package main

import (
	"flag"
	"fmt"

	"github.com/ssun3/bsos/pkg/driver"
)

func main() {
	var (
		endpoint = flag.String("endpoint", "defaultValue", "Endpoint out gRPC server would run at")
		token    = flag.String("token", "defaultValue", "Token of the storage provider")
		region   = flag.String("region", "ams3", "Region where the the volumes are going to be provisioned")
	)

	flag.Parse()

	fmt.Println(*endpoint, *token, *region)

	drv := driver.NewDriver(driver.InputParams{
		Name:     driver.DefaultName,
		Endpoint: *endpoint,
		Region:   *region,
		Token:    *token,
	})

	if err := drv.Run(); err != nil {
		fmt.Printf("Error %s, running the driver", err.Error())
	}

}
