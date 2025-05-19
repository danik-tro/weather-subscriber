package background_job

import (
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

type BackgroundJobService interface {
	Start() error
	Stop() error
	AddJob(schedule string, job func()) error
}

type CronBackgroundJobService struct {
	cron    *cron.Cron
	mu      sync.RWMutex
	running bool
}

func NewCronBackgroundJobService() *CronBackgroundJobService {
	return &CronBackgroundJobService{
		cron: cron.New(cron.WithSeconds()),
	}
}

func (s *CronBackgroundJobService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("background job service is already running")
	}

	s.cron.Start()
	s.running = true
	return nil
}

func (s *CronBackgroundJobService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return fmt.Errorf("background job service is not running")
	}

	ctx := s.cron.Stop()
	<-ctx.Done()
	s.running = false
	return nil
}

func (s *CronBackgroundJobService) AddJob(schedule string, job func()) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return fmt.Errorf("background job service is not running")
	}

	_, err := s.cron.AddFunc(schedule, job)
	if err != nil {
		return fmt.Errorf("failed to add job: %w", err)
	}

	return nil
}
