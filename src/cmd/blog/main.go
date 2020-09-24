package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hatajoe/blog/ent"
)

var (
	dbUser string = os.Getenv("DB_USER")
	dbPass string = os.Getenv("DB_PASS")
	dbHost string = os.Getenv("DB_HOST")
	dbPort string = os.Getenv("DB_PORT")
	dbName string = os.Getenv("DB_NAME")
)

func main() {
	client, err := ent.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName))
	if err != nil {
		log.Fatalf("failed connecting to mysql: %v", err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalf("failed closing mysql client: %v", err)
		}
	}()

	ctx := context.Background()
	// Run migration.
	err = client.Schema.Create(ctx)
	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
