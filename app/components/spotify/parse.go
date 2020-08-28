package spotify

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

func ParseUriPattern(uri string) error {

	parsed := strings.Split(uri, ":")
	if len(parsed) == 3 {

		switch sType := parsed[1]; sType {
		case "track":
			trackId := strings.TrimSpace(parsed[2])
			log.Println(fmt.Sprintf("Getting Spotify Track : [%s]", trackId))

			track, err := GetTrack(trackId, "")
			if err != nil {
				return err
			}

			log.Println(track)
		case "episode":
			episodeId := strings.TrimSpace(parsed[2])
			log.Println(fmt.Sprintf("Getting Spotify Episode : [%s]", episodeId))

			track, err := GetTrack(episodeId, "")
			if err != nil {
				return err
			}

			log.Println(track)
			break
		default:
			return errors.New("only spotify [tracks|episodes] items are allowed")
		}
	}

	return errors.New("could not parse spotify uri pattern")
}