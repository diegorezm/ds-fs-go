package main

import (
	"bytes"
	"fmt"
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
	s := newStore()
	defer teardown(t, s)
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("test_%d", i)
		data := []byte("some jpg")
		err := s.Write(key, bytes.NewReader(data))
		assert.NoError(t, err, "Error while writing file with key: "+key)
		r, err := s.Read(key)
		assert.NoError(t, err, "Error while reading file with key: "+key)
		b, err := io.ReadAll(r)
		assert.Equal(t, b, data, "Expected parsed data to be equal to original data.")
		err = s.Delete(key)
		assert.NoError(t, err, "Error while deleting file with key: "+key)
		assert.False(t, s.Has(key), "File still exists.")
	}
}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFun: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
