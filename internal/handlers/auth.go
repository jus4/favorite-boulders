
package handlers

import (
  "net/http"
  "time"
	"crypto/rand"
  "github.com/joho/godotenv"
	"encoding/base64"
  "os"
  "encoding/json"
  "context"
  "io"
  "log"
  "fmt"
  "github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
  "github.com/labstack/echo-contrib/session"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func OauthGoogleLogin(c echo.Context) error {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  var googleOauthConfig = &oauth2.Config{
  	RedirectURL:  "http://localhost:8000/auth/google/callback",
  	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
  	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
  	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
  	Endpoint:     google.Endpoint,
  }

  oauthState := generateStateOauthCookie(c)
  u := googleOauthConfig.AuthCodeURL(oauthState)
  // fmt.Printf("Visit the URL for the auth dialog: %v", u)
  return c.Redirect(http.StatusTemporaryRedirect, u) 
}


func generateStateOauthCookie(ctx echo.Context) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	ctx.SetCookie(&cookie)

	return state
}

func OauthGoogleCallback(c echo.Context) error  {
	// Read oauthState from Cookie
	oauthState, _ := c.Cookie("oauthstate")

	if c.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return nil
	}

	data, err := getUserDataFromGoogle(c.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return nil
  }
  var user map[string]interface{}
  json.Unmarshal(data, &user)
  email := user["email"].(string)

  sess, _ := session.Get("session", c)
  sess.Values["user_id"] = email
  sess.Save(c.Request(), c.Response())
	c.Redirect(http.StatusTemporaryRedirect, "/")

  return nil
}


func getUserDataFromGoogle(code string) ([]byte, error) {
  var googleOauthConfig = &oauth2.Config{
  	RedirectURL:  "http://localhost:8000/auth/google/callback",
  	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
  	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
  	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
  	Endpoint:     google.Endpoint,
  }
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}


func Logout(c echo.Context) error{
    sess, _ := session.Get("session", c)
    sess.Values = map[interface{}]interface{}{}
    sess.Options.MaxAge = -1
    sess.Save(c.Request(), c.Response())

    return c.Redirect(http.StatusFound, "/")
}
