package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sawka/dashborg-go-sdk/pkg/dash"
)

type SimpleQueryModel struct {
	DB *sqlx.DB
}

type Company struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	CompanyUrl  string `json:"companyurl"`
	LogoUrl     string `json:"logourl"`
	Rank        int    `json:"rank"`
	Overview    string `json:"overview"`
	Sector      string `json:"sector"`
	JobsCreated int    `json:"jobscreated"`
	Batch       string `json:"batch"`
	HQLocation  string `json:"hqlocation"`
	JobsUrl     string `json:"jobsurl"`
}

func (m *SimpleQueryModel) GetRows(req dash.Request) (interface{}, error) {
	var rtn []Company
	err := m.DB.SelectContext(req.Context(), &rtn, "SELECT * from companies")
	if err != nil {
		return nil, err
	}
	return rtn, nil
}

func main() {
	db, err := sqlx.Open("sqlite3", "./data/yc-companies.db")
	if err != nil {
		fmt.Printf("Error opening sqlitedb: %v\n", err)
		return
	}
	model := &SimpleQueryModel{
		DB: db,
	}
	config := &dash.Config{AnonAcc: true}
	client, err := dash.ConnectClient(config)
	if err != nil {
		fmt.Printf("Error connecting client: %v\n", err)
		return
	}
	app := client.AppClient().NewApp("simplequery")
	app.SetAppTitle("Simple Query")
	app.SetOfflineAccess(false)
	app.WatchHtmlFile("./panels/simplequery.html", nil)
	app.Runtime().PureHandler("get-rows", model.GetRows)
	err = client.AppClient().WriteAndConnectApp(app)
	if err != nil {
		fmt.Printf("Error writing app: %v\n", err)
		return
	}
	client.WaitForShutdown()
}
