package avatar

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"panel.go/internal/interfaces/handler"
)

const avatarUrl = "https://avatars.laravel.cloud"

type AvatarResponse struct {
	Body        []byte
	ContentType string
	Status      int
}

var avatarCache = map[string]AvatarResponse{}

func Get(options *handler.Options) handler.HandlerFunc {
	return func(c *fiber.Ctx) error {
		avatar := c.Params("avatar")
		if avatar == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Avatar parameter is required")
		}

		if avatarResponse, ok := avatarCache[avatar]; ok {
			c.Set("Content-Type", avatarResponse.ContentType)
			c.Set("Content-Length", strconv.Itoa(len(avatarResponse.Body)))
			c.Set("Content-Disposition", "inline")
			c.Set("Cache-Control", "public, max-age=31536000")
			c.Set("Expires", time.Now().Add(365*24*time.Hour).Format(http.TimeFormat))
			c.Set("Last-Modified", time.Now().Format(http.TimeFormat))
			c.Set("ETag", fmt.Sprintf("%x", sha1.Sum(avatarResponse.Body)))
			return c.Send(avatarResponse.Body)
		}

		url := fmt.Sprintf("%s/%s?vibe=ocean", avatarUrl, avatar)
		fmt.Printf("Fetching avatar from: %s\n", url)

		response, err := http.Get(url)
		if err != nil {
			fmt.Printf("HTTP request error: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching avatar")
		}
		defer response.Body.Close()

		fmt.Printf("Response status: %d, Content-Length: %s, Content-Type: %s\n",
			response.StatusCode, response.Header.Get("Content-Length"), response.Header.Get("Content-Type"))

		if response.StatusCode != http.StatusOK {
			return c.Status(response.StatusCode).SendString("Avatar not found")
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error reading avatar")
		}

		if len(body) == 0 {
			return c.Status(fiber.StatusNoContent).SendString("Empty avatar response")
		}

		if contentType := response.Header.Get("Content-Type"); contentType != "" {
			c.Set("Content-Type", contentType)
		}

		avatarCache[avatar] = AvatarResponse{
			Body:        body,
			ContentType: response.Header.Get("Content-Type"),
			Status:      response.StatusCode,
		}

		return c.Send(body)
	}
}
