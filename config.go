package main

import (
	"MuseBot/utils"
	"crypto/rsa"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Parsed config with embedded fields
type configParsed struct {
	config
	TravisPublicKey *rsa.PublicKey
	LogLevel        logrus.Level
}

// Raw config config
type config struct {
	// Telegram config
	TgPollInterval time.Duration `arg:"env:TG_POLL_INTERVAL"`
	TgToken        string        `arg:"env:TG_TOKEN"`
	// GitHub Config
	GitHubPullUrl   string `arg:"env:GITHUB_PULL_URL"`
	GitHubCommitUrl string `arg:"env:GITHUB_COMMIT_URL"`
	GitHubHookPath  string `arg:"env:GITHUB_HOOK_PATH"`
	GitHubHookToken string `arg:"env:GITHUB_HOOK_TOKEN"`
	// MU Node config
	MUNodeUrl string `arg:"env:MU_NODE_URL"`
	// Travis config
	TravisHookPath  string `arg:"env:TRAVIS_HOOK_PATH"`
	TravisConfigUrl string `arg:"env:TRAVIS_CONFIG_URL"`
	TravisPublicKey string `arg:"env:TRAVIS_PUBLIC_KEY" help:"Custom travis public key. By default it will downloaded from TRAVIS_CONFIG_URL"`
	// Main config
	STPath    string `arg:"env:ST_PATH"`
	WebListen string `arg:"env:WEB_LISTEN"`
	LogLevel  string `arg:"env:LOG_LEVEL"`
}

// VersionId return program version string on human format
func (config) Version() string {
	return fmt.Sprintf("VersionId: %v, commit: %v, built at: %v", version, commit, date)
}

// Description return program description string
func (config) Description() string {
	return "Really fast sync tool for S3"
}

// GetConfig parse config, set default values, check input values and return configParsed struct
func GetConfig() (cfg configParsed) {
	rawConfig := config{}
	rawConfig.TgPollInterval = time.Millisecond * 500
	rawConfig.STPath = "/storage/storage.db"
	rawConfig.WebListen = ":8080"
	rawConfig.GitHubPullUrl = "https://github.com/musescore/MuseScore/pull/"
	rawConfig.GitHubCommitUrl = "https://github.com/musescore/MuseScore/commit/"
	rawConfig.GitHubHookPath = "/webhook/github/"
	rawConfig.MUNodeUrl = "https://musescore.org/node/"
	rawConfig.TravisConfigUrl = "https://api.travis-ci.org/config"
	rawConfig.TravisHookPath = "/webhook/travis/"
	rawConfig.LogLevel = "INFO"

	arg.MustParse(&rawConfig)
	cfg.config = rawConfig

	{
		if cfg.config.TravisPublicKey == "" {
			key, err := utils.TravisGetPubKey(cfg.TravisConfigUrl)
			if err != nil {
				log.Fatalf("Failed to get Travis public key, error: %s", err)
			}
			cfg.config.TravisPublicKey = key
		}

		key, err := utils.TravisParsePubKey(cfg.config.TravisPublicKey)
		if err != nil {
			log.Fatalf("Failed to parse Travis pub key, error: %s", err)
		}
		cfg.TravisPublicKey = key
	}

	{
		switch strings.ToUpper(cfg.config.LogLevel) {
		case "DEBUG":
			cfg.LogLevel = logrus.DebugLevel
		case "INFO":
			cfg.LogLevel = logrus.InfoLevel
		case "WARN":
			cfg.LogLevel = logrus.WarnLevel
		case "ERROR":
			cfg.LogLevel = logrus.ErrorLevel
		case "FATAL":
			cfg.LogLevel = logrus.FatalLevel
		default:
			log.Fatalf("Unmnown LOG_LEVEL: %s", cfg.config.LogLevel)
		}
	}

	return cfg
}
