package middleware

import (
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo-contrib/session"
  "context"
  "net/http"
  "fmt"
)

type contextKey string

const emailKey contextKey = "email"

const UserIDKey contextKey = "user_id"

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
       sess, err := session.Get("session", c)
       if err != nil {
        return next(c)
       }

       if userID, ok := sess.Values["user_id"].(string); ok {
          fmt.Println("User ID from session:", userID)
          c.Set("user_id", userID)
          ctx := context.WithValue(c.Request().Context(), UserIDKey, userID)
          c.SetRequest(c.Request().WithContext(ctx))
       } else {
           fmt.Println("No user ID in session")
       }

       return next(c)
    }
}

type UserIDResult struct {
    ID string
    OK bool
}

func GetUserID(ctx context.Context) UserIDResult {
    val := ctx.Value(UserIDKey)
    userID, ok := val.(string)
    if (userID == "") {
      return UserIDResult{ID: "", OK: false}
    }
    return UserIDResult{ID: userID, OK: ok}
}

