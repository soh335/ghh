package main

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
	"io"
	"strings"
)

type TomlHook struct {
	Name   *string                `toml:"name"`
	Events []string               `toml:"events"`
	Active *bool                  `toml:"active"`
	Config map[string]interface{} `toml:"config"`
}

func createHookFromToml(reader io.Reader) (*github.Hook, error) {
	tomlHook := new(TomlHook)
	if _, err := toml.DecodeReader(reader, tomlHook); err != nil {
		return nil, err
	}

	hook := new(github.Hook)
	if tomlHook.Name != nil {
		hook.Name = tomlHook.Name
	}
	if tomlHook.Active != nil {
		hook.Active = tomlHook.Active
	}
	if tomlHook.Events != nil {
		hook.Events = tomlHook.Events
	}
	if tomlHook.Config != nil {
		hook.Config = tomlHook.Config
	}

	return hook, nil
}

func writeTomlHook(hook *github.Hook, hookSchema *HookSchema, writer io.Writer) error {
	// https://github.com/BurntSushi/toml/blob/master/encode.go#L136
	replacer := strings.NewReplacer(
		"\t", "\\t",
		"\n", "\\n",
		"\r", "\\r",
		"\"", "\\\"",
		"\\", "\\\\",
	)

	strVal := func(str string) string {
		return "\"" + replacer.Replace(str) + "\""
	}

	writeVal := func(key string, val interface{}) {
		switch val.(type) {
		case string:
			val = strVal(val.(string))
		}
		fmt.Fprintf(
			writer,
			"%v = %v\n",
			key,
			val,
		)
	}

	if hook.ID != nil {
		writeVal("id", *hook.ID)
	}
	if hook.Name != nil {
		writeVal("name", *hook.Name)
	}
	if hook.Active != nil {
		writeVal("active", *hook.Active)
	}
	if hook.CreatedAt != nil {
		writeVal("created_at", hook.CreatedAt.String())
	}
	if hook.UpdatedAt != nil {
		writeVal("updated_at", hook.UpdatedAt.String())
	}

	// events
	fmt.Fprintf(writer, "%v = [\n", "events")
	writtenEvents := make(map[string]bool)
	for _, event := range hook.Events {
		writtenEvents[event] = true
		fmt.Fprintf(writer, "  %v,\n", strVal(event))
	}
	for _, event := range hookSchema.SupportedEvents {
		if writtenEvents[event] {
			continue
		}
		fmt.Fprintf(writer, "# %v,\n", strVal(event))
	}
	fmt.Fprintf(writer, "]\n")

	// config
	fmt.Fprintf(writer, "[config]\n")
	writtenConfig := make(map[string]bool)
	for key, val := range hook.Config {
		writtenConfig[key] = true
		writeVal(key, val)
	}
	for _, schema := range hookSchema.Schema {
		if writtenConfig[schema[1]] {
			continue
		}
		// schema type
		// https://github.com/github/github-services/blob/32456c8c3815e08bd26bfb4e68066d2200f2c21b/lib/service.rb#L158
		switch schema[0] {
		case "string":
			writeVal("# "+schema[1], "")
		case "password":
			writeVal("# "+schema[1], "*****")
		case "boolean":
			// can i detect default value ?
			writeVal("# "+schema[1], false)
		default:
			return errors.New("invalid config type: " + schema[0])
		}
	}
	return nil
}
