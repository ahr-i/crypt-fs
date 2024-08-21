package setting

import (
	"encoding/json"
	"os"

	"github.com/ahr-i/crypt-fs/aes/src/log/logIPFS"
)

const settingFilePath string = "./setting/setting.json"

func Init() {
	err := readSettingFile()
	if err != nil {
		logIPFS.Error(err)

		os.Exit(1)
	}
	logIPFS.Info("Successfully finished initializing setting.")
}

func readSettingFile() error {
	file, err := os.ReadFile(settingFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &Setting)
	if err != nil {
		return err
	}

	return nil
}
