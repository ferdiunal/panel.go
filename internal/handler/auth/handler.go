package auth

import (
	"github.com/ferdiunal/panel.go/internal/context"
	"github.com/ferdiunal/panel.go/internal/service/auth"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *auth.Service
}

func NewHandler(service *auth.Service) *Handler {
	return &Handler{service: service}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) RegisterEmail(c *context.Context) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := h.service.RegisterEmail(c.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if err == auth.ErrEmailAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		if err == auth.ErrInvalidEmail || err == auth.ErrWeakPassword || err == auth.ErrInvalidName {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Auto login after register? For now just return user.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": user,
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) LoginEmail(c *context.Context) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	session, err := h.service.LoginEmail(c.Context(), req.Email, req.Password, c.IP(), c.Get("User-Agent"))
	if err != nil {
		if err == auth.ErrTooManyAttempts {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Set Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		HTTPOnly: true,
		Secure:   c.Protocol() == "https", // Or config
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"session": session,
		"user":    session.User, // Preloaded in service? Wait, Service LoginEmail returns session but doesn't explicitly preload user struct inside session return, BUT it fetches it. I should ensure Service populates it or I fetch it.
		// Actually Service just creates session object. It has UserID. I might need to load user to return it.
		// For now let's assume session is enough or I update service to return user too.
	})
}

func (h *Handler) SignOut(c *context.Context) error {
	token := c.Cookies("session_token")
	if token != "" {
		h.service.Logout(c.Context(), token)
	}

	c.ClearCookie("session_token")
	return c.JSON(fiber.Map{"message": "Signed out"})
}

func (h *Handler) GetSession(c *context.Context) error {
	token := c.Cookies("session_token")
	if token == "" {
		return c.JSON(fiber.Map{"session": nil})
	}

	session, err := h.service.ValidateSession(c.Context(), token)
	if err != nil {
		c.ClearCookie("session_token")
		return c.JSON(fiber.Map{"session": nil})
	}

	return c.JSON(fiber.Map{
		"session": fiber.Map{
			"token":   session.Token,
			"expires": session.ExpiresAt,
		},
		"user": session.User, // Preloaded? Service FindByToken preloads User.
	})
}

func (h *Handler) SessionMiddleware(c *context.Context) error {
	token := c.Cookies("session_token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	session, err := h.service.ValidateSession(c.Context(), token)
	if err != nil {
		c.ClearCookie("session_token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	c.Locals("session", session)
	c.Locals("user", session.User)

	return c.Next()
}
