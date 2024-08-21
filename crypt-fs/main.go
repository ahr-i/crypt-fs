package main

import (
	"os"

	"github.com/ahr-i/crypt-fs/fuse/fuse"
	"github.com/ahr-i/crypt-fs/fuse/setting"
	"github.com/ahr-i/crypt-fs/fuse/src/log/logIPFS"
)

func initialization() {
	logIPFS.Info("Initializing...")

	setting.Init()

	logIPFS.Info("Initialization successful.")
}

func checkArgs() {
	if len(os.Args) == 3 {
		return
	}

	logIPFS.Error("Correct usage: 'go run main.go [mount folder] [source folder]'")
	os.Exit(1)
}

func main() {
	initialization()
	checkArgs()

	fuse.Execute(os.Args[1], os.Args[2])
}
