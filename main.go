package main

import (
	"gitlab.com/bobymcbobs/url-redirector/pkg/common"
)

func main() {
	common.Logger().Printf("launching url-redirector (%v, %v, %v, %v)\n", common.GetAppBuildVersion(), common.GetAppBuildHash(), common.GetAppBuildDate(), common.GetAppBuildMode())
	common.CheckForConfigYAML()
	common.PrintEnvConfig()
	common.HandleWebserver()
}
