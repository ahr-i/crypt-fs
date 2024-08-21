package main

import (
	"os"

	"github.com/ahr-i/crypt-fs/aes/aes"
	"github.com/ahr-i/crypt-fs/aes/setting"
	"github.com/ahr-i/crypt-fs/aes/src/log/logIPFS"
)

func initialization() {
	logIPFS.Info("Initializing...")

	setting.Init()

	logIPFS.Info("Initialization successful.")
}

func checkArgs() {
	if len(os.Args) == 2 {
		return
	}

	logIPFS.Error("Correct Usage: 'go run main.go [target folder]'")
	os.Exit(1)
}

func main() {
	initialization()
	checkArgs()

	if err := aes.EncryptFolder(os.Args[1]); err != nil {
		logIPFS.Error(err)
	}
	logIPFS.Info("Success.")
}
