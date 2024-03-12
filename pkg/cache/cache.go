package cache

import (
	"fmt"
	"log"
	"module/internal/domain"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
)

func loadDeliveries(db *sqlx.DB, c *cache.Cache) error {
	query := `SELECT "delivery_id", "name", "phone", "zip", "city", "address", "region", "email" FROM deliveries`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var d domain.DeliveryCache
		err := rows.Scan(&d.ID, &d.Name, &d.Phone, &d.Zip, &d.City, &d.Address, &d.Region, &d.Email)
		if err != nil {
			return err
		}
		c.Set(fmt.Sprintf("delivery_%d", d.ID), d, cache.DefaultExpiration)
	}

	return rows.Err()
}
func loadPayments(db *sqlx.DB, c *cache.Cache) error {
	query := `SELECT "payments_id", "transaction", "request_id", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee" FROM payments`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.PaymentCache
		err := rows.Scan(&p.ID, &p.Transaction, &p.RequestID, &p.Currency, &p.Provider, &p.Amount, &p.PaymentDt, &p.Bank, &p.DeliveryCost, &p.GoodsTotal, &p.CustomFee)
		if err != nil {
			return err
		}
		c.Set(fmt.Sprintf("payment_%d", p.ID), p, cache.DefaultExpiration)
	}

	return rows.Err()
}

func loadItems(db *sqlx.DB, c *cache.Cache) error {
	query := `SELECT "chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status" FROM items`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var i domain.ItemCache
		err := rows.Scan(&i.ID, &i.TrackNumber, &i.Price, &i.RID, &i.Name, &i.Sale, &i.Size, &i.TotalPrice, &i.NmID, &i.Brand, &i.Status)
		if err != nil {
			return err
		}
		c.Set(fmt.Sprintf("item_%d", i.ID), i, cache.DefaultExpiration)
	}

	return rows.Err()
}

func loadOrders(db *sqlx.DB, c *cache.Cache) error {
	query := `SELECT "order_uid", "track_number", "entry", "delivery", "payment", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard" FROM orders`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var o domain.OrderCache
		var dateCreatedStr string
		err := rows.Scan(&o.UID, &o.TrackNumber, &o.Entry, &o.Delivery, &o.Payment, &o.Locale, &o.InternalSignature, &o.CustomerID, &o.DeliveryService, &o.ShardKey, &o.SmID, &dateCreatedStr, &o.OofShard)
		if err != nil {
			return err
		}
		o.DateCreated, err = time.Parse(time.RFC3339, dateCreatedStr)
		if err != nil {
			return err
		}
		c.Set(fmt.Sprintf("order_%s", o.UID), o, cache.DefaultExpiration)
	}

	return rows.Err()
}

func loadOrderItems(db *sqlx.DB, c *cache.Cache) error {
	query := `SELECT "order_id", "item_id" FROM orders_items`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var oi domain.OrderItemCache
		err := rows.Scan(&oi.OrderID, &oi.ItemID)
		if err != nil {
			return err
		}
		c.Set(fmt.Sprintf("order_item_%s_%d", oi.OrderID, oi.ItemID), oi, cache.DefaultExpiration)
	}

	return rows.Err()
}
func RefreshCache(db *sqlx.DB, c *cache.Cache) error {
	if err := loadDeliveries(db, c); err != nil {
		log.Printf("Error loading deliveries into cache: %v", err)
		return err
	}

	if err := loadPayments(db, c); err != nil {
		log.Printf("Error loading payments into cache: %v", err)
		return err
	}

	if err := loadItems(db, c); err != nil {
		log.Printf("Error loading items into cache: %v", err)
		return err
	}

	if err := loadOrders(db, c); err != nil {
		log.Printf("Error loading orders into cache: %v", err)
		return err
	}

	if err := loadOrderItems(db, c); err != nil {
		log.Printf("Error loading order items into cache: %v", err)
		return err
	}
	return nil
}
