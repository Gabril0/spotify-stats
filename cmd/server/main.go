package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"log"
	"os"
	"spotify-stats/internal/auth"
	spotifydata "spotify-stats/internal/spotify-data"
	"spotify-stats/internal/stats"
)

func main() {
	app := fiber.New()
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Please copy .env.example and fill in the credentials")
	}

	clientID := os.Getenv("CLIENT_ID")
	port := os.Getenv("PORT")
	clientSecret := os.Getenv("CLIENT_SECRET")

	a := auth.New(clientID, clientSecret, port)
	dataService := spotifydata.NewService(a)
	dataHandler := spotifydata.NewHandler(dataService)
	statsService := stats.NewService(dataService)
	statsHandler := stats.NewHandler(statsService)

	app.Get("/api/callback", a.SpotifyCallback)
	app.Get("/api/web-auth", a.SpotifyAuthWebRedirect)
	app.Get("/api/token", a.GetToken)
	app.Get("/api/me", dataHandler.GetUserInfo)
	app.Get("/api/top/artists", dataHandler.GetTopArtists)
	app.Get("/api/top/tracks", dataHandler.GetTopTracks)
	app.Get("/api/recently-played", dataHandler.GetRecentlyPlayed)
	app.Get("/api/saved-tracks", dataHandler.GetSavedTracks)
	app.Get("/api/stats/taste-evolution", statsHandler.GetTasteEvolution)
	app.Get("/api/stats/listening-patterns", statsHandler.GetListeningPatterns)
	app.Get("/api/stats/album-colors", statsHandler.GetAlbumColors)
	log.Fatal(app.Listen(":" + port))
}
