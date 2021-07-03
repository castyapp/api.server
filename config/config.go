package config

import (
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

var validBuckets = []string{
	"avatars",
	"subtitles",
	"posters",
}

type ConfMap struct {
	Debug     bool            `hcl:"debug"`
	Env       string          `hcl:"env"`
	Grpc      GrpcConfig      `hcl:"grpc,block"`
	HTTP      HTTPConfig      `hcl:"http,block"`
	S3        S3Config        `hcl:"s3,block"`
	Sentry    SentryConfig    `hcl:"sentry,block"`
	Recaptcha RecaptchaConfig `hcl:"recaptcha,block"`
}

type HTTPConfig struct {
	Rules HTTPConfigRules `hcl:"rules,label"`
}

type HTTPConfigRules struct {
	AccessControlAllowOrigin string `hcl:"access_control_allow_origin"`
}

type GrpcConfig struct {
	Host string `hcl:"host"`
	Port int    `hcl:"port"`
}

type S3Config struct {
	Endpoint           string `hcl:"endpoint"`
	AccessKey          string `hcl:"access_key"`
	SecretKey          string `hcl:"secret_key"`
	UseHTTPS           bool   `hcl:"use_https"`
	InsecureSkipVerify bool   `hcl:"insecure_skip_verify"`
}

type SentryConfig struct {
	Enabled bool   `hcl:"enabled"`
	Dsn     string `hcl:"dsn"`
}

type RecaptchaConfig struct {
	Enabled bool   `hcl:"enabled"`
	Type    string `hcl:"type"`
	Secret  string `hcl:"secret"`
}

var Map = new(ConfMap)

func Load(filename string) (err error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	obj, err := hcl.Parse(string(d))
	if err != nil {
		return err
	}
	// Build up the result
	if err := hcl.DecodeObject(&Map, obj); err != nil {
		return err
	}
	return
}

func IsValidBucketName(bucketname string) bool {
	for _, bk := range validBuckets {
		if bucketname == bk {
			return true
		}
	}
	return false
}
