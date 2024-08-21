package fuse

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ahr-i/crypt-fs/fuse/src/log/logIPFS"
)

func loadKey(keyPath string) ([]byte, error) {
	keyHex, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	key, err := hex.DecodeString(string(keyHex))
	if err != nil {
		return key, nil
	}

	logIPFS.Info(fmt.Sprintf("Successfully loaded the AES key: %s", hex.EncodeToString(key)))
	return key, nil
}
