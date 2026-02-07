package auth

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/middleware"
	"github.com/ferdiunal/panel.go/pkg/service/auth"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service        *auth.Service
	accountLockout *middleware.AccountLockout
	environment    string
}

func NewHandler(service *auth.Service, accountLockout *middleware.AccountLockout, environment string) *Handler {
	return &Handler{
		service:        service,
		accountLockout: accountLockout,
		environment:    environment,
	}
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

	// SECURITY: Check if account is locked
	if h.accountLockout != nil && h.accountLockout.IsLocked(req.Email) {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Account temporarily locked due to too many failed login attempts. Please try again later.",
		})
	}

	// Get IP with fallback to X-Forwarded-For
	ip := c.IP()
	if forwarded := c.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}

	session, err := h.service.LoginEmail(c.Context(), req.Email, req.Password, ip, c.Get("User-Agent"))
	if err != nil {
		// SECURITY: Record failed login attempt
		if h.accountLockout != nil {
			h.accountLockout.RecordFailedAttempt(req.Email)
			remaining := h.accountLockout.GetRemainingAttempts(req.Email)
			if remaining > 0 {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":              err.Error(),
					"remaining_attempts": remaining,
				})
			}
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// SECURITY: Reset failed attempts on successful login
	if h.accountLockout != nil {
		h.accountLockout.ResetAttempts(req.Email)
	}

	// SECURITY: Set secure session cookie
	// Use __Host- prefix for additional security (requires Secure=true, Path=/, no Domain)
	// In test environment, use regular cookie name to allow HTTP
	cookieName := "__Host-session_token"
	secure := true
	if h.environment == "test" {
		cookieName = "session_token"
		secure = false
	}

	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Strict", // Strict for admin panels to prevent CSRF
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"session": fiber.Map{
			"token":   session.Token,
			"expires": session.ExpiresAt,
		},
		"user": session.User,
	})
}

func (h *Handler) SignOut(c *context.Context) error {
	cookieName := "__Host-session_token"
	if h.environment == "test" {
		cookieName = "session_token"
	}

	token := c.Cookies(cookieName)
	if token != "" {
		h.service.Logout(c.Context(), token)
	}

	c.ClearCookie(cookieName)
	return c.JSON(fiber.Map{"message": "Signed out"})
}

func (h *Handler) GetSession(c *context.Context) error {
	cookieName := "__Host-session_token"
	if h.environment == "test" {
		cookieName = "session_token"
	}

	token := c.Cookies(cookieName)
	if token == "" {
		return c.JSON(fiber.Map{"session": nil})
	}

	session, err := h.service.ValidateSession(c.Context(), token)
	if err != nil {
		c.ClearCookie(cookieName)
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
	cookieName := "__Host-session_token"
	if h.environment == "test" {
		cookieName = "session_token"
	}

	token := c.Cookies(cookieName)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	session, err := h.service.ValidateSession(c.Context(), token)
	if err != nil {
		c.ClearCookie(cookieName)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	c.Locals("session", session)
	c.Locals("user", session.User)

	return c.Next()
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func (h *Handler) ForgotPassword(c *context.Context) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.ForgotPassword(c.Context(), req.Email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Always return success for security (don't reveal if email exists)
	return c.JSON(fiber.Map{
		"message": "If an account exists with this email, a password reset link has been sent.",
	})
}
