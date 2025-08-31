package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"panel.go/internal/ent"
	"panel.go/internal/ent/user"
	"panel.go/internal/interfaces/repository"
	"panel.go/internal/interfaces/resource"
	"panel.go/internal/resource/user_resource"
)

type UserRepository struct {
	repository.BaseRepository
	resource resource.ResourceInterface[ent.User, *user_resource.UserResource]
}

type UserCreatePayload struct {
	Email string `json:"email" validate:"required,email,max=150,omitempty"`
	Name  string `json:"name" validate:"required,min=3,max=150,omitempty"`
}

type UserUpdatePayload struct {
	ID    uuid.UUID `json:"id" validate:"required,uuid"`
	Email *string   `json:"email,omitempty" validate:"omitempty,email,max=150"`
	Name  *string   `json:"name,omitempty" validate:"omitempty,min=3,max=150"`
}

type UserRespositoryInterface interface {
	repository.BaseRepositoryInterface[UserCreatePayload, UserUpdatePayload, *user_resource.UserResource]
	FindByEmail(ctx context.Context, email string) (*user_resource.UserResource, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	UniqueEmailWithID(ctx context.Context, email string, id uuid.UUID) (bool, error)
}

func NewUserRepository(ent *ent.Client) UserRespositoryInterface {
	return &UserRepository{BaseRepository: repository.BaseRepository{Ent: ent}, resource: user_resource.NewResource()}
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*user_resource.UserResource, error) {
	users, err := r.Ent.User.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return r.resource.Collection(users), nil
}

func (r *UserRepository) UniqueEmailWithID(ctx context.Context, email string, id uuid.UUID) (bool, error) {
	return r.Ent.User.Query().Where(user.Email(email), user.IDNEQ(id)).Exist(ctx)
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.Ent.User.Query().Where(user.Email(email)).Exist(ctx)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user_resource.UserResource, error) {
	exists, err := r.ExistsByEmail(ctx, email)
	if err != nil || !exists {
		return nil, err
	}

	user, _ := r.Ent.User.Query().Where(user.Email(email)).First(ctx)

	return r.resource.Resource(user), nil
}

func (r *UserRepository) FindOne(ctx context.Context, id uuid.UUID) (*user_resource.UserResource, error) {
	exists, err := r.Exists(ctx, id)
	if err != nil || !exists {
		return nil, err
	}

	user, _ := r.Ent.User.Get(ctx, id)
	return r.resource.Resource(user), nil
}

func (r *UserRepository) Create(ctx context.Context, payload UserCreatePayload) (*user_resource.UserResource, error) {
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

	user, err := tx.User.Create().
		SetName(payload.Name).
		SetEmail(payload.Email).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return response
	return r.resource.Resource(user), nil
}

func (r *UserRepository) Update(ctx context.Context, payload UserUpdatePayload) (*user_resource.UserResource, error) {
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

	update := tx.User.UpdateOneID(payload.ID)
	update.SetNillableName(payload.Name)
	update.SetNillableEmail(payload.Email)

	user, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.resource.Resource(user), nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
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

	err = tx.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

func (r *UserRepository) Count(ctx context.Context) (int, error) {
	return r.Ent.User.Query().Count(ctx)
}

func (r *UserRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.Ent.User.Query().Where(user.ID(id)).Exist(ctx)
}
