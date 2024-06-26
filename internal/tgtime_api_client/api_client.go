package tgtime_api_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robertobadjio/tgtime-aggregator/internal/config"
	"net/http"
	"time"
)

type tgTimeClient struct {
	HTTPClient *http.Client
	baseURL    string
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

/*type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}*/

var authData Tokens

func NewTgTimeClient() *tgTimeClient {
	cfg := config.New()
	return &tgTimeClient{
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
		baseURL: cfg.ApiURL,
	}
}

func (otc *tgTimeClient) sendRequest(request *http.Request, v interface{}) error {
	request.Header.Add("Content-Type", "application/json")

	if request.URL.String() != otc.baseURL+"/login" {
		if authData.AccessTokenExpires <= uint64(time.Now().Unix()) {
			c := NewTgTimeClient()
			data, err := c.Login()
			if err != nil {
				panic(err)
			}
			// TODO: Если не истек refresh token
			//data := Login()
			authData.AccessToken = data.AccessToken
			authData.RefreshToken = data.RefreshToken
			authData.AccessTokenExpires = data.AccessTokenExpires
			authData.RefreshTokenExpires = data.RefreshTokenExpires
		}

		request.Header.Add("Token", authData.AccessToken)
	}

	response, err := otc.HTTPClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(response.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", response.StatusCode)
	}

	/*fullResponse := successResponse{
		Data: v,
	}*/

	if err = json.NewDecoder(response.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}
