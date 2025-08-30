package account_resource

import (
	"time"

	"github.com/google/uuid"
	"panel.go/internal/ent"
	"panel.go/internal/interfaces/resource"
)

type AccountResource struct {
	ID                    uuid.UUID  `json:"id"`
	Password              *string    `json:"password,omitempty"`
	Provider              string     `json:"provider,omitempty"`
	ProviderID            string     `json:"provider_id,omitempty"`
	AccessToken           string     `json:"access_token,omitempty"`
	RefreshToken          string     `json:"refresh_token,omitempty"`
	IdToken               string     `json:"id_token,omitempty"`
	AccessTokenExpiresAt  *time.Time `json:"access_token_expires_at,omitempty"`
	RefreshTokenExpiresAt *time.Time `json:"refresh_token_expires_at,omitempty"`
}

type initResource struct {
}

func NewResource() resource.ResourceInterface[ent.Account, *AccountResource] {
	return &initResource{}
}

func (initResource) Resource(account *ent.Account) *AccountResource {
	return &AccountResource{
		ID:                    account.ID,
		Password:              &account.Password,
		Provider:              account.Provider,
		ProviderID:            account.ProviderID,
		AccessToken:           account.AccessToken,
		RefreshToken:          account.RefreshToken,
		IdToken:               account.IDToken,
		AccessTokenExpiresAt:  &account.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: &account.RefreshTokenExpiresAt,
	}
}

func (r *initResource) Collection(accounts []*ent.Account) []*AccountResource {
	response := []*AccountResource{}
	for _, account := range accounts {
		response = append(response, r.Resource(account))
	}
	return response
}
