package orm

import (
	"context"

	"github.com/ferdiunal/panel.go/pkg/domain/verification"
	"gorm.io/gorm"
)

type VerificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

func (r *VerificationRepository) Create(ctx context.Context, v *verification.Verification) error {
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *VerificationRepository) FindByToken(ctx context.Context, token string) (*verification.Verification, error) {
	var v verification.Verification
	if err := r.db.WithContext(ctx).First(&v, "token = ?", token).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *VerificationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&verification.Verification{}, "id = ?", id).Error
}

func (r *VerificationRepository) DeleteByIdentifier(ctx context.Context, identifier string) error {
	return r.db.WithContext(ctx).Delete(&verification.Verification{}, "identifier = ?", identifier).Error
}
