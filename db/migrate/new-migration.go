package main

import (
	"fmt"
	"log"
	"os"
)

func createMigration(name string) {
	os.MkdirAll("db/migrations", 0755)

	fileName := fmt.Sprintf("db/migrations/%s.sql", name)
	_, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Migration %s created successfully\n", name)
}
