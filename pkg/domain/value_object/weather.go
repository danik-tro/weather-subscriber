package domain

type Weather struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Description string  `json:"description"`
}

type WeatherEvent struct {
	City             string
	Country          string
	Temperature      float64
	Humidity         float64
	Description      string
	Email            string
	UnsubscribeToken string
}
