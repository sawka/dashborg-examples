package main

import (
	"fmt"
	"math/rand"
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

type PanelState struct {
	SelectedAccId string
	CreateData    CreateAccountFormData `json:"create" mapstructure:"create"`
}

type PanelModel struct {
	AccModel *AccModel
}

func (pm *PanelModel) RootHandler(req *dash.PanelRequest) error {
	if !req.CheckAuth(dash.AuthPassword{"hello"}) {
		return nil
	}
	err := req.SetHtmlFromFile("panels/accdemo.html")
	if err != nil {
		return err
	}
	return nil
}

func (pm *PanelModel) Upgrade(req *dash.PanelRequest, panelState interface{}, accId string) error {
	if accId == "" {
		return fmt.Errorf("No AccountId Selected")
	}
	pm.AccModel.Upgrade(accId)
	req.InvalidateData("/accounts/.*")
	return nil
}

func (pm *PanelModel) Downgrade(req *dash.PanelRequest, panelState interface{}, accId string) error {
	if accId == "" {
		return fmt.Errorf("No AccountId Selected")
	}
	pm.AccModel.Downgrade(accId)
	req.InvalidateData("/accounts/.*")
	return nil
}

func (pm *PanelModel) GetAccountsList(req *dash.PanelRequest) (interface{}, error) {
	accList := pm.AccModel.CopyAccList()
	return accList, nil
}

func (pm *PanelModel) GetAccount(req *dash.PanelRequest, panelState interface{}, accId string) (interface{}, error) {
	acc := pm.AccModel.AccById(accId)
	if acc == nil {
		return nil, nil
	}
	return acc, nil
}

func (pm *PanelModel) RefreshAccList(req *dash.PanelRequest) error {
	req.SetData("$state.selaccid", nil)
	req.InvalidateData("/accounts/.*")
	return nil
}

func (pm *PanelModel) RegenAccList(req *dash.PanelRequest) error {
	pm.AccModel.RegenAccounts()
	req.SetData("$state.selaccid", nil)
	req.InvalidateData("/accounts/.*")
	return nil
}

func (pm *PanelModel) RemoveAccount(req *dash.PanelRequest, panelState interface{}, accId string) error {
	if accId == "" {
		return fmt.Errorf("No AccountId Selected")
	}
	pm.AccModel.RemoveAcc(accId)
	req.InvalidateData("/accounts/.*")
	req.SetData("$state.selaccid", nil)
	return nil
}

func (pm *PanelModel) CreateAccount(req *dash.PanelRequest, panelState PanelState) error {
	errors := panelState.CreateData.Validate()
	if len(errors) > 0 {
		req.SetData("$state.create.errors", errors)
		return nil
	}
	req.SetData("$state.create.errors", nil)
	newAccId := pm.AccModel.CreateAcc(panelState.CreateData.Name, panelState.CreateData.Email)
	req.SetData("$state.createAccountModal", false)
	req.SetData("$state.selaccid", newAccId)
	req.InvalidateData("/accounts/.*")
	return nil
}

func main() {
	rand.Seed(time.Now().Unix())
	cfg := &dash.Config{ProcName: "acc-demo", AnonAcc: true, AutoKeygen: true}
	dash.StartProcClient(cfg)
	defer dash.WaitForClear()

	accModel := MakeAccModel()
	panel := &PanelModel{AccModel: accModel}
	dash.RegisterPanelHandlerEx("acc-demo", "/", panel.RootHandler)
	dash.RegisterPanelHandlerEx("acc-demo", "/acc/upgrade", panel.Upgrade)
	dash.RegisterPanelHandlerEx("acc-demo", "/acc/downgrade", panel.Downgrade)
	dash.RegisterDataHandlerEx("acc-demo", "/accounts/list", panel.GetAccountsList)
	dash.RegisterDataHandlerEx("acc-demo", "/accounts/get", panel.GetAccount)
	dash.RegisterPanelHandlerEx("acc-demo", "/acc/refresh-accounts", panel.RefreshAccList)
	dash.RegisterPanelHandlerEx("acc-demo", "/acc/regen-acclist", panel.RegenAccList)
	dash.RegisterPanelHandlerEx("acc-demo", "/acc/remove", panel.RemoveAccount)
	dash.RegisterPanelHandlerEx("acc-demo", "/acc/create-account", panel.CreateAccount)

	select {}
}
