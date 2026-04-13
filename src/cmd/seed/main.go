package main

import (
	"fmt"
	"log"
	"os"

	commondb "bitsflow/common/db"
	internaldb "bitsflow/internal/db"
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
	if err != nil {
		log.Fatal(err)
	}

	file := "cmd/seed/seed_seguimiento.sql"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	sql, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	if r := db.Exec(string(sql)); r.Error != nil {
		log.Fatal("seed error:", r.Error)
	}

	fmt.Println("Seed ejecutado correctamente.")
}
