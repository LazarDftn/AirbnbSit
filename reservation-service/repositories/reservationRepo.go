package repositories

import (
	"fmt"
	"log"
	"os"
	"reservation-service/domain"

	// NoSQL: module containing Cassandra api client
	"github.com/gocql/gocql"
)

// NoSQL: ReservationRepo struct encapsulating Cassandra api client
type ReservationRepo struct {
	session *gocql.Session
	logger  *log.Logger
}

// NoSQL: Constructor which reads db configuration from environment and creates a keyspace
func New(logger *log.Logger) (*ReservationRepo, error) {
	db := os.Getenv("CASS_DB")

	// Connect to default keyspace
	cluster := gocql.NewCluster(db)
	cluster.Keyspace = "system"
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	// Create 'reservation' keyspace
	err = session.Query(
		fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
					WITH replication = {
						'class' : 'SimpleStrategy',
						'replication_factor' : %d
					}`, "reservation", 1)).Exec()
	if err != nil {
		logger.Println(err)
	}
	session.Close()

	// Connect to reservation keyspace
	cluster.Keyspace = "reservation"
	cluster.Consistency = gocql.One
	session, err = cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	// Return repository with logger and DB session
	return &ReservationRepo{
		session: session,
		logger:  logger,
	}, nil
}

// Disconnect from database
func (rr *ReservationRepo) CloseSession() {
	rr.session.Close()
}

// Create tables
func (rr *ReservationRepo) CreateTables() {

	// table for reservations that are relevant to an accommodation
	err := rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(location text, accomm_id text, reservation_id UUID, guest_email text, host_email text, price int, 
					num_of_People int, start_date timestamp, end_date timestamp, 
					PRIMARY KEY ((accomm_id), guest_email, reservation_id)) 
					WITH CLUSTERING ORDER BY (guest_email ASC, reservation_id ASC)`,
			"reservation_by_accommodation")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	// table for reservations that are relevant to a guest (he wants to see all of the reservations he made)
	/*err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
					(guest_email text, reservation_id UUID, accomm_id text, price int, num_of_People int, start_date timestamp, end_date timestamp,
					PRIMARY KEY ((guest_email), reservation_id, accomm_id))
					WITH CLUSTERING ORDER BY (reservation_id ASC, accomm_id ASC)`,
			"reservation_by_guest")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	// table for reservations that are relevant to a host (host wants to see reservations of all his accommodations)
	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
					(host_email text, reservation_id UUID, accomm_id text, price int, num_of_People int, start_date timestamp, end_date timestamp,
					PRIMARY KEY ((host_email), reservation_id, accomm_id))
					WITH CLUSTERING ORDER BY (reservation_id ASC, accomm_id ASC)`,
			"reservation_by_host")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}*/

	/* Additional information about the price of an accommodation is stored here.
	   The reason why it is stored in the Reservation database (and service)
	   is because the price can vary at different dates depending on when guests want to
	   reserve. Instead of sending the info from the client side about pricing every time
	   the guest makes a reservation, it's better to have the info available in this service.
	*/
	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(accomm_id text PRIMARY KEY, price int, pay_per text)`,
			"price_by_accommodation")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	// table for price variations of an accommodation (host wants to have different pricings for different periods)
	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(variation_id UUID, location text, accomm_id text, percentage int, start_date timestamp, end_date timestamp, 
					PRIMARY KEY ((location), accomm_id, start_date, end_date, variation_id)) 
					WITH CLUSTERING ORDER BY (accomm_id ASC, start_date DESC, end_date DESC, variation_id ASC)`,
			"variation_by_accomm_and_interval")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	// table for hosting availability of an accommodation (host wants can't host an accommodation at certain time periods)
	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(availability_id UUID, accomm_id text, start_date timestamp, end_date timestamp,
					name text, location text, min_capacity int, max_capacity int, 
					PRIMARY KEY ((location), accomm_id, start_date, end_date)) 
					WITH CLUSTERING ORDER BY (accomm_id ASC, start_date DESC, end_date DESC)`,
			"availability_by_accomm")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

}

func (rr *ReservationRepo) InsertAvailability(availability *domain.Availability) (string, error) {

	scanner := rr.session.Query(`SELECT * FROM availability_by_accomm WHERE location = ? AND accomm_id = ? 
	AND start_date <= ? AND end_date >= ?
	ALLOW FILTERING`,
		availability.Location, availability.AccommID, availability.EndDate, availability.StartDate).Iter().Scanner()

	for scanner.Next() {
		return "Can't change availability because there is another availability during this period", nil
	}

	scanner = rr.session.Query(`SELECT * FROM reservation_by_accommodation WHERE location = ? AND accomm_id = ? 
	AND start_date <= ? AND end_date >= ?
	ALLOW FILTERING`,
		availability.Location, availability.AccommID, availability.EndDate, availability.StartDate).Iter().Scanner()

	for scanner.Next() {
		return "Can't change availability because there are reservations during this period", nil
	}

	availability.AvailabilityID, _ = gocql.RandomUUID()
	err := rr.session.Query(
		`INSERT INTO availability_by_accomm (location, accomm_id, start_date, end_date, name, min_capacity, max_capacity, availability_id) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		availability.Location, availability.AccommID, availability.StartDate, availability.EndDate, availability.Name,
		availability.MinCapacity, availability.MaxCapacity, availability.AvailabilityID).Exec()
	if err != nil {
		rr.logger.Println(err)
		return "", err
	}
	return "Changed", nil
}

