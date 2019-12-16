package main

import (
	"gitlab.com/bobymcbobs/url-redirector/src/common"
)

func main() {
	common.Logger().Println("Warming up")
	common.CheckForConfigYAML()
	common.PrintEnvConfig()
	common.HandleWebserver()
}
