package spotifydata

import (
	"encoding/json"
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

type Service struct {
	auth             *auth.Auth
	client           *httpclient.Client
	defaultTimeRange TimeRange
	defaultLimit     string
}

func NewService(a *auth.Auth) *Service {
	return &Service{auth: a, client: httpclient.New(), defaultTimeRange: MediumTerm, defaultLimit: defaultLimit}
}

func (s *Service) QuickGet(path string) (string, error) {
	token, err := s.auth.ValidToken()
	if err != nil {
		return "", err
	}

	headers := map[string]string{"Authorization": "Bearer " + token}

	return s.client.MakeRequest(httpclient.GET, spotifyAPIBaseURL+path, "", headers)
}

func (s *Service) topQuery(opts TopOptions) string {
	if opts.TimeRange == "" {
		opts.TimeRange = s.defaultTimeRange
	}
	if opts.Limit == "" {
		opts.Limit = s.defaultLimit
	}

	params := url.Values{}
	params.Set("time_range", string(opts.TimeRange))
	params.Set("limit", opts.Limit)
	return "?" + params.Encode()
}

func (s *Service) FetchUserProfile() (UserProfile, error) {
	resp, err := s.QuickGet("me")
	if err != nil {
		return UserProfile{}, err
	}

	var profile UserProfile
	err = json.Unmarshal([]byte(resp), &profile)
	if err != nil {
		return UserProfile{}, err
	}

	return profile, nil
}

func (s *Service) FetchTopArtists(opts TopOptions) (TopArtistsResponse, error) {
	resp, err := s.QuickGet("me/top/artists" + s.topQuery(opts))
	if err != nil {
		return TopArtistsResponse{}, err
	}

	var artists TopArtistsResponse
	err = json.Unmarshal([]byte(resp), &artists)
	if err != nil {
		return TopArtistsResponse{}, err
	}

	return artists, nil
}

func (s *Service) FetchTopTracks(opts TopOptions) (TopTracksResponse, error) {
	resp, err := s.QuickGet("me/top/tracks" + s.topQuery(opts))
	if err != nil {
		return TopTracksResponse{}, err
	}

	var tracks TopTracksResponse
	err = json.Unmarshal([]byte(resp), &tracks)
	if err != nil {
		return TopTracksResponse{}, err
	}

	return tracks, nil
}

func (s *Service) FetchRecentlyPlayed(limit string) (RecentlyPlayedResponse, error) {
	if limit == "" {
		limit = s.defaultLimit
	}

	params := url.Values{}
	params.Set("limit", limit)

	resp, err := s.QuickGet("me/player/recently-played?" + params.Encode())
	if err != nil {
		return RecentlyPlayedResponse{}, err
	}

	var recent RecentlyPlayedResponse
	err = json.Unmarshal([]byte(resp), &recent)
	if err != nil {
		return RecentlyPlayedResponse{}, err
	}

	return recent, nil
}
