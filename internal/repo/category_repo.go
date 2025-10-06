package repo

import (
	"github.com/aborgas90/expense-tracker-api/internal/model"
	"gorm.io/gorm"
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo{
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) GetCategoriesByUserID(userID uint) ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepo) CreateCategory(category *model.Category)error {
	return r.db.Create(category).Error
}

func (r *CategoryRepo) UpdateCategory(categoryID uint, userID uint, typeCategory string) error {
	return r.db.Model(&model.Category{}).
		Where("id = ? AND user_id = ?", categoryID, userID).
		Update("type", typeCategory).Error
}

func (r *CategoryRepo) DeleteCategory(categoryID string, userId uint) error {
	if err := r.db.Where("id = ? AND user_id = ?", categoryID, userId).
		Delete(&model.Category{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepo) FindByIDAndUser(id string, userID uint) (*model.Category, error) {
	var category model.Category
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

