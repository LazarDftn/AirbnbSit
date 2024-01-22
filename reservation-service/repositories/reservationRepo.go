package repositories

import (
	"context"
	"fmt"
	"log"
	"os"
	"reservation-service/domain"
	"time"

	// NoSQL: module containing Cassandra and Mongo api clients
	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NoSQL: ReservationRepo struct encapsulating Cassandra and Mongo api client
type ReservationRepo struct {
	session *gocql.Session
	logger  *log.Logger
	cli     *mongo.Client
}

// NoSQL: Constructor which reads db configuration from environment and creates a keyspace
func New(logger *log.Logger, ctx context.Context) (*ReservationRepo, error) {
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

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://reservations_mongo_db:27017/"))
	if err != nil {
		return nil, err
	}

	// Return repository with logger and DB session
	return &ReservationRepo{
		session: session,
		logger:  logger,
		cli:     client,
	}, nil
}

// Disconnect from database
func (rr *ReservationRepo) CloseSession(ctx context.Context) {
	rr.session.Close()
	rr.cli.Disconnect(ctx)
}

func (rr *ReservationRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := rr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		rr.logger.Println(err)
	}

	// Print available databases
	databases, err := rr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		rr.logger.Println(err)
	}
	fmt.Println(databases)
}

// Create tables
func (rr *ReservationRepo) CreateTables() {

	// table for reservations that are relevant to an accommodation
	err := rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(location text, accomm_id text, reservation_id UUID, guest_email text, host_email text, price int, 
					num_of_People int, start_date timestamp, end_date timestamp, 
					PRIMARY KEY ((location), accomm_id, end_date, reservation_id)) 
					WITH CLUSTERING ORDER BY (accomm_id ASC, end_date DESC, reservation_id ASC)`,
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
					PRIMARY KEY ((location), accomm_id, end_date, start_date, variation_id)) 
					WITH CLUSTERING ORDER BY (accomm_id ASC, end_date DESC, start_date DESC, variation_id ASC)`,
			"variation_by_accomm_and_interval")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	// table for hosting availability of an accommodation (host wants can't host an accommodation at certain time periods)
	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(availability_id UUID, accomm_id text, start_date timestamp, end_date timestamp,
					name text, location text, min_capacity int, max_capacity int, 
					PRIMARY KEY ((location), accomm_id, end_date, start_date)) 
					WITH CLUSTERING ORDER BY (accomm_id ASC, end_date DESC, start_date DESC)`,
			"availability_by_accomm")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(availability_id UUID, accomm_id text, start_date timestamp, end_date timestamp,
					name text, location text, min_capacity int, max_capacity int, 
					PRIMARY KEY ((location), end_date, start_date)) 
					WITH CLUSTERING ORDER BY (end_date DESC, start_date DESC)`,
			"availability_by_search_parameters")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

}

