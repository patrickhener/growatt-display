package api

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/olekukonko/tablewriter"
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
	data := [][]string{}

	plantList, err := api.PlantList()
	if err != nil {
		return err
	}

	for _, p := range plantList {
		data = append(data, []string{fmt.Sprintf("Plant '%s'", p.PlantName), "Total Energy Last Month", p.TotalEnergyLastMonth})
		data = append(data, []string{"", "Total Energy Last Week", p.TotalEnergyLastWeek})
		data = append(data, []string{"", "Total Energy Yesterday", p.TotalEnergyYesterday})
		data = append(data, []string{"", "Total Energy Today", p.TotalEnergyToday})
		data = append(data, []string{"", "Total Energy This Week", p.TotalEnergyThisWeek})
		data = append(data, []string{"", "Total Energy This Month", p.TotalEnergyThisMonth})
		data = append(data, []string{"", "Total Energy All Time", p.TotalEnergyAllTime})
		for _, d := range p.Devices {
			data = append(data, []string{fmt.Sprintf("Collector '%s'", d.Alias), "Current Power", d.CurrentPower})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.AppendBulk(data)
	table.Render()

	return nil
}

func (api *GrowattAPI) GetUrl(page string) string {
	return fmt.Sprintf("%s%s", api.URI, page)
}
