# Crypt FS  
Encrypted files can be read using FUSE.

## 0. Encryption
```
cd aes_encryption_tool
go run main.go [target folder]
```

## 1. Start FUSE
```
cd crypt-fs
go run main.go [mount folder] [source folder]
```