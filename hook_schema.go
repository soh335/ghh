package main

import (
	"fmt"
	"github.com/google/go-github/github"
)

type HookSchema struct {
	Name            *string    `json:"name,omitempty"`
	Events          []string   `json:"events,omitempty"`
	SupportedEvents []string   `json:"supported_events,omitempty"`
	Schema          [][]string `json:"schema,omitempty"`
}

func (h HookSchema) String() string {
	return github.Stringify(h)
}

func GetHookSchemas(c *github.Client) ([]HookSchema, *github.Response, error) {
	req, err := c.NewRequest("GET", "hooks", nil)
	if err != nil {
		return nil, nil, err
	}
	hookSchemas := new([]HookSchema)
	resp, err := c.Do(req, hookSchemas)
	if err != nil {
		return nil, resp, err
	}

	return *hookSchemas, resp, err
}

func GetHookSchema(c *github.Client, name string) (*HookSchema, *github.Response, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("hooks/%v", name), nil)
	if err != nil {
		return nil, nil, err
	}
	hookSchema := new(HookSchema)
	resp, err := c.Do(req, hookSchema)
	if err != nil {
		return nil, resp, err
	}

	return hookSchema, resp, err
}
