package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type LoginResponse struct {
	Back LoginBack `json:"back"`
}

type LoginBack struct {
	Msg         string      `json:"msg"`
	Data        []LoginData `json:"data"`
	DeviceCount string      `json:"deviceCount"`
	Success     bool        `json:"success"`
	User        LoginUser   `json:"user"`
}

type LoginData struct {
	PlantID   string `json:"plantId"`
	PlantName string `json:"plantName"`
}

type LoginUser struct {
	ID          int    `json:"id"`
	RightLevel  int    `json:"rightlevel"`
	AgentCode   string `json:"agentCode"`
	AccountName string `json:"accountName"`
	Email       string `json:"email"`
	PhoneNum    string `json:"phoneNum"`
}

func (api *GrowattAPI) Login() error {
	loginData := url.Values{
		"userName": {api.Username},
		"password": {api.HashedPassword},
	}
	req, err := http.NewRequest(http.MethodPost, api.GetUrl("newTwoLoginAPI.do"), strings.NewReader(loginData.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := api.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var resJson LoginResponse
	if err := json.Unmarshal(resBodyBytes, &resJson); err != nil {
		return err
	}

	if !resJson.Back.Success {
		return fmt.Errorf("something went wrong logging in: %+v", resJson)
	}

	api.UserID = resJson.Back.User.ID
	api.UserRightLevel = resJson.Back.User.RightLevel

	return nil
}
