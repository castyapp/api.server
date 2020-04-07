package subtitle

import (
	"bytes"
	"fmt"
	"github.com/CastyLab/api.server/app/components/strings"
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

func Save(sFile *multipart.FileHeader) (file *os.File, err error) {

	subtitle, err := sFile.Open()
	if err != nil {
		return nil, err
	}

	buf, err := ConvertToVTT(subtitle)
	if err != nil {
		return nil, err
	}

	subtitleName := strings.RandomNumber(20)

	file, err = os.Create(fmt.Sprintf("./storage/uploads/subtitles/%s.vtt", subtitleName))
	if err != nil {
		return
	}

	if _, err := file.Write(buf.Bytes()); err != nil {
		return nil, err
	}

	return file, nil
}