package utils

import (
	"fmt"
	"github.com/ip2location/ip2location-go/v9"
)

func GetInfo(ip string) (error, string, string, string) {
	db, err := ip2location.OpenDB("config/IP2LOCATION-LITE-DB3.BIN")

	if err != nil {
		fmt.Print(err)
		return err, "", "", ""
	}
	//ip := "31.42.47.99"
	results, err := db.Get_all(ip)

	if err != nil {
		fmt.Print(err)
		return err, "", "", ""
	}

	db.Close()

	return err, results.Country_short, results.Region, results.City
}
