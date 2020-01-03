package validators

import (
	"errors"
	"fmt"
)

func Access(field string, rule string, message string, value interface{}) error {

	val := value.(string)

	if val == "" {
		return nil
	}

	switch val {
	case "1", "2", "3":
		return nil
	}

	return errors.New(fmt.Sprintf("The %s field format is invalid", field))
}