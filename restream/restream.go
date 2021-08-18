package restream

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
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
		log.Printf("Received Code=%+v starting token exchange\n", e.Code)

		tokens, err := requestTokens(e.Code, "authorization_code")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if tokens.AccessToken == "" {
			http.Error(w, "No Access token received", http.StatusNoContent)
			return
		}
		log.Printf("Tokens Received from provider using code=%+v, starting store\n", e.Code)
		w.WriteHeader(http.StatusCreated)
		res, err := storeTokens(&tokens)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotModified)
			return
		}
		log.Printf("Tokens Stored using code=%+v exiting\n", e.Code)
		msg := ResponseData{res}
		json.NewEncoder(w).Encode(msg)
	}
}

func RefreshTokens() error {
	tokens, err := getTokens()
	if err != nil {
		log.Fatalf("RefreshTokens.getTokens error=%+v\n", err)
		return err
	}
	ert, err := time.Parse("2006-01-02T15:04:05.999999999Z07:00", tokens.RefreshTokenExpiresAt)
	if err != nil {
		log.Fatalf("RefreshTokens.time.Parse Refresh Token error=%+v\n", err)
		return err
	}
	eat, err := time.Parse("2006-01-02T15:04:05.999999999Z07:00", tokens.AccessTokenExpiresAt)
	if err != nil {
		log.Fatalf("RefreshTokens.time.Parse Access Token error=%+v\n", err)
		return err
	}
	if eat.Sub(time.Now()) > 0*time.Second {
		log.Printf("Access Token is still valid")
		return nil
	}
	if ert.Sub(time.Now()) <= 0*time.Second {
		log.Fatalf("Refresh Token is expired")
		return errors.New("The received refresh token is already expired")
	}
	tr, err := requestTokens(tokens.RefreshToken, "refresh_token")
	if err != nil {
		log.Fatalf("RefreshTokens.requestTokens error=%+v\n", err)
		return err
	}
	log.Printf("refresh got:%+v\n", tr)
	_, err = storeTokens(&tr)
	if err != nil {
		log.Fatalf("RefreshTokens.storeTokens error=%+v\n", err)
		return err
	}
	return nil
}

func GetAccessToken() (string, error) {
	if err := RefreshTokens(); err != nil {
		log.Fatalf("GetAccessToken.RefreshTokens error=%+v\n", err)
		return "", err
	}
	tokens, err := getTokens()
	if err != nil {
		log.Fatalf("GetAccessToken.getTokens error=%+v\n", err)
		return "", err
	}
	if tokens.AccessToken == "" {
		log.Fatalf("GetAccessToken.AccessToken error= Access Token is empty")
		return "", errors.New("The Access Token received is blank")
	}
	return tokens.AccessToken, nil
}
