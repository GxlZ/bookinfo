package main

import (
	"bookinfo/bookdetails-service/svc/server"

	"bookinfo/bookdetails-service/global"
	"runtime"
	"time"
	"fmt"
	_ "github.com/mkevac/debugcharts"
)

func main() {

	global.SetPid(global.ProjectRealPath + "/runtime/pid")

	go func(){
		for  {
			printMem()
			time.Sleep(time.Second)
		}
	}()

	server.Run()
}


func printMem() {
	mem := runtime.MemStats{}

	runtime.ReadMemStats(&mem)
	fmt.Println("==============================")
	fmt.Println("mem.Sys: ", mem.Sys/1024/1024)
	fmt.Println("mem.Alloc: ", mem.Alloc/1024/1024)
	fmt.Println("mem.HeapIdle: ", mem.HeapIdle/1024/1024)
	fmt.Println("mem.TotalAlloc: ", mem.TotalAlloc/1024/1024)
	fmt.Println("mem.HeapAlloc: ", mem.HeapAlloc/1024/1024)
	fmt.Println("mem.HeapSys: ", mem.HeapSys/1024/1024)
	fmt.Println("mem.HeapObjects: ", mem.HeapObjects)
}