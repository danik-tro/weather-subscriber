package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	domain "github.com/danik-tro/weather-subscriber/pkg/domain/value_object"
)

type RedisWeatherCache struct {
	client *redis.Client
	prefix string
}

func NewWeatherCache(redisAddress, password string, db int, prefix string) (WeatherCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisWeatherCache{
		client: client,
		prefix: prefix,
	}, nil
}

func (r *RedisWeatherCache) generateKey(city string) string {
	return fmt.Sprintf("%s:weather:%s", r.prefix, city)
}

func (r *RedisWeatherCache) GetWeather(ctx context.Context, city string) (*domain.Weather, error) {
	key := r.generateKey(city)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get weather data from Redis: %w", err)
	}

	var weather domain.Weather
	if err := json.Unmarshal([]byte(data), &weather); err != nil {
		return nil, fmt.Errorf("failed to deserialize weather data: %w", err)
	}

	return &weather, nil
}

func (r *RedisWeatherCache) SetWeather(ctx context.Context, city string, weather *domain.Weather, ttl time.Duration) error {
	key := r.generateKey(city)

	data, err := json.Marshal(weather)
	if err != nil {
		return fmt.Errorf("failed to serialize weather data: %w", err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set weather data in Redis: %w", err)
	}

	return nil
}

func (r *RedisWeatherCache) Close() error {
	return r.client.Close()
}
