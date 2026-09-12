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

// GetTasteEvolution godoc
// @Summary Compare top artists/tracks across time ranges
// @Tags stats
// @Produce json
// @Success 200 {object} TasteEvolutionResponse
// @Router /api/stats/taste-evolution [get]
func (h *Handler) GetTasteEvolution(c fiber.Ctx) error {
	result, err := h.service.TasteEvolution()
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// GetListeningPatterns godoc
// @Summary Get hour-of-day and day-of-week listening histograms
// @Tags stats
// @Produce json
// @Success 200 {object} ListeningPatternsResponse
// @Router /api/stats/listening-patterns [get]
func (h *Handler) GetListeningPatterns(c fiber.Ctx) error {
	result, err := h.service.ListeningPatterns()
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// GetAlbumColors godoc
// @Summary Get the average color of each top track's album artwork
// @Tags stats
// @Produce json
// @Param time_range query string false "short_term, medium_term, or long_term"
// @Param limit query string false "Number of items to return"
// @Success 200 {object} AlbumColorsResponse
// @Router /api/stats/album-colors [get]
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
