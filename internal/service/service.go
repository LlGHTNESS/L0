package service

import (
	"encoding/json"
	"log"
	"module/internal/database"
	"module/internal/model"
	cachex "module/pkg/cache"

	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"github.com/xeipuuv/gojsonschema"
)

func HandleMessage(db *sqlx.DB, data []byte, c *cache.Cache) {
	var order model.OrderData
	err := json.Unmarshal(data, &order)
	if err != nil {
		log.Printf("Error decoding JSON from NATS message: %v\n", err)
		return
	}

	isValid, err := validateJSON(data)
	if err != nil {
		log.Printf("Validation error: %s\n", err)
	} else if isValid {
		var existingID string
		err = db.Get(&existingID, "SELECT order_uid FROM orders WHERE order_uid=$1", order.OrderUID)
		if err == nil && existingID != "" {
			log.Printf("The order with ID %s is already in the database. \n", order.OrderUID)
			return
		}

		orderID, err := database.DbInsertOrder(db, order)
		if err != nil {
			log.Printf("Error writing an order to the database: %v", err)
			return
		}
		log.Printf("The order %s has been successfully added to the database.", orderID)

		err = cachex.RefreshCache(db, c)
		if err != nil {
			log.Fatalf("Error refreshing cache: %v", err)
		}
	} else {
		log.Println("JSON failed validation")
	}

}
func validateJSON(input []byte) (bool, error) {
	schemaLoader := gojsonschema.NewReferenceLoader("file:///app/configs/schema.json")
	documentLoader := gojsonschema.NewBytesLoader(input)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, err
	}
	return result.Valid(), nil
}
