package main

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	file, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())
	io.WriteString(file, `{"hoge":"fuga"}`)

	m, err := ReadConfig(file.Name())
	if err != nil {
		t.Error(err)
	}
	var keys []string
	for key, _ := range m {
		keys = append(keys, key)
	}
	if keys[0] != "hoge" {
		t.Error("should key count is 1")
	}
}

func TestWriteConfig(t *testing.T) {
	file, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())
	m := make(map[string]string)
	m["yap"] = "dap"
	WriteConfig(file.Name(), m)
	byt, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Error(err)
	}
	expect := `{"yap":"dap"}
`
	if string(byt) != expect {
		t.Error("should", expect, "but", string(byt))
	}
}
