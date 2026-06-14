package main

import (
	"log"
	"os"
	"spotify-stats/internal/auth"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)


func main(){
	var app *fiber.App = fiber.New()
	err := godotenv.Load()

  if err != nil {
      log.Fatal("Please copy .env.example and fill in the credentials")
  }
	
	var clientID string = os.Getenv("CLIENT_ID")
	var port string = os.Getenv("PORT")
	var clientSecret string = os.Getenv("CLIENT_SECRET")

	var a *auth.Auth = auth.New(clientID, clientSecret, port )

	app.Get("/api/callback", a.SpotifyCallback)
	app.Get("/api/web-auth", a.SpotifyAuthWebRedirect)
	log.Fatal(app.Listen(":" + port))
}



