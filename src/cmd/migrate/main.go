package main

import (
	"fmt"
	"log"

	commondb "bitsflow/common/db"
	internaldb "bitsflow/internal/db"
	"bitsflow/internal/models"
)

func main() {
	cfg := commondb.DBClientConfig{
		Hostname:     "dpg-d76je5ffte5s73elf39g-a.virginia-postgres.render.com",
		Port:         "5432",
		DatabaseName: "salvia_pruebas_q9xt",
		UserName:     "root",
		Password:     "k07wAhHOlrkswVlOAIms3jvyymiLt2k4",
	}
	db, err := internaldb.NewGormDB(cfg)
	if err != nil { log.Fatal(err) }

	tables := []string{
		"salvia.answer", "salvia.repeater_entry", "salvia.visibility_condition",
		"salvia.question", "salvia.repeater_group", "salvia.form_submission",
	}
	for _, t := range tables {
		if r := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", t)); r.Error != nil {
			log.Printf("drop %s: %v", t, r.Error)
		} else {
			fmt.Println("Dropped:", t)
		}
	}

	ms := []interface{}{
		&models.RepeaterGroup{}, &models.Question{}, &models.VisibilityCondition{},
		&models.FormSubmission{}, &models.RepeaterEntry{}, &models.Answer{},
	}
	for _, m := range ms {
		if err := db.AutoMigrate(m); err != nil {
			log.Printf("WARN %T: %v", m, err)
		} else {
			fmt.Printf("Migrated: %T\n", m)
		}
	}
	fmt.Println("Done!")
}
