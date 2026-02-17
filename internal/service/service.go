package service

import (
	"database/sql"
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
