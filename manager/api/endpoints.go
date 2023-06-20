package api

import (
	"fmt"
	"github.com/labstack/gommon/bytes"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/bytes"

	"phoenixManager/files"
	"phoenixManager/nats"
)

const maxImageSize = 10 * bytes.MiB

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

func imageInstructions(c echo.Context) error {
	return c.String(http.StatusOK, "Get image\nOptional get params\n type [webp|png|jpg]\n width\n height")
}

func getImage(c echo.Context) error {
	id := c.Param("id")
	if len(id) < 32 {
		return c.String(http.StatusBadRequest, "not a valid id") // too short
	}
	imType := strings.ToLower(c.QueryParam("type"))
	imWidth, _ := strconv.Atoi(c.QueryParam("width")) // sets to 0 if invalid
	imHeight, _ := strconv.Atoi(c.QueryParam("height"))

	if imType == "" {
		imType = "webp"
	}
	if !files.ValidImage.MatchString(imType) {
		return c.String(http.StatusBadRequest, fmt.Sprintf("invalid image type '%s'", imType))
	}
	if imWidth < 0 {
		imWidth = 0
	}
	if imHeight < 0 {
		imHeight = 0
	}

	im := files.GetImage(id, imWidth, imHeight, imType)
	if im == nil {
		return c.String(http.StatusInternalServerError, "error getting image")
	}
	return c.File(im.Path)
}

func getImageList(c echo.Context) error {
	return c.JSON(http.StatusOK, files.GetImages())
}

func imageUploadInstructions(c echo.Context) error {
	return c.String(http.StatusOK, "Upload image\nform params\n image")
}

func imageDeleteInstructions(c echo.Context) error {
	return c.String(http.StatusOK, "Delete image\nform params\n imageid")
}

func imageUpload(c echo.Context) error {
	img, err := c.FormFile("image")
	if err != nil {
		return jsonError(c, http.StatusBadRequest, "image not provided")
	}
	if img.Size > maxImageSize {
		return jsonError(c, http.StatusRequestEntityTooLarge, fmt.Sprintf("%s is bigger then %s", bytes.Format(img.Size), bytes.Format(maxImageSize)))
	}

	f, err := img.Open()
	if err != nil {
		log.Println("issue reading file ", err)
		return jsonError(c, http.StatusBadRequest, "problem reading file")
	}

	image, err := files.AddImage(f, img.Filename)
	if err != nil {
		log.Println("error saving image ", err)
		return jsonError(c, http.StatusInternalServerError, "error saving image")
	}
	return c.JSON(http.StatusOK, image)
}

func imageDelete(c echo.Context) error {
	img := c.FormValue("imageid")
	if img == "" {
		return jsonError(c, http.StatusBadRequest, "'imageid' not provided")
	}

	im := files.GetImage(img, 0, 0, "webp")
	if im == nil {
		return jsonError(c, http.StatusNotFound, "image does not exist")
	}

	err := files.RemoveImage(img)
	if err != nil {
		return jsonError(c, http.StatusInternalServerError, "could not remove image "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "true"})
}

func imageAuthPerm(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		co, err := c.Cookie("PHPSESSID")
		if err != nil {
			return jsonError(c, http.StatusUnauthorized, "no session id provided")
		}
		if co.Value != "" { //todo do something
			return next(c)
		}
		return jsonError(c, http.StatusUnauthorized, "insufficient perms")
	}
}
