package storage

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/CastyLab/api.server/config"
	"github.com/minio/minio-go"
)

var Client *minio.Client

func Configure() (err error) {
	log.Println(config.Map.Secrets)
	Client, err = minio.NewV4(
		config.Map.Secrets.ObjectStorage.Endpoint,
		config.Map.Secrets.ObjectStorage.AccessKey,
		config.Map.Secrets.ObjectStorage.SecretKey,
		config.Map.Secrets.ObjectStorage.UseHttps,
	)
	if err != nil {
		return err
	}
	if config.Map.Secrets.ObjectStorage.InsecureSkipVerify {
		Client.SetCustomTransport(&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		})
	}
	return nil
}
