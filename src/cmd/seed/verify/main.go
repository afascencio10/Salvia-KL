package main

import (
	"fmt"
	"log"

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

	fmt.Println("\n=== REPEATER ENTRIES SUB4 ===")
	rows, _ := db.Raw(`
		SELECT re.iteration, a.question_id, a.value
		FROM salvia.repeater_entry re
		JOIN salvia.answer a ON a.repeater_entry_id = re.id
		WHERE re.form_submission_id = 'ed09276b-47b5-4ea5-9a43-49afb45f9d7b'
		ORDER BY re.iteration, a.question_id
	`).Rows()
	defer rows.Close()
	cur := 0
	for rows.Next() {
		var iter int
		var qid, val string
		rows.Scan(&iter, &qid, &val)
		if iter != cur {
			fmt.Printf("\n  Iteración %d:\n", iter)
			cur = iter
		}
		fmt.Printf("    q=%s val=%q\n", qid[:8], val)
	}
}
