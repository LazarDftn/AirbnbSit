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
					(accomm_id text, reservation_id UUID, guest_email text, host_email text, price int, num_of_People int, start_date timestamp, end_date timestamp, 
					PRIMARY KEY ((accomm_id), reservation_id, guest_email)) 
					WITH CLUSTERING ORDER BY (reservation_id ASC, guest_email ASC)`,
			"reservation_by_accommodation")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	// table for reservations that are relevant to a guest (he wants to see all of the reservations he made)
	err = rr.session.Query(
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
	}

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
					(variation_id UUID, accomm_id text, percentage int, start_date timestamp, end_date timestamp, 
					PRIMARY KEY ((accomm_id), start_date, end_date, variation_id)) 
					WITH CLUSTERING ORDER BY (start_date DESC, end_date DESC, variation_id ASC)`,
			"variation_by_accomm_and_interval")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	// table for hosting availability of an accommodation (host wants can't host an accommodation at certain time periods)
	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(availability_id UUID, accomm_id text, start_date timestamp, end_date timestamp, 
					PRIMARY KEY ((accomm_id), start_date, end_date, availability_id)) 
					WITH CLUSTERING ORDER BY (start_date DESC, end_date DESC, availability_id ASC)`,
			"availability_by_accomm_and_interval")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

}

func (rr *ReservationRepo) InsertAvailability(availability *domain.Availability) error {
	availability.AvailabilityID, _ = gocql.RandomUUID()
	err := rr.session.Query(
		`INSERT INTO availability_by_accomm_and_interval (accomm_id, start_date, end_date, availability_id) 
		VALUES (?, ?, ?)`,
		availability.AccommID, availability.StartDate, availability.EndDate, availability.AvailabilityID).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}
	return nil
}

func (rr *ReservationRepo) CheckExistingForAvailability(availability *domain.Availability) {
	/* method that will check if there is existing reservations or Host unavailabilities
	   when he tries to set a period of accommodation unavailability */
}

func (rr *ReservationRepo) InsertPriceVariation(variation *domain.PriceVariation) error {
	variation.VariationID, _ = gocql.RandomUUID()
	err := rr.session.Query(
		`INSERT INTO variation_by_accomm_and_interval (accomm_id, start_date, end_date, variation_id) 
		VALUES (?, ?, ?, ?)`,
		variation.AccommID, variation.StartDate, variation.EndDate, variation.VariationID).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}
	return nil
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

	// there is no OR statement in cql so the only way I found to check the time overlap is to make 4 queries
	scanner := rr.session.Query(`SELECT * FROM reservation_by_accommodation WHERE accomm_id = ? AND start_date >= ? AND end_date <= ?
	ALLOW FILTERING`,
		reservation.AccommID).Iter().Scanner()

	for scanner.Next() {
		return "That period is occupied"
	}

	scanner = rr.session.Query(`SELECT * FROM reservation_by_accommodation WHERE accomm_id = ? AND start_date <= ? AND end_date >= ?
	ALLOW FILTERING`,
		reservation.AccommID, reservation.StartDate, reservation.EndDate).Iter().Scanner()

	for scanner.Next() {
		return "That period is occupied"
	}

	scanner = rr.session.Query(`SELECT * FROM reservation_by_accommodation WHERE accomm_id = ? AND start_date >= ? AND start_date <= ?
	ALLOW FILTERING`,
		reservation.AccommID, reservation.StartDate, reservation.EndDate).Iter().Scanner()

	for scanner.Next() {
		return "That period is occupied"
	}

	scanner = rr.session.Query(`SELECT * FROM reservation_by_accommodation WHERE accomm_id = ? AND end_date >= ? AND end_date <= ?
	ALLOW FILTERING`,
		reservation.AccommID, reservation.StartDate, reservation.EndDate).Iter().Scanner()

	for scanner.Next() {
		return "That period is occupied"
	}

	scanner2 := rr.session.Query(`SELECT * FROM price_by_accommodation WHERE accomm_id = ?`,
		reservation.AccommID).Iter().Scanner()

	var prices []domain.AccommPrice
	for scanner2.Next() {
		var price domain.AccommPrice
		err := scanner2.Scan(&price.AccommID, &price.PayPer, &price.Price)
		if err != nil {
			rr.logger.Println(err)
			return "error"
		}
		prices = append(prices, price)
	}

	var priceNum int

	if prices[0].PayPer == "per guest" {
		// check if accommodation payment is per person, then multiply the price with the number of people in the reservation
		priceNum = prices[0].Price * reservation.NumOfPeople
	} else {
		priceNum = prices[0].Price
	}

	reservation.ReservationID, _ = gocql.RandomUUID()
	err := rr.session.Query(
		`INSERT INTO reservation_by_accommodation (accomm_id, reservation_id, start_date, end_date, num_of_people,
			guest_email, host_email, price) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		reservation.AccommID, reservation.ReservationID, reservation.StartDate, reservation.EndDate,
		reservation.NumOfPeople, reservation.GuestEmail, reservation.HostEmail, priceNum).Exec()
	if err != nil {
		rr.logger.Println(err)
		return "error"
	}
	return ""
}

func (rr *ReservationRepo) GetReservationsByAccomm(id string) (domain.Reservations, error) {

	scanner := rr.session.Query(`SELECT accomm_id, reservation_id, guest_email, host_email, price, 
	num_of_People, start_date, end_date FROM reservation_by_accommodation WHERE accomm_id = ?`,
		id).Iter().Scanner()

	var foundReservations domain.Reservations
	for scanner.Next() {
		var res domain.Reservation
		err := scanner.Scan(&res.AccommID, &res.ReservationID, &res.GuestEmail, &res.HostEmail, &res.Price,
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
