package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type PlantListResponse struct {
	Back PlantListBack `json:"back"`
}

type PlantListBack struct {
	TotalData PlantListTotalData `json:"totalData"`
	Data      []PlantListData    `json:"data"`
	Success   bool
}

type PlantListTotalData struct {
	CO2Sum          string `json:"CO2Sum"`
	CurrentPowerSum string `json:"currentPowerSum"`
	IsHaveStorage   string `json:"isHaveStorage"`
	TodayEnergySum  string `json:"todayEnergySum"`
	TotalEnergySum  string `json:"totalEnergySum"`
}

type PlantListData struct {
	IsHaveStorage string `json:"isHaveStorage"`
	CurrentPower  string `json:"currentPower"`
	TotalEnergy   string `json:"totalEnergy"`
	TodayEnergy   string `json:"todayEnergy"`
	PlantID       string `json:"plantId"`
	PlantName     string `json:"plantName"`
}

func (api *GrowattAPI) PlantList() (PlantListResponse, error) {
	var resJson PlantListResponse
	data := url.Values{"userId": {fmt.Sprint(api.UserID)}}
	req, err := http.NewRequest(http.MethodGet, api.GetUrl("PlantListAPI.do"), strings.NewReader(data.Encode()))
	if err != nil {
		return resJson, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := api.HttpClient.Do(req)
	if err != nil {
		return resJson, err
	}
	defer res.Body.Close()

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return resJson, err
	}

	if err := json.Unmarshal(resBodyBytes, &resJson); err != nil {
		return resJson, err
	}

	if !resJson.Back.Success {
		return resJson, fmt.Errorf("something went wrong pulling plant list info: %+v", resJson)
	}

	return resJson, nil
}
