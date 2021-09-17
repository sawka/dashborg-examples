package main

import (
	"fmt"
	"regexp"

	"github.com/sawka/dashborg-go-sdk/pkg/dash"
)

func TestClick(req *dash.AppRequest) error {
	req.SetData("$.clicked", true)
	return nil
}

var ZIP_RE = regexp.MustCompile("^[0-9]{5}$")

func ClickZipCode(req *dash.AppRequest, data string) error {
	if data == "" || !ZIP_RE.MatchString(data) {
		req.SetData("$.ziperror", "Please Enter a Valid Zip Code")
		return nil
	}
	req.SetData("$.ziperror", nil)
	req.SetData("append:$.tabledata", data)
	req.SetData("$state.zipcode", "")
	return nil
}

func main() {
	config := &dash.Config{AutoKeygen: true, AnonAcc: true}
	client, err := dash.ConnectClient(config)
	if err != nil {
		fmt.Printf("Error connecting client: %v\n", err)
		return
	}
	app := client.AppClient().NewApp("bulma")
	app.SetAppTitle("Bulma UI Demo")
	app.SetHtmlFromFile("panels/bulma.html")
	app.Runtime().Handler("test-click", TestClick)
	app.Runtime().Handler("test-zip", ClickZipCode)
	err = client.AppClient().WriteAndConnectApp(app)
	if err != nil {
		fmt.Printf("Error writing app: %v\n", err)
		return
	}

	select {}
}
