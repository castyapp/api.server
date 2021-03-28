package storage

import (
	"crypto/tls"
	"net/http"

	"github.com/castyapp/api.server/config"
	"github.com/minio/minio-go"
)

var Client *minio.Client

func Configure() (err error) {
	Client, err = minio.NewV4(
		config.Map.S3.Endpoint,
		config.Map.S3.AccessKey,
		config.Map.S3.SecretKey,
		config.Map.S3.UseHttps,
	)
	if err != nil {
		return err
	}
	if config.Map.S3.InsecureSkipVerify {
		Client.SetCustomTransport(&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		})
	}
	return nil
}
