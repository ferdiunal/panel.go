package orm

import (
	"context"

	"github.com/ferdiunal/panel.go/internal/domain/session"
	"github.com/ferdiunal/panel.go/shared/uuid"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, s *session.Session) error {
	s.ID = uuid.NewUUID().String()
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *SessionRepository) FindByID(ctx context.Context, id string) (*session.Session, error) {
	var s session.Session
	if err := r.db.WithContext(ctx).Preload("User").First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*session.Session, error) {
	var s session.Session
	if err := r.db.WithContext(ctx).Preload("User").First(&s, "token = ?", token).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&session.Session{}, "id = ?", id).Error
}

func (r *SessionRepository) DeleteByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Delete(&session.Session{}, "token = ?", token).Error
}

func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Delete(&session.Session{}, "user_id = ?", userID).Error
}
