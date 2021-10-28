package main

import (
	"log"

	"github.com/dgraph-io/badger"
	"github.com/sawka/dashborg-go-sdk/pkg/dash"
)

type Model struct {
	DB *badger.DB
}

func (m *Model) GetTables(withKeysCount bool) []badger.TableInfo {
	return m.DB.Tables(withKeysCount)
}

func (m *Model) GetSizes() interface{} {
	lsm, vlog := m.DB.Size()
	return map[string]int64{
		"lsm":  lsm,
		"vlog": vlog,
	}
}

func (m *Model) GetStringValue(key string) (string, error) {
	var strVal string
	err := m.DB.View(func(tx *badger.Txn) error {
		item, err := tx.Get([]byte(key))
		if err != nil {
			return err
		}
		bytesVal, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		strVal = string(bytesVal)
		return nil
	})
	return strVal, err
}

func (m *Model) SetStringValue(key string, val string) error {
	err := m.DB.Update(func(tx *badger.Txn) error {
		err := tx.Set([]byte(key), []byte(val))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (m *Model) RunValueLogGC(discardRatio float64) error {
	err := m.DB.RunValueLogGC(discardRatio)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	db, err := badger.Open(badger.DefaultOptions("./data/badgerdb"))
	if err != nil {
		log.Printf("Cannot open badgerdb: %v\n", err)
		return
	}
	m := &Model{DB: db}
	cfg := &dash.Config{AnonAcc: true, AutoKeygen: true}
	client, err := dash.ConnectClient(cfg)
	if err != nil {
		panic(err)
	}
	app := client.AppClient().NewApp("badgerdb")
	app.SetAppTitle("Sample BadgerDB Admin")
	app.WatchHtmlFile("panels/badgerdb.html", nil)
	app.Runtime().PureHandler("get-tables", m.GetTables)
	app.Runtime().PureHandler("get-sizes", m.GetSizes)
	app.Runtime().PureHandler("get-string-value", m.GetStringValue)
	app.Runtime().PureHandler("set-string-value", m.SetStringValue)
	app.Runtime().PureHandler("run-value-log-gc", m.RunValueLogGC)
	err = client.AppClient().WriteAndConnectApp(app)
	if err != nil {
		panic(err)
	}
	client.WaitForShutdown()
}
