package restream

import (
	"encoding/json"
	"net/http"
)

type ExchangeBody struct {
	Code string `json:"code"`
}

type ResponseData struct {
	Message string `json:"msg"`
}

func ExchangeTokens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		e := ExchangeBody{}

		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tokens, err := requestTokens(e.Code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		res, err := StoreTokens(&tokens)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotModified)
			return
		}
		
		msg := ResponseData{res}
		json.NewEncoder(w).Encode(msg)
	}
}
