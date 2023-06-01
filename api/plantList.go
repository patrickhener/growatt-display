package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/patrickhener/growatt-display/utils"
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
	CO2Sum string `json:"CO2Sum"`
}

type PlantListData struct {
	IsHaveStorage string `json:"isHaveStorage"`
	CurrentPower  string `json:"currentPower"`
	TotalEnergy   string `json:"totalEnergy"`
	TodayEnergy   string `json:"todayEnergy"`
	PlantID       string `json:"plantId"`
	PlantName     string `json:"plantName"`
}

type ResponseDevice struct {
	Alias        string
	CurrentPower string
}

type ResponseStats struct {
	PlantName            string
	TotalEnergyLastMonth string
	TotalEnergyLastWeek  string
	TotalEnergyYesterday string
	TotalEnergyToday     string
	TotalEnergyThisWeek  string
	TotalEnergyThisMonth string
	TotalEnergyAllTime   string
	CurrentPower         string
	Devices              []ResponseDevice
}

type CurrentData struct {
	TotalEnergyThisWeek  string
	TotalEnergyThisMonth string
}
type PreviousData struct {
	TotalEnergyLastMonth string
	TotalEnergyLastWeek  string
	TotalEnergyYesterday string
}

func (api *GrowattAPI) PlantList() ([]ResponseStats, error) {
	var resStats []ResponseStats

	// Init Plant List Request to fetch plantId
	var resJson PlantListResponse
	data := url.Values{"userId": {fmt.Sprint(api.UserID)}}
	req, err := http.NewRequest(http.MethodGet, api.GetUrl("PlantListAPI.do"), strings.NewReader(data.Encode()))
	if err != nil {
		return resStats, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := api.HttpClient.Do(req)
	if err != nil {
		return resStats, err
	}
	defer res.Body.Close()

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return resStats, err
	}

	if err := json.Unmarshal(resBodyBytes, &resJson); err != nil {
		return resStats, err
	}

	if !resJson.Back.Success {
		return resStats, fmt.Errorf("something went wrong pulling plant list info: %+v", resJson)
	}

	for _, p := range resJson.Back.Data {
		var devices []ResponseDevice

		// current year and week
		year, week := time.Now().ISOWeek()

		// Current Data (actual Month)
		curData, err := api.getCurrentMonthData(p.PlantID)
		if err != nil {
			return nil, err
		}

		// Previous Data (previous Month)
		prevData, err := api.getLastMonthData(p.PlantID)
		if err != nil {
			return nil, err
		}

		// Get last month
		lastMonthTotalEnergy, err := api.getMonthEnergy(prevData.Back.Data)
		if err != nil {
			return nil, err
		}

		// Get last week
		lastWeekTotalEnergy, err := api.getWeekEnergy(prevData, curData, week-1, year)
		if err != nil {
			return nil, err
		}

		// Get yesterday
		yesterdayTotalEnergy, err := api.getYesterday(time.Now().AddDate(0, 0, -1), p.PlantID)
		if err != nil {
			return nil, err
		}

		// Get this week
		thisWeekTotalEnergy, err := api.getWeekEnergy(prevData, curData, week, year)
		if err != nil {
			return nil, err
		}

		// Get this month
		thisMonthTotalEnergy, err := api.getMonthEnergy(curData.Back.Data)
		if err != nil {
			return nil, err
		}

		// General plant info
		plantInfo, err := api.PlantInfo(p.PlantID)
		if err != nil {
			return resStats, err
		}

		for _, d := range plantInfo.DeviceList {
			devices = append(devices, ResponseDevice{
				Alias:        d.DeviceAlias,
				CurrentPower: strings.TrimRight(d.PowerStr, "kW"),
			})
		}

		resStats = append(resStats, ResponseStats{
			TotalEnergyLastMonth: lastMonthTotalEnergy,
			TotalEnergyLastWeek:  lastWeekTotalEnergy,
			TotalEnergyYesterday: yesterdayTotalEnergy,
			TotalEnergyToday:     strings.TrimRight(p.TodayEnergy, " kWh"),
			TotalEnergyThisWeek:  thisWeekTotalEnergy,
			TotalEnergyThisMonth: thisMonthTotalEnergy,
			TotalEnergyAllTime:   strings.TrimRight(p.TotalEnergy, " kWh"),
			PlantName:            p.PlantName,
			Devices:              devices,
		})
	}

	return resStats, nil
}

func (api *GrowattAPI) getLastMonthData(id string) (PlantDetailResponse, error) {
	// Last Month (all)
	lastMonth := time.Now().AddDate(0, -1, 0)
	lastMonthData, err := api.PlantDetail(id, string(utils.Month), lastMonth)
	if err != nil {
		return PlantDetailResponse{}, err
	}

	return lastMonthData, nil
}

func (api *GrowattAPI) getCurrentMonthData(id string) (PlantDetailResponse, error) {
	// Fetch This month
	thisMonthData, err := api.PlantDetail(id, string(utils.Month), time.Now())
	if err != nil {
		return PlantDetailResponse{}, err
	}

	return thisMonthData, nil
}

func (api *GrowattAPI) getMonthEnergy(data map[string]string) (string, error) {
	var counter float64

	for _, e := range data {
		eFloat, err := strconv.ParseFloat(e, 32)
		if err != nil {
			return "", err
		}

		// Calc to total Energy
		counter = counter + eFloat
	}
	return fmt.Sprintf("%.1f", counter), nil
}

func (api *GrowattAPI) getYesterday(timeYesterday time.Time, id string) (string, error) {

	yesterday, err := api.PlantDetail(id, string(utils.Day), timeYesterday)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(yesterday.Back.PlantData.Energy, " kWh"), nil
}

func (api *GrowattAPI) getWeekEnergy(prevData, curData PlantDetailResponse, week, year int) (string, error) {
	var weekEnergy float64

	// Get week start end
	start, end, throughMonth, curPrev, err := utils.GetWeekStartEnd(year, week)
	if err != nil {
		return "", err
	}

	if throughMonth {
		// Previous months last day
		lengthPrevMonth := utils.DaysInMonth(start.Month(), year)

		// From start.Day() to Month end in prevData
		out1, err := utils.AddKilowatts(start.Day(), lengthPrevMonth, prevData.Back.Data)
		if err != nil {
			return "", err
		}
		out2, err := utils.AddKilowatts(1, end.Day(), curData.Back.Data)
		if err != nil {
			return "", err
		}
		weekEnergy = out1 + out2
	} else {
		switch curPrev {
		case "cur":
			weekEnergy, err = utils.AddKilowatts(start.Day(), end.Day(), curData.Back.Data)
			if err != nil {
				return "", err
			}
		case "prev":
			weekEnergy, err = utils.AddKilowatts(start.Day(), end.Day(), prevData.Back.Data)
			if err != nil {
				return "", err
			}
		}
	}

	return fmt.Sprintf("%.1f", weekEnergy), nil
}
