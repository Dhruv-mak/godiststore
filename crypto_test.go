package main

import (
	"bytes"
	"testing"
)

func TestCopyEncryptDecrypt(t *testing.T) {
	payload := "some text data"
	src := bytes.NewReader([]byte(payload))
	dst := new(bytes.Buffer)
	key := newEncryptionKey()
	_, err := copyEncrypt(key, src, dst)
	if err != nil {
		t.Error(err)
	}

	out := new(bytes.Buffer)
	nw, err := copyDecrypt(key, dst, out)
	if err != nil {
		t.Error(err)
	}

	if nw != 16*len(payload) {
		t.Fail()
	}

	if out.String() != payload {
		t.Errorf("have %s, want %s", out.String(), payload)
	}

}
