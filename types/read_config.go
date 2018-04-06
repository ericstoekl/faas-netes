// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package types

import (
	"os"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"
)

// OsEnv implements interface to wrap os.Getenv
type OsEnv struct {
}

// Getenv wraps os.Getenv
func (OsEnv) Getenv(key string) string {
	return os.Getenv(key)
}

// HasEnv provides interface for os.Getenv
type HasEnv interface {
	Getenv(key string) string
}

// ReadConfig constitutes config from env variables
type ReadConfig struct {
}

func parseIntOrDurationValue(val string, fallback time.Duration) time.Duration {
	if len(val) > 0 {
		parsedVal, parseErr := strconv.Atoi(val)
		if parseErr == nil && parsedVal >= 0 {
			return time.Duration(parsedVal) * time.Second
		}
	}

	duration, durationErr := time.ParseDuration(val)
	if durationErr != nil {
		return fallback
	}

	return duration
}

func parseBoolValue(val string, fallback bool) bool {
	if len(val) > 0 {
		return val == "true"
	}
	return fallback
}

// Read fetches config from environmental variables.
func (ReadConfig) Read(hasEnv HasEnv) BootstrapConfig {
	cfg := BootstrapConfig{}

	enableProbe := parseBoolValue(hasEnv.Getenv("enable_function_readiness_probe"), true)

	readTimeout := parseIntOrDurationValue(hasEnv.Getenv("read_timeout"), time.Second*10)
	writeTimeout := parseIntOrDurationValue(hasEnv.Getenv("write_timeout"), time.Second*10)
	imagePullPolicy := hasEnv.Getenv("image_pull_policy")

	cfg.ReadTimeout = readTimeout
	cfg.WriteTimeout = writeTimeout
	cfg.ImagePullPolicy = v1.PullPolicy(imagePullPolicy)

	cfg.EnableFunctionReadinessProbe = enableProbe

	return cfg
}

// BootstrapConfig for the process.
type BootstrapConfig struct {
	EnableFunctionReadinessProbe bool
	ReadTimeout                  time.Duration
	WriteTimeout                 time.Duration
	ImagePullPolicy              v1.PullPolicy
}
