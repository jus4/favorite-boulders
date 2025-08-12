package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/jus4/favorite-boulders/internal/store/queries"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
  "time"
)

var c = cache.New(5*time.Minute, 10*time.Minute)

func GetSectors(ctx echo.Context) error {

	if cachedSectors, found := c.Get("sectors"); found {
		// Return cached data immediately
		return ctx.JSON(http.StatusOK, cachedSectors)
	}

  err := godotenv.Load()
  if err != nil {
    log.Print("Error loading .env file")
  }
  conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
  if err != nil {
    fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
    os.Exit(1)
  }
  defer conn.Close(context.Background())

  q := queries.New(conn)
  sectors, err := q.GetSectors(context.Background())
  if err != nil {
    return nil
  }

	// Cache the sectors data
	c.Set("sectors", sectors, cache.DefaultExpiration)
  
  return ctx.JSON(http.StatusOK, sectors)
}
