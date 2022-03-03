package gorm

import (
	"github.com/porter-dev/porter/internal/models"
	"github.com/porter-dev/porter/internal/repository"
	"gorm.io/gorm"
)

// APITokenRepository uses gorm.DB for querying the database
type APITokenRepository struct {
	db *gorm.DB
}

// NewAPITokenRepository returns a APITokenRepository which uses
// gorm.DB for querying the database
func NewAPITokenRepository(db *gorm.DB) repository.APITokenRepository {
	return &APITokenRepository{db}
}

func (repo *APITokenRepository) CreateAPIToken(a *models.APIToken) (*models.APIToken, error) {
	if err := repo.db.Create(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

func (repo *APITokenRepository) ListAPITokensByProjectID(projectID uint) ([]*models.APIToken, error) {
	tokens := []*models.APIToken{}

	if err := repo.db.Where("project_id = ?", projectID).Find(&tokens).Error; err != nil {
		return nil, err
	}

	return tokens, nil
}

func (repo *APITokenRepository) ReadAPIToken(uid string) (*models.APIToken, error) {
	token := &models.APIToken{}

	if err := repo.db.Where("id = ?", uid).First(&token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

func (repo *APITokenRepository) UpdateAPIToken(
	token *models.APIToken,
) (*models.APIToken, error) {
	if err := repo.db.Save(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}
