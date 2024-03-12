package domain

import "time"

type PaymentCache struct {
	ID           int
	Transaction  string
	RequestID    string
	Currency     string
	Provider     string
	Amount       int
	PaymentDt    int
	Bank         string
	DeliveryCost int
	GoodsTotal   int
	CustomFee    int
}

type ItemCache struct {
	ID          int
	TrackNumber string
	Price       int
	RID         string
	Name        string
	Sale        int
	Size        string
	TotalPrice  int
	NmID        int
	Brand       string
	Status      int
}

type OrderCache struct {
	UID               string
	TrackNumber       string
	Entry             string
	Delivery          int
	Payment           int
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	ShardKey          string
	SmID              int
	DateCreated       time.Time
	OofShard          string
}

type OrderItemCache struct {
	OrderID string
	ItemID  int
}

type DeliveryCache struct {
	ID      int
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}
