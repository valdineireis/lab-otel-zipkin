package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"go.opentelemetry.io/otel"
)

type CepRequest struct {
	Cep string `json:"cep"`
}

func CepHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("service-a")
	ctx, span := tracer.Start(r.Context(), "CepHandler")
	defer span.End()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
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

	// Validate CEP (8 digits)
	isValid := validateCep(req.Cep)
	if !isValid {
		http.Error(w, `{"message": "invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	// Forward to Service B
	resp, err := forwardToServiceB(ctx, req.Cep)
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	// Copiar o conteúdo de resp.Body para um buffer
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from Service B", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Escrever a resposta no ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func validateCep(cep string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}

func forwardToServiceB(ctx context.Context, cep string) (*http.Response, error) {
	tracer := otel.Tracer("service-a")
	ctx, span := tracer.Start(ctx, "forwardToServiceB")
	defer span.End()

	url := "http://service-b:8081/process"
	payload := map[string]string{"cep": cep}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
