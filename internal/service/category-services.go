package service

import (
	"time"

	"github.com/aborgas90/expense-tracker-api/internal/dto/category"
	"github.com/aborgas90/expense-tracker-api/internal/model"
	"github.com/aborgas90/expense-tracker-api/internal/repo"
)

type CategoryService struct {
	repo *repo.CategoryRepo
}

func NewCategoryService(r *repo.CategoryRepo) *CategoryService {
	return &CategoryService{repo: r}
}

func (s *CategoryService) GetCategoriesByUserID(userID uint) ([]category.CategoryResponse, error) {
	categories, err := s.repo.GetCategoriesByUserID(userID)
	if err != nil {
		return nil, err
	}

	// mapping
	var res []category.CategoryResponse
	for _, c := range categories {
		res = append(res, category.CategoryResponse{
			ID:        c.ID,
			UserID:    c.UserID,
			Type:      c.Type,
			CreatedAt: c.CreatedAt.Format(time.RFC3339),
		})
	}

	return res, nil
}

func (s *CategoryService) CreateCategory(userID uint, req *category.CategoryRequest) (*category.CategoryResponse, error) {
	newCategory := &model.Category{
		UserID: userID,
		Type:   req.TypeCategory,
	}

	if err := s.repo.CreateCategory(newCategory); err != nil {
		return nil, err
	}

	return &category.CategoryResponse{
		ID:        newCategory.ID,
		UserID:    newCategory.UserID,
		Type:      newCategory.Type,
		CreatedAt: newCategory.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *CategoryService) UpdateCategory(categoryID string, userID uint, req *category.CategoryRequest) (*category.CategoryResponse, error) {
	existing, err := s.repo.FindByIDAndUser(categoryID, userID)
	if err != nil {
		return nil, err
	}

	existing.Type = req.TypeCategory

	if err := s.repo.UpdateCategory(userID, existing.ID, existing.Type); err != nil {
		return nil, err
	}

	return &category.CategoryResponse{
		ID:        existing.ID,
		UserID:    existing.UserID,
		Type:      existing.Type,
		CreatedAt: existing.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *CategoryService) DeleteCategory(categoryID string, userId uint) error {
	return s.repo.DeleteCategory(categoryID, userId)
}
