package auth

import (
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"net/url"
	"spotify-stats/internal/httpclient"
	"time"
)

const spotifyBaseURL = "https://accounts.spotify.com/"

type Auth struct {
	clientID        string
	clientSecret    string
	port            string
	userState       string
	userCode        string
	token           string
	tokenExpireTime time.Time
	client          *httpclient.Client
}

func New(clientID string, clientSecret string, port string) *Auth {
	return &Auth{
		clientID:     clientID,
		clientSecret: clientSecret,
		port:         port,
		client:       httpclient.New(),
	}
}

func (a *Auth) SpotifyAuthWebRedirect(c fiber.Ctx) error {
	params := url.Values{}
	params.Set("client_id", a.clientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", "http://127.0.0.1:"+a.port+"/api/callback")
	params.Set("scope", "user-top-read user-read-recently-played")
	params.Set("state", "spotify island auth request")

	loginURL := spotifyBaseURL + "authorize?" + params.Encode()

	return c.JSON(fiber.Map{"url": loginURL})
}

func (a *Auth) SpotifyCallback(c fiber.Ctx) error {
	a.userState = c.Query("state")
	a.userCode = c.Query("code")

	if a.userCode == "" || a.userState == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Not all parameters were returned",
			"code":  a.userCode,
			"state": a.userState,
		})
	}
	return c.JSON(fiber.Map{"message": "Authorized successfully!! Yay ꉂ(˵˃ ᗜ ˂˵)"})
}

func (a *Auth) GetToken(c fiber.Ctx) error {
	token, err := a.ValidToken()
	if err != nil {
		return err
	}
	return c.JSON(token)
}

func (a *Auth) ValidToken() (string, error) {
	if a.token != "" && time.Now().Before(a.tokenExpireTime) {
		return a.token, nil
	}

	token, err := a.ObtainNewToken()
	if err != nil {
		return "", err
	}

	a.token = token.AccessToken
	a.tokenExpireTime = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	return a.token, nil
}

func (a *Auth) ObtainNewToken() (TokenResponse, error) {
	body := url.Values{}
	body.Set("grant_type", "authorization_code")
	body.Set("code", a.userCode)
	body.Set("redirect_uri", "http://127.0.0.1:"+a.port+"/api/callback")
	body.Set("client_id", a.clientID)
	body.Set("client_secret", a.clientSecret)

	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}

	resp, err := a.client.MakeRequest(httpclient.POST, spotifyBaseURL+"api/token", body.Encode(), headers)
	if err != nil {
		return TokenResponse{}, err
	}

	var token TokenResponse
	err = json.Unmarshal([]byte(resp), &token)
	if err != nil {
		return TokenResponse{}, err
	}
	return token, nil
}
