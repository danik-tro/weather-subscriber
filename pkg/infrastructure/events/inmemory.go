package events

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain"
)

type Publisher struct {
	eventChan chan domain.Event
	handlers  map[domain.EventType][]domain.EventHandler
	mu        sync.RWMutex
	workers   int
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewPublisher(workers int, bufferSize int) *Publisher {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	if bufferSize <= 0 {
		bufferSize = 100
	}

	ctx, cancel := context.WithCancel(context.Background())

	p := &Publisher{
		eventChan: make(chan domain.Event, bufferSize),
		handlers:  make(map[domain.EventType][]domain.EventHandler),
		workers:   workers,
		ctx:       ctx,
		cancel:    cancel,
	}
	return p
}

func (p *Publisher) Register(eventType domain.EventType, handler domain.EventHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.handlers[eventType]; !exists {
		p.handlers[eventType] = []domain.EventHandler{}
	}

	p.handlers[eventType] = append(p.handlers[eventType], handler)
}

func (p *Publisher) Trigger(event domain.Event) []error {
	p.mu.RLock()
	handlers, exists := p.handlers[event.Type]
	p.mu.RUnlock()

	if !exists {
		return nil
	}

	var errors []error

	for _, handler := range handlers {
		if err := handler(event.Context, event); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (p *Publisher) TriggerAsync(event domain.Event) {
	select {
	case p.eventChan <- event:
	case <-p.ctx.Done():
		log.Printf("Publisher is closed, dropping event %s", event.Type)
	default:
		go func() {
			ctx, cancel := context.WithTimeout(p.ctx, 5*time.Second)
			defer cancel()

			done := make(chan struct{})
			go func() {
				p.processEvent(event)
				close(done)
			}()

			select {
			case <-done:
			case <-ctx.Done():
				log.Printf("Timeout processing event %s", event.Type)
			}
		}()
	}
}

func (p *Publisher) Start() {
	p.wg.Add(p.workers)

	for i := range p.workers {
		go p.worker(i)
	}
}

func (p *Publisher) worker(id int) {
	defer p.wg.Done()

	log.Printf("Event worker %d started", id)

	for {
		select {
		case evt := <-p.eventChan:
			p.processEvent(evt)
		case <-p.ctx.Done():
			log.Printf("Event worker %d stopping", id)
			return
		}
	}
}

func (p *Publisher) processEvent(evt domain.Event) {

	p.mu.RLock()
	handlers, exists := p.handlers[evt.Type]
	p.mu.RUnlock()

	if !exists {
		return
	}

	for _, handler := range handlers {
		if err := handler(evt.Context, evt); err != nil {
			log.Printf("Error handling event %s: %v", evt.Type, err)
		}
	}
}

func (p *Publisher) Close() {
	p.cancel()
	p.wg.Wait()
	close(p.eventChan)
	log.Println("Event publisher closed")
}
