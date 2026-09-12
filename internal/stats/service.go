package stats

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"sort"
	"time"

	"spotify-stats/internal/httpclient"
	spotifydata "spotify-stats/internal/spotify-data"
)

const evolutionSampleLimit = "50"

type Service struct {
	dataService *spotifydata.Service
	client      *httpclient.Client
}

func NewService(dataService *spotifydata.Service) *Service {
	return &Service{dataService: dataService, client: httpclient.New()}
}

// TasteEvolution compares the user's top artists/tracks across Spotify's
// three time ranges to show how their taste has shifted over time.
func (s *Service) TasteEvolution() (TasteEvolutionResponse, error) {
	longArtists, err := s.dataService.FetchTopArtists(spotifydata.TopOptions{TimeRange: spotifydata.LongTerm, Limit: evolutionSampleLimit})
	if err != nil {
		return TasteEvolutionResponse{}, err
	}
	mediumArtists, err := s.dataService.FetchTopArtists(spotifydata.TopOptions{TimeRange: spotifydata.MediumTerm, Limit: evolutionSampleLimit})
	if err != nil {
		return TasteEvolutionResponse{}, err
	}
	shortArtists, err := s.dataService.FetchTopArtists(spotifydata.TopOptions{TimeRange: spotifydata.ShortTerm, Limit: evolutionSampleLimit})
	if err != nil {
		return TasteEvolutionResponse{}, err
	}

	longTracks, err := s.dataService.FetchTopTracks(spotifydata.TopOptions{TimeRange: spotifydata.LongTerm, Limit: evolutionSampleLimit})
	if err != nil {
		return TasteEvolutionResponse{}, err
	}
	mediumTracks, err := s.dataService.FetchTopTracks(spotifydata.TopOptions{TimeRange: spotifydata.MediumTerm, Limit: evolutionSampleLimit})
	if err != nil {
		return TasteEvolutionResponse{}, err
	}
	shortTracks, err := s.dataService.FetchTopTracks(spotifydata.TopOptions{TimeRange: spotifydata.ShortTerm, Limit: evolutionSampleLimit})
	if err != nil {
		return TasteEvolutionResponse{}, err
	}

	longArtistMap := artistIDNameMap(longArtists.Items)
	mediumArtistMap := artistIDNameMap(mediumArtists.Items)
	shortArtistMap := artistIDNameMap(shortArtists.Items)

	longTrackMap := trackIDNameMap(longTracks.Items)
	mediumTrackMap := trackIDNameMap(mediumTracks.Items)
	shortTrackMap := trackIDNameMap(shortTracks.Items)

	return TasteEvolutionResponse{
		LongToMedium: TasteEvolutionComparison{
			From:                 spotifydata.LongTerm,
			To:                   spotifydata.MediumTerm,
			ArtistOverlapPercent: overlapPercent(longArtistMap, mediumArtistMap),
			TrackOverlapPercent:  overlapPercent(longTrackMap, mediumTrackMap),
			NewArtistsInToRange:  newItems(longArtistMap, mediumArtistMap),
		},
		MediumToShort: TasteEvolutionComparison{
			From:                 spotifydata.MediumTerm,
			To:                   spotifydata.ShortTerm,
			ArtistOverlapPercent: overlapPercent(mediumArtistMap, shortArtistMap),
			TrackOverlapPercent:  overlapPercent(mediumTrackMap, shortTrackMap),
			NewArtistsInToRange:  newItems(mediumArtistMap, shortArtistMap),
		},
		LongToShort: TasteEvolutionComparison{
			From:                 spotifydata.LongTerm,
			To:                   spotifydata.ShortTerm,
			ArtistOverlapPercent: overlapPercent(longArtistMap, shortArtistMap),
			TrackOverlapPercent:  overlapPercent(longTrackMap, shortTrackMap),
			NewArtistsInToRange:  newItems(longArtistMap, shortArtistMap),
		},
	}, nil
}

// ListeningPatterns builds hour-of-day and day-of-week histograms from the
// user's recently played tracks.
func (s *Service) ListeningPatterns() (ListeningPatternsResponse, error) {
	recent, err := s.dataService.FetchRecentlyPlayed("50")
	if err != nil {
		return ListeningPatternsResponse{}, err
	}

	var hourHistogram [24]int
	dayHistogram := map[string]int{
		"Monday": 0, "Tuesday": 0, "Wednesday": 0, "Thursday": 0,
		"Friday": 0, "Saturday": 0, "Sunday": 0,
	}

	sampleSize := 0
	for _, item := range recent.Items {
		playedAt, err := time.Parse(time.RFC3339, item.PlayedAt)
		if err != nil {
			continue
		}
		hourHistogram[playedAt.Hour()]++
		dayHistogram[playedAt.Weekday().String()]++
		sampleSize++
	}

	busiestHour := 0
	for hour, count := range hourHistogram {
		if count > hourHistogram[busiestHour] {
			busiestHour = hour
		}
	}

	busiestDay := busiestDayOf(dayHistogram)

	return ListeningPatternsResponse{
		HourOfDayHistogram: hourHistogram,
		DayOfWeekHistogram: dayHistogram,
		BusiestHour:        busiestHour,
		BusiestDay:         busiestDay,
		SampleSize:         sampleSize,
	}, nil
}

