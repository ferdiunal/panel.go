package notification

import (
	"context"
	"fmt"
	"sync"
)

type Dispatcher struct {
	manager *NotificationManager
}

func NewDispatcher(manager *NotificationManager) *Dispatcher {
	return &Dispatcher{
		manager: manager,
	}
}

func (d *Dispatcher) Send(ctx context.Context, notifiable Notifiable, notification Notification) error {
	return d.manager.Send(ctx, notifiable, notification)
}

func (d *Dispatcher) SendVia(ctx context.Context, channels []string, notifiable Notifiable, notification Notification) error {
	return d.manager.SendVia(ctx, channels, notifiable, notification)
}

func (d *Dispatcher) SendToMany(ctx context.Context, notifiables []Notifiable, notification Notification) error {
	var wg sync.WaitGroup
	errorCh := make(chan error, len(notifiables))

	for _, notifiable := range notifiables {
		wg.Add(1)
		go func(n Notifiable) {
			defer wg.Done()
			if err := d.Send(ctx, n, notification); err != nil {
				errorCh <- fmt.Errorf("failed to send to %T: %w", n, err)
			}
		}(notifiable)
	}

	wg.Wait()
	close(errorCh)

	var errors []error
	for err := range errorCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification sending failed for %d recipients: %v", len(errors), errors[0])
	}

	return nil
}

func (d *Dispatcher) SendToManyVia(ctx context.Context, channels []string, notifiables []Notifiable, notification Notification) error {
	var wg sync.WaitGroup
	errorCh := make(chan error, len(notifiables))

	for _, notifiable := range notifiables {
		wg.Add(1)
		go func(n Notifiable) {
			defer wg.Done()
			if err := d.SendVia(ctx, channels, n, notification); err != nil {
				errorCh <- fmt.Errorf("failed to send to %T: %w", n, err)
			}
		}(notifiable)
	}

	wg.Wait()
	close(errorCh)

	var errors []error
	for err := range errorCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification sending failed for %d recipients: %v", len(errors), errors[0])
	}

	return nil
}

func (d *Dispatcher) GetAvailableChannels() []string {
	return d.manager.GetAvailableChannels()
}