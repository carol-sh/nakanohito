// Package config reads config.json.
package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Config is everything nakanohito needs to know, read from one file.
type Config struct {
	MediaMTX MediaMTX `json:"mediamtx"`
	Ingest   Ingest   `json:"ingest"`
	OBS      OBS      `json:"obs"`
}

// MediaMTX is the RTMP server binary and the addresses it listens on.
type MediaMTX struct {
	Bin  string `json:"bin"`
	RTMP string `json:"rtmp"`
	API  string `json:"api"`
}

// Ingest is where a remote encoder publishes the feed.
type Ingest struct {
	Path     string `json:"path"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// OBS is the OBS Studio app that puts the feed on air.
type OBS struct {
	App string `json:"app"`
}

// Load reads the config at path and fails on unknown or empty fields.
func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	var c Config
	if err := d.Decode(&c); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if err := c.check(); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return &c, nil
}

func (c *Config) check() error {
	required := []struct{ key, value string }{
		{"mediamtx.bin", c.MediaMTX.Bin},
		{"mediamtx.rtmp", c.MediaMTX.RTMP},
		{"mediamtx.api", c.MediaMTX.API},
		{"ingest.path", c.Ingest.Path},
		{"ingest.user", c.Ingest.User},
		{"ingest.password", c.Ingest.Password},
		{"obs.app", c.OBS.App},
	}
	var errs []error
	for _, r := range required {
		if r.value == "" {
			errs = append(errs, fmt.Errorf("%s is empty", r.key))
		}
	}
	return errors.Join(errs...)
}
