package service

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	_err "panel.go/internal/errors"
	"panel.go/internal/repository"
	"panel.go/internal/resource/session_resource"
	"panel.go/shared/encrypt"
)

type AuthService struct {
	AccountRepository repository.AccountRespositoryInterface
	UserRepository    repository.UserRespositoryInterface
	SessionRepository repository.SessionRespositoryInterface
	Encrypt           encrypt.Crypt
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email,max=150"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type RegisterPayload struct {
	Name     string `json:"name" validate:"required,min=3,max=150"`
	Email    string `json:"email" validate:"required,email,max=150"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type TokenBundle struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	ExpiresAt string `json:"expires_at"`
}

type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func NewAuthService(
	accountRepository repository.AccountRespositoryInterface,
	userRepository repository.UserRespositoryInterface,
	sessionRepository repository.SessionRespositoryInterface,
	encrypt encrypt.Crypt,
) *AuthService {
	return &AuthService{AccountRepository: accountRepository, UserRepository: userRepository, SessionRepository: sessionRepository, Encrypt: encrypt}
}

func (s *AuthService) Login(c *fiber.Ctx, email string, password string) (*TokenResponse, error) {
	user, err := s.UserRepository.FindByEmail(c.Context(), email)
	if err != nil {
		return nil, _err.ErrAuthentication
	}

	account, err := s.AccountRepository.FindByUserIDWithPassword(c.Context(), user.ID)
	if err != nil {
		return nil, _err.ErrAuthentication
	}

	if account.Password == nil {
		return nil, _err.ErrAuthentication
	}

	err = bcrypt.CompareHashAndPassword([]byte(*account.Password), []byte(password))

	if err != nil {
		return nil, _err.ErrAuthentication
	}

	session, err := s.SessionRepository.Create(c.Context(), repository.SessionCreatePayload{
		UserID:    user.ID,
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
	})

	if err != nil {
		return nil, _err.ErrAuthentication
	}

	accessToken, err := s.encryptedSession(session)
	if err != nil {
		return nil, _err.ErrAuthentication
	}

	return &TokenResponse{
		AccessToken: *accessToken,
		ExpiresAt:   session.ExpiresAt,
	}, nil
}

func (s *AuthService) encryptedSession(session *session_resource.SessionResource) (*string, error) {

	tokenBundle := TokenBundle{
		Token:     session.Token,
		SessionID: session.ID.String(),
		IPAddress: session.IPAddress,
		UserAgent: session.UserAgent,
		ExpiresAt: session.ExpiresAt.UTC().Format(time.RFC3339Nano),
	}

	_json, err := json.Marshal(tokenBundle)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.Encrypt.Encrypt(string(_json))

	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}
