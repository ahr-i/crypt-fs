package logIPFS

import (
	"log"

	log2 "github.com/ipfs/go-log/v2"
)

func newLogger(system string) *log2.ZapEventLogger {
	logger := log2.Logger(system)

	levelSetting("info")

	return logger
}

func levelSetting(level string) {
	lvl, err := log2.LevelFromString(level)
	if err != nil {
		log.Println(err)
	}

	log2.SetAllLoggers(lvl)
}
