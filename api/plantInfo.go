package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PlantInfoResponse struct {
	Co2Reduction string            `json:"Co2Reduction"`
	DeviceList   []PlantInfoDevice `json:"deviceList"`
}

type PlantInfoDevice struct {
	PowerStr    string `json:"powerStr"`
	DeviceAlias string `json:"deviceAilas"` // Yes typo is intended
}

func (api *GrowattAPI) PlantInfo(id string) (PlantInfoResponse, error) {
	var resJson PlantInfoResponse

	finalURL := fmt.Sprintf("%s?op=getAllDeviceList&plantId=%s&pageNum=1&pageSize=1", api.GetUrl("newTwoPlantAPI.do"), id)

	req, err := http.NewRequest(http.MethodGet, finalURL, nil)
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

	return resJson, nil
}
