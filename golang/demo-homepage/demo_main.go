package main

import (
	"fmt"

	"github.com/sawka/dashborg-go-sdk/pkg/dash"
)

func main() {
	cfg := &dash.Config{AnonAcc: true, AutoKeygen: true}
	client, err := dash.ConnectClient(cfg)
	if err != nil {
		fmt.Printf("Error connecting DashborgCloudClient: %v\n", err)
		return
	}
	app := client.AppClient().NewApp("default")
	app.SetOfflineAccess(true)
	app.SetAppTitle("Demo Homepage")
	app.WatchHtmlFile("./panels/demo-homepage.html", nil)
	err = client.AppClient().WriteApp(app)
	if err != nil {
		fmt.Printf("Error writing app: %v\n", err)
		return
	}
	select {}
}
