package subtitle

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"

	"github.com/asticode/go-astisub"
	"github.com/castyapp/api.server/app/components/strings"
	"github.com/castyapp/api.server/storage"
	"github.com/minio/minio-go"
)

// Convert and return subtitle files to WebVTT
func ConvertToVTT(file multipart.File) (buffer *bytes.Buffer, err error) {

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	vttSubtitle, err := astisub.ReadFromSRT(bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	buffer = new(bytes.Buffer)
	if err := vttSubtitle.WriteToWebVTT(buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}

func Save(sFile *multipart.FileHeader) (string, error) {

	subtitleName := strings.RandomNumber(20)

	subtitle, err := sFile.Open()
	if err != nil {
		return "", err
	}

	buf, err := ConvertToVTT(subtitle)
	if err != nil {
		return "", err
	}

	_, err = storage.Client.PutObject("subtitles", fmt.Sprintf("%s.vtt", subtitleName), buf, -1, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	return subtitleName, nil
}
