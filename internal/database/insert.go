package database

import (
	"log"
	"module/internal/model"

	"github.com/jmoiron/sqlx"
)

func DbInsertOrder(db *sqlx.DB, order model.OrderData) (orderID string, err error) {
	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	deliveryID, err := InsertDelivery(tx, order.Delivery)
	if err != nil {
		tx.Rollback()
		log.Printf("Error inserting the delivery: %v\n", err)
		return
	}

	paymentID, err := InsertPayment(tx, order.Payment)
	if err != nil {
		tx.Rollback()
		log.Printf("Error when inserting a payment: %v\n", err)
		return
	}

	orderID, err = InsertOrder(tx, order, deliveryID, paymentID)
	if err != nil {
		return "", err
	}

	for _, item := range order.Items {
		itemID, err := InsertItem(tx, item, orderID)
		if err != nil {
			return "", err
		}

		orderItemData := model.ItemOrderData{
			ChrtID:   itemID,
			OrderUID: order.OrderUID,
		}

		_, err = InsertItemOrder(tx, orderItemData)
		if err != nil {

			log.Printf("Error linking an order item (ChrtID: %d) with the order (OrderUID: '%s'): %v", itemID, order.OrderUID, err)
			return "", err
		}
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return orderID, nil
}

func InsertDelivery(tx *sqlx.Tx, delivery model.DeliveryData) (int64, error) {
	var deliveryID int64

	err := tx.QueryRow(`
        INSERT INTO deliveries (name, phone, zip, city, address, region, email)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING delivery_id
    `, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email).Scan(&deliveryID)

	if err != nil {
		return 0, err
	}

	return deliveryID, nil
}

func InsertPayment(tx *sqlx.Tx, payment model.PaymentData) (int64, error) {
	var paymentID int64

	err := tx.QueryRow(`
        INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING payments_id
    `, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee).Scan(&paymentID)

	if err != nil {
		return 0, err
	}

	return paymentID, nil
}

func InsertOrder(tx *sqlx.Tx, order model.OrderData, deliveryID, paymentID int64) (string, error) {
	var orderID string
	err := tx.QueryRow(`
        INSERT INTO orders (order_uid, track_number, entry, delivery, payment, locale,
            internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING order_uid
    `, order.OrderUID, order.TrackNumber, order.Entry, deliveryID, paymentID, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey,
		order.SmID, order.DateCreated, order.OofShard).Scan(&orderID)

	if err != nil {
		return "", err
	}

	return orderID, nil
}

func InsertItem(tx *sqlx.Tx, item model.ItemData, orderID string) (int64, error) {
	var itemID int64

	err := tx.QueryRow(`
        INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING chrt_id
    `, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status).Scan(&itemID)

	if err != nil {
		return 0, err
	}

	return itemID, nil
}
func InsertItemOrder(tx *sqlx.Tx, orderitem model.ItemOrderData) (int64, error) {
	var chrtID int64

	err := tx.QueryRow(`
        INSERT INTO orders_items (item_id, order_id)
        VALUES ($1, $2)
        RETURNING item_id
    `, orderitem.ChrtID, orderitem.OrderUID).Scan(&chrtID)

	if err != nil {

		log.Printf("Couldn't insert (ChrtID: %d, OrderUID: '%s'): %v\n", orderitem.ChrtID, orderitem.OrderUID, err)
		return 0, err
	}

	return chrtID, nil
}
