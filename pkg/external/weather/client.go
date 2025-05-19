package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/danik-tro/weather-subscriber/pkg/domain"
)

type WeatherAPIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewWeatherAPIClient(apiKey string) WeatherClient {
	return &WeatherAPIClient{
		baseURL:    "https://api.weatherapi.com/v1",
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *WeatherAPIClient) GetCurrentWeather(ctx context.Context, city string) (*WeatherData, error) {
	endpoint := fmt.Sprintf("%s/current.json", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	// TODO: add logging
	if err != nil {
		// TODO: add logging
		return nil, domain.ErrInternalServerError
	}

	q := url.Values{}
	q.Add("key", c.apiKey)
	q.Add("q", city)
	q.Add("aqi", "no")
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// TODO: add logging
		return nil, domain.ErrInternalServerError
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// TODO: add logging
		return nil, domain.ErrInternalServerError
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var weatherData WeatherData
		if err := json.Unmarshal(body, &weatherData); err != nil {
			// TODO: add logging
			return nil, domain.ErrInternalServerError
		}
		return &weatherData, nil
	case http.StatusBadRequest:
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			// TODO: add logging
			return nil, domain.ErrInternalServerError
		}

		if errorResp.Error.Code == 1006 {
			return nil, domain.ErrCityNotFound
		}

		// TODO: add logging
		return nil, domain.ErrBadRequest
	default:
		return nil, domain.ErrBadRequest
	}
}
