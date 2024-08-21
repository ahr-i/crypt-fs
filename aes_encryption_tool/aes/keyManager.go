package aes

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ahr-i/crypt-fs/aes/src/log/logIPFS"
)

func getKey(keyPath string) ([]byte, error) {
	var key []byte
	var err error

	logIPFS.Info("Checking the key.")
	if checkKey(keyPath) {
		logIPFS.Info("The existing AES key is present.")
		key, err = loadKey(keyPath)
		if err != nil {
			return nil, err
		}
	} else {
		logIPFS.Info("The existing AES key is not present.")
		logIPFS.Info("Generating a new AES key.")
		key, err = generateKey()
		if err != nil {
			return nil, err
		}

		logIPFS.Info(fmt.Sprintf("Saving the generated key. [%s]", keyPath))
		if err := saveKey(keyPath, key); err != nil {
			return nil, err
		}
	}

	logIPFS.Info("Successfully loaded the AES key.")
	return key, nil
}

func checkKey(keyPath string) bool {
	_, err := os.Stat(keyPath)

	return !os.IsNotExist(err)
}

func loadKey(keyPath string) ([]byte, error) {
	keyHex, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	return hex.DecodeString(string(keyHex))
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func saveKey(keyPath string, key []byte) error {
	keyHex := hex.EncodeToString(key)

	return os.WriteFile(keyPath, []byte(keyHex), 0600)
}
