package validators

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func MediaSourceUri(field string, _ string, _ string, value interface{}) error {
	val := value.(string)
	if val == "" {
		return nil
	}
	if uri, err := url.ParseRequestURI(val); err == nil {
		if strings.Contains(uri.Host, "spotify") {
			return nil
		}
		return nil
	}
	return errors.New(fmt.Sprintf("The %s field format is invalid", field))
}