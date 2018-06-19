package main

import (
	// This Service
	"bookinfo/bookcomments-service/svc/server"
	"bookinfo/bookcomments-service/lib"
	"bookinfo/bookcomments-service/global"
)

func main() {
	// Update addresses if they have been overwritten by flags
	lib.SetPid(global.ProjectRealPath + "/runtime/pid")

	server.Run()
}
