package middleware

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	secret := os.Getenv("JWT_SECRET")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			var tokenString string

			// Try Authorization header
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader != "" {
				tokenString = authHeader[len("Bearer "):]
			} else {

				// Try cookie
				cookie, err := c.Cookie("access_token")
				if err != nil {
					return c.JSON(http.StatusUnauthorized, echo.Map{
						"error": "missing token",
					})
				}

				tokenString = cookie.Value
			}

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid claims"})
			}

			c.Set("user_id", claims["user_id"])

			return next(c)
		}
	}
}