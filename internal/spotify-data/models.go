package spotifydata

type TopOptions struct {
	TimeRange TimeRange `json:"time_range"`
	Limit     string    `json:"limit"`
}

type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type UserProfile struct {
	DisplayName string  `json:"display_name"`
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	Country     string  `json:"country"`
	Images      []Image `json:"images"`
	Followers   struct {
		Total int `json:"total"`
	} `json:"followers"`
}

type Artist struct {
	Name       string   `json:"name"`
	ID         string   `json:"id"`
	Genres     []string `json:"genres"`
	Popularity int      `json:"popularity"`
	Images     []Image  `json:"images"`
}

type Album struct {
	Name   string  `json:"name"`
	Images []Image `json:"images"`
}

type Track struct {
	Name       string   `json:"name"`
	ID         string   `json:"id"`
	Popularity int      `json:"popularity"`
	Album      Album    `json:"album"`
	Artists    []Artist `json:"artists"`
}

type TopArtistsResponse struct {
	Items []Artist `json:"items"`
}

type TopTracksResponse struct {
	Items []Track `json:"items"`
}

type RecentlyPlayedResponse struct {
	Items []struct {
		Track    Track  `json:"track"`
		PlayedAt string `json:"played_at"`
	} `json:"items"`
}
