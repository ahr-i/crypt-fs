package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ahr-i/crypt-fs/aes/setting"
	"github.com/ahr-i/crypt-fs/aes/src/log/logIPFS"
)

func EncryptFolder(folderPath string) error {
	logIPFS.Info("Attempting to load the AES key.")
	key, err := getKey(setting.Setting.KeyPath)
	if err != nil {
		return err
	}

	logIPFS.Info("Starting encryption.")
	if err = encrypt(folderPath, key); err != nil {
		return err
	}

	return nil
}

func encrypt(folderPath string, key []byte) error {
	return filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return encryptFile(path, key)
	})
}

func encryptFile(filePath string, key []byte) error {
	logIPFS.Info(fmt.Sprintf("%s ...", filePath))
	plaintext, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	encryptedFilePath := filePath + ".enc"

	return os.WriteFile(encryptedFilePath, ciphertext, 0644)
}
