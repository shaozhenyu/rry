/**************************************************************************

Copyright:YOYA

Author: shaozhenyu

Date:2017-06-22

Description:  AES test file

**************************************************************************/

package aes

import (
	"testing"
)

func Test_getKey(t *testing.T) {
	_, err := getKey("abc")
	if err == nil {
		t.Error("get key error")
	}

	value, err := getKey("abcdefghijklmnopqrstuvw")
	if err != nil {
		t.Error(err)
	}
	if string(value) != string("abcdefghijklmnop") {
		t.Error("get key the first 16 error")
	}

	value, err = getKey("aaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err != nil {
		t.Error(err)
	}
	if len(value) != 24 {
		t.Error("get key the first 24 error")
	}

	value, err = getKey("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err != nil {
		t.Error(err)
	}
	if len(value) != 32 {
		t.Error("get key the first 32 error")
	}
}

func Test_AES(t *testing.T) {
	text := "abcdefg"
	c, err := Encrypter([]byte(text))
	if err != nil {
		t.Error(err)
	}
	newT, err := Decrypter(c)
	if err != nil {
		t.Error(err)
	}

	if text != string(newT) {
		t.Error("aes encrypter error")
	}
}
