package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"zombie_locator/internal/http"
	"zombie_locator/internal/logger"
	"zombie_locator/internal/repository/zombie"
	"zombie_locator/internal/service/locator"
	"zombie_locator/internal/service/observer"
	"zombie_locator/internal/storage/broker"
	"zombie_locator/internal/storage/db"
	"zombie_locator/internal/utils/shema_registry"
)

var (
	kafkaBroker        = "localhost:9092"
	kafkaConsumerGroup = "zombie-tracker"

	httpAddr = "127.0.0.1:8000"

	postgresqlURL = "postgresql://test:secret@localhost:5432/zombies?sslmode=disable"
	redisURL      = "localhost:9851"
)

func main() {
	appLog, err := logger.NewAppLogger()
	if err != nil {
		log.Fatalf("unable to create logger: %s", err)
	}
	// Set up a postgresql database connection.
	dbConnect, err := db.NewPostgresConnection(postgresqlURL)
	if err != nil {
		appLog.Fatal("unable to connect to database", err)
	}

	// Set up a tile38 database connection.
	tile38Client, err := db.NewTile38Connection(redisURL)
	if err != nil {
		appLog.Fatal("tile38 database is not reachable", err)
	}

	registry := shema_registry.NewRegistry([]int{1})
	zRepo := zombie.NewZombieRepository(dbConnect, tile38Client)

	appLog.Info("init observer service")
	// dead-letter queue producers init here
	locationDLQProducer := broker.NewKafkaProducer(appLog, kafkaBroker, "zombie-location-dql")
	statusDLQProducer := broker.NewKafkaProducer(appLog, kafkaBroker, "zombie-status-dql")

	// consumers for different type of events
	locationConsumer := broker.NewKafkaConsumer(appLog, locationDLQProducer, []string{kafkaBroker}, kafkaConsumerGroup, "zombie_locations")
	zombieStatusConsumer := broker.NewKafkaConsumer(appLog, statusDLQProducer, []string{kafkaBroker}, kafkaConsumerGroup, "captured_zombies")
	zombieObserver := observer.NewObserver(appLog, zRepo, registry, locationConsumer, zombieStatusConsumer)

	// Set up HTTP handler and router
	appLog.Info("init http service")
	appHTTPServer := http.NewServer(appLog, httpAddr, locator.NewLocatorService(appLog, zRepo))

	// Start the HTTP handler and Kafka consumers in parallel.
	appLog.Info("starting services")
	go func() {
		if err = appHTTPServer.Run(); err != nil {
			appLog.Fatal("error start appHTTPServer", err)
		}
	}()
	zombieObserver.Run()

	// Wait for OS termination signal
	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	<-wait
	if err = appHTTPServer.Shutdown(); err != nil {
		appLog.Error("unable to shutdown http server", err)
	}
	// Since kafka consumers are not disconnected, the kafka broker will wait for a timeout before it is able to reassign partitions to new instances.
	// We consider this acceptabke in the context of this technical test.

	if err = zombieObserver.Shutdown(); err != nil {
		appLog.Error("unable to shutdown zombie observer", err)
	}
}
