package main

import (
	"bytes"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "something"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "1af17e73721dbe0c40011b82ed4bb1a7dbe3ce29"
	expectedPathName := "1af17/e7372/1dbe0/c4001/1b82e/d4bb1/a7dbe/3ce29"
	if pathKey.Pathname != expectedPathName {
		t.Errorf("have %s, want %s", pathKey.Pathname, expectedPathName)
	}
	if pathKey.Filename != expectedOriginalKey {
		t.Errorf("have %s, want %s", pathKey.Filename, expectedOriginalKey)
	}
}

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "somefolder"
	data := []byte("some text data")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	key := "somefolder"
	data := []byte("some text data")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("expected key %s to be present", key)
	}
	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)

	if string(b) != string(data) {
		t.Errorf("have %s, want %s", b, data)
	}

}
