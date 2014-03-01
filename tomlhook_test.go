package main

import (
	"bytes"
	"github.com/google/go-github/github"
	"io/ioutil"
	"testing"
)

func TestWriteTomlHook(t *testing.T) {
	id := 12345
	hook := &github.Hook{
		ID:     &id,
		Events: []string{"yap"},
	}
	hookSchemaName := "hoge"
	hookSchema := &HookSchema{
		Name:            &hookSchemaName,
		Events:          []string{"foo", "bar"},
		SupportedEvents: []string{"foo", "bar", "yap", "dap"},
		Schema: [][]string{
			[]string{"boolean", "bool"},
			[]string{"password", "pass"},
			[]string{"string", "str"},
		},
	}

	expect := `id = 12345
events = [
  "yap",
# "foo",
# "bar",
# "dap",
]
[config]
# bool = false
# pass = "*****"
# str = ""
`
	compare(t, hook, hookSchema, expect)
}

func compare(t *testing.T, hook *github.Hook, hookSchema *HookSchema, expect string) {
	b := bytes.NewBuffer([]byte{})
	err := writeTomlHook(hook, hookSchema, b)
	if err != nil {
		t.Error(err)
	}

	byt, err := ioutil.ReadAll(b)
	if err != nil {
		t.Error(err)
	}

	if string(byt) != expect {
		t.Error("expect:\n", expect, "but got\n", string(byt))
	}
}
