package stats

import (
	"github.com/gofiber/fiber/v3"

	spotifydata "spotify-stats/internal/spotify-data"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetTasteEvolution(c fiber.Ctx) error {
	result, err := h.service.TasteEvolution()
	if err != nil {
		return err
	}

	return c.JSON(result)
}

func (h *Handler) GetListeningPatterns(c fiber.Ctx) error {
	result, err := h.service.ListeningPatterns()
	if err != nil {
		return err
	}

	return c.JSON(result)
}

func (h *Handler) GetAlbumColors(c fiber.Ctx) error {
	opts := spotifydata.TopOptions{}
	err := c.Bind().Query(&opts)
	if err != nil {
		return err
	}

	result, err := h.service.AlbumArtColors(opts)
	if err != nil {
		return err
	}

	return c.JSON(result)
}
