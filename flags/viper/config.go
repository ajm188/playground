package main

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrBadDSN       = errors.New("bad DSN format")
	ErrMissingValue = errors.New("missing required value")

	discoveryFlagRegexp = regexp.MustCompile(`^discovery-(?P<impl>\w+)-(?P<flag>.+)$`)
)

type FileConfig struct {
	Defaults Config
	Clusters map[string]*Config
}

type Config struct {
	Name                 string
	Code                 string
	DiscoveryImpl        string
	DiscoveryFlagsByImpl map[string]map[string]string
}

func (c Config) Merge(override *Config) Config {
	if override == nil {
		return c
	}

	merged := Config{
		Name:                 c.Name,
		Code:                 c.Code,
		DiscoveryImpl:        c.DiscoveryImpl,
		DiscoveryFlagsByImpl: map[string]map[string]string{},
	}

	if override.Name != "" {
		merged.Name = override.Name
	}

	if override.Code != "" {
		merged.Code = override.Code
	}

	if override.DiscoveryImpl != "" {
		merged.DiscoveryImpl = override.DiscoveryImpl
	}

	mergeFlagsByImpl(&merged.DiscoveryFlagsByImpl, c.DiscoveryFlagsByImpl)
	mergeFlagsByImpl(&merged.DiscoveryFlagsByImpl, override.DiscoveryFlagsByImpl)

	return merged
}

func mergeFlagsByImpl(base *map[string]map[string]string, override map[string]map[string]string) {
	if base == nil || (*base) == nil {
		*base = map[string]map[string]string{}
	}

	for impl, flags := range override {
		_, ok := (*base)[impl]
		if !ok {
			(*base)[impl] = map[string]string{}
		}

		for k, v := range flags {
			(*base)[impl][k] = v
		}
	}
}

func (c *Config) Set(value string) error {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.StringVar(&c.Name, "name", "", "")
	fs.StringVar(&c.Code, "code", "", "")
	fs.StringVar(&c.DiscoveryImpl, "discovery", "", "")

	args := strings.Split(value, ",")
	strictArgs := make([]string, 0, len(args))

	for _, arg := range args {
		if !strings.HasPrefix(arg, "discovery-") {
			strictArgs = append(strictArgs, "-"+arg)
			continue
		}

		if !strings.Contains(arg, "=") {
			return fmt.Errorf("%w: %s is missing value", ErrBadDSN, arg)
		}

		parts := strings.Split(arg, "=")
		name := parts[0]
		val := strings.Join(parts[1:], "=")

		match := discoveryFlagRegexp.FindStringSubmatch(name)
		if match == nil {
			// not a discovery flag
			continue
		}

		var impl, flag string

		for i, g := range discoveryFlagRegexp.SubexpNames() {
			switch g {
			case "impl":
				impl = match[i]
			case "flag":
				flag = match[i]
			}
		}

		if c.DiscoveryFlagsByImpl[impl] == nil {
			c.DiscoveryFlagsByImpl[impl] = map[string]string{}
		}

		c.DiscoveryFlagsByImpl[impl][flag] = val
	}

	return fs.Parse(strictArgs)
}

func (c *Config) String() string {
	return ""
}

func (c *Config) Type() string {
	return "Config"
}

type ConfigMap map[string]*Config

func (cf ConfigMap) Set(value string) error {
	cfg := &Config{}

	if err := cfg.Set(value); err != nil {
		return err
	}

	if cfg.Name == "" {
		return fmt.Errorf("%w: per-cluster values must include `name` attribute", ErrMissingValue)
	}

	finalCfg := cfg
	base, ok := cf[cfg.Name]
	if ok {
		merged := base.Merge(cfg)
		finalCfg = &merged
	}

	cf[cfg.Name] = finalCfg

	return nil
}

func (cf ConfigMap) String() string {
	return ""
}

func (cf ConfigMap) Type() string {
	return "ClustersFlag"
}
