package restream

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type ResponseData struct {
	Message string `json:"msg"`
}

func ExchangeTokens(c echo.Context) error {
	code := c.FormValue("code")
	logger := log.WithFields(log.Fields{"source": "restream.ExchangeTokens()", "code": code})
	logger.Debugln("token exchange started")

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("Content-Type", "application/json")

	tokens, err := requestTokens(code, "authorization_code")
	if err != nil {
		logger.Errorln(err)
		c.Error(err)
		return err
	}
	log.WithField("tokens", tokens).Infoln("received tokens from code exchange")

	if tokens.AccessToken == "" {
		logger.Errorln("empty access token received")
		c.String(http.StatusInternalServerError, "no access token received form exchange.")
		return err
	}

	c.Response().WriteHeader(http.StatusCreated)
	logger.Debugln("started tokens store process")
	res, err := storeTokens(&tokens)
	if err != nil {
		logger.Errorln(err)
		c.Error(err)
		return err
	}
	logger.Debugln("tokens stored successfully")
	msg := ResponseData{res}
	logger.Debugln("token exchange finished")
	return c.JSON(http.StatusOK, msg)
}

func RefreshTokens(ttr time.Duration) error {
	logger := log.WithFields(log.Fields{"source": "restream.RefreshTokens()"})
	logger.Debugln("token refresh started")
	tokens, err := getTokens()
	if err != nil {
		logger.Errorln(err)
		return err
	}
	log.WithField("tokens", tokens).Infoln("received tokens from database")
	ert, err := time.Parse("2006-01-02T15:04:05.999999999Z07:00", tokens.RefreshTokenExpiresAt)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	log.WithField("refreshToken", tokens.RefreshTokenExpiresAt).Debugln("refresh token is still valid")
	eat, err := time.Parse("2006-01-02T15:04:05.999999999Z07:00", tokens.AccessTokenExpiresAt)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	if time.Until(eat) > ttr {
		logger.Printf("Access Token is still valid at next ttr")
		return nil
	}

	if time.Until(ert) <= 0*time.Second {
		logger.Fatalf("Refresh Token is expired")
		return errors.New("The received refresh token is already expired")
	}
	tr, err := requestTokens(tokens.RefreshToken, "refresh_token")
	if err != nil {
		logger.Fatalf("RefreshTokens.requestTokens error=%+v\n", err)
		return err
	}
	logger.Printf("refresh got:%+v\n", tr)
	_, err = storeTokens(&tr)
	if err != nil {
		logger.Fatalf("RefreshTokens.storeTokens error=%+v\n", err)
		return err
	}
	logger.Debugln("token refresh finished")
	return nil
}

func GetAccessToken() (string, error) {
	logger := log.WithFields(log.Fields{"source": "restream.GetAccessToken()"})
	logger.Debugln("token get started")
	tokens, err := getTokens()
	if err != nil {
		logger.Errorln(err)
		return "", err
	}
	if tokens.AccessToken == "" {
		logger.Errorln("received access token is empty")
		return "", errors.New("he Access Token received is blank")
	}
	logger.Debugln("token get finished")
	return tokens.AccessToken, nil
}
