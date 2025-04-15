package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
  "github.com/jus4/favorite-boulders/internal/handlers"
)

func main() {
  fmt.Printf("Start app")
  e := echo.New()
  e.GET("/", handlers.HomePage)
  e.POST("/api/get-routes/", handlers.GetRoutesByName)
  e.Static("/static", "static") 
	e.Logger.Fatal(e.Start(":1323"))
}
