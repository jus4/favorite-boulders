package helpers

import (
  "fmt"
	"github.com/jackc/pgx/v5/pgxpool"
  "os"
  "context"
  "strconv"
  "time"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetTimeStamp()string {
  version := strconv.FormatInt(time.Now().Unix(), 10)
  return version
}

func Text(v *string) pgtype.Text {
  if v == nil {
    return pgtype.Text{String: "", Valid: false }
  }
  return pgtype.Text{String: *v, Valid: true}
}

func GetDbPool() *pgxpool.Pool {
  dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
  if err != nil {
    fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
    os.Exit(1)
  }

  return dbpool
}
