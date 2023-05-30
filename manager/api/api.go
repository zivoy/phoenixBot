package api

// todo switch to echo
import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"phoenixManager/nats"
)

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
	req := &nats.DiscordVerifyRequest{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if req.RSICode == "" {
		rest.Error(w, "code is required", http.StatusBadRequest)
		return
	}
	if req.DiscordName == "" && req.DiscordID == "" {
		rest.Error(w, "discord_name or discord_name are required", http.StatusBadRequest)
		return
	}

	user, err := nats.Gateway.VerifyUser(req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Error != "" {
		_ = w.WriteJson(map[string]string{
			"success": "false",
			"error":   err.Error(),
		})
		return
	}

	_ = w.WriteJson(map[string]string{
		"success":    "true",
		"discord_id": user.DiscordID,
	})
}
