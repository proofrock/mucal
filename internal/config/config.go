// Copyright 2026 Mano
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	TimeZone    string     `yaml:"time_zone"`
	AutoRefresh int        `yaml:"auto_refresh"`
	Calendars   []Calendar `yaml:"calendars"`
}

// Calendar represents a single CalDAV calendar configuration
type Calendar struct {
	Name         string `yaml:"name"`
	URL          string `yaml:"url"`
	UserID       string `yaml:"user_id"`
	PasswordFile string `yaml:"password_file"`
	Color        string `yaml:"color"`
}

// LoadConfig loads and validates the configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate timezone
	if c.TimeZone == "" {
		return fmt.Errorf("time_zone is required")
	}
	if _, err := time.LoadLocation(c.TimeZone); err != nil {
		return fmt.Errorf("invalid time_zone: %w", err)
	}

	// Validate auto_refresh
	if c.AutoRefresh <= 0 {
		return fmt.Errorf("auto_refresh must be positive")
	}

	// Validate calendars
	if len(c.Calendars) == 0 {
		return fmt.Errorf("at least one calendar is required")
	}

	for i, cal := range c.Calendars {
		if err := cal.Validate(); err != nil {
			return fmt.Errorf("calendar %d (%s): %w", i, cal.Name, err)
		}
	}

	return nil
}

// Validate validates a single calendar configuration
func (c *Calendar) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.URL == "" {
		return fmt.Errorf("url is required")
	}
	if c.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if c.PasswordFile == "" {
		return fmt.Errorf("password_file is required")
	}
	if c.Color == "" {
		return fmt.Errorf("color is required")
	}
	// Basic color validation (should start with #)
	if !strings.HasPrefix(c.Color, "#") || len(c.Color) != 7 {
		return fmt.Errorf("color must be in hex format (#RRGGBB)")
	}

	return nil
}

// GetPassword reads the password from the configured password file
func (c *Calendar) GetPassword() (string, error) {
	data, err := os.ReadFile(c.PasswordFile)
	if err != nil {
		return "", fmt.Errorf("failed to read password file %s: %w", c.PasswordFile, err)
	}

	// Trim any whitespace/newlines
	password := strings.TrimSpace(string(data))
	if password == "" {
		return "", fmt.Errorf("password file %s is empty", c.PasswordFile)
	}

	return password, nil
}

// GetLocation returns the time.Location for the configured timezone
func (c *Config) GetLocation() (*time.Location, error) {
	return time.LoadLocation(c.TimeZone)
}

// Sanitize returns a sanitized version of the config (without credentials)
// for sending to the frontend
func (c *Config) Sanitize() map[string]interface{} {
	cals := make([]map[string]interface{}, len(c.Calendars))
	for i, cal := range c.Calendars {
		cals[i] = map[string]interface{}{
			"name":  cal.Name,
			"color": cal.Color,
		}
	}

	return map[string]interface{}{
		"timezone":    c.TimeZone,
		"autoRefresh": c.AutoRefresh,
		"calendars":   cals,
	}
}