// AlbumArtColors computes the average color of each unique album's artwork
// among the user's top tracks, as a lightweight stand-in for the "vibe"
// audio-features Spotify used to provide (now deprecated).
func (s *Service) AlbumArtColors(opts spotifydata.TopOptions) (AlbumColorsResponse, error) {
	tracks, err := s.dataService.FetchTopTracks(opts)
	if err != nil {
		return AlbumColorsResponse{}, err
	}

	seenAlbums := map[string]bool{}
	var colors []AlbumColor

	for _, track := range tracks.Items {
		if track.Album.ID == "" || seenAlbums[track.Album.ID] || len(track.Album.Images) == 0 {
			continue
		}
		seenAlbums[track.Album.ID] = true

		img := smallestImage(track.Album.Images)
		hex, err := s.averageColor(img.URL)
		if err != nil {
			// Skip albums whose artwork can't be fetched/decoded rather than
			// failing the whole request.
			continue
		}

		colors = append(colors, AlbumColor{
			AlbumName: track.Album.Name,
			ImageURL:  img.URL,
			HexColor:  hex,
		})
	}

	return AlbumColorsResponse{Colors: colors}, nil
}

func (s *Service) averageColor(imageURL string) (string, error) {
	body, err := s.client.MakeRequest(httpclient.GET, imageURL, "", nil)
	if err != nil {
		return "", err
	}

	decoded, _, err := image.Decode(bytes.NewReader([]byte(body)))
	if err != nil {
		return "", err
	}

	bounds := decoded.Bounds()
	const step = 4 // sample every 4th pixel to keep this fast
	var rSum, gSum, bSum, count uint64

	for y := bounds.Min.Y; y < bounds.Max.Y; y += step {
		for x := bounds.Min.X; x < bounds.Max.X; x += step {
			r, g, b, _ := decoded.At(x, y).RGBA()
			rSum += uint64(r >> 8)
			gSum += uint64(g >> 8)
			bSum += uint64(b >> 8)
			count++
		}
	}

	if count == 0 {
		return "", fmt.Errorf("no pixels sampled from image %s", imageURL)
	}

	return fmt.Sprintf("#%02X%02X%02X", rSum/count, gSum/count, bSum/count), nil
}

func artistIDNameMap(artists []spotifydata.Artist) map[string]string {
	m := make(map[string]string, len(artists))
	for _, a := range artists {
		m[a.ID] = a.Name
	}
	return m
}

func trackIDNameMap(tracks []spotifydata.Track) map[string]string {
	m := make(map[string]string, len(tracks))
	for _, t := range tracks {
		m[t.ID] = t.Name
	}
	return m
}

// overlapPercent returns the Jaccard similarity (as a 0-100 percentage)
// between two ID->name sets.
func overlapPercent(a, b map[string]string) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	intersection := 0
	for id := range a {
		if _, ok := b[id]; ok {
			intersection++
		}
	}

	union := len(a) + len(b) - intersection
	if union == 0 {
		return 0
	}

	return math.Round(float64(intersection)/float64(union)*10000) / 100
}

// newItems returns names present in `to` but absent from `from`.
func newItems(from, to map[string]string) []string {
	items := make([]string, 0)
	for id, name := range to {
		if _, ok := from[id]; !ok {
			items = append(items, name)
		}
	}
	sort.Strings(items)
	return items
}

func smallestImage(images []spotifydata.Image) spotifydata.Image {
	smallest := images[0]
	for _, img := range images {
		if img.Width > 0 && img.Width < smallest.Width {
			smallest = img
		}
	}
	return smallest
}

func busiestDayOf(dayHistogram map[string]int) string {
	busiestDay := ""
	busiestCount := -1
	// Iterate in a fixed order so ties resolve deterministically.
	for _, day := range []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"} {
		if dayHistogram[day] > busiestCount {
			busiestCount = dayHistogram[day]
			busiestDay = day
		}
	}
	return busiestDay
}
