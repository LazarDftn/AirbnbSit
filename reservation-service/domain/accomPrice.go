package domain

type AccommPrice struct {
	AccommID string `json:"accommId"`
	Price    int    `json:"price"`
	PayPer   int    `bson:"payPer" json:"payPer"` //nacin placanja, 0=ukupno, 1=po osobi
}
