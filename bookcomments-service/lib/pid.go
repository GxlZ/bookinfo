package lib

import (
	"os"
	"io/ioutil"
	"strconv"
)

func SetPid(filePath string) int {
	pid := os.Getpid()

	ioutil.WriteFile(
		filePath,
		[]byte(strconv.Itoa(pid)),
		os.ModePerm,
	)

	return pid
}