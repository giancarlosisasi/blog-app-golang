package services

import (
	"blog-app/internal/models"
	"blog-app/internal/repositories"

	"github.com/rs/zerolog/log"
)

type PostService struct {
	postRepository repositories.PostRepository
}

func NewPostService(postRepository repositories.PostRepository) *PostService {
	return &PostService{
		postRepository: postRepository,
	}
}

func (s *PostService) CreatePost(post models.Post) (*models.Post, error) {
	postCreated, err := s.postRepository.CreatePost(post)
	if err != nil {
		return nil, err
	}

	return postCreated, nil
}

func (s *PostService) GetPostBySlug(slug string) (*models.Post, error) {
	post, err := s.postRepository.GetPostBySlug(slug)

	if err != nil {
		log.Error().Err(err).Msg("postRepository.GetPostBySlug")
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetPostByID(id string) (*models.Post, error) {
	post, err := s.postRepository.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetPostsByAuthorID(authorId string) ([]models.Post, error) {
	posts, err := s.postRepository.GetPostsByAuthorID(authorId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostService) UpdatePostByID(authorId string, id string, data *models.UpdatePostData) (*models.Post, error) {
	post, err := s.postRepository.UpdatePostByID(authorId, id, data)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) DeletePostByID(id string) error {
	return nil
}
