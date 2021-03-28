package models

import (
	"errors"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/castyapp/api.server/app/components/spotify"
	rnd "github.com/castyapp/api.server/app/components/strings"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/knadh/go-get-youtube/youtube"
)

type MediaFile struct {
	Id     string        `json:"id"`
	Title  string        `json:"title"`
	Length time.Duration `json:"length"`
	Size   int           `json:"size"`
}

type MediaSource struct {
	u           string
	trackType   string
	accessToken string
	proto       *proto.MediaSource
}

func (m *MediaSource) Proto() *proto.MediaSource {
	return m.proto
}

func NewMediaSource(uri string, accessToken string) *MediaSource {
	return &MediaSource{u: uri, accessToken: accessToken}
}

func (m *MediaSource) parseYoutube() error {
	video, err := youtube.Get(m.u)
	if err != nil {
		return err
	}
	m.proto.Type = proto.MediaSource_YOUTUBE
	m.proto.Banner = video.Thumbnail_url
	m.proto.Title = video.Title
	m.proto.Artist = video.Author
	m.proto.Length = int64(video.Length_seconds)
	return nil
}

func (m *MediaSource) parseSpotify(id string) error {
	switch m.trackType {
	case "track":
		track, err := spotify.GetTrack(id, m.accessToken)
		if err != nil {
			return err
		}
		m.proto.Type = proto.MediaSource_SPOTIFY
		m.proto.Title = track.Name
		m.proto.Length = int64(track.DurationMs / 1000)
		m.proto.Banner = track.Album.Images[0].URL
		m.proto.Artist = ""
		for i, artist := range track.Artists {
			m.proto.Artist += artist.Name
			if (len(track.Artists) - 1) != i {
				m.proto.Artist += ", "
			}
		}
		return nil
	case "episode":
		episode, err := spotify.GetEpisode(id, m.accessToken)
		if err != nil {
			return err
		}
		m.proto.Type = proto.MediaSource_SPOTIFY
		m.proto.Title = episode.Name
		m.proto.Length = int64(episode.DurationMs / 1000)
		m.proto.Banner = episode.Images[0].URL
		m.proto.Artist = episode.Show.Name
		return nil
	}
	return errors.New("could not parse spotify")
}

func (m *MediaSource) parseDownloadUri() error {
	response, err := http.Get(m.u)
	if err != nil {
		m.proto.Type = proto.MediaSource_UNKNOWN
		return errors.New("type is unknown")
	}
	contentType := response.Header.Get("Content-Type")
	switch contentType {
	case "video/mp4", "application/octet-stream":
		m.proto.Type = proto.MediaSource_DOWNLOAD_URI
		m.proto.Title = getMovieTitle(m.u)
		duration, err := getMovieDuration(m.u)
		if err == nil {
			m.proto.Length = duration
		}
		return nil
	case "audio/x-mpegurl", "application/vnd.apple.mpegurl":
		m.proto.Type = proto.MediaSource_M3U8
		m.proto.Title = rnd.Random(20)
		return nil
	default:
		m.proto.Type = proto.MediaSource_UNKNOWN
		return errors.New("type is unknown")
	}
}

func getMovieTitle(uri string) string {
	fileUrl, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	path := fileUrl.Path
	segments := strings.Split(path, "/")
	return segments[len(segments)-1]
}

func getMovieDuration(uri string) (duration int64, err error) {
	command := exec.Command("ffprobe", "-show_streams", uri)
	output, err := command.Output()
	if err != nil {
		return 0, err
	}
	strOutput := string(output)
	dRegex := regexp.MustCompile(`(?m)duration=([0-9]+).([0-9]+)`)
	if len(dRegex.FindStringIndex(strOutput)) > 0 {
		strDuration := strings.ReplaceAll(dRegex.FindString(strOutput), "duration=", "")
		splits := strings.Split(strDuration, ".")
		durationInt, err := strconv.Atoi(splits[0])
		if err != nil {
			return 0, err
		}
		return int64(durationInt), nil
	}
	err = errors.New("could not get the duration parameter")
	return
}

func (m *MediaSource) Parse() error {
	m.proto = &proto.MediaSource{Uri: m.u}
	u, err := url.ParseRequestURI(m.u)
	if err != nil {
		return err
	}
	switch domain := strings.ReplaceAll(u.Hostname(), "www.", ""); domain {
	case "youtube.com", "yt.com":
		return m.parseYoutube()
	case "spotify", "open.spotify.com":
		parsed := strings.Split(u.Path, "/")
		switch strings.TrimSpace(parsed[1]) {
		case "track":
			m.trackType = "track"
			break
		case "episode":
			m.trackType = "episode"
			break
		default:
			return errors.New("could not parse spotify")
		}
		return m.parseSpotify(strings.TrimSpace(parsed[2]))
	default:
		return m.parseDownloadUri()
	}
}

func (m *MediaSource) IsUnknown() bool {
	return m.proto.Type == proto.MediaSource_UNKNOWN
}

func (m *MediaSource) IsYoutube() bool {
	return m.proto.Type == proto.MediaSource_YOUTUBE
}

func (m *MediaSource) IsTorrent() bool {
	return m.proto.Type == proto.MediaSource_TORRENT
}

func (m *MediaSource) IsSoundCloud() bool {
	return m.proto.Type == proto.MediaSource_SOUND_CLOUD
}

func (m *MediaSource) IsSpotify() bool {
	return m.proto.Type == proto.MediaSource_SPOTIFY
}

func (m *MediaSource) IsDownloadUri() bool {
	return m.proto.Type == proto.MediaSource_DOWNLOAD_URI
}
