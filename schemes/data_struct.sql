-- Active: 1709711485375@@127.0.0.1@5432@postgres
BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS deliveries(
    "delivery_id" SERIAL PRIMARY KEY,
    "name" VARCHAR(50) NOT NULL,
    "phone" VARCHAR(25) NOT NULL,
    "zip" VARCHAR(20) NOT NULL,
    "city" VARCHAR(50) NOT NULL,
    "address" VARCHAR(255) NOT NULL,
    "region" VARCHAR(255) NOT NULL,
    "email" VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS payments(
    "payments_id" SERIAL PRIMARY KEY,
    "transaction" VARCHAR(50) NOT NULL,
    "request_id" VARCHAR(50) NOT NULL,
    "currency" VARCHAR(3) NOT NULL,
    "provider" VARCHAR(50) NOT NULL,
    "amount" INTEGER NOT NULL,
    "payment_dt" INTEGER NOT NULL,
    "bank" VARCHAR(50) NOT NULL,
    "delivery_cost" INTEGER NOT NULL,
    "goods_total" INTEGER NOT NULL,
    "custom_fee" INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS items(
    "chrt_id" INTEGER PRIMARY KEY,
    "track_number" VARCHAR(50) NOT NULL,
    "price" INTEGER NOT NULL,
    "rid" VARCHAR(50) NOT NULL,
    "name" VARCHAR(50) NOT NULL,
    "sale" INTEGER NOT NULL,
    "size" VARCHAR(20) NOT NULL,
    "total_price" INTEGER NOT NULL,
    "nm_id" INTEGER NOT NULL,
    "brand" VARCHAR(50) NOT NULL,
    "status" INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS orders(
    "order_uid"  VARCHAR(50) PRIMARY KEY,
    "track_number"  VARCHAR(50) NOT NULL,
    "entry"  VARCHAR(50) NOT NULL,
    "delivery" INTEGER REFERENCES deliveries(delivery_id) ON DELETE CASCADE ON UPDATE CASCADE,
    "payment" INTEGER REFERENCES payments(payments_id) ON DELETE CASCADE ON UPDATE CASCADE,
    "locale"  VARCHAR(2) NOT NULL,
    "internal_signature"  VARCHAR(50) NOT NULL,
    "customer_id"  VARCHAR(50) NOT NULL,
    "delivery_service"  VARCHAR(50) NOT NULL,
    "shardkey"  VARCHAR(10) NOT NULL,
    "sm_id" INTEGER NOT NULL,
    "date_created" TIMESTAMP NOT NULL,
    "oof_shard"  VARCHAR(10) NOT NULL
);
CREATE TABLE IF NOT EXISTS orders_items(
    "order_id" VARCHAR(50) REFERENCES orders(order_uid) ON DELETE CASCADE ON UPDATE CASCADE,
    "item_id" INTEGER REFERENCES items(chrt_id) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (order_id, item_id)
);
COMMIT TRANSACTION;