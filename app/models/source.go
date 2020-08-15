package models

import (
	"errors"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/knadh/go-get-youtube/youtube"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MediaFile struct {
	Id      string         `json:"id"`
	Title   string         `json:"title"`
	Length  time.Duration  `json:"length"`
	Size    int            `json:"size"`
}

func (mf *MediaFile) Download() {

}

type MediaSource struct {
	u string
	proto *proto.MediaSource
}

func (m *MediaSource) Proto() *proto.MediaSource {
	return m.proto
}

func NewMediaSource(uri string) *MediaSource {
	return &MediaSource{u: uri}
}

//var (
//	cErr error
//	premiumSize int
//	client *torrent.Client
//)
//
//func init() {
//	client, cErr = newTorrentClient()
//	premiumSize, _ = strconv.Atoi(os.Getenv("PREMIUM_SIZE"))
//}

func (m *MediaSource) parseYoutube() error {
	video, err := youtube.Get(m.u)
	if err != nil {
		return err
	}
	m.proto.Type = proto.MediaSource_YOUTUBE
	m.proto.Banner = video.Thumbnail_url
	m.proto.Title = video.Title
	m.proto.Length = int64(video.Length_seconds)
	return nil
}

func (m *MediaSource) parseDownloadUri() error {

	response, err := http.Get(m.u)
	if err != nil {
		m.proto.Type = proto.MediaSource_UNKNOWN
		return errors.New("type is unknown")
	}

	contentType := response.Header.Get("Content-Type")
	switch contentType {
	// valid types
	case "video/mp4", "application/octet-stream":
		m.proto.Type = proto.MediaSource_DOWNLOAD_URI
		m.proto.Title = getMovieTitle(m.u)
		duration, err := getMovieDuration(m.u)
		if err == nil {
			m.proto.Length = duration
		}
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

//func (m *MediaSource) parseTorrent() error {
//	if cErr != nil {
//		return cErr
//	}
//	t, err := client.AddMagnet(m.u)
//	if err != nil {
//		return err
//	}
//	m.Type = proto.MediaSource_TORRENT
//	select {
//	case <-t.GotInfo():
//
//		for index, file := range t.Files() {
//			m.Files = append(m.Files, MediaFile{
//				Id: strconv.Itoa(index),
//				Title: file.Path(),
//				Size: int(file.Length()),
//				file: file,
//			})
//		}
//
//		if m.calcSize() > uint(premiumSize) {
//			m.Premium = true
//			m.PremiumText = "This source has more than 1GB files, You have to purchase premium to continue!"
//		}
//
//	}
//	return nil
//}

//func newTorrentClient() (*torrent.Client, error) {
//	config := torrent.NewDefaultClientConfig()
//	config.DataDir = "storage/uploads/media/"
//	config.Logger = tLog.Default
//	return torrent.NewClient(config)
//}

func (m *MediaSource) Parse() error {

	m.proto = &proto.MediaSource{Uri: m.u}

	u, err := url.ParseRequestURI(m.u)
	if err != nil {
		return err
	}

	domain := strings.ReplaceAll(u.Hostname(), "www.", "")

	switch domain {
	case "youtube.com", "yt.com":
		return m.parseYoutube()
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

func (m *MediaSource) IsDownloadUri() bool {
	return m.proto.Type == proto.MediaSource_DOWNLOAD_URI
}
