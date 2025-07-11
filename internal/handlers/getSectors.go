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
)

func GetSectors(c echo.Context) error {
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
  
  return c.JSON(http.StatusOK, sectors)
}
