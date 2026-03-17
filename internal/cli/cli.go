package cli

import (
	"log"
	"strconv"
	"strings"

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

	switch {

	case cmd == "stop":
		svc.Stop()

	case cmd == "status":
		svc.Status()

	case cmd == "logs":
		svc.Logs()

	case cmd == "todos":
		log.Printf("Fetching todos from Notion...")
		if err := notionClient.GetTodos(); err != nil {
			log.Printf("Error fetching todos: %v", err)
		}

	case cmd == "logs-all":
		svc.LogsAll()

	case strings.HasPrefix(cmd, "mark-"):
		n, err := strconv.Atoi(strings.TrimPrefix(cmd, "mark-"))
		if err != nil {
			log.Printf("Invalid mark command: %v", err)
			return
		}
		log.Printf("Marking todo %d in Notion...", n)
		if err := notionClient.MarkTodo(n); err != nil {
			log.Printf("Error marking todo: %v", err)
		}

	case strings.HasPrefix(cmd, "unmark-"):
		n, err := strconv.Atoi(strings.TrimPrefix(cmd, "unmark-"))
		if err != nil {
			log.Printf("Invalid unmark command: %v", err)
			return
		}
		log.Printf("Unmarking todo %d in Notion...", n)
		if err := notionClient.UnmarkTodo(n); err != nil {
			log.Printf("Error unmarking todo: %v", err)
		}

	default:
		svc.Start(cmd)

	}
}
