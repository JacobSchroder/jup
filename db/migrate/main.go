package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

	namePtr := flag.String("name", "", "The name of the migration")
	flag.Parse()

	fmt.Println(flag.Args())

	if len(flag.Args()) != 1 {
		log.Fatal("Usage: migrate [-name=<migration_name>] <new|up>")
		os.Exit(1)
	}

	action := flag.Args()[0]

	switch action {
	case "new":
		suffix := ""
		if *namePtr != "" {
			suffix = "_" + *namePtr
		}

		migrationName := fmt.Sprintf("%s%s", time.Now().Format("20060102150405"), suffix)
		if len(migrationName) > 250 {
			log.Fatalf("Migration name must be less than 250 characters\nprovided: %s\n", migrationName)
			os.Exit(1)
		}

		createMigration(migrationName)
	case "up":
		Migrate(*namePtr)
	}

}
