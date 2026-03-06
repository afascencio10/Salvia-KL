package utils

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"time"
)

// Task models a scheduled job
type Task struct {
	Name            string `json:"name"`
	ID              string `json:"id"`
	ScheduleType    string `json:"schedule_type"` // "interval", "daily", "weekly"
	IntervalSeconds int    `json:"interval_seconds,omitempty"`
	Hour            int    `json:"hour,omitempty"`    // for daily/weekly
	Minute          int    `json:"minute,omitempty"`  // for daily/weekly
	Weekday         string `json:"weekday,omitempty"` // for weekly
	Action          func() `json:"-"`
	Enabled         bool   `json:"enabled"`
}

var tasks []Task

// Inicializa las tareas junto con la función que debe ejecutar cada una
func Load(taskFuncs map[string]func()) {

	/* ejemplo de tasks.json
	[
		{
		"name": "CleanupTempFiles",
		"id": "cleanup_01",
		"schedule_type": "interval",
		"interval_seconds": 300
		},
		{
		"name": "DailyReport",
		"id": "report_01",
		"schedule_type": "daily",
		"hour": 2,
		"minute": 30
		},
		{
		"name": "WeeklySync",
		"id": "sync_01",
		"schedule_type": "weekly",
		"weekday": "Monday",
		"hour": 4,
		"minute": 0
		}
	]
	*/

	data, err := fs.ReadFile(ConfigAssets, "config/tasks.json")
	if err != nil {
		fmt.Println("Error: TaskScheduler leyendo tasks.json: ", err)
	}

	if err := json.Unmarshal(data, &tasks); err != nil {
		fmt.Println("Error: TaskScheduler JSON inválido: ", err)
	}

	for _, t := range tasks {
		if t.Enabled {
			fun, found := taskFuncs[t.ID]
			if found {
				t.Action = fun
			}
			if t.Action != nil {
				switch t.ScheduleType {
				case "interval":
					startInterval(t)
				case "daily":
					go scheduleDaily(t)
				case "weekly":
					go scheduleWeekly(t)
				default:
					fmt.Println("Error: TaskScheduler Tipo de schedule desconocido para ", t.ID)
				}
			} else {
				fmt.Println("Error: TaskScheduler tarea sin función a ejecutar: ", t.ID)
			}
		}

	}

	select {} // bloquea para que las goroutines sigan corriendo
}

func executeTask(t Task) {
	fmt.Printf("[%s] Ejecutando %s (%s)\n", time.Now().Format(time.RFC3339), t.Name, t.ID)
	go t.Action()
}

func startInterval(t Task) {
	ticker := time.NewTicker(time.Duration(t.IntervalSeconds) * time.Second)
	go func() {
		for range ticker.C {
			executeTask(t)
		}
	}()
}

func scheduleDaily(t Task) {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), t.Hour, t.Minute, 0, 0, now.Location())
		if next.Before(now) {
			next = next.AddDate(0, 0, 1)
		}
		time.Sleep(time.Until(next))
		executeTask(t)
	}
}

func scheduleWeekly(t Task) {
	weekdays := map[string]time.Weekday{
		"Sunday":    time.Sunday,
		"Monday":    time.Monday,
		"Tuesday":   time.Tuesday,
		"Wednesday": time.Wednesday,
		"Thursday":  time.Thursday,
		"Friday":    time.Friday,
		"Saturday":  time.Saturday,
	}
	target := weekdays[t.Weekday]
	for {
		now := time.Now()
		offset := (int(target) - int(now.Weekday()) + 7) % 7
		next := time.Date(now.Year(), now.Month(), now.Day()+offset, t.Hour, t.Minute, 0, 0, now.Location())
		if next.Before(now) {
			next = next.AddDate(0, 0, 7)
		}
		time.Sleep(time.Until(next))
		executeTask(t)
	}
}
