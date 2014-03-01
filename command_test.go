package main

import (
	"github.com/gonuts/commander"
	"testing"
)

func TestvalidateSubUsage(t *testing.T) {
	cmd := &commander.Command{
		Run:       runListCmd,
		UsageLine: "cmd <foo> <bar>",
	}

	var err error
	err = validateSubUsage(cmd, []string{})
	if err == nil {
		t.Error("empty arguments. so should raise error")
	}

	err = validateSubUsage(cmd, []string{"hoge"})
	if err == nil {
		t.Error("not enough arugments. so should raise error")
	}

	err = validateSubUsage(cmd, []string{"hoge", "fuga"})
	if err != nil {
		t.Error("arguments count is correct. so should not raise error")
	}

	err = validateSubUsage(cmd, []string{"hoge", "fuga", "yap"})
	if err == nil {
		t.Error("arguments count is over. so should raise error")
	}
}
