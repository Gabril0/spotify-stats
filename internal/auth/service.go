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

// SpotifyAuthWebRedirect godoc
// @Summary Get the Spotify login/authorize URL
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/web-auth [get]
func (a *Auth) SpotifyAuthWebRedirect(c fiber.Ctx) error {
	params := url.Values{}
	params.Set("client_id", a.clientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", "http://127.0.0.1:"+a.port+"/api/callback")
	params.Set("scope", "user-top-read user-read-recently-played user-library-read")
	params.Set("state", "spotify island auth request")

	loginURL := spotifyBaseURL + "authorize?" + params.Encode()

	return c.JSON(fiber.Map{"url": loginURL})
}

// SpotifyCallback godoc
// @Summary OAuth redirect target used by Spotify
// @Tags auth
// @Produce json
// @Param code query string false "Authorization code"
// @Param state query string false "State parameter"
// @Success 200 {object} map[string]string
// @Router /api/callback [get]
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

// GetToken godoc
// @Summary Get a valid (cached or refreshed) access token
// @Tags auth
// @Produce json
// @Success 200 {string} string
// @Router /api/token [get]
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
