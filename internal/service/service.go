package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RohanDSkaria/time-it/internal/model"
	"github.com/RohanDSkaria/time-it/internal/repository"
)

type Service struct {
	Repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Start(task string) error {
	if err := s.Stop(); err != nil {
		return err
	}

	now := time.Now().Unix()

	return s.Repo.SetCurrentEntry(task, now)
}

func (s *Service) Stop() error {
	currentEntry, err := s.Repo.GetCurrentEntry()
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	now := time.Now().Unix()

	entry := model.Entry{
		Task:     currentEntry.Task,
		Start:    currentEntry.Start,
		Duration: now - currentEntry.Start,
	}

	if err := s.Repo.AddEntry(entry); err != nil {
		return err
	}

	if err := s.Repo.DeleteCurrentEntry(); err != nil {
		return err
	}

	return nil
}

func (s *Service) Status() error {
	currentEntry, err := s.Repo.GetCurrentEntry()
	if err == sql.ErrNoRows {
		fmt.Println("No task is currently being timed.")
		return nil
	}
	if err != nil {
		return err
	}

	fmt.Printf("Current task: %s\n", currentEntry.Task)
	fmt.Printf("Started at: %s\n", time.Unix(currentEntry.Start, 0).Format(time.RFC1123))
	fmt.Printf("Elapsed time: %s\n", time.Since(time.Unix(currentEntry.Start, 0)).Truncate(time.Second))

	return nil
}

func (s *Service) Logs() error {
	entries, err := s.Repo.GetAllEntries()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("No entries found.")
		return nil
	}

	cutoff := time.Now().Add(-24 * time.Hour).Unix()

	idx := 0
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].Start < cutoff {
			idx = i + 1
			break
		}
	}

	for i := idx; i < len(entries); i++ {
		entry := entries[i]

		fmt.Printf("Task: %s\n", entry.Task)
		fmt.Printf("Start: %s\n", time.Unix(entry.Start, 0).Format(time.RFC1123))
		fmt.Printf("Duration: %s\n", time.Duration(entry.Duration)*time.Second)
		fmt.Println()
	}

	return nil
}

func (s *Service) LogsAll() error {
	entries, err := s.Repo.GetAllEntries()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("No entries found.")
		return nil
	}

	for _, entry := range entries {
		fmt.Printf("Task: %s\n", entry.Task)
		fmt.Printf("Start: %s\n", time.Unix(entry.Start, 0).Format(time.RFC1123))
		fmt.Printf("Duration: %s\n", time.Duration(entry.Duration)*time.Second)
		fmt.Println()
	}

	return nil
}
