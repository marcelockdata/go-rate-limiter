package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func ZipCodeHandler(w http.ResponseWriter, r *http.Request) {
	zipcode := chi.URLParam(r, "zipcode")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://viacep.com.br/ws/%s/json", zipcode), nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Fail to create the request: %v", err), http.StatusInternalServerError)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Fail to make the request: %v", err), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	ctx_err := ctx.Err()
	if ctx_err != nil {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			http.Error(w, fmt.Sprintf("Max timeout reached: %v", err), http.StatusRequestTimeout)
			return
		}
	}

	resp_json, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Fail to read the response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp_json)
}
