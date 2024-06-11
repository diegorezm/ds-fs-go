package main

import (
	"bytes"
	"io"
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
	opts := StoreOps{
		PathTransformFun: CASPathTransformFunc,
	}
	s := NewStore(opts)
	data := []byte("some jpg")
	key := "testing"
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}
	b, err := io.ReadAll(r)
	assert.Equal(t, b, data, "Expected parsed data to be equal to original data.")
	s.Delete(key)
	assert.False(t, s.Has(key), "Expected 'hasKey' to be false.")
}
