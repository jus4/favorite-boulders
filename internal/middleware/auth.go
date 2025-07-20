package middleware

import (
	"context"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type contextKey string

const emailKey contextKey = "email"

const UserIDKey contextKey = "user_email"

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return next(c)
		}

		if userID, ok := sess.Values["user_email"].(string); ok {
			c.Set("user_email", userID)
			ctx := context.WithValue(c.Request().Context(), UserIDKey, userID)
			c.SetRequest(c.Request().WithContext(ctx))
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
	if userID == "" {
		return UserIDResult{ID: "", OK: false}
	}
	return UserIDResult{ID: userID, OK: ok}
}
