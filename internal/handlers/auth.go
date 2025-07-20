package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/jus4/favorite-boulders/internal/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func getAuthConfig() *oauth2.Config {
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	return googleOauthConfig
}

func OauthGoogleLogin(c echo.Context) error {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	googleOauthConfig := getAuthConfig()

	oauthState := generateStateOauthCookie(c)
	u := googleOauthConfig.AuthCodeURL(oauthState)

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

func OauthGoogleCallback(c echo.Context) error {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}
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
	var userData map[string]interface{}
	json.Unmarshal(data, &userData)
	email := userData["email"].(string)

	// get create user
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	userRepository := model.NewUserRepository(dbpool)
	user := model.User{}
	u, err := userRepository.GetUserByEmail(context.Background(), email)
	if u != nil {
		user.Email = u.Email
		user.ID = u.ID
	} else {
		newUser, _ := userRepository.CreateUser(context.Background(), email)
		user.Email = newUser.Email
		user.ID = newUser.ID
	}

	sess, _ := session.Get("session", c)
	sess.Values["user_email"] = user.Email
	sess.Values["user_uuid"] = user.ID.String()
	sess.Save(c.Request(), c.Response())
	c.Redirect(http.StatusTemporaryRedirect, "/")

	return nil
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	googleOauthConfig := getAuthConfig()
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

func Logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Values = map[interface{}]interface{}{}
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, "/")
}
