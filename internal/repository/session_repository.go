package repository

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"panel.go/internal/ent"
	"panel.go/internal/ent/session"
	"panel.go/internal/interfaces/repository"
	"panel.go/internal/interfaces/resource"
	"panel.go/internal/resource/session_resource"
	_uuid "panel.go/shared/uuid"
)

type SessionRepository struct {
	repository.BaseRepository
	resource resource.ResourceInterface[ent.Session, resource.Response]
}

type SessionCreatePayload struct {
	UserID         uuid.UUID  `json:"user_id"`
	IPAddress      string     `json:"ip_address"`
	UserAgent      string     `json:"user_agent"`
	ImpersonatorID *uuid.UUID `json:"impersonator_id"`
}

type SessionUpdatePayload struct {
	ID uuid.UUID `json:"id"`
	SessionCreatePayload
}

type SessionRespositoryInterface interface {
	repository.BaseRepositoryInterface[SessionCreatePayload, SessionUpdatePayload, resource.Response]
	FindByToken(ctx context.Context, token string) (resource.Response, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (resource.Response, error)
	Touch(ctx context.Context, id uuid.UUID) (resource.Response, error)
}

func NewSessionRepository(ent *ent.Client) SessionRespositoryInterface {
	return &SessionRepository{BaseRepository: repository.BaseRepository{Ent: ent}, resource: session_resource.NewResource()}
}

func (r *SessionRepository) FindAll(ctx context.Context) ([]resource.Response, error) {
	sessions, err := r.Ent.Session.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	return r.resource.Collection(sessions), nil
}

func (r *SessionRepository) FindByToken(ctx context.Context, token string) (resource.Response, error) {
	session, err := r.Ent.Session.Query().Where(session.Token(token)).First(ctx)
	if err != nil {
		return nil, err
	}
	return r.resource.Resource(session), nil
}

func (r *SessionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (resource.Response, error) {
	session, err := r.Ent.Session.Query().Where(session.UserID(userID)).First(ctx)
	if err != nil {
		return nil, err
	}
	return r.resource.Resource(session), nil
}

func (r *SessionRepository) FindOne(ctx context.Context, id uuid.UUID) (resource.Response, error) {

	exists, err := r.Exists(ctx, id)
	if err != nil || !exists {
		return nil, err
	}

	session, _ := r.Ent.Session.Get(ctx, id)

	return r.resource.Resource(session), nil
}

func (r *SessionRepository) Create(ctx context.Context, payload SessionCreatePayload) (resource.Response, error) {
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

	token := _uuid.NewUUID()

	session, err := tx.Session.Create().
		SetUserID(payload.UserID).
		SetIPAddress(payload.IPAddress).
		SetUserAgent(payload.UserAgent).
		SetToken(token.String()).
		SetExpiresAt(time.Now().Add(time.Hour * 2)).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return response
	return r.resource.Resource(session), nil
}

func (r *SessionRepository) Touch(ctx context.Context, id uuid.UUID) (resource.Response, error) {
	exists, err := r.Exists(ctx, id)
	if err != nil || !exists {
		return nil, err
	}

	tx, err := r.Ent.Tx(ctx)

	defer func() {
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Failed to rollback transaction after panic: %v", rollbackErr)
			}
			panic(p)
		}
	}()

	session, err := tx.Session.UpdateOneID(id).Save(ctx)

	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.resource.Resource(session), nil
}

func (r *SessionRepository) Update(ctx context.Context, payload SessionUpdatePayload) (resource.Response, error) {
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

	update := tx.Session.UpdateOneID(payload.ID)
	update.SetNillableImpersonatorID(payload.ImpersonatorID).
		SetIPAddress(payload.IPAddress).
		SetUserAgent(payload.UserAgent).
		SetExpiresAt(time.Now().Add(time.Hour * 2))
	session, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return response
	return r.resource.Resource(session), nil

}

func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
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

	err = tx.Session.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

func (r *SessionRepository) Count(ctx context.Context) (int, error) {
	return r.Ent.Session.Query().Count(ctx)
}

func (r *SessionRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.Ent.Session.Query().Where(session.ID(id)).Exist(ctx)
}
