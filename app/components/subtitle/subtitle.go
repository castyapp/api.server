package subtitle

import (
	"bytes"
	"fmt"
	"github.com/CastyLab/api.server/app/components/strings"
	"github.com/CastyLab/api.server/config"
	"github.com/asticode/go-astisub"
	"io/ioutil"
	"mime/multipart"
	"os"
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

	subtitleFileName := fmt.Sprintf("%s/uploads/subtitles/%s.vtt", config.Map.StoragePath, subtitleName)
	file, err := os.Create(subtitleFileName)
	if err != nil {
		return "", err
	}

	if _, err := file.Write(buf.Bytes()); err != nil {
		return "", err
	}

	return subtitleName, nil
}