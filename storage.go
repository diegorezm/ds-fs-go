package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type PathTransformFunc func(string) PathKey

type PathKey struct {
	Pathname string
	Filename string
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.Pathname, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}

type StoreOps struct {
	PathTransformFun PathTransformFunc
}

type Store struct {
	opts StoreOps
}

func NewStore(opts StoreOps) *Store {
	return &Store{
		opts: opts,
	}
}

var DefaultPathTransformFunc = func(key string) string {
	return key
}

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:]) // convert byte[20] to slice
	blockSize := 20                        // two dept
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err
}

// FIX: does not delete properly
func (s *Store) Delete(key string) error {
	pathKey := s.opts.PathTransformFun(key)
	defer func() {
		fmt.Printf("Deleted [%s] from system disk\n", pathKey.Filename)
	}()
	return os.RemoveAll(pathKey.FullPath())
}

func (s *Store) Has(key string) bool {
	pathKey := s.opts.PathTransformFun(key)
	_, err := os.Stat(pathKey.FullPath())
	if err != nil {
		return false
	}
	return true
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.opts.PathTransformFun(key)
	return os.Open(pathKey.FullPath())
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.opts.PathTransformFun(key)
	if err := os.MkdirAll(pathKey.Pathname, os.ModePerm); err != nil {
		return err
	}

	pathAndFileName := pathKey.FullPath()
	f, err := os.Create(pathAndFileName)

	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("Written (%d) bytes to disk: %s", n, pathAndFileName)
	return nil
}
