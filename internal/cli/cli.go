package cli

import (
	"log"

	"github.com/RohanDSkaria/time-it/internal/db"
	"github.com/RohanDSkaria/time-it/internal/notion"
	"github.com/RohanDSkaria/time-it/internal/repository"
	"github.com/RohanDSkaria/time-it/internal/service"
)

func Run(args []string) {
	if len(args) == 1 {
		return
	}

	db, err := db.Open()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	repo := repository.New(db)
	svc := service.New(repo)
	notionClient := notion.New()
	cmd := args[1]

	switch cmd {

	case "stop":
		svc.Stop()

	case "status":
		svc.Status()

	case "logs":
		svc.Logs()

	case "todos":
		log.Printf("Fetching todos from Notion...")
		if err := notionClient.GetTodos(); err != nil {
			log.Printf("Error fetching todos: %v", err)
		}

	default:
		svc.Start(cmd)

	}
}
