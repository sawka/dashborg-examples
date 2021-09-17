package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sawka/dashborg-go-sdk/pkg/dash"
)

type DataPoint struct {
	X     int
	Y     int
	Val   int
	Color string
}

const NUM_DATA = 30

var Colors []string = []string{"red", "green", "blue", "purple"}

func RegenData(req *dash.AppRequest) error {
	rtn := make([]DataPoint, 0)
	for i := 0; i < NUM_DATA; i++ {
		point := DataPoint{
			X:     rand.Intn(50),
			Y:     rand.Intn(50),
			Val:   rand.Intn(50) + 1,
			Color: Colors[rand.Intn(len(Colors))],
		}
		rtn = append(rtn, point)
	}
	req.SetData("$.data", rtn)
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

	app := client.AppClient().NewApp("d3-test")
	app.SetAppTitle("D3 Demo")
	app.WatchHtmlFile("panels/d3-test.html", nil)
	app.SetInitRequired(true)
	app.Runtime().Handler("regen-data", RegenData)
	app.Runtime().SetInitHandler(RegenData)

	err = client.AppClient().WriteAndConnectApp(app)
	if err != nil {
		fmt.Printf("Error connecting app: %v\n", err)
		return
	}
	client.WaitForShutdown()
}