func (rr *ReservationRepo) InsertAvailability(availability *domain.Availability) (string, error) {

	scanner := rr.session.Query(`SELECT start_date, end_date FROM availability_by_accomm WHERE location = ? AND accomm_id = ?
	 AND end_date >= ?`,
		availability.Location, availability.AccommID, time.Now()).Iter().Scanner()

	for scanner.Next() {
		var av domain.Availability
		err := scanner.Scan(&av.StartDate, &av.EndDate)
		if err != nil {
			rr.logger.Println(err)
		}
		if av.StartDate.Before(availability.EndDate) && av.EndDate.After(availability.StartDate) {
			return "Can't change availability because there is another availability during this period", nil
		}
	}

	scanner = rr.session.Query(`SELECT start_date, end_date FROM reservation_by_accommodation WHERE location = ? AND accomm_id = ?
	 AND end_date >= ?`,
		availability.Location, availability.AccommID, time.Now()).Iter().Scanner()

	for scanner.Next() {
		var res domain.Reservation
		err := scanner.Scan(&res.StartDate, &res.EndDate)
		if err != nil {
			rr.logger.Println(err)
		}
		if res.StartDate.Before(availability.EndDate) && res.EndDate.After(availability.StartDate) {
			return "Can't change availability because there are reservations during this period", nil
		}
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

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	avCollection := rr.getAvCollection()

	result, err := avCollection.InsertOne(ctx, &availability)
	if err != nil {
		rr.logger.Println(err)
		return "", err
	}
	rr.logger.Println(result.InsertedID)

	return "Changed", nil
}

func (rr *ReservationRepo) DeleteAvailability(availability *domain.Availability) string {

	avCollection := rr.getAvCollection()

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	_, err := avCollection.DeleteOne(ctx, bson.M{"accommId": availability.AccommID, "startDate": availability.StartDate})

	if err != nil {
		rr.logger.Println(err)
		return "Database error"
	}

	scanner := rr.session.Query(`SELECT start_date, end_date FROM reservation_by_accommodation WHERE location = ? AND accomm_id = ?
	 AND end_date >= ?`,
		availability.Location, availability.AccommID, time.Now()).Iter().Scanner()

	for scanner.Next() {
		var res domain.Reservation
		err := scanner.Scan(&res.StartDate, &res.EndDate)
		if err != nil {
			rr.logger.Println(err)
		}
		if res.StartDate.Before(availability.EndDate) && res.EndDate.After(availability.StartDate) {
			return "Can't change availability because there are reservations during this period"
		}
	}

	err = rr.session.Query(
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

	isAvailable := false

	scanner := rr.session.Query(`SELECT start_date, end_date FROM availability_by_accomm WHERE location = ? AND 
	accomm_id = ? AND end_date >= ?`,
		variation.Location, variation.AccommID, time.Now()).Iter().Scanner()

	for scanner.Next() {
		var av domain.Availability
		err := scanner.Scan(&av.StartDate, &av.EndDate)
		if err != nil {
			rr.logger.Println(err)
		}
		if av.StartDate.Before(variation.StartDate.AddDate(0, 0, 1)) && av.EndDate.After(variation.EndDate.AddDate(0, 0, -1)) {
			isAvailable = true
			break
		}
	}

	if !isAvailable {
		return "Can't change price because the accommodation is unavailable during this period", nil
	}

	scanner = rr.session.Query(`SELECT start_date, end_date FROM reservation_by_accommodation WHERE location = ? AND accomm_id = ?
	 AND end_date >= ?`,
		variation.Location, variation.AccommID, time.Now()).Iter().Scanner()

	for scanner.Next() {
		var res domain.Reservation
		err := scanner.Scan(&res.StartDate, &res.EndDate)
		if err != nil {
			rr.logger.Println(err)
		}
		if res.StartDate.Before(variation.EndDate) && res.EndDate.After(variation.StartDate) {
			return "Can't change price because there are reservations during this period", nil
		}
	}

	scanner = rr.session.Query(`SELECT start_date, end_date FROM variation_by_accomm_and_interval WHERE location = ? AND accomm_id = ?
	 AND end_date >= ?`,
		variation.Location, variation.AccommID, time.Now()).Iter().Scanner()

	for scanner.Next() {
		var pv domain.PriceVariation
		err := scanner.Scan(&pv.StartDate, &pv.EndDate)
		if err != nil {
			rr.logger.Println(err)
		}
		if pv.StartDate.Before(variation.EndDate) && pv.EndDate.After(variation.StartDate) {
			return "Can't change price because there is another price change during this period", nil
		}
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
	WHERE location = ? AND accomm_id = ? AND end_date >= ?`,
		location, id, time.Now()).Iter().Scanner()

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

	scanner := rr.session.Query(`SELECT start_date, end_date FROM reservation_by_accommodation WHERE location = ? AND accomm_id = ?
	 AND end_date >= ?`,
		pv.Location, pv.AccommID, time.Now()).Iter().Scanner()

	for scanner.Next() {
		var res domain.Reservation
		err := scanner.Scan(&res.StartDate, &res.EndDate)
		if err != nil {
			rr.logger.Println(err)
		}
		if res.StartDate.Before(pv.EndDate) && res.EndDate.After(pv.StartDate) {
			return "Can't change price because there are reservations during this period"
		}
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

	isAvailable := false

	scanner := rr.session.Query(`SELECT start_date, end_date FROM availability_by_accomm WHERE location = ? AND 
	accomm_id = ? AND end_date >= ?`,
		reservation.Location, reservation.AccommID, time.Now()).Iter().Scanner()

	for scanner.Next() {
		var av domain.Availability
		err := scanner.Scan(&av.StartDate, &av.EndDate)
		if err != nil {
			rr.logger.Println(err)
		}
		if av.StartDate.Before(reservation.StartDate.AddDate(0, 0, 1)) && av.EndDate.After(reservation.EndDate.AddDate(0, 0, -1)) {
			isAvailable = true
			break
		}
	}

	if !isAvailable {
		return "Accommodation is unavailable during this period"
	}

	scanner = rr.session.Query(`SELECT start_date, end_date FROM reservation_by_accommodation WHERE location = ? AND accomm_id = ?
	 AND end_date >= ?`,
		reservation.Location, reservation.AccommID, time.Now()).Iter().Scanner()

	for scanner.Next() {
		var res domain.Reservation
		err := scanner.Scan(&res.StartDate, &res.EndDate)
		if err != nil {
			rr.logger.Println(err)
		}
		if res.StartDate.Before(reservation.EndDate) && res.EndDate.After(reservation.StartDate) {
			return "That period is occupied"
		}
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
		return err.Error()
	}
	return ""
}

func (rr *ReservationRepo) GetReservationsByAccomm(location string, id string) (domain.Reservations, error) {

	scanner := rr.session.Query(`SELECT location, accomm_id, reservation_id, guest_email, host_email, price, 
	num_of_People, start_date, end_date FROM reservation_by_accommodation WHERE location = ? AND accomm_id = ?`,
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
	WHERE location = ? AND accomm_id = ? AND end_date >= ?`,
		res.Location, res.AccommID, time.Now()).Iter().Scanner()

	for scanner.Next() {
		var variation domain.PriceVariation
		err := scanner.Scan(&variation.AccommID, &variation.Percentage, &variation.StartDate, &variation.EndDate)
		if err != nil {
			rr.logger.Println(err)
			//return nil, err
		}
		if variation.StartDate.Before(res.EndDate) && variation.EndDate.After(res.StartDate) {
			variations = append(variations, variation)
		}
	}

	if len(variations) > 0 {
		return variations
	} else {
		return nil
	}
}

func (rr *ReservationRepo) getAvCollection() *mongo.Collection {
	avDatabase := rr.cli.Database("mongoDemo")
	avCollection := avDatabase.Collection("availabilities")
	return avCollection
}

func (rr *ReservationRepo) SearchAccommodations(av domain.Availability) ([]domain.Availability, error) {

	var availabilites []domain.Availability

	/* filter pretrage je po defaultu filter sa svim parametrima a ispod
	   su if slucajevi kada korisnik u pretragu ne unese neki od parametara
	   (ili nijedan od parametara osim perioda bukiranja koji je obavezan)
	*/
	filter := bson.M{
		"location":    av.Location,
		"minCapacity": bson.M{"$lte": av.MinCapacity},
		"maxCapacity": bson.M{"$gte": av.MinCapacity},
		"startDate":   bson.M{"$lte": av.StartDate},
		"endDate":     bson.M{"$gte": av.EndDate},
	}

	if av.Location == "" && av.MinCapacity != 0 {

		filter = bson.M{
			"minCapacity": bson.M{"$lte": av.MinCapacity},
			"maxCapacity": bson.M{"$gte": av.MinCapacity},
			"startDate":   bson.M{"$lte": av.StartDate},
			"endDate":     bson.M{"$gte": av.EndDate},
		}
	}

	if av.MinCapacity == 0 && av.Location != "" {

		filter = bson.M{
			"location":  av.Location,
			"startDate": bson.M{"$lte": av.StartDate},
			"endDate":   bson.M{"$gte": av.EndDate},
		}
	}

	if av.Location == "" && av.MinCapacity == 0 {

		filter = bson.M{
			"startDate": bson.M{"$lte": av.StartDate},
			"endDate":   bson.M{"$gte": av.EndDate},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := rr.getAvCollection()

	coll, err := collection.Find(ctx, filter)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err := coll.All(ctx, &availabilites); err != nil {
		log.Fatal(err)
		return nil, err
	}

	if len(availabilites) > 0 {
		return availabilites, nil
	} else {
		return nil, nil
	}
}
