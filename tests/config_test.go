package tests

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/castyapp/api.server/config"
)

var (
	configFileName = "./config_test.hcl"
	defaultConfig  = &config.ConfigMap{
		Debug:    false,
		Metrics:  false,
		Env:      "dev",
		Timezone: "America/California",
		S3: config.S3Config{
			Endpoint:           "127.0.0.1:9000",
			AccessKey:          "secret-access-key",
			SecretKey:          "secret-key",
			InsecureSkipVerify: true,
			UseHttps:           false,
		},
		Http: config.HttpConfig{
			Rules: config.HttpConfigRules{
				AccessControlAllowOrigin: "*",
			},
		},
		Grpc: config.GrpcConfig{
			Host: "localhost",
			Port: 55283,
		},
		Sentry: config.SentryConfig{
			Enabled: false,
			Dsn:     "sentry.dsn.here",
		},
		Recaptcha: config.RecaptchaConfig{
			Enabled: false,
			Type:    "hcaptcha",
			Secret:  "hcaptcha-secret-token",
		},
	}
)

func TestLoadConfig(t *testing.T) {
	if err := config.Load(filepath.Join(configFileName)); err != nil {
		t.Fatalf("err: %s", err)
	}
	if !reflect.DeepEqual(defaultConfig, config.Map) {
		t.Fatalf("bad: %#v", config.Map)
	}
}
