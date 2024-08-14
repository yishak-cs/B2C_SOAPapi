package main

import (
	"github.com/yishak-cs/B2C_SOAPapi/client"
	"github.com/yishak-cs/B2C_SOAPapi/service"
)

func main() {
	service.StartService()
	client.StartClient()

}
