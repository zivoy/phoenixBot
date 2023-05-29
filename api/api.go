package api

// idealy this should be its own program, bot can should only handle interactions with discord

import (
	"fmt"
	"log"
	"net/http"
	"phoenixDiscordBot/discord"
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
)

type discordVerify struct {
	DiscordName string `json:"discord_name,omitempty"`
	DiscordID   string `json:"discord_id,omitempty"`
	RSICode     string `json:"code"`
}

func StartApi() {
	api := rest.NewApi()
	api.Use(
		&rest.RecoverMiddleware{
			EnableResponseStackTrace: true,
		},
		&rest.ContentTypeCheckerMiddleware{},
	)
	router, err := rest.MakeRouter(
		rest.Post("/verify", verifyPost),
	)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	go func() {
		fmt.Print("Api is now running")
		log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
	}()
}

func verifyPost(w rest.ResponseWriter, r *rest.Request) {
	req := &discordVerify{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if req.RSICode == "" {
		rest.Error(w, "code is required", 400)
		return
	}
	if req.DiscordName == "" && req.DiscordID == "" {
		rest.Error(w, "discord_name or discord_name are required", 400)
		return
	}

	if req.DiscordID == "" {
		parts := strings.SplitN(req.DiscordName, "#", 2)
		if len(parts) < 2 {
			rest.Error(w, "discriminator not provided", 400)
			return
		}
		fmt.Println(parts)

		req.DiscordID, err = discord.FindUser(parts[0], parts[1])
		if err != nil {
			rest.Error(w, "cannot find user", 404)
			return
		}
	}

	if err = discord.VerifyUser(req.DiscordID, req.RSICode); err != nil {
		_ = w.WriteJson(map[string]string{
			"success": "false",
			"error":   err.Error(),
		})
		return
	}

	_ = w.WriteJson(map[string]string{
		"success":    "true",
		"discord_id": req.DiscordID,
	})
}