func (rr *ReservationRepo) DeleteAvailability(availability *domain.Availability) string {

	scanner := rr.session.Query(`SELECT * FROM reservation_by_accommodation WHERE location = ? AND accomm_id = ?
	 AND start_date <= ? AND end_date >= ?
	ALLOW FILTERING`,
		availability.Location, availability.AccommID, availability.EndDate, availability.StartDate).Iter().Scanner()

	for scanner.Next() {
		return "Can't change availability because there are reservations during this period"
	}

	err := rr.session.Query(
		`DELETE FROM availability_by_accomm WHERE location = ? AND accomm_id = ? AND start_date = ? AND end_date = ?`,
		availability.Location, availability.AccommID, availability.StartDate, availability.EndDate).Exec()
	if err != nil {
		rr.logger.Println(err)
		return "Database error"
	}
	return ""
}

func (rr *ReservationRepo) GetAvailability(location string, id string) ([]*domain.Availability, error) {

	scanner := rr.session.Query(`SELECT availability_id, accomm_id, start_date, end_date,
	name, location, min_capacity, max_capacity
	FROM availability_by_accomm WHERE location = ? AND accomm_id = ?`,
		location, id).Iter().Scanner()

	var availabilities []*domain.Availability
	for scanner.Next() {
		var av domain.Availability
		err := scanner.Scan(&av.AvailabilityID, &av.AccommID, &av.StartDate, &av.EndDate,
			&av.Name, &av.Location, &av.MinCapacity, &av.MaxCapacity)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		availabilities = append(availabilities, &av)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return availabilities, nil
}

func (rr *ReservationRepo) InsertPriceVariation(variation *domain.PriceVariation) (string, error) {

	var foundAvailability = false

	scanner := rr.session.Query(`SELECT start_date, end_date FROM availability_by_accomm WHERE location = ? AND 
	accomm_id = ? AND start_date <= ? AND end_date >= ?
	ALLOW FILTERING`,
		variation.Location, variation.AccommID, variation.StartDate, variation.EndDate).Iter().Scanner()

	for scanner.Next() {
		foundAvailability = true
	}

	if !foundAvailability {
		return "Can't change price because the accommodation is unavailable during this period", nil
	}

	scanner = rr.session.Query(`SELECT * FROM reservation_by_accommodation WHERE location = ? AND accomm_id = ? 
	AND start_date <= ? AND end_date >= ?
	ALLOW FILTERING`,
		variation.AccommID, variation.EndDate, variation.StartDate).Iter().Scanner()

	for scanner.Next() {
		return "Can't change price because there are reservations during this period", nil
	}

	scanner = rr.session.Query(`SELECT * FROM variation_by_accomm_and_interval WHERE location = ? AND accomm_id = ? 
	AND start_date <= ? AND end_date >= ?
	ALLOW FILTERING`,
		variation.AccommID, variation.EndDate, variation.StartDate).Iter().Scanner()

	for scanner.Next() {
		return "Can't change price because there is another price change during this period", nil
	}

	variation.VariationID, _ = gocql.RandomUUID()

	err := rr.session.Query(
		`INSERT INTO variation_by_accomm_and_interval (location, accomm_id, percentage, start_date, end_date, variation_id) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		variation.Location, variation.AccommID, variation.Percentage, variation.StartDate, variation.EndDate, variation.VariationID).Exec()
	if err != nil {
		rr.logger.Println(err)
		return "database error", err
	}
	return "Created", nil
}

func (rr *ReservationRepo) GetVariationsByAccommId(location string, id string) ([]domain.PriceVariation, error) {
	scanner := rr.session.Query(`SELECT variation_id, location, accomm_id, percentage, start_date, end_date FROM variation_by_accomm_and_interval 
	WHERE location = ? AND accomm_id = ?`,
		location, id).Iter().Scanner()

	var foundVariations []domain.PriceVariation
	for scanner.Next() {
		var pv domain.PriceVariation
		err := scanner.Scan(&pv.VariationID, &pv.Location, &pv.AccommID, &pv.Percentage, &pv.StartDate, &pv.EndDate)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		foundVariations = append(foundVariations, pv)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return foundVariations, nil
}

func (rr *ReservationRepo) DeletePriceVariation(pv *domain.PriceVariation) string {

	scanner := rr.session.Query(`SELECT * FROM reservation_by_accommodation WHERE location = ? AND accomm_id = ?
	 AND start_date <= ? AND end_date >= ?
	ALLOW FILTERING`,
		pv.Location, pv.AccommID, pv.EndDate, pv.StartDate).Iter().Scanner()

	for scanner.Next() {
		return "Can't change price because there are reservations during this period"
	}

	err := rr.session.Query(
		`DELETE FROM variation_by_accomm_and_interval WHERE location = ? AND accomm_id = ? AND start_date = ? AND end_date = ?`,
		pv.Location, pv.AccommID, pv.StartDate, pv.EndDate).Exec()
	if err != nil {
		rr.logger.Println(err)
		return "Database error"
	}
	return ""
}

func (rr *ReservationRepo) InsertAccommodationPrice(price *domain.AccommPrice) error {
	err := rr.session.Query(
		`INSERT INTO price_by_accommodation (accomm_id, price, pay_per) 
		VALUES (?, ?, ?)`,
		price.AccommID, price.Price, price.PayPer).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}
	return nil
}

func (rr *ReservationRepo) GetPriceByAccomm(id string) (*domain.AccommPrice, error) {

	scanner := rr.session.Query(`SELECT accomm_id, price, pay_per FROM price_by_accommodation WHERE accomm_id = ?`,
		id).Iter().Scanner()

	var foundPrices []domain.AccommPrice
	for scanner.Next() {
		var price domain.AccommPrice
		err := scanner.Scan(&price.AccommID, &price.Price, &price.PayPer)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		foundPrices = append(foundPrices, price)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}

	if len(foundPrices) == 0 {
		var price *domain.AccommPrice
		return price, nil
	}

	return &foundPrices[0], nil
}

// create the reservationID in the handler before calling this method
func (rr *ReservationRepo) InsertReservation(reservation *domain.Reservation) string {

	var isAvailable = false

	scanner := rr.session.Query(`SELECT start_date, end_date FROM availability_by_accomm WHERE location = ? AND 
	accomm_id = ? AND start_date <= ? AND end_date >= ?
	ALLOW FILTERING`,
		reservation.Location, reservation.AccommID, reservation.StartDate, reservation.EndDate).Iter().Scanner()

	for scanner.Next() {
		isAvailable = true
	}

	if !isAvailable {
		return "Accommodation is unavailable during this period"
	}

	scanner = rr.session.Query(`SELECT * FROM reservation_by_accommodation WHERE location = ? AND accomm_id = ? 
	AND start_date <= ? AND end_date >= ?
	ALLOW FILTERING`,
		reservation.Location, reservation.AccommID, reservation.EndDate, reservation.StartDate).Iter().Scanner()

	for scanner.Next() {
		return "That period is occupied"
	}

	reservation.ReservationID, _ = gocql.RandomUUID()
	err := rr.session.Query(
		`INSERT INTO reservation_by_accommodation (location, accomm_id, reservation_id, start_date, end_date, num_of_people,
			guest_email, host_email, price) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		reservation.Location, reservation.AccommID, reservation.ReservationID, reservation.StartDate, reservation.EndDate,
		reservation.NumOfPeople, reservation.GuestEmail, reservation.HostEmail, reservation.Price).Exec()
	if err != nil {
		rr.logger.Println(err)
		return "error"
	}
	return ""
}

func (rr *ReservationRepo) GetReservationsByAccomm(location string, id string) (domain.Reservations, error) {

	scanner := rr.session.Query(`SELECT location, accomm_id, reservation_id, guest_email, host_email, price, 
	num_of_People, start_date, end_date FROM reservation_by_accommodation WHERE location = ? AND accomm_id = ? 
	ALLOW FILTERING`,
		location, id).Iter().Scanner()

	var foundReservations domain.Reservations
	for scanner.Next() {
		var res domain.Reservation
		err := scanner.Scan(&res.Location, &res.AccommID, &res.ReservationID, &res.GuestEmail, &res.HostEmail, &res.Price,
			&res.NumOfPeople, &res.StartDate, &res.EndDate)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		foundReservations = append(foundReservations, &res)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return foundReservations, nil
}

func (rr *ReservationRepo) CheckPrice(res domain.Reservation) []domain.PriceVariation {

	var variations []domain.PriceVariation

	scanner := rr.session.Query(`SELECT accomm_id, percentage, start_date, end_date FROM variation_by_accomm_and_interval 
	WHERE location = ? AND accomm_id = ? AND start_date <= ? AND end_date >= ? 
	ALLOW FILTERING`,
		res.Location, res.AccommID, res.EndDate, res.StartDate).Iter().Scanner()

	for scanner.Next() {
		var variation domain.PriceVariation
		err := scanner.Scan(&variation.AccommID, &variation.Percentage, &variation.StartDate, &variation.EndDate)
		if err != nil {
			rr.logger.Println(err)
			//return nil, err
		}
		variations = append(variations, variation)
	}

	if len(variations) > 0 {
		return variations
	} else {
		return nil
	}
}
