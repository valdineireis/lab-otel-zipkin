package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}
}

type CepRequest struct {
	Cep string `json:"cep"`
}

type Response struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func CepHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(r.Context(), "CepHandler")
	defer span.End()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body) // Substituindo ioutil.ReadAll por io.ReadAll
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req CepRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate CEP
	if len(req.Cep) != 8 {
		http.Error(w, `{"message": "invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	// Fetch location from ViaCEP
	city, err := fetchLocation(ctx, req.Cep)
	if err != nil {
		http.Error(w, `{"message": "can not find zipcode"}`, http.StatusNotFound)
		return
	}

	// Fetch temperature from WeatherAPI
	tempC, err := fetchTemperature(ctx, city)
	if err != nil {
		http.Error(w, "Failed to fetch temperature", http.StatusInternalServerError)
		return
	}

	// Convert temperatures
	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	// Respond with data
	resp := Response{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func fetchLocation(ctx context.Context, cep string) (string, error) {
	tracer := otel.Tracer("service-b")
	_, span := tracer.Start(ctx, "fetchLocation")
	defer span.End()

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if _, ok := data["erro"]; ok {
		return "", fmt.Errorf("CEP not found")
	}

	return data["localidade"].(string), nil
}

func fetchTemperature(ctx context.Context, city string) (float64, error) {
	tracer := otel.Tracer("service-b")
	_, span := tracer.Start(ctx, "fetchTemperature")
	defer span.End()

	apiKey := os.Getenv("WEATHERAPI_KEY") // Load the key from the environment
	if apiKey == "" {
		return 0, fmt.Errorf("WeatherAPI key not set")
	}

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, city)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data["current"].(map[string]interface{})["temp_c"].(float64), nil
}
