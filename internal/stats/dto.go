package stats

import spotifydata "spotify-stats/internal/spotify-data"

// TasteEvolutionComparison expresses how a user's top artists/tracks changed
// between two Spotify time ranges (e.g. their all-time favorites vs. what
// they've been playing in the last 4 weeks).
type TasteEvolutionComparison struct {
	From                 spotifydata.TimeRange `json:"from"`
	To                   spotifydata.TimeRange `json:"to"`
	ArtistOverlapPercent float64               `json:"artist_overlap_percent"`
	TrackOverlapPercent  float64               `json:"track_overlap_percent"`
	NewArtistsInToRange  []string              `json:"new_artists_in_to_range"`
}

type TasteEvolutionResponse struct {
	LongToMedium  TasteEvolutionComparison `json:"long_to_medium"`
	MediumToShort TasteEvolutionComparison `json:"medium_to_short"`
	LongToShort   TasteEvolutionComparison `json:"long_to_short"`
}

// ListeningPatternsResponse summarizes when (hour of day / day of week) a
// user tends to listen, based on their recently played tracks.
type ListeningPatternsResponse struct {
	HourOfDayHistogram [24]int        `json:"hour_of_day_histogram"`
	DayOfWeekHistogram map[string]int `json:"day_of_week_histogram"`
	BusiestHour        int            `json:"busiest_hour"`
	BusiestDay         string         `json:"busiest_day"`
	SampleSize         int            `json:"sample_size"`
}

// AlbumColor is the computed average color of an album's artwork.
type AlbumColor struct {
	AlbumName string `json:"album_name"`
	ImageURL  string `json:"image_url"`
	HexColor  string `json:"hex_color"`
}

type AlbumColorsResponse struct {
	Colors []AlbumColor `json:"colors"`
}
