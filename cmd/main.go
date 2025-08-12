package main

import (
	"github.com/labstack/echo/v4"
  "github.com/jus4/favorite-boulders/internal/handlers"
  "github.com/labstack/echo-contrib/session"
  "github.com/gorilla/sessions"
  "github.com/jus4/favorite-boulders/internal/middleware"
  "github.com/joho/godotenv"
  "log"
  "net/http"
  "os"
	"time"
	"github.com/patrickmn/go-cache"
  "io"
)

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Print("Error loading .env file")
  }
  e := echo.New()
  e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))))
  e.Use(middleware.AuthMiddleware)
  e.GET("/", handlers.HomePage)
  e.GET("/auth/google/login", handlers.OauthGoogleLogin) 
  e.GET("/auth/google/callback", handlers.OauthGoogleCallback) 
  e.GET("/auth/logout", handlers.Logout)
  e.GET("/proxy/wmts/*", func(c echo.Context) error {
    var tileCache = cache.New(30*time.Minute, 1*time.Hour)
	  tilePath := c.Param("*")
	  cacheKey := tilePath
    mapApiKey := os.Getenv("MAP_API_KEY")

  	// Check cache
  	if cached, found := tileCache.Get(cacheKey); found {
  		cachedData := cached.(map[string]interface{})
  		body := cachedData["body"].([]byte)
  		contentType := cachedData["contentType"].(string)
  
  		return c.Blob(http.StatusOK, contentType, body)
  	}

	  // Construct full upstream URL
	  targetURL := "https://avoin-karttakuva.maanmittauslaitos.fi/avoin/wmts/" + tilePath + "?api-key=" + mapApiKey

	  // Forward the request
	  resp, err := http.Get(targetURL)
	  if err != nil {
	  	return c.String(http.StatusBadGateway, "Failed to fetch tile")
	  }
	  defer resp.Body.Close()

  	// Read full body
  	bodyBytes, err := io.ReadAll(resp.Body)
  	if err != nil {
  		return c.String(http.StatusInternalServerError, "Failed to read tile response")
  	}
  
  	// Cache the response
  	tileCache.Set(cacheKey, map[string]interface{}{
  		"body":        bodyBytes,
  		"contentType": resp.Header.Get("Content-Type"),
  	}, cache.DefaultExpiration)


	  // Serve to client
	  return c.Blob(resp.StatusCode, resp.Header.Get("Content-Type"), bodyBytes)

	  // return c.Stream(resp.StatusCode, resp.Header.Get("Content-Type"), resp.Body)
  })
  e.POST("/api/get-routes/", handlers.GetRoutesByName)
  e.POST("/api/search-route/", handlers.SearchRouteByName)
  e.POST("/api/select-climb/", handlers.SearchEditClimb)
  e.POST("/api/edit-route/", handlers.EditRoute) // TODO fix to :id
  e.POST("/api/update-route/:id", handlers.UpdateRoute)
  e.GET("/add-modify/", handlers.AddModify)
  e.GET("/api/get-sectors/", handlers.GetSectors)
  e.GET("/api/routes-by-sector/:id", handlers.GetRoutesBySectorId)
  e.Static("/static", "static") 


  // Start the app
  port := os.Getenv("PORT")
  if port == "" {
      port = "8080"
  }

  e.Logger.Fatal(e.Start(":" + port))

}
