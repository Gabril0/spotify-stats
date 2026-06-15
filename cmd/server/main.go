package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"log"
	"os"
	"spotify-stats/internal/auth"
	"spotify-stats/internal/spotify-data"
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
	data := spotifydata.New(a)

	app.Get("/api/callback", a.SpotifyCallback)
	app.Get("/api/web-auth", a.SpotifyAuthWebRedirect)
	app.Get("/api/token", a.GetToken)
	app.Get("/api/me", data.GetUserInfo)
	app.Get("/api/top/artists", data.GetTopArtists)
	app.Get("/api/top/tracks", data.GetTopTracks)
	app.Get("/api/recently-played", data.GetRecentlyPlayed)
	log.Fatal(app.Listen(":" + port))
}
