package repository

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"panel.go/internal/ent"
	"panel.go/internal/ent/account"
	"panel.go/internal/interfaces/repository"
	"panel.go/internal/interfaces/resource"
	"panel.go/internal/resource/account_resource"
)

type AccountRepository struct {
	repository.BaseRepository
	resource resource.ResourceInterface[ent.Account, resource.Response]
}

type AccountCreatePayload struct {
	Password              *string    `json:"password,omitempty"`
	Provider              *string    `json:"provider,omitempty"`
	ProviderID            *string    `json:"provider_id,omitempty"`
	AccessToken           *string    `json:"access_token,omitempty"`
	RefreshToken          *string    `json:"refresh_token,omitempty"`
	IdToken               *string    `json:"id_token,omitempty"`
	AccessTokenExpiresAt  *time.Time `json:"access_token_expires_at,omitempty"`
	RefreshTokenExpiresAt *time.Time `json:"refresh_token_expires_at,omitempty"`
	Scopes                []string   `json:"scopes"`
	UserID                uuid.UUID  `json:"user_id"`
}

type AccountUpdatePayload struct {
	ID uuid.UUID `json:"id"`
	AccountCreatePayload
}

type AccountCreatePasswordPayload struct {
	UserID   uuid.UUID `json:"user_id"`
	Password string    `json:"password"`
}

type AccountRespositoryInterface interface {
	repository.BaseRepositoryInterface[AccountCreatePayload, AccountUpdatePayload, resource.Response]
	CreatePassword(ctx context.Context, payload AccountCreatePasswordPayload) (resource.Response, error)
}

func NewAccountRepository(ent *ent.Client) AccountRespositoryInterface {
	return &AccountRepository{BaseRepository: repository.BaseRepository{Ent: ent}, resource: account_resource.NewResource()}
}

func (r *AccountRepository) FindAll(ctx context.Context) ([]resource.Response, error) {
	accounts, err := r.Ent.Account.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	return r.resource.Collection(accounts), nil
}

func (r *AccountRepository) FindOne(ctx context.Context, id uuid.UUID) (resource.Response, error) {

	exists, err := r.Exists(ctx, id)
	if err != nil || !exists {
		return nil, err
	}

	account, _ := r.Ent.Account.Get(ctx, id)

	return r.resource.Resource(account), nil
}

func (r *AccountRepository) Create(ctx context.Context, payload AccountCreatePayload) (resource.Response, error) {
	tx, err := r.Ent.Tx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Failed to rollback transaction after panic: %v", rollbackErr)
			}
			panic(p)
		}
	}()

	account, err := tx.Account.Create().
		SetNillableAccessToken(payload.AccessToken).
		SetNillableRefreshToken(payload.RefreshToken).
		SetNillableIDToken(payload.IdToken).
		SetNillableAccessTokenExpiresAt(payload.AccessTokenExpiresAt).
		SetNillableRefreshTokenExpiresAt(payload.RefreshTokenExpiresAt).
		SetScopes(payload.Scopes).
		SetNillableProvider(payload.Provider).
		SetNillableProviderID(payload.ProviderID).
		SetUserID(payload.UserID).
		SetNillablePassword(payload.Password).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return response
	return r.resource.Resource(account), nil
}

func (r *AccountRepository) Update(ctx context.Context, payload AccountUpdatePayload) (resource.Response, error) {
	exists, err := r.Exists(ctx, payload.ID)
	if err != nil || !exists {
		return nil, err
	}

	tx, err := r.Ent.Tx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Failed to rollback transaction after panic: %v", rollbackErr)
			}
			panic(p)
		}
	}()

	update := tx.Account.UpdateOneID(payload.ID)

	update.SetNillableAccessToken(payload.AccessToken).
		SetNillableRefreshToken(payload.RefreshToken).
		SetNillableIDToken(payload.IdToken).
		SetNillableAccessTokenExpiresAt(payload.AccessTokenExpiresAt).
		SetNillableRefreshTokenExpiresAt(payload.RefreshTokenExpiresAt).
		SetScopes(payload.Scopes).
		SetNillableProvider(payload.Provider).
		SetNillableProviderID(payload.ProviderID).
		SetNillablePassword(payload.Password)

	account, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return response
	return r.resource.Resource(account), nil

}

func (r *AccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	exists, err := r.Exists(ctx, id)
	if err != nil || !exists {
		return err
	}

	tx, err := r.Ent.Tx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Failed to rollback transaction after panic: %v", rollbackErr)
			}
			panic(p)
		}
	}()

	err = tx.Account.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

func (r *AccountRepository) Count(ctx context.Context) (int, error) {
	return r.Ent.Account.Query().Count(ctx)
}

func (r *AccountRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.Ent.Account.Query().Where(account.ID(id)).Exist(ctx)
}

func (r *AccountRepository) CreatePassword(ctx context.Context, payload AccountCreatePasswordPayload) (resource.Response, error) {
	exists, err := r.Exists(ctx, payload.UserID)
	if err != nil || !exists {
		return nil, err
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	password := string(encryptedPassword)

	return r.Create(ctx, AccountCreatePayload{
		UserID:   payload.UserID,
		Password: &password,
	})
}
