package restream

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

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

func requestTokens(payload string, grant_type string) (TokenResponse, error) {
	logger := log.WithFields(log.Fields{"source": "restream.requestTokens()", "payload": payload, "type": grant_type})
	logger.Debugln("token request started")
	tr := TokenResponse{}
	endpoint := os.Getenv("RESTREAM_TOKEN_ENDPOINT")
	redirect_uri := os.Getenv("RESTREAM_REDIRECT_URI")
	client_id := os.Getenv("RESTREAM_CLIENT_ID")
	secret := os.Getenv("RESTREAM_SECRET")

	data := url.Values{}
	if grant_type == "authorization_code" {
		data.Set("grant_type", "authorization_code")
		data.Set("redirect_uri", redirect_uri)
		data.Set("code", payload)
	}
	if grant_type == "refresh_token" {
		data.Set("grant_type", "refresh_token")
		data.Set("refresh_token", payload)
	}
	client := &http.Client{}
	logger.WithFields(log.Fields{"method": http.MethodPost, "endpoint": endpoint, "body": data.Encode()}).Infoln("started request")
	r, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalf("http.NewRequest() error=%+v\n", err)
		return tr, err
	}

	r.SetBasicAuth(client_id, secret)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(r)
	if err != nil {
		logger.Errorln(err)
		return tr, err
	}
	if res.StatusCode == http.StatusBadRequest {
		logger.WithFields(log.Fields{"status-code": res.StatusCode, "body": res.Body, "header": res.Header}).Errorln("request rejected")
		return tr, err
	}
	defer res.Body.Close()
	logger.Debugln("started decode of res.body")
	if err := json.NewDecoder(res.Body).Decode(&tr); err != nil {
		logger.Errorln(err)
		return tr, err
	}
	logger.Debugln("token request finished")
	return tr, nil
}
