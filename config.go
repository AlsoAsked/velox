package velox

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	ref           string = "ref"
	defaultBranch string = "master"
	gitlabBaseURL string = "https://gitlab.com"
)

type Config struct {
	Velox map[string][]string `mapstructure:"velox"`

	// Version
	Roadrunner map[string]string `mapstructure:"roadrunner"`

	// GitHub configuration
	GitHub *CodeHosting `mapstructure:"github"`

	// GitLab configuration
	GitLab *CodeHosting `mapstructure:"gitlab"`

	// Log contains log configuration
	Log map[string]string `mapstructure:"log"`
}

type Token struct {
	Token string `mapstructure:"token"`
}

type Endpoint struct {
	BaseURL string `mapstructure:"endpoint"`
}

type CodeHosting struct {
	BaseURL *Endpoint                `mapstructure:"endpoint"`
	Token   *Token                   `mapstructure:"token"`
	Plugins map[string]*PluginConfig `mapstructure:"plugins"`
}

type PluginConfig struct {
	Ref        string   `mapstructure:"ref"`
	Owner      string   `mapstructure:"owner"`
	Repo       string   `mapstructure:"repository"`
	Replace    string   `mapstructure:"replace"`
	BuildFlags []string `mapstructure:"build-flags"`
}

func (c *Config) Validate() error { //nolint:gocognit,gocyclo
	if _, ok := c.Roadrunner[ref]; !ok {
		c.Roadrunner[ref] = defaultBranch
	}

	if (c.GitLab != nil && len(c.GitLab.Plugins) == 0) && (c.GitHub != nil && len(c.GitHub.Plugins) == 0) {
		return errors.New("no plugins specified in the configuration")
	}

	if c.GitHub != nil {
		for k, v := range c.GitHub.Plugins {
			if v.Owner == "" {
				return fmt.Errorf("no owner specified for the plugin: %s", k)
			}

			if v.Ref == "" {
				return fmt.Errorf("no ref specified for the plugin: %s", k)
			}

			if v.Repo == "" {
				return fmt.Errorf("no repository specified for the plugin: %s", k)
			}
		}

		if c.GitHub.Token == nil || c.GitHub.Token.Token == "" {
			return errors.New("github.token should not be empty, create a token with any permissions: https://github.com/settings/tokens")
		}
	}

	if c.GitLab != nil {
		for k, v := range c.GitLab.Plugins {
			if v.Owner == "" {
				return fmt.Errorf("no owner specified for the plugin: %s", k)
			}

			if v.Ref == "" {
				return fmt.Errorf("no ref specified for the plugin: %s", k)
			}

			if v.Repo == "" {
				return fmt.Errorf("no repository specified for the plugin: %s", k)
			}
		}

		if c.GitLab.BaseURL == nil {
			c.GitLab.BaseURL = &Endpoint{BaseURL: gitlabBaseURL}
		}

		if c.GitLab.Token == nil || c.GitLab.Token.Token == "" {
			return errors.New("gitlab.token should not be empty, create a token with at least [api, read_api] permissions: https://gitlab.com/-/profile/personal_access_tokens")
		}
	}

	if len(c.Log) == 0 {
		c.Log = map[string]string{"level": "debug", "mode": "development"}
	}

	return nil
}
