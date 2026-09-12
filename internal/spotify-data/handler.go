package spotifydata

import "github.com/gofiber/fiber/v3"

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// GetUserInfo godoc
// @Summary Get current user's Spotify profile
// @Tags data
// @Produce json
// @Success 200 {object} UserProfile
// @Router /api/me [get]
func (h *Handler) GetUserInfo(c fiber.Ctx) error {
	profile, err := h.service.FetchUserProfile()
	if err != nil {
		return err
	}

	return c.JSON(profile)
}

// GetTopArtists godoc
// @Summary Get the user's top artists
// @Tags data
// @Produce json
// @Param time_range query string false "short_term, medium_term, or long_term"
// @Param limit query string false "Number of items to return"
// @Success 200 {object} TopArtistsResponse
// @Router /api/top/artists [get]
func (h *Handler) GetTopArtists(c fiber.Ctx) error {
	opts := TopOptions{}
	err := c.Bind().Query(&opts)
	if err != nil {
		return err
	}

	artists, err := h.service.FetchTopArtists(opts)
	if err != nil {
		return err
	}

	return c.JSON(artists)
}

// GetTopTracks godoc
// @Summary Get the user's top tracks
// @Tags data
// @Produce json
// @Param time_range query string false "short_term, medium_term, or long_term"
// @Param limit query string false "Number of items to return"
// @Success 200 {object} TopTracksResponse
// @Router /api/top/tracks [get]
func (h *Handler) GetTopTracks(c fiber.Ctx) error {
	opts := TopOptions{}
	err := c.Bind().Query(&opts)
	if err != nil {
		return err
	}

	tracks, err := h.service.FetchTopTracks(opts)
	if err != nil {
		return err
	}

	return c.JSON(tracks)
}

// GetRecentlyPlayed godoc
// @Summary Get the user's recently played tracks
// @Tags data
// @Produce json
// @Param limit query string false "Number of items to return"
// @Success 200 {object} RecentlyPlayedResponse
// @Router /api/recently-played [get]
func (h *Handler) GetRecentlyPlayed(c fiber.Ctx) error {
	opts := TopOptions{}
	err := c.Bind().Query(&opts)
	if err != nil {
		return err
	}

	recent, err := h.service.FetchRecentlyPlayed(opts.Limit)
	if err != nil {
		return err
	}

	return c.JSON(recent)
}

// GetSavedTracks godoc
// @Summary Get the user's saved (liked) tracks
// @Tags data
// @Produce json
// @Param limit query string false "Number of items to return"
// @Success 200 {object} SavedTracksResponse
// @Router /api/saved-tracks [get]
func (h *Handler) GetSavedTracks(c fiber.Ctx) error {
	opts := TopOptions{}
	err := c.Bind().Query(&opts)
	if err != nil {
		return err
	}

	saved, err := h.service.FetchSavedTracks(opts.Limit)
	if err != nil {
		return err
	}

	return c.JSON(saved)
}
