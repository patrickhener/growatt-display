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

		// Current Data
		curData, err := api.GetCurrentData(p.PlantID)
		if err != nil {
			return nil, err
		}

		// Previous Data
		prevData, err := api.GetPreviousData(p.PlantID)
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
			TotalEnergyLastMonth: prevData.TotalEnergyLastMonth,
			TotalEnergyLastWeek:  prevData.TotalEnergyLastWeek,
			TotalEnergyYesterday: prevData.TotalEnergyYesterday,
			TotalEnergyToday:     strings.TrimRight(p.TodayEnergy, " kWh"),
			TotalEnergyThisWeek:  curData.TotalEnergyThisWeek,
			TotalEnergyThisMonth: curData.TotalEnergyThisMonth,
			TotalEnergyAllTime:   strings.TrimRight(p.TotalEnergy, " kWh"),
			PlantName:            p.PlantName,
			Devices:              devices,
		})
	}

	return resStats, nil
}

func (api *GrowattAPI) GetPreviousData(id string) (PreviousData, error) {
	var prevData PreviousData

	// Yesterday
	timeYesterday := time.Now().AddDate(0, 0, -1)
	yesterday, err := api.PlantDetail(id, string(utils.Day), timeYesterday)
	if err != nil {
		return prevData, err
	}

	// Last Week (monday - sunday)
	year, week := time.Now().ISOWeek()
	week = week - 1
	// Find Last Monday
	start, _ := utils.DateRange(week, year)
	mon := start.Day()
	sun := mon + 6

	// Last Month (all)
	lastMonth := time.Now().AddDate(0, -1, 0)
	lastMonthData, err := api.PlantDetail(id, string(utils.Month), lastMonth)
	if err != nil {
		return prevData, err
	}

	// Go over data map and fill both counters
	lastWeekEnergy := 0.0
	lastMonthEnergy := 0.0

	for d, e := range lastMonthData.Back.Data {
		dInt, err := strconv.Atoi(d)
		if err != nil {
			return prevData, err
		}
		eFloat, err := strconv.ParseFloat(e, 32)
		if err != nil {
			return prevData, err
		}

		// Conditonally calc to Week Energy
		for i := mon; i <= sun; i++ {
			if dInt == i {
				lastWeekEnergy = lastWeekEnergy + eFloat
			}
		}

		// Calc to total Energy
		lastMonthEnergy = lastWeekEnergy + eFloat
	}

	prevData.TotalEnergyLastWeek = fmt.Sprintf("%.1f", lastWeekEnergy)
	prevData.TotalEnergyLastMonth = fmt.Sprintf("%.1f", lastMonthEnergy)
	prevData.TotalEnergyYesterday = strings.TrimRight(yesterday.Back.PlantData.Energy, " kWh")

	return prevData, nil

}

func (api *GrowattAPI) GetCurrentData(id string) (CurrentData, error) {
	var curData CurrentData

	// This Week (monday - sunday)
	year, week := time.Now().ISOWeek()
	// Find This Monday
	start, _ := utils.DateRange(week, year)
	mon := start.Day()
	sun := mon + 6
	// Fetch This month
	thisMonthData, err := api.PlantDetail(id, string(utils.Month), time.Now())
	if err != nil {
		return curData, err
	}

	// Go over data map and fill both counters
	thisWeekEnergy := 0.0
	thisMonthEnergy := 0.0
	for d, e := range thisMonthData.Back.Data {
		dInt, err := strconv.Atoi(d)
		if err != nil {
			return curData, err
		}
		eFloat, err := strconv.ParseFloat(e, 32)
		if err != nil {
			return curData, err
		}

		// Conditonally calc to Week Energy
		for i := mon; i <= sun; i++ {
			if dInt == i {
				thisWeekEnergy = thisWeekEnergy + eFloat
			}
		}

		// Calc to total Energy
		thisMonthEnergy = thisMonthEnergy + eFloat
	}

	curData.TotalEnergyThisMonth = fmt.Sprintf("%.1f", thisMonthEnergy)
	curData.TotalEnergyThisWeek = fmt.Sprintf("%.1f", thisWeekEnergy)

	return curData, nil
}
