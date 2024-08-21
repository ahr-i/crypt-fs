package fuse

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
	"github.com/ahr-i/crypt-fs/fuse/setting"
	"github.com/ahr-i/crypt-fs/fuse/src/log/logIPFS"
)

func Execute(mountpoint string, source string) error {
	key, err := loadKey(setting.Setting.KeyPath)
	if err != nil {
		return err
	}

	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("decryptfs"),
		fuse.Subtype("decryptfs"),
		fuse.AllowOther(),
	)
	if err != nil {
		return err
	}
	defer c.Close()

	filesys := &FS{source, key}
	go func() {
		if err := fs.Serve(c, filesys); err != nil {
			log.Fatalf("Failed to serve FUSE filesystem: %v", err)
		}
	}()

	err = catchSignal(mountpoint)
	if err != nil {
		return err
	}

	return nil
}

func catchSignal(mountpoint string) error {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc

	if err := fuse.Unmount(mountpoint); err != nil {
		return err
	}

	return nil
}

// Root is the root directory of the file system.
func (f *FS) Root() (fs.Node, error) {
	return &Dir{f.root, f.key}, nil
}

// Dir implements both Node and Handle for directories.
type Dir struct {
	path string
	key  []byte
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	//a.Mode = os.ModeDir | 0555
	a.Mode = os.ModeDir | 0777

	return nil
}

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	p := filepath.Join(d.path, name)
	fi, err := os.Lstat(p)
	if err != nil {
		return nil, fuse.ENOENT
	}

	if fi.IsDir() {
		return &Dir{p, d.key}, nil
	}

	return &File{p, fi, d.key}, nil
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	f, err := os.Open(d.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fis, err := f.Readdir(0)
	if err != nil {
		return nil, err
	}

	var dirs []fuse.Dirent
	for _, fi := range fis {
		var de fuse.Dirent
		de.Name = fi.Name()

		if fi.IsDir() {
			de.Type = fuse.DT_Dir
		} else {
			de.Type = fuse.DT_File
		}
		dirs = append(dirs, de)
	}

	return dirs, nil
}

// File implements both Node and Handle for files.
type File struct {
	path string
	fi   os.FileInfo
	key  []byte
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = uint64(f.fi.Sys().(*syscall.Stat_t).Ino)
	a.Mode = f.fi.Mode()
	a.Size = uint64(f.fi.Size())

	return nil
}

func (f *File) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	start := time.Now()

	logIPFS.Info(fmt.Sprintf("Decrypting file:", f.path))
	file, err := os.Open(f.path)

	if err != nil {
		return err
	}
	defer file.Close()

	ciphertext, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(f.key)
	if err != nil {
		return err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := aead.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	fuseutil.HandleRead(req, resp, plaintext)

	duration := time.Since(start)
	logIPFS.Info(fmt.Sprintf("Decryption and read of file %s took %s\n", f.path, duration))

	return nil
}
