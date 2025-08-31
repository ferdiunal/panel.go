package service

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_err "panel.go/internal/errors"
	"panel.go/internal/repository"
	"panel.go/internal/resource/session_resource"
	"panel.go/internal/resource/user_resource"
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

func (s *AuthService) VerifyToken(c *fiber.Ctx, token string) error {
	cookie := c.Cookies("access_token")
	if cookie == "" {
		return _err.ErrTokenExpired
	}

	decryptedCookie, err := s.Encrypt.Decrypt(cookie)
	if err != nil {
		return _err.ErrTokenExpired
	}

	var tokenBundle TokenBundle
	if err := json.Unmarshal([]byte(decryptedCookie), &tokenBundle); err != nil {
		return _err.ErrTokenExpired
	}

	session, err := s.SessionRepository.FindOne(c.Context(), uuid.MustParse(tokenBundle.SessionID))
	if err != nil {
		return _err.ErrTokenExpired
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = s.SessionRepository.Delete(c.Context(), session.ID)
		return _err.ErrTokenExpired
	}

	user, err := s.UserRepository.FindOne(c.Context(), session.UserID)
	if err != nil {
		return _err.ErrTokenExpired
	}

	c.Locals("user", user)

	return nil
}

func (s *AuthService) Logout(c *fiber.Ctx) error {
	cookie := c.Cookies("access_token")
	if cookie == "" {
		return _err.ErrTokenExpired
	}

	decryptedCookie, err := s.Encrypt.Decrypt(cookie)
	if err != nil {
		return _err.ErrTokenExpired
	}

	var tokenBundle TokenBundle
	if err := json.Unmarshal([]byte(decryptedCookie), &tokenBundle); err != nil {
		return _err.ErrTokenExpired
	}

	_ = s.SessionRepository.Delete(c.Context(), uuid.MustParse(tokenBundle.SessionID))

	return nil
}

func (s *AuthService) Login(c *fiber.Ctx, body *LoginPayload) (*TokenResponse, error) {
	user, err := s.UserRepository.FindByEmail(c.Context(), body.Email)
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

	err = bcrypt.CompareHashAndPassword([]byte(*account.Password), []byte(body.Password))

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

func (s *AuthService) Register(c *fiber.Ctx, body *RegisterPayload) (bool, error) {

	exists, err := s.UserRepository.ExistsByEmail(c.Context(), body.Email)
	if err != nil {
		return false, _err.ErrRegister
	}

	if exists {
		return false, _err.ErrUserExists
	}

	user, err := s.UserRepository.Create(c.Context(), repository.UserCreatePayload{
		Name:  body.Name,
		Email: body.Email,
	})

	if err != nil {
		return false, _err.ErrRegister
	}

	password, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, _err.ErrRegister
	}

	provider := "email"
	_password := string(password)

	_, err = s.AccountRepository.Create(c.Context(), repository.AccountCreatePayload{
		UserID:   user.ID,
		Provider: &provider,
		Password: &_password,
	})

	if err != nil {
		return false, _err.ErrRegister
	}

	return true, nil
}

func (s *AuthService) UpdateAccount(c *fiber.Ctx, body *repository.UserCreatePayload) (bool, error) {
	user := c.Locals("user").(*user_resource.UserResource)
	exists, err := s.UserRepository.Exists(c.Context(), user.ID)

	if err != nil || !exists {
		return false, _err.ErrUpdateAccount
	}

	existsEmail, err := s.UserRepository.UniqueEmailWithID(c.Context(), body.Email, user.ID)

	if err != nil || existsEmail {
		return false, _err.ErrUniqueEmail
	}

	_, err = s.UserRepository.Update(c.Context(), repository.UserUpdatePayload{
		ID:    user.ID,
		Name:  &body.Name,
		Email: &body.Email,
	})

	if err != nil {
		return false, _err.ErrUpdateAccount
	}
	return true, nil
}
