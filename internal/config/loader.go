package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
	"os"
	"strings"
)

var (
	k         = koanf.New(".")
	ymlParser = yaml.Parser()
)

type Loader[C any] struct {
	envVarsPrefix         string
	defaultConfigFilePath string
}

func NewLoader[C any](envVarsPrefix string, defaultConfigFilePath string) *Loader[C] {
	return &Loader[C]{envVarsPrefix: envVarsPrefix, defaultConfigFilePath: defaultConfigFilePath}
}

func NewDefaultLoader[C any]() *Loader[C] {
	return &Loader[C]{envVarsPrefix: "PROP_", defaultConfigFilePath: "config.yaml"}
}

func (l Loader[C]) LoadConfig() (*C, error) {
	confPath, err := l.loadConfPathFromFlag()
	if err != nil {
		return nil, err
	}

	if err := k.Load(file.Provider(confPath), ymlParser); err != nil {
		return nil, err
	}

	if err := k.Load(env.Provider(l.envVarsPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, l.envVarsPrefix)), "_", ".", -1)
	}), nil); err != nil {
		return nil, err
	}

	conf := new(C)
	err = k.Unmarshal("", conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (l Loader[C]) loadConfPathFromFlag() (string, error) {
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.StringP("config", "c", l.defaultConfigFilePath, "the file path for the config file")
	err := f.Parse(os.Args[1:])
	if err != nil {
		return "", err
	}

	return f.GetString("config")
}
