package spotify

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

func ParseURIPattern(uri string) error {

	parsed := strings.Split(uri, ":")
	if len(parsed) == 3 {

		switch sType := parsed[1]; sType {
		case "track":
			trackID := strings.TrimSpace(parsed[2])
			log.Println(fmt.Sprintf("Getting Spotify Track : [%s]", trackID))

			track, err := GetTrack(trackID, "")
			if err != nil {
				return err
			}

			log.Println(track)
		case "episode":
			episodeID := strings.TrimSpace(parsed[2])
			log.Println(fmt.Sprintf("Getting Spotify Episode : [%s]", episodeID))

			track, err := GetTrack(episodeID, "")
			if err != nil {
				return err
			}

			log.Println(track)
		default:
			return errors.New("only spotify [tracks|episodes] items are allowed")
		}
	}

	return errors.New("could not parse spotify uri pattern")
}
