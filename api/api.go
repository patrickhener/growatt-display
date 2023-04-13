package api

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
)

type GrowattAPI struct {
	URI            string
	Username       string
	HashedPassword string
	HttpClient     *http.Client
	UserID         int
	UserRightLevel int
}

func New(url, username, password string) *GrowattAPI {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Jar: jar,
	}
	return &GrowattAPI{
		URI:            url,
		Username:       username,
		HashedPassword: password,
		HttpClient:     client,
	}
}

func (api *GrowattAPI) Display() error {
	plantList, err := api.PlantList()
	if err != nil {
		return err
	}

	for _, p := range plantList.Back.Data {
		fmt.Printf("Plant '%s':\n", p.PlantName)
		fmt.Printf("\tTotal Energy Today: %s\n", p.TodayEnergy)
		fmt.Printf("\tTotal Energy All Time: %s\n", p.TotalEnergy)
		plantInfo, err := api.PlantInfo(p.PlantID)
		if err != nil {
			return err
		}
		fmt.Printf("\tTotal Co² reduction: %s kg\n", plantInfo.Co2Reduction)
		fmt.Println("")
		fmt.Println("Data collectors:")
		fmt.Println("")
		for _, d := range plantInfo.DeviceList {
			fmt.Printf("Collector '%s'\n", d.DeviceAlias)
			fmt.Printf("\tCurrent Power: %s\n", d.PowerStr)
			fmt.Printf("")
		}
	}

	return nil
}

func (api *GrowattAPI) GetUrl(page string) string {
	return fmt.Sprintf("%s%s", api.URI, page)
}
