package observer

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"zombie_locator/internal/entities"
	"zombie_locator/internal/logger"
	"zombie_locator/internal/repository/zombie"
	"zombie_locator/internal/storage/broker"
	"zombie_locator/internal/utils/shema_registry"

	"go.uber.org/zap"
)

var (
	UnsupportedConsumerType = errors.New("unsupported event type")
)

type Observer struct {
	ctx              context.Context
	cancel           context.CancelFunc
	log              logger.AppLogger
	repo             zombie.Zombier
	registry         shema_registry.SchemaRegistry
	statusConsumer   broker.Consumer
	locationConsumer broker.Consumer
}

func NewObserver(
	log logger.AppLogger,
	repo zombie.Zombier,
	registry shema_registry.SchemaRegistry,
	locationConsumer, statusConsumer broker.Consumer) *Observer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Observer{
		ctx:              ctx,
		cancel:           cancel,
		log:              log.With(zap.String("service", "observer")),
		repo:             repo,
		registry:         registry,
		locationConsumer: locationConsumer,
		statusConsumer:   statusConsumer,
	}
}

func (o *Observer) Run() {
	go func() {
		if err := o.locationConsumer.Run(o.ctx, o.ZombieLocationUpdate); err != nil {
			o.log.Error("failed to run location consumer", err)
		}
	}()
	go func() {
		if err := o.statusConsumer.Run(o.ctx, o.ZombieCapturedUpdate); err != nil {
			o.log.Error("failed to run status consumer", err)
		}
	}()
}

// ZombieLocationUpdate processes Kafka messages containing location updates..
func (o *Observer) ZombieLocationUpdate(ctx context.Context, payload []byte) error {
	log := o.log.With(zap.String("method", "ZombieLocationUpdate"))
	data, err := o.registry.DecodeZombieLocationStreamEvent(payload)
	if err != nil {
		// unsupported message structure
		log.Error("unsupported message structure", err)
		return fmt.Errorf("unsupported message structure: %w", err)
	}
	// we expect that map contains only one element
	for v := range data {
		switch v {
		case 1:
			// process v1
			return o.zombieLocationUpdateV1(ctx, log, data[v])
		default:
			// unsupported version
			err = fmt.Errorf("unsupported version: %d", v)
			log.Error("error update zombie location", err)
			return err
		}
	}
	return nil
}

// ZombieCapturedUpdate processes Kafka messages containing captured zombies data updates.
func (o *Observer) ZombieCapturedUpdate(ctx context.Context, payload []byte) error {
	log := o.log.With(zap.String("method", "ZombieCapturedUpdate"))
	data, err := o.registry.DecodeZombieCapturedStreamEvent(payload)
	if err != nil {
		// unsupported message structure
		log.Error("unsupported message structure", err)
		return fmt.Errorf("unsupported message structure: %w", err)
	}
	// we expect that map contains only one element
	for v := range data {
		switch v {
		case 1:
			// process v1
			return o.zombieCapturedUpdateV1(ctx, log, data[v])
		default:
			// unsupported version
			err = fmt.Errorf("unsupported version: %d", v)
			log.Error("error update zombie status", err)
			return fmt.Errorf("error update zombie status: %w", err)
		}
	}
	return nil
}

func (o *Observer) zombieCapturedUpdateV1(ctx context.Context, log logger.AppLogger, payload any) error {
	zC, ok := payload.(entities.ZombieCapturedV1)
	if !ok {
		payloadType := reflect.TypeOf(payload).String()
		log.Error("unsupported event type", UnsupportedConsumerType, zap.String("type", payloadType))
		return fmt.Errorf("unsupported type for zombie capture update v1 :%s", payloadType)
	}
	if err := o.repo.CapturedZombie(ctx, zC.ZombieID, zC.UpdatedAt); err != nil {
		log.Error("failed to update zombie status", err)
		return fmt.Errorf("failed to update zombie status: %w", err)
	}
	return nil
}

func (o *Observer) zombieLocationUpdateV1(ctx context.Context, log logger.AppLogger, payload any) error {
	zL, ok := payload.(entities.ZombieLocationV1)
	if !ok {
		payloadType := reflect.TypeOf(payload).String()
		log.Error("unsupported event type", UnsupportedConsumerType, zap.String("type", payloadType))
		return fmt.Errorf("unsupported type for zombie location update v1 :%s", payloadType)
	}
	if err := o.repo.LocatedZombie(ctx, zL.ZombieID, zL.Latitude, zL.Longitude, zL.UpdatedAt); err != nil {
		log.Error("failed to update zombie location", err)
		return fmt.Errorf("failed store zombie location: %w", err)
	}
	return nil
}

func (o *Observer) Shutdown() error {
	var wg sync.WaitGroup
	wg.Add(2)
	o.cancel()
	go func() {
		if err := o.locationConsumer.Shutdown(); err != nil {
			o.log.Error("failed to shutdown location consumer", err)
		}
		wg.Done()
	}()

	go func() {
		if err := o.statusConsumer.Shutdown(); err != nil {
			o.log.Error("failed to shutdown status consumer", err)
		}
		wg.Done()
	}()
	wg.Wait()
	return nil
}
