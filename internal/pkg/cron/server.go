package cron

import (
	"context"
	"fmt"

	"github.com/robfig/cron"
)

type Server struct {
	cron *cron.Cron
}

func NewServer() *Server {
	return &Server{
		cron: cron.New(),
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.cron.Start()
	<-ctx.Done()
	s.cron.Stop()
	return ctx.Err()
}

func (s *Server) Stop(ctx context.Context) error {
	s.cron.Stop()
	return nil
}

func (s *Server) AddTask(schedule string, job func()) error {
	err := s.cron.AddFunc(schedule, job)
	if err != nil {
		return fmt.Errorf("failed to add task: %v", err)
	}
	return nil
}
