package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sawka/dashborg-go-sdk/pkg/dash"
)

var EmailStrs = []string{"mike", "matt", "michelle", "pat", "jim", "marc", "andrew", "alan"}
var EmailDomains = []string{"nomail.net", "gmail.com", "hotmail.com", "outlook.com", "msn.com"}

func RandomWord(list []string) string {
	return list[rand.Int31n(int32(len(list)))]
}

// RandomWord(EmailStrs) + strconv.Itoa(int(rand.Int31n(70)+10)) + "@nomail.com",

type CreditModel struct {
	Users []*UserType
}

type PanelState struct {
	SearchEmail string
}

type UserType struct {
	UserId     int
	Email      string
	NumCredits int
}

func MakeCreditModel() *CreditModel {
	rtn := &CreditModel{}
	for i := 0; i < 50; i++ {
		email := RandomWord(EmailStrs) + strconv.Itoa(int(rand.Int31n(70)+10)) + "@" + RandomWord(EmailDomains)
		user := &UserType{UserId: i + 1, Email: email, NumCredits: int(rand.Int31n(26))}
		rtn.Users = append(rtn.Users, user)
	}
	return rtn
}

func (m *CreditModel) RootHandler(req *dash.PanelRequest) error {
	if !req.CheckAuth(dash.AuthPassword{Password: "hello"}) {
		return nil
	}
	req.SetHtmlFromFile("panels/creditcount.html")
	return nil
}

type updateCreditsParams struct {
	UserId     int
	NumCredits int
	Index      int
}

func (m *CreditModel) UpdateCredits(req *dash.PanelRequest, state *PanelState, params *updateCreditsParams) error {
	if params.UserId <= 0 || params.NumCredits < 0 || params.NumCredits > 100 {
		return fmt.Errorf("Invalid Params")
	}
	var foundUser *UserType
	for i := 0; i < len(m.Users); i++ {
		user := m.Users[i]
		if user.UserId == params.UserId {
			user.NumCredits = params.NumCredits
			foundUser = user
			break
		}
	}
	if foundUser == nil {
		return fmt.Errorf("User not found")
	}
	req.SetData("$state.updateemail", nil)
	userPath := fmt.Sprintf("$.users[%d].NumCredits", params.Index)
	req.SetData(userPath, params.NumCredits)
	// req.InvalidateData("/get-users")
	return nil
}

type SortSpec struct {
	Column string `json:"column"`
	Asc    bool   `json:"asc"`
}

type getUsersParams struct {
	SearchEmail string
	SortSpec    SortSpec
}

func (m *CreditModel) GetUsers(req *dash.PanelRequest, state *PanelState, params *getUsersParams) (interface{}, error) {
	var rtn []*UserType

	if params.SearchEmail == "" {
		rtn = m.Users
	} else {
		for _, user := range m.Users {
			if strings.Contains(user.Email, params.SearchEmail) {
				rtn = append(rtn, user)
			}
		}
	}

	if params.SortSpec.Column == "" {
		return rtn, nil
	}
	sort.Slice(rtn, func(i int, j int) bool {
		v1 := rtn[i]
		v2 := rtn[j]
		var rtn bool
		if params.SortSpec.Column == "Email" {
			rtn = v1.Email < v2.Email
		} else if params.SortSpec.Column == "NumCredits" {
			rtn = v1.NumCredits < v2.NumCredits
		} else {
			rtn = true
		}
		if !params.SortSpec.Asc {
			rtn = !rtn
		}
		return rtn
	})
	return rtn, nil
}

func main() {
	rand.Seed(time.Now().Unix())
	config := &dash.Config{AutoKeygen: true, AnonAcc: true}
	dash.StartProcClient(config)
	defer dash.WaitForClear()

	model := MakeCreditModel()
	dash.RegisterPanelHandler("creditcount", "/", model.RootHandler)
	dash.RegisterDataHandlerEx("creditcount", "/get-users", model.GetUsers)
	dash.RegisterPanelHandlerEx("creditcount", "/update-credits", model.UpdateCredits)

	select {}

}
