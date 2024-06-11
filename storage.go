package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const defaultRootFolderName = "dsfs"

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

func (p PathKey) FullPathWithRoot(root string) string {
	return fmt.Sprintf("%s/%s", root, p.FullPath())
}

type PathTransformFunc func(string) PathKey

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

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

type StoreOpts struct {
	// root folder
	Root             string
	PathTransformFun PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFun == nil {
		opts.PathTransformFun = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}

	return &Store{
		StoreOpts: opts,
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

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFun(key)
	fullPath := fmt.Sprintf("%s/%s", s.Root, pathKey.FirstPathName())
	return os.RemoveAll(fullPath)
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFun(key)
	_, err := os.Stat(pathKey.FullPathWithRoot(s.Root))
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFun(key)
	return os.Open(pathKey.FullPathWithRoot(s.Root))
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFun(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.Pathname)

	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(pathKey.FullPathWithRoot(s.Root))

	if err != nil {
		return err
	}

	defer f.Close()

	// tests only work if i copy
	// the io reader buffer
	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}
	return err
}
