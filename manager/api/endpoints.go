package api

import (
	"github.com/labstack/echo/v4"
	"net/http"

	"phoenixManager/nats"
)

func jsonError(c echo.Context, status int, error string) error {
	return c.JSON(status, map[string]string{
		"success": "false",
		"error":   error,
	})
}

func home(c echo.Context) error {
	return c.HTML(http.StatusOK, "<h1>web endpoint</h1>")
}

func verifyPost(c echo.Context) error {
	req := &nats.DiscordVerifyRequest{}
	if err := c.Bind(req); err != nil {
		return jsonError(c, http.StatusBadRequest, err.Error())
	}

	if req.RSICode == "" {
		return jsonError(c, http.StatusBadRequest, "code is required")
	}

	if req.DiscordName == "" && req.DiscordID == "" {
		return jsonError(c, http.StatusBadRequest, "discord_name or discord_id are required")
	}

	user, err := nats.Gateway.VerifyUser(req)
	if err != nil {
		return jsonError(c, http.StatusInternalServerError, err.Error())
	}

	if user.Error != "" {
		return jsonError(c, http.StatusOK, user.Error)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"success":    "true",
		"discord_id": user.DiscordID,
	})
}
