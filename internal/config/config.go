package config

import (
	"context"
	"fmt"
	"github.com/chn555/property-service/pkg/db/mongo"
	"log"

	"github.com/go-playground/validator"
)

var (
	v = validator.New()
)

type MainConfig struct {
	MongoConfig           mongo.Config
	MongoEventStateConfig mongo.EventStateConfig
}

func LoadConfig(ctx context.Context) (*MainConfig, error) {
	log.Println("beginning loading configurations")
	c, err := NewDefaultLoader[MainConfig]().LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error: failed to load configurations: %w", err)
	}

	err = validateConfig(c)
	if err != nil {
		return nil, fmt.Errorf("error: failed to validate configurations: %w", err)
	}

	log.Println("successfully loaded configurations")
	return c, nil
}

func validateConfig(config any) error {
	return v.Struct(config)
}
