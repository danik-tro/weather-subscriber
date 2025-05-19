package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	domain_errors "github.com/danik-tro/weather-subscriber/pkg/domain"
	usecase "github.com/danik-tro/weather-subscriber/pkg/domain/usecases"
	domain "github.com/danik-tro/weather-subscriber/pkg/domain/value_object"
)

type MockGetWeatherUseCase struct {
	mock.Mock
}

func (m *MockGetWeatherUseCase) GetWeather(ctx context.Context, city string) (domain.Weather, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(domain.Weather), args.Error(1)
}

func setupRouter(uc usecase.GetWeatherUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/weather", GetWeatherHandler(uc))
	return r
}

func TestGetWeatherHandler_Success(t *testing.T) {
	mockUC := new(MockGetWeatherUseCase)

	want := domain.Weather{Temperature: 20.5, Humidity: 55, Description: "Sunny"}
	mockUC.On("GetWeather", mock.Anything, "Kyiv").Return(want, nil).Once()

	router := setupRouter(mockUC)
	req := httptest.NewRequest(http.MethodGet, "/api/weather?city=Kyiv", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t,
		`{"temperature":20.5,"humidity":55,"description":"Sunny"}`,
		w.Body.String(),
	)
	mockUC.AssertExpectations(t)
}

func TestGetWeatherHandler_CityNotFound(t *testing.T) {
	mockUC := new(MockGetWeatherUseCase)
	mockUC.On("GetWeather", mock.Anything, "Nowhere").
		Return(domain.Weather{}, domain_errors.ErrCityNotFound).Once()

	router := setupRouter(mockUC)
	req := httptest.NewRequest(http.MethodGet, "/api/weather?city=Nowhere", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"error":"city not found"}`, w.Body.String())
	mockUC.AssertExpectations(t)
}

func TestGetWeatherHandler_MissingCity(t *testing.T) {
	mockUC := new(MockGetWeatherUseCase)
	router := setupRouter(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/api/weather", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"query parameter 'city' is required"}`, w.Body.String())
}
