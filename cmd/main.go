package main

import (
	"github.com/labstack/echo/v4"
  "github.com/jus4/favorite-boulders/internal/handlers"
  "github.com/labstack/echo-contrib/session"
  "github.com/gorilla/sessions"
  "github.com/jus4/favorite-boulders/internal/middleware"
  "github.com/joho/godotenv"
  "log"
  "os"
)

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  e := echo.New()
  e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))))
  e.Use(middleware.AuthMiddleware)
  e.GET("/", handlers.HomePage)
  e.GET("/auth/google/login", handlers.OauthGoogleLogin) 
  e.GET("/auth/google/callback", handlers.OauthGoogleCallback) 
  e.GET("/auth/logout", handlers.Logout)
  e.POST("/api/get-routes/", handlers.GetRoutesByName)
  e.Static("/static", "static") 
	e.Logger.Fatal(e.Start(":8000"))
}
