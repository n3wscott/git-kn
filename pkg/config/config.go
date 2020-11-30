package config

import (
	"github.com/n3wscott/git-kn/pkg/ghutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type C interface {
	GitHubClient() (*ghutil.GithubClient, error)
	Orgs() []string
	Fork() string
}

var once = sync.Once{}
var cfg C
var cfgErr error

func Config() C {
	once.Do(func() {
		cfg, cfgErr = Force()
	})
	if cfgErr != nil {
		panic(cfgErr)
	}
	return cfg
}

func Force() (C, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Printf("Couldn't detect home dir, using cwd: %s", err)
		homeDir = "."
	}
	// Walk the file tree from here backwards looking for a .bujo file.
	viper.SetDefault("orgs", "knative,knative-sandbox")
	viper.SetConfigName(".kngit") // .yaml is implicit
	viper.SetEnvPrefix("KNGIT")
	viper.AutomaticEnv()

	if override := os.Getenv("KNGIT_CONFIG_PATH"); override != "" {
		viper.AddConfigPath(override)
	}

	viper.AddConfigPath(homeDir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	return &gitConfig{
		orgs:      viper.GetString("orgs"),
		fork:      viper.GetString("fork"),
		tokenPath: viper.GetString("tokenPath"), // Optional. Or use GITHUB_TOKEN.
	}, nil
}

type gitConfig struct {
	orgs      string
	fork      string
	tokenPath string
}

func (c *gitConfig) Orgs() []string {
	orgs := strings.Split(c.orgs, ",")
	for i, org := range orgs {
		orgs[i] = strings.TrimSpace(org)
	}
	return orgs
}

func (c *gitConfig) Fork() string {
	return c.fork
}

func (c *gitConfig) GitHubClient() (*ghutil.GithubClient, error) {
	return ghutil.NewGithubClient(c.tokenPath)
}
