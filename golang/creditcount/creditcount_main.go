package main

import (
	"fmt"
	"math/rand"
	"reflect"
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

type AppState struct {
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

func (m *CreditModel) UpdateCredits(req *dash.AppRequest, userId int, numCredits int) error {
	if userId <= 0 || numCredits < 0 || numCredits > 100 {
		return fmt.Errorf("Invalid Params")
	}
	var foundUser *UserType
	for i := 0; i < len(m.Users); i++ {
		user := m.Users[i]
		if user.UserId == userId {
			user.NumCredits = numCredits
			foundUser = user
			break
		}
	}
	if foundUser == nil {
		return fmt.Errorf("User not found")
	}
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

func (m *CreditModel) GetUsers(req dash.Request, state *AppState, params getUsersParams) (interface{}, error) {
	var rtn []*UserType
	if params.SearchEmail == "" {
		rtn = make([]*UserType, len(m.Users))
		copy(rtn, m.Users)
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

func (m *CreditModel) TestHandler(req *dash.AppRequest) error {
	fmt.Printf("TEST HANDLER\n")
	req.SetData("@rtn", map[string]interface{}{"x": 5})
	return nil
}

func main() {
	rand.Seed(time.Now().Unix())

	config := &dash.Config{AutoKeygen: true, AnonAcc: true}
	client, err := dash.ConnectClient(config)
	if err != nil {
		fmt.Printf("Error connecting client: %v\n", err)
		return
	}
	model := MakeCreditModel()
	app := client.AppClient().NewApp("creditcount")
	app.SetOfflineAccess(true)
	app.SetAllowedRoles("admin", "public")
	app.WatchHtmlFile("panels/creditcount.html", nil)
	app.Runtime().SetAppStateType(reflect.TypeOf(&AppState{}))
	app.Runtime().PureHandler("get-users", model.GetUsers)
	app.Runtime().Handler("update-credits", model.UpdateCredits)
	app.Runtime().Handler("test", model.TestHandler)
	err = client.AppClient().WriteAndConnectApp(app)
	if err != nil {
		fmt.Printf("Error connecting app: %v\n", err)
		return
	}
	client.WaitForShutdown()
}
