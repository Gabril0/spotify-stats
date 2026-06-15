package spotifydata

import (
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"net/url"
	"spotify-stats/internal/auth"
	"spotify-stats/internal/httpclient"
)

const spotifyAPIBaseURL = "https://api.spotify.com/v1/"

type TimeRange string

const (
	ShortTerm  TimeRange = "short_term"
	MediumTerm TimeRange = "medium_term"
	LongTerm   TimeRange = "long_term"
)

const defaultLimit = "20"

type SpotifyData struct {
	auth             *auth.Auth
	client           *httpclient.Client
	defaultTimeRange TimeRange
	defaultLimit     string
}

func New(a *auth.Auth) *SpotifyData {
	return &SpotifyData{auth: a, client: httpclient.New(), defaultTimeRange: MediumTerm, defaultLimit: defaultLimit}
}

func (s *SpotifyData) QuickGet(path string) (string, error) {
	token, err := s.auth.ValidToken()
	if err != nil {
		return "", err
	}

	headers := map[string]string{"Authorization": "Bearer " + token}

	return s.client.MakeRequest(httpclient.GET, spotifyAPIBaseURL+path, "", headers)
}

func (s *SpotifyData) topQuery(c fiber.Ctx) (string, error) {
	opts := TopOptions{}
	if err := c.Bind().Body(&opts); err != nil {
		return "", err
	}

	if opts.TimeRange == "" {
		opts.TimeRange = s.defaultTimeRange
	}
	if opts.Limit == "" {
		opts.Limit = s.defaultLimit
	}

	params := url.Values{}
	params.Set("time_range", string(opts.TimeRange))
	params.Set("limit", opts.Limit)
	return "?" + params.Encode(), nil
}

func (s *SpotifyData) GetUserInfo(c fiber.Ctx) error {
	resp, err := s.QuickGet("me")
	if err != nil {
		return err
	}

	var profile UserProfile
	err = json.Unmarshal([]byte(resp), &profile)
	if err != nil {
		return err
	}

	return c.JSON(profile)
}

func (s *SpotifyData) GetTopArtists(c fiber.Ctx) error {
	query, err := s.topQuery(c)
	if err != nil {
		return err
	}

	resp, err := s.QuickGet("me/top/artists" + query)
	if err != nil {
		return err
	}

	var artists TopArtistsResponse
	err = json.Unmarshal([]byte(resp), &artists)
	if err != nil {
		return err
	}

	return c.JSON(artists)
}

func (s *SpotifyData) GetTopTracks(c fiber.Ctx) error {
	query, err := s.topQuery(c)
	if err != nil {
		return err
	}

	resp, err := s.QuickGet("me/top/tracks" + query)
	if err != nil {
		return err
	}

	var tracks TopTracksResponse
	err = json.Unmarshal([]byte(resp), &tracks)
	if err != nil {
		return err
	}

	return c.JSON(tracks)
}

func (s *SpotifyData) GetRecentlyPlayed(c fiber.Ctx) error {
	opts := TopOptions{}
	if err := c.Bind().Body(&opts); err != nil {
		return err
	}
	if opts.Limit == "" {
		opts.Limit = s.defaultLimit
	}

	params := url.Values{}
	params.Set("limit", opts.Limit)

	resp, err := s.QuickGet("me/player/recently-played?" + params.Encode())
	if err != nil {
		return err
	}

	var recent RecentlyPlayedResponse
	err = json.Unmarshal([]byte(resp), &recent)
	if err != nil {
		return err
	}

	return c.JSON(recent)
}
