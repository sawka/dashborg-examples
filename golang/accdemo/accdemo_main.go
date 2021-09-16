package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"time"

	"github.com/sawka/dashborg-go-sdk/pkg/dash"
)

var EmailRe = regexp.MustCompile("^[a-zA-Z0-9.-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type CreateAccountFormData struct {
	Name  string `json:"name" mapstructure:"name"`
	Email string `json:"email" mapstructure:"email"`
}

func (d CreateAccountFormData) Validate() map[string]string {
	errors := make(map[string]string)
	if d.Name == "" {
		errors["name"] = "Name must not be empty"
	} else if len(d.Name) > 40 {
		errors["name"] = "Name can only be 40 characters"
	}
	if d.Email == "" {
		errors["email"] = "Email must not be empty"
	} else if len(d.Email) > 40 {
		errors["email"] = "Email can only be 40 characters"
	} else if !EmailRe.MatchString(d.Email) {
		errors["email"] = "Email format is not correct"
	}
	return errors
}

type AppState struct {
	SelectedAccId string
	CreateData    CreateAccountFormData `json:"create" mapstructure:"create"`
}

type AppModel struct {
	AccModel *AccModel
}

func (pm *AppModel) Upgrade(req *dash.AppRequest, accId string) error {
	if accId == "" {
		return fmt.Errorf("No AccountId Selected")
	}
	pm.AccModel.Upgrade(accId)
	req.InvalidateData(".*")
	return nil
}

func (pm *AppModel) Downgrade(req *dash.AppRequest, accId string) error {
	if accId == "" {
		return fmt.Errorf("No AccountId Selected")
	}
	pm.AccModel.Downgrade(accId)
	req.InvalidateData(".*")
	return nil
}

func (pm *AppModel) GetAccountsList(req dash.Request) (interface{}, error) {
	accList := pm.AccModel.CopyAccList()
	return accList, nil
}

func (pm *AppModel) GetAccount(req dash.Request, accId string) (interface{}, error) {
	acc := pm.AccModel.AccById(accId)
	if acc == nil {
		return nil, nil
	}
	return acc, nil
}

func (pm *AppModel) RefreshAccList(req *dash.AppRequest) error {
	req.SetData("$state.selaccid", nil)
	req.InvalidateData(".*")
	return nil
}

func (pm *AppModel) RegenAccList(req *dash.AppRequest) error {
	pm.AccModel.RegenAccounts()
	req.SetData("$state.selaccid", nil)
	req.InvalidateData(".*")
	return nil
}

func (pm *AppModel) RemoveAccount(req *dash.AppRequest, accId string) error {
	if accId == "" {
		return fmt.Errorf("No AccountId Selected")
	}
	pm.AccModel.RemoveAcc(accId)
	req.InvalidateData(".*")
	req.SetData("$state.selaccid", nil)
	return nil
}

func (pm *AppModel) CreateAccount(req *dash.AppRequest, appState AppState) error {
	errors := appState.CreateData.Validate()
	if len(errors) > 0 {
		req.SetData("$state.create.errors", errors)
		return nil
	}
	req.SetData("$state.create.errors", nil)
	newAccId := pm.AccModel.CreateAcc(appState.CreateData.Name, appState.CreateData.Email)
	req.SetData("$state.createAccountModal", false)
	req.SetData("$state.selaccid", newAccId)
	req.InvalidateData(".*")
	return nil
}

func main() {
	rand.Seed(time.Now().Unix())
	cfg := &dash.Config{AnonAcc: true, AutoKeygen: true}
	client, err := dash.ConnectClient(cfg)
	if err != nil {
		fmt.Printf("Error connecting DashborgCloudClient: %v\n", err)
		return
	}
	accModel := MakeAccModel()
	panel := &AppModel{AccModel: accModel}

	app := client.AppClient().NewApp("acc-demo")
	app.WatchHtmlFile("panels/accdemo.html", nil)
	app.Runtime().SetAppStateType(reflect.TypeOf(AppState{}))
	app.Runtime().PureHandler("get-accounts-list", panel.GetAccountsList)
	app.Runtime().PureHandler("get-account", panel.GetAccount)
	app.Runtime().Handler("acc-upgrade", panel.Upgrade)
	app.Runtime().Handler("acc-downgrade", panel.Downgrade)
	app.Runtime().Handler("regen-acclist", panel.RegenAccList)
	app.Runtime().Handler("acc-remove", panel.RemoveAccount)
	app.Runtime().Handler("refresh-accounts", panel.RefreshAccList)
	app.Runtime().Handler("create-account", panel.CreateAccount)
	err = client.AppClient().WriteAndConnectApp(app)
	if err != nil {
		fmt.Printf("Error writing app: %v\n", err)
		return
	}

	client.WaitForShutdown()
}
