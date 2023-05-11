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

func New(url, username, password string) (*GrowattAPI, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return &GrowattAPI{}, err
	}
	client := &http.Client{
		Jar: jar,
	}
	return &GrowattAPI{
		URI:            url,
		Username:       username,
		HashedPassword: password,
		HttpClient:     client,
	}, nil
}

func (api *GrowattAPI) Display() error {
	empty := []string{"", "", "", ""}
	data := [][]string{}

	plantList, err := api.PlantList()
	if err != nil {
		return err
	}

	for _, p := range plantList {
		data = append(data, []string{fmt.Sprintf("Plant '%s'", p.PlantName), "Total Energy Last Month", p.TotalEnergyLastMonth, "kWh"})
		data = append(data, []string{"", "Total Energy Last Week", p.TotalEnergyLastWeek, "kWh"})
		data = append(data, []string{"", "Total Energy Yesterday", p.TotalEnergyYesterday, "kWh"})
		data = append(data, empty)
		data = append(data, []string{"", "Total Energy Today", p.TotalEnergyToday, "kWh"})
		data = append(data, []string{"", "Total Energy This Week", p.TotalEnergyThisWeek, "kWh"})
		data = append(data, []string{"", "Total Energy This Month", p.TotalEnergyThisMonth, "kWh"})
		data = append(data, empty)
		data = append(data, []string{"", "Total Energy All Time", p.TotalEnergyAllTime, "kWh"})
		data = append(data, empty)
		for _, d := range p.Devices {
			data = append(data, []string{fmt.Sprintf("Collector '%s'", d.Alias), "Current Power", d.CurrentPower, "kW"})
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
