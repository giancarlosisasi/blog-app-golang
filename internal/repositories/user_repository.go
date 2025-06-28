package repositories

import (
	database "blog-app/internal/database/queries"
	"blog-app/internal/models"
	"blog-app/internal/utils"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserPostgresRepository struct {
	ctx    context.Context
	dbConn *pgx.Conn
}

func NewUserPostgresRepository(ctx context.Context, dbConn *pgx.Conn) *UserPostgresRepository {
	return &UserPostgresRepository{
		ctx:    ctx,
		dbConn: dbConn,
	}
}

func (r *UserPostgresRepository) RegisterUser(email string, hashedPassword string, username string) (*models.UserCreated, error) {
	u, err := database.New(r.dbConn).CreateUser(r.ctx, database.CreateUserParams{
		Email:        email,
		PasswordHash: hashedPassword,
		Username:     username,
	})

	if err != nil {
		return nil, err
	}

	return &models.UserCreated{
		BaseUser: models.BaseUser{
			ID:       u.ID,
			Email:    u.Email,
			Username: u.Username,
		},
	}, nil
}

func (r *UserPostgresRepository) CheckIfEmailExists(email string) (bool, error) {
	_, err := r.GetUserByEmail(email)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (r *UserPostgresRepository) GetUserByEmail(email string) (*models.BaseUser, error) {
	user, err := database.New(r.dbConn).GetUserByEmail(r.ctx, email)
	if err != nil {
		return nil, err
	}

	return &models.BaseUser{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

func (r *UserPostgresRepository) GetUserByID(id string) (*models.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, utils.ErrInvalidUUID
	}

	uuid := pgtype.UUID{
		Valid: true,
		Bytes: uid,
	}
	user, err := database.New(r.dbConn).GetUserByID(r.ctx, uuid)
	if err != nil {
		return nil, utils.ErrResourceNotFoundInDB
	}

	return &models.User{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		FirstName: &user.FirstName.String,
		LastName:  &user.LastName.String,
		AvatarUrl: &user.AvatarUrl.String,
		Bio:       &user.Bio.String,
		IsActive:  user.IsActive.Bool,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

func (r *UserPostgresRepository) GetUserByEmailWithPassword(email string) (*models.UserWithPasswordHash, error) {
	user, err := database.New(r.dbConn).GetUserByEmailWithPassword(r.ctx, email)
	if err != nil {
		return nil, err
	}

	return &models.UserWithPasswordHash{
		PasswordHash: user.PasswordHash,
		BaseUser: models.BaseUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (r UserPostgresRepository) UpdateUserProfileById(id string, updateProfile models.UserProfile) (*models.UserProfile, error) {
	uuid, err := utils.FromStrToPGUUID(id)
	if err != nil {
		return nil, err
	}

	profile := database.UpdateUserParams{
		ID: uuid,
	}

	if updateProfile.AvatarURL != nil && *updateProfile.AvatarURL != "" {
		profile.AvatarUrl = utils.FromStrToPGText(*updateProfile.AvatarURL)
	}

	if updateProfile.FirstName != nil && *updateProfile.FirstName != "" {
		profile.FirstName = utils.FromStrToPGText(*updateProfile.FirstName)
	}

	if updateProfile.LastName != nil && *updateProfile.LastName != "" {
		profile.LastName = utils.FromStrToPGText(*updateProfile.LastName)
	}

	if updateProfile.Bio != nil && *updateProfile.Bio != "" {
		profile.Bio = utils.FromStrToPGText(*updateProfile.Bio)
	}

	user, err := database.New(r.dbConn).UpdateUser(r.ctx, profile)

	return &models.UserProfile{
		ID:        id,
		FirstName: &user.FirstName.String,
		LastName:  &user.LastName.String,
		Bio:       &user.Bio.String,
		AvatarURL: &user.AvatarUrl.String,
	}, nil
}
