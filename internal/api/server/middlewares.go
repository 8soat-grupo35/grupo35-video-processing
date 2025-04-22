package server

import (
	"fmt"
	"net/http"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/external"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/gateways"
	"github.com/labstack/echo/v4"
)

func AuthenticationMiddleware(cfg external.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header is missing"})
			}

			authGateway := gateways.NewAuthGateway(cfg.CognitoUserPoolClientId, cfg.Region)

			claims, err := authGateway.ValidateToken(authHeader)

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			fmt.Println("claims", claims)

			email, err := gateways.GetCognitoUserEmail(authHeader)

			if err != nil {
				fmt.Println(err)
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			userId, ok := (*claims)["sub"].(string)

			if !ok {
				fmt.Println("cant get user id")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			fmt.Println("user id", userId)
			fmt.Println("user email", *email)

			c.Set("user_id", userId)
			c.Set("user_email", *email)

			return next(c)
		}
	}
}
