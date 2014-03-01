package main

import (
	"code.google.com/p/goauth2/oauth"
	"errors"
	"fmt"
	"github.com/gonuts/commander"
	"github.com/google/go-github/github"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var supportCmd = &commander.Command{
	Run:       runSupportCmd,
	UsageLine: "support",
}

var listCmd = &commander.Command{
	Run:       runListCmd,
	UsageLine: "list <owner> <repo>",
}

var showCmd = &commander.Command{
	Run:       runShowCmd,
	UsageLine: "show <owner> <repo> <id>",
}

var createCmd = &commander.Command{
	Run:       runCreateCmd,
	UsageLine: "create <owner> <repo> <type>",
}

var editCmd = &commander.Command{
	Run:       runEditCmd,
	UsageLine: "edit <owner> <repo> <id>",
}

var deleteCmd = &commander.Command{
	Run:       runDeleteCmd,
	UsageLine: "delete <owner> <repo> <id>",
}

var testCmd = &commander.Command{
	Run:       runTestCmd,
	UsageLine: "test <owner> <repo> <id>",
}

var configCmd = &commander.Command{
	Run:       runConfigCmd,
	UsageLine: "config",
}

func Client() (*github.Client, error) {
	config, err := ReadConfig(ConfigPath)
	if err != nil {
		return nil, err
	}
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: config["token"]},
	}
	client := github.NewClient(t.Client())
	return client, nil
}

func runSupportCmd(cmd *commander.Command, args []string) error {
	c, err := Client()
	if err != nil {
		return err
	}
	hooks, _, err := GetHookSchemas(c)
	if err != nil {
		return err
	}
	for _, hook := range hooks {
		fmt.Println(*hook.Name)
	}
	return nil
}

func runListCmd(cmd *commander.Command, args []string) error {
	if err := validateSubUsage(cmd, args); err != nil {
		return err
	}
	//TODO pager
	c, err := Client()
	if err != nil {
		return err
	}
	hooks, _, err := c.Repositories.ListHooks(args[0], args[1], &github.ListOptions{})

	if err != nil {
		return err
	}

	var dummyHookSchema HookSchema
	for _, hook := range hooks {
		fmt.Println("---------------------------------")
		writeTomlHook(&hook, &dummyHookSchema, os.Stdout)
	}

	return nil
}

func runShowCmd(cmd *commander.Command, args []string) error {
	if err := validateSubUsage(cmd, args); err != nil {
		return err
	}
	id, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	c, err := Client()
	if err != nil {
		return err
	}
	hook, _, err := c.Repositories.GetHook(args[0], args[1], id)

	if err != nil {
		return err
	}

	var dummyHookSchema HookSchema
	writeTomlHook(hook, &dummyHookSchema, os.Stdout)

	return nil
}

func runCreateCmd(cmd *commander.Command, args []string) error {
	if err := validateSubUsage(cmd, args); err != nil {
		return err
	}
	c, err := Client()
	if err != nil {
		return err
	}
	hookSchema, _, err := GetHookSchema(c, args[2])
	if err != nil {
		return err
	}

	hook := new(github.Hook)
	hook.Name = hookSchema.Name

	active := true
	hook.Active = &active

	events := make([]string, 0)
	for _, event := range hookSchema.Events {
		events = append(events, event)
	}
	hook.Events = events

	hook.Config = make(map[string]interface{})

	newHook, err := edit(
		hook,
		hookSchema,
	)

	if err != nil {
		return err
	}

	if _, _, err := c.Repositories.CreateHook(args[0], args[1], newHook); err != nil {
		return err
	}

	return nil
}

func runEditCmd(cmd *commander.Command, args []string) error {
	if err := validateSubUsage(cmd, args); err != nil {
		return err
	}
	c, err := Client()
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	hook, _, err := c.Repositories.GetHook(args[0], args[1], id)
	if err != nil {
		return err
	}
	hookSchema, _, err := GetHookSchema(c, *hook.Name)
	if err != nil {
		return err
	}

	newHook, err := edit(
		&github.Hook{
			Name:   hook.Name,
			Events: hook.Events,
			Active: hook.Active,
			Config: hook.Config,
		},
		hookSchema,
	)
	if err != nil {
		return err
	}

	if _, _, err := c.Repositories.EditHook(args[0], args[1], *hook.ID, newHook); err != nil {
		return err
	}

	return nil
}

func runDeleteCmd(cmd *commander.Command, args []string) error {
	if err := validateSubUsage(cmd, args); err != nil {
		return err
	}
	id, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}

	c, err := Client()
	if err != nil {
		return err
	}
	_, err = c.Repositories.DeleteHook(args[0], args[1], id)
	if err != nil {
		return err
	}
	return nil
}

func runTestCmd(cmd *commander.Command, args []string) error {
	if err := validateSubUsage(cmd, args); err != nil {
		return err
	}
	id, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	c, err := Client()
	if err != nil {
		return err
	}
	_, err = c.Repositories.TestHook(args[0], args[1], id)
	if err != nil {
		return err
	}
	return nil
}

func runConfigCmd(cmd *commander.Command, args []string) error {
	switch len(args) {
	case 1:
		config, err := ReadConfig(ConfigPath)
		if err != nil {
			return err
		}
		val := config[args[0]]
		fmt.Println(val)
	case 2:
		config, err := ReadConfig(ConfigPath)
		if config == nil {
			config = make(map[string]string)
		}
		config[args[0]] = args[1]
		err = WriteConfig(ConfigPath, config)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid argument count")
	}
	return nil
}

func validateSubUsage(cmd *commander.Command, args []string) error {
	actions := strings.Split(cmd.UsageLine, " ")
	expectVal := 0
	r := regexp.MustCompile("<.*?>")
	for _, action := range actions {
		if r.MatchString(action) {
			expectVal++
		}
	}
	if len(args) != expectVal {
		return errors.New("invalid argument count\nUsage: " + cmd.UsageLine)
	}
	return nil
}
