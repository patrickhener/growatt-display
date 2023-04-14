package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/patrickhener/growatt-display/utils"
)

type PlantDetailBack struct {
	PlantData PlantDetailPlantData `json:"plantData"`
	Data      map[string]string    `json:"data"`
}

type PlantDetailPlantData struct {
	Energy string `json:"currentEnergy"`
}

type PlantDetailResponse struct {
	Back PlantDetailBack `json:"back"`
}

func (api *GrowattAPI) PlantDetail(id string, timespan string, time time.Time) (PlantDetailResponse, error) {
	var resJson PlantDetailResponse

	processedDate := utils.GetDateString(utils.Timespan(timespan), time)

	finalURL := fmt.Sprintf("%s?plantId=%s&type=%s&date=%s", api.GetUrl("PlantDetailAPI.do"), id, timespan, processedDate)

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
