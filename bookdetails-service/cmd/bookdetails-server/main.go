package main

import (
	"bookinfo/bookdetails-service/svc/server"

	"bookinfo/bookdetails-service/global"
)

func main() {

	global.SetPid(global.ProjectRealPath + "/runtime/pid")

	server.Run()
}
