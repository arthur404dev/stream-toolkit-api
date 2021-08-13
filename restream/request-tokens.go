package restream

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type TokenResponse struct {
	AccessToken              string   `json:"accessToken"`
	AccessTokenExpiresIn     int64    `json:"accessTokenExpiresIn"`
	AccessTokenExpiresAt     string   `json:"accessTokenExpiresAt"`
	AccessTokenExpiresEpoch  int64    `json:"accessTokenExpiresEpoch"`
	RefreshToken             string   `json:"refreshToken"`
	RefreshTokenExpiresIn    int64    `json:"refreshTokenExpiresIn"`
	RefreshTokenExpiresAt    string   `json:"refreshTokenExpiresAt"`
	RefreshTokenExpiresEpoch int64    `json:"refreshTokenExpiresEpoch"`
	Scopes                   []string `json:"scopeJson"`
	TokenType                string   `json:"tokenType"`
}

func requestTokens(c string) (TokenResponse, error) {
	endpoint := os.Getenv("RESTREAM_TOKEN_ENDPOINT")
	redirect_uri := os.Getenv("RESTREAM_REDIRECT_URI")
	client_id := os.Getenv("RESTREAM_CLIENT_ID")
	secret := os.Getenv("RESTREAM_SECRET")

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirect_uri)
	data.Set("code", c)

	client := &http.Client{}
	r, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	r.SetBasicAuth(client_id, secret)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var tr TokenResponse

	if err := json.NewDecoder(res.Body).Decode(&tr); err != nil {
		log.Fatal(err)
		return tr, err
	}

	return tr, nil
}
