package services

import (
	"blog-app/internal/models"
	"blog-app/internal/repositories"
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
	return nil, nil
}

func (s *PostService) GetPostByID(id string) (*models.Post, error) {
	return nil, nil
}

func (s *PostService) GetPostsByAuthorID(authorId string) ([]*models.Post, error) {
	return nil, nil
}

func (s *PostService) UpdatePostByID(id string) (*models.Post, error) {
	return nil, nil
}

func (s *PostService) DeletePostByID(id string) error {
	return nil
}
