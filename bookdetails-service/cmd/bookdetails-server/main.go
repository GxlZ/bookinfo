package main

import (
	"bookinfo/bookdetails-service/svc/server"

	"bookinfo/bookdetails-service/global"
	_ "github.com/mkevac/debugcharts"
)

func main() {

	global.SetPid(global.ProjectRealPath + "/runtime/pid")

	server.Run()
}