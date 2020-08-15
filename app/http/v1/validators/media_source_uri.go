package validators

import (
	"errors"
	"fmt"
	"net/url"
)

func MediaSourceUri(field string, _ string, _ string, value interface{}) error {

	val := value.(string)

	if val == "" {
		return nil
	}

	if _, err := url.ParseRequestURI(val); err == nil {
		return nil
	}

	return errors.New(fmt.Sprintf("The %s field format is invalid", field))
}