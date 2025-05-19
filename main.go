package main

import (
	"context"
	"fmt"
	"log"

	config "github.com/danik-tro/weather-subscriber/pkg"
	domain "github.com/danik-tro/weather-subscriber/pkg/domain"
	entity "github.com/danik-tro/weather-subscriber/pkg/domain/entity"
	"github.com/danik-tro/weather-subscriber/pkg/external/weather"
	"github.com/danik-tro/weather-subscriber/pkg/infrastructure/background_job"
	"github.com/danik-tro/weather-subscriber/pkg/infrastructure/db"
	smtp "github.com/danik-tro/weather-subscriber/pkg/infrastructure/email_service"
	events "github.com/danik-tro/weather-subscriber/pkg/infrastructure/events"
	"github.com/danik-tro/weather-subscriber/pkg/infrastructure/usecases"
	"github.com/danik-tro/weather-subscriber/pkg/presenter/http"

	_ "github.com/danik-tro/weather-subscriber/docs"
)

const (
	workers     = 10
	bufferSize  = 100
	redisPrefix = "weather_app"
)

// @title			Weather Service
// @version		1.0
// @description	Weather Service
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:8080
// @BasePath		/api
// @schemes		http https
func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	publisher := events.NewPublisher(workers, bufferSize)

	weatherClient := weather.NewWeatherAPIClient(config.WeatherAPIKey)
	weatherCache, err := weather.NewWeatherCache(config.RedisAddress, config.RedisPassword, config.RedisDB, redisPrefix)

	if err != nil {
		log.Fatal(err)
	}

	gormDb, err := db.NewGormConnection(config)
	if err != nil {
		log.Fatal(err)
	}

	repository := db.NewGormRepository(gormDb)

	if config.DBAutoMigrate {
		err = repository.EnsureSchema()
		if err != nil {
			log.Fatal(err)
		}
	}

	weatherService := weather.NewWeatherService(weatherClient, weatherCache)

	getWeatherUC := usecases.NewGetWeatherUseCase(*weatherService)
	subscribeUC := usecases.NewSubscribeWeatherUseCase(repository, publisher, *weatherService)
	confirmUC := usecases.NewConfirmSubscription(repository)
	unsubscribeUC := usecases.NewUnsubscribeUseCase(repository)
	checkTokensUC := usecases.NewCheckTokens(repository)

	router := http.NewRouter(*config, subscribeUC, getWeatherUC, confirmUC, unsubscribeUC, checkTokensUC)

	smtpConfig := smtp.SMTPConfig{
		Host:     config.SMTPHost,
		Port:     config.SMTPPort,
		Username: config.SMTPUsername,
		Password: config.SMTPPassword,
		From:     config.SMTPFrom,
	}
	emailService := smtp.NewSmtpService(smtpConfig)

	handler := events.Handler{
		EmailService:   emailService,
		WeatherService: weatherService,
		Repository:     repository,
		Config:         *config,
	}

	publisher.Register(domain.UserSubscribed, handler.UserSubscribed())
	publisher.Register(domain.WeatherEvent, handler.WeatherEvent())

	backgroundJobService := background_job.NewCronBackgroundJobService()
	if err := backgroundJobService.Start(); err != nil {
		log.Fatal(err)
	}
	defer backgroundJobService.Stop()

	if err := backgroundJobService.AddJob("0 0 12 * * *", func() {
		event := domain.Event{
			Type: domain.WeatherEvent,
		}
		if err := handler.FetchAndUpdateWeatherSubscribers(entity.FrequencyDaily)(context.Background(), event); err != nil {
			log.Printf("Failed to process daily weather updates: %v", err)
		}
	}); err != nil {
		log.Fatal(err)
	}

	if err := backgroundJobService.AddJob("0 0 * * * *", func() {
		event := domain.Event{
			Type: domain.WeatherEvent,
		}
		if err := handler.FetchAndUpdateWeatherSubscribers(entity.FrequencyHourly)(context.Background(), event); err != nil {
			log.Printf("Failed to process hourly weather updates: %v", err)
		}
	}); err != nil {
		log.Fatal(err)
	}

	publisher.Start()
	defer publisher.Close()

	if err := router.Run(fmt.Sprintf("%s:%d", config.AppHost, config.AppPort)); err != nil {
		log.Fatal(err)
	}
}
