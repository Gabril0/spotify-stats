package auth

import (
	"io"
	"net/http"
	"net/url"
	"spotify-stats/internal/httpclient"
	"github.com/gofiber/fiber/v3"
)

const spotifyBaseURL = "https://accounts.spotify.com/"

type Auth struct {
	clientID string
	clientSecret string
	port string
	userState string
	userCode string
	client *httpclient.Client
}

func New(clientID string, clientSecret string, port string) *Auth {
	var client *httpclient.Client = httpclient.New()
	return &Auth{clientID: clientID, clientSecret: clientSecret, port: port, client: client}
}

//func SpotifyAuthWeb(c fiber.Ctx) error {
// resp, err := http.Get(spotifyBaseURL + "authorize")
// 	if err != nil{
// 		return err
// 	}
// 	defer resp.Body.Close()
//
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil{
// 		return err
// 	}
// 	return c.SendString(string(body))
// }

func (a *Auth) SpotifyAuthWebRedirect(c fiber.Ctx) error {
	params := url.Values{}
	params.Set("client_id", a.clientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", "http://127.0.0.1:" + a.port + "/api/callback")
	params.Set("scope", "user-top-read user-read-recently-played")
	params.Set("state", "spotify island auth request")

	loginURL := spotifyBaseURL + "authorize?" + params.Encode()

	return c.JSON(fiber.Map{"url": loginURL})
}

func (a *Auth) SpotifyCallback(c fiber.Ctx) error {
	a.userState = c.Query("state")
	a.userCode = c.Query("code") 

	if a.userCode == "" || a.userState == ""{
		return c.SendString("Error: Not all parameters were returned, \n code:" + a.userCode + "\n state:" + a.userState)
	}
	return c.SendString("Authorized successfully!! Yay ꉂ(˵˃ ᗜ ˂˵)")
}

func (a *Auth) SpotifyGetToken(c fiber.Ctx) error {
	var tokenURL string = spotifyBaseURL + "api/token"

	resp, err := http.Get(tokenURL)
	if err != nil {
		return c.SendString(string(err.Error()))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil{
		return c.SendString(string(err.Error()))
	}
	return c.SendString(string(body))
}

func (a *Auth) GetToken(){
	a.client.MakeRequest(spotifyBaseURL + "api/token", "", http.MethodPost)
	
}
