package repositories

import (
	database "blog-app/internal/database/queries"
	"blog-app/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type PostPostgresRepository struct {
	ctx    context.Context
	dbConn *pgx.Conn
}

func NewPostPostgresRepository(ctx context.Context, dbConn *pgx.Conn) *PostPostgresRepository {
	return &PostPostgresRepository{
		ctx:    ctx,
		dbConn: dbConn,
	}
}

func (r *PostPostgresRepository) CreatePost(post models.Post) (*models.Post, error) {
	queries := database.New(r.dbConn)

	// Parse the string AuthorID to UUID
	authorUUID, err := uuid.Parse(post.AuthorID)
	if err != nil {
		return nil, err
	}

	// .Format("2006-01-02T15:04:05 -070000")

	postData := database.CreatePostParams{
		Title:   post.Title,
		Slug:    post.Slug,
		Content: post.Content,
		Excerpt: pgtype.Text{
			String: post.Excerpt,
			Valid:  true,
		},
		FeaturedImageUrl: pgtype.Text{
			String: *post.FeatureImageURL,
			Valid:  true,
		},
		Status: pgtype.Text{
			String: post.Status,
			Valid:  true,
		},
		AuthorID: pgtype.UUID{
			Bytes: authorUUID,
			Valid: true,
		},
	}
	if post.Status == "published" {
		postData.PublishedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	postCreated, err := queries.CreatePost(r.ctx, postData)
	log.Info().Msg(fmt.Sprintf("postData is: %+v", postData))
	if err != nil {
		return nil, err
	}

	// Convert the created post back to our model
	var publishedAt *time.Time
	if postCreated.PublishedAt.Valid {
		publishedAt = &postCreated.PublishedAt.Time
	}

	return &models.Post{
		Title:           postCreated.Title,
		Slug:            postCreated.Slug,
		Content:         postCreated.Content,
		Excerpt:         postCreated.Excerpt.String,
		FeatureImageURL: &postCreated.FeaturedImageUrl.String,
		Status:          postCreated.Status.String,
		AuthorID:        postCreated.AuthorID.String(), // Keep the original string
		PublishedAt:     publishedAt,
	}, nil
}

func (r *PostPostgresRepository) GetPostBySlug(slug string) (*models.Post, error) {
	queries := database.New(r.dbConn)
	post, err := queries.GetPostBySlug(r.ctx, slug)
	if err != nil {
		return nil, err
	}

	id := post.ID.String()

	return &models.Post{
		ID:              &id,
		Title:           post.Title,
		Slug:            post.Slug,
		Content:         post.Content,
		Excerpt:         post.Excerpt.String,
		FeatureImageURL: &post.FeaturedImageUrl.String,
		Status:          post.Status.String,
		PublishedAt:     &post.PublishedAt.Time,
		AuthorID:        post.AuthorID.String(),
	}, nil
}

func (r *PostPostgresRepository) GetPostByID(id string) (*models.Post, error) {
	return nil, nil
}

func (r *PostPostgresRepository) GetPostsByAuthorID(authorId string) ([]*models.Post, error) {
	return nil, nil
}

func (r *PostPostgresRepository) UpdatePostByID(id string) (*models.Post, error) {
	return nil, nil
}

func (r *PostPostgresRepository) DeletePostByID(id string) error {
	return nil
}
