package internal

import (
	"fmt"
	"github.com/madrabit/mini-market/search/internal/common"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	FullTextSearching(req SearchRequest) (SearchResponse, error)
	DropDownHint(query string, limit int64) (SuggestResponse, error)
}

type Validator interface {
	Validate(request any) error
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{repo, validator}
}

func (s *Service) FullTextSearching(req SearchRequest) (SearchResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return SearchResponse{}, &common.RequestValidationError{Message: err.Error()}
	}
	resp, err := s.repo.FullTextSearching(req)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("search service: failed to search full text query: %w", err)
	}
	return resp, nil
}

func (s *Service) DropDownHint(query string, limit int64) (SuggestResponse, error) {
	hint, err := s.repo.DropDownHint(query, limit)
	if err != nil {
		return SuggestResponse{}, fmt.Errorf("search service: failed to get drop down hint: %w", err)
	}
	return hint, nil
}
