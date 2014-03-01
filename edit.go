package main

import (
	"errors"
	"github.com/google/go-github/github"
	"io/ioutil"
	"os"
	"os/exec"
)

func edit(hook *github.Hook, hookSchema *HookSchema) (*github.Hook, error) {
	file, err := ioutil.TempFile("", os.Args[0])
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())

	if err := writeTomlHook(hook, hookSchema, file); err != nil {
		return nil, err
	}
	oldStat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if err := execEditCmd(file.Name()); err != nil {
		return nil, err
	}

	newFile, err := os.Open(file.Name())
	if err != nil {
		return nil, err
	}
	newStat, err := newFile.Stat()
	if oldStat.ModTime().Unix() == newStat.ModTime().Unix() {
		return nil, errors.New("not update")
	}
	newHook, err := createHookFromToml(newFile)
	if err != nil {
		return nil, err
	}
	return newHook, nil
}

func execEditCmd(filePath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return errors.New("cant detect editor cmd")
	}
	editorPath, err := exec.LookPath(editor)
	if err != nil {
		return err
	}
	// http://stackoverflow.com/questions/12088138/trying-to-launch-an-external-editor-from-within-a-go-program/12089980#12089980
	cmd := exec.Command(editorPath, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
