package spotifydata

import "github.com/gofiber/fiber/v3"

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetUserInfo(c fiber.Ctx) error {
	profile, err := h.service.FetchUserProfile()
	if err != nil {
		return err
	}

	return c.JSON(profile)
}

func (h *Handler) GetTopArtists(c fiber.Ctx) error {
	opts := TopOptions{}
	err := c.Bind().Body(&opts)
	if err != nil {
		return err
	}

	artists, err := h.service.FetchTopArtists(opts)
	if err != nil {
		return err
	}

	return c.JSON(artists)
}

func (h *Handler) GetTopTracks(c fiber.Ctx) error {
	opts := TopOptions{}
	err := c.Bind().Body(&opts)
	if err != nil {
		return err
	}

	tracks, err := h.service.FetchTopTracks(opts)
	if err != nil {
		return err
	}

	return c.JSON(tracks)
}

func (h *Handler) GetRecentlyPlayed(c fiber.Ctx) error {
	opts := TopOptions{}
	err := c.Bind().Body(&opts)
	if err != nil {
		return err
	}

	recent, err := h.service.FetchRecentlyPlayed(opts.Limit)
	if err != nil {
		return err
	}

	return c.JSON(recent)
}
