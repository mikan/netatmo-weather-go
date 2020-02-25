package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mikan/netatmo-weather-go"
)

func main() {
	clientID := flag.String("c", "", "netatmo client id")
	clientSecret := flag.String("s", "", "netatmo client secret")
	username := flag.String("u", "", "netatmo user name")
	password := flag.String("p", "", "netatmo password")
	deviceID := flag.String("d", "", "device id (MAC address)")
	moduleID := flag.String("m", "", "module id (MAC address)")
	minutes := flag.Int("a", -1, "how many minutes ago")
	flag.Parse()
	if *clientID == "" || *clientSecret == "" || *username == "" || *password == "" {
		flag.Usage()
		os.Exit(2)
	}
	client, err := netatmo.NewClient(context.Background(), *clientID, *clientSecret, *username, *password)
	if err != nil {
		panic(err)
	}
	if len(*deviceID) == 0 {
		stations(client)
		return
	}
	if len(*moduleID) == 0 {
		moduleID = deviceID
	}
	if *minutes > 0 {
		measureRange(client, *deviceID, *moduleID, *minutes)
	} else {
		measureNewest(client, *deviceID, *moduleID)
	}
}

func stations(client *netatmo.Client) {
	devices, user, err := client.GetStationsData()
	if err != nil {
		panic(err)
	}
	if err := printStationsData(devices, *user, os.Stdout); err != nil {
		panic(err)
	}
}

func measureRange(client *netatmo.Client, device, module string, minutes int) {
	end := time.Now().UTC()
	begin := end.Add(-time.Duration(minutes) * time.Minute)
	values, err := client.GetMeasureByTimeRange(device, module, begin.Unix(), end.Unix())
	if err != nil {
		panic(err)
	}
	if err := printMeasures(values, os.Stdout); err != nil {
		panic(err)
	}
}

func measureNewest(client *netatmo.Client, device, module string) {
	value, err := client.GetMeasureByNewest(device, module)
	if err != nil {
		panic(err)
	}
	if value != nil {
		if err := printMeasures([]netatmo.Measure{*value}, os.Stdout); err != nil {
			panic(err)
		}
	} else {
		fmt.Println("No Data")
	}
}
