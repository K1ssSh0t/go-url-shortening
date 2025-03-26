package service

import (
	"context"

	"url-shortener/internal/models"
	"url-shortener/internal/repository"

	"github.com/rs/zerolog/log"
)

type URLService struct {
	repo *repository.URLRepository
}

func NewURLService(repo *repository.URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) CreateShortURL(ctx context.Context, originalURL string) (*models.URL, error) {
	// Validate URL
	url := &models.URL{
		OriginalURL: originalURL,
	}
	
	if err := url.Validate(); err != nil {
		return nil, err
	}
	
	// Attempt to create short URL
	err := s.repo.Create(ctx, url)
	if err != nil {
		return nil, err
	}
	
	return url, nil
}

func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (*models.URL, error) {
	url, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	
	if url == nil {
		return nil, nil
	}
	
	// Increment access count
	err = s.repo.IncrementAccessCount(ctx, shortCode)
	if err != nil {
		// Log error but don't stop the process
		log.Error().Err(err).Msg("Failed to increment access count")
	}
	
	return url, nil
}

func (s *URLService) UpdateShortURL(ctx context.Context, shortCode string, newURL string) (*models.URL, error) {
	// Validate new URL
	url := &models.URL{
		OriginalURL: newURL,
	}
	
	if err := url.Validate(); err != nil {
		return nil, err
	}
	
	// Update URL
	updatedURL, err := s.repo.Update(ctx, shortCode, newURL)
	if err != nil {
		return nil, err
	}
	
	return updatedURL, nil
}

func (s *URLService) DeleteShortURL(ctx context.Context, shortCode string) error {
	return s.repo.Delete(ctx, shortCode)
}

func (s *URLService) GetURLStats(ctx context.Context, shortCode string) (*models.URL, error) {
	return s.repo.FindByShortCode(ctx, shortCode)
}