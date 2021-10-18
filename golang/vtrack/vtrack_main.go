package main

import (
	"fmt"
	"strings"

	"github.com/sawka/dashborg-go-sdk/pkg/dash"
)

type Vehicle struct {
	Id      int    `json:"id"`
	Make    string `json:"make"`
	Model   string `json:"model"`
	Year    int    `json:"year"`
	License string `json:"license"`
}

func (v *Vehicle) TrimFields() {
	v.Make = strings.TrimSpace(v.Make)
	v.Model = strings.TrimSpace(v.Model)
	v.License = strings.TrimSpace(v.License)
}

func (v *Vehicle) Validate() error {
	v.TrimFields()
	if v.Make == "" {
		return fmt.Errorf("Please enter a valid Make")
	}
	if v.Model == "" {
		return fmt.Errorf("Please enter a valid Model")
	}
	if v.Year < 1900 || v.Year > 2050 {
		return fmt.Errorf("Please enter a valid Year")
	}
	if v.License == "" {
		return fmt.Errorf("Please enter a valid License")
	}
	return nil
}

type Model struct {
	Vehicles []Vehicle
	NextId   int
}

func (m *Model) GetAllVehicles() (interface{}, error) {
	return m.Vehicles, nil
}

func (m *Model) AddVehicle(newVehicle Vehicle) error {
	err := newVehicle.Validate()
	if err != nil {
		return err
	}
	newVehicle.Id = m.NextId
	m.NextId++
	m.Vehicles = append(m.Vehicles, newVehicle)
	return nil
}

func (m *Model) EditVehicle(vehicle Vehicle) error {
	err := vehicle.Validate()
	if err != nil {
		return err
	}
	if vehicle.Id == 0 {
		return fmt.Errorf("edit-vehicle requires an id")
	}
	found := false
	for idx, v := range m.Vehicles {
		if v.Id == vehicle.Id {
			found = true
			m.Vehicles[idx] = vehicle
			continue
		}
	}
	if !found {
		return fmt.Errorf("Vehicle with id=%d not found\n", vehicle.Id)
	}
	return nil
}

func (m *Model) DeleteVehicle(id int) error {
	var newList []Vehicle
	for _, v := range m.Vehicles {
		if v.Id == id {
			continue
		}
		newList = append(newList, v)
	}
	m.Vehicles = newList
	return nil
}

func main() {
	cfg := &dash.Config{AnonAcc: true, AutoKeygen: true}
	client, err := dash.ConnectClient(cfg)
	if err != nil {
		panic(err)
	}
	m := Model{
		NextId:   2,
		Vehicles: []Vehicle{Vehicle{Id: 1, Make: "Honda", Model: "Accord", Year: 1998, License: "2VEH123"}},
	}
	app := client.AppClient().NewApp("vtrack")
	app.WatchHtmlFile("panels/vtrack.html", nil)
	app.SetInitRequired(true)
	app.SetPagesEnabled(true)
	app.Runtime().PureHandler("add-vehicle", m.AddVehicle)
	app.Runtime().PureHandler("get-all-vehicles", m.GetAllVehicles)
	app.Runtime().PureHandler("delete-vehicle", m.DeleteVehicle)
	app.Runtime().PureHandler("edit-vehicle", m.EditVehicle)
	err = client.AppClient().WriteAndConnectApp(app)
	if err != nil {
		panic(err)
	}
	client.WaitForShutdown()
}
