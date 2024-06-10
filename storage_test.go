package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpicture"
	pathname := CASPathTransformFunc(key)
	expectedPath := "6804429f74181a63c50c/3d81d733a12f14a353ff"
	originalPath := "6804429f74181a63c50c3d81d733a12f14a353ff"
	slicedPath := strings.Split(expectedPath, "/")
	assert.Equal(t, pathname.Pathname, expectedPath)
	assert.Equal(t, pathname.Filename, originalPath)
	assert.Len(t, slicedPath, 2)
}

func TestStore(t *testing.T) {
	b, err := createDefaultFakeStore("test")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, string(b.content), string(b.originalData))
	fmt.Printf("content: %s\n", b)
}

func TestDelete(t *testing.T) {
	opts := StoreOps{
		PathTransformFun: CASPathTransformFunc,
	}
	s := NewStore(opts)
	b, err := createDefaultFakeStore()
	if err != nil {
		t.Error(err)
	}
	if err = s.Delete(b.key); err != nil {
		t.Error(err)
	}
	assert.False(t, s.Has(b.key))
}

type fakeStore struct {
	key          string
	content      []byte
	originalData []byte
}

func createFakeStore(key string, data []byte) (fakeStore, error) {
	opts := StoreOps{
		PathTransformFun: CASPathTransformFunc,
	}
	s := NewStore(opts)
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		return fakeStore{}, err
	}
	r, err := s.Read(key)
	if err != nil {
		return fakeStore{}, err
	}
	b, err := io.ReadAll(r)
	return fakeStore{key: key, content: b, originalData: data}, err
}

func createDefaultFakeStore(fakeKey ...string) (fakeStore, error) {
	var key string
	data := []byte("jpg")
	if len(fakeKey) > 0 {
		key = fakeKey[0]
	} else {
		key = randomBase64String(8)
	}
	b, err := createFakeStore(key, data)
	if err != nil {
		return fakeStore{}, err
	}
	return b, err
}

func randomBase64String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/float64(1.33333333333))))
	rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:l]
}
