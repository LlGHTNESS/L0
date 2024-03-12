package main

import (
	"fmt"
	"html/template"
	"log"
	"module/internal/database"
	"module/internal/service"
	"module/internal/transport/rest"
	cachex "module/pkg/cache"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
	stan "github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
)

func main() {

	if err := godotenv.Load("configs/config.env"); err != nil {
		log.Fatalln("Error loading .env file")
	}

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"))

	db, err := sqlx.Connect("postgres", dbinfo)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	sqlContent, err := database.ReadSQLFile("/app/schemes/data_struct.sql")
	if err != nil {
		log.Fatalf("Error reading SQL file: %v\n", err)
	}

	if err := database.ExecuteSQLCommands(db, sqlContent); err != nil {
		log.Fatalf("Error executing SQL commands: %v\n", err)
	}

	natsConnection := fmt.Sprintf("nats://%s", os.Getenv("NATS_URL"))
	sc, err := stan.Connect(
		os.Getenv("NATS_CLUSTER_ID"),
		os.Getenv("NATS_CLIENT_ID"),
		stan.NatsURL(natsConnection),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer sc.Close()

	c := cache.New(5*time.Minute, 10*time.Minute)

	sub, err := sc.Subscribe("main", func(m *stan.Msg) {
		service.HandleMessage(db, m.Data, c)
	}, stan.StartWithLastReceived())

	if err != nil {
		log.Fatalln(err)
	}
	defer sub.Unsubscribe()

	err = cachex.RefreshCache(db, c)
	if err != nil {
		log.Fatalf("Error refreshing cache: %v", err)
	}

	tmpl, err := template.ParseFiles("/app/web/index.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	http.HandleFunc("/", rest.ServeTemplate(tmpl))

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		rest.SearchOrderData(w, r, c, tmpl)
	})

	fmt.Println("The server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
