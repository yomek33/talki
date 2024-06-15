package stores

import (
	"errors"
	"log"

	"github.com/yomek33/talki/internal/models"
	"gorm.io/gorm"
)

const (
	ErrMaterialCannotBeNil  = "material cannot be nil"
	ErrMismatchedMaterialID = "material ID does not match the provided ID"
)

type MaterialStore interface {
	CreateMaterial(material *models.Material) (uint, error)
	GetMaterialByID(id uint, UserUID string) (*models.Material, error)
	UpdateMaterial(id uint, material *models.Material) error
	DeleteMaterial(id uint, UserUID string) error
	GetAllMaterials(searchQuery string, UserUID string) ([]models.Material, error)
	UpdateMaterialStatus(id uint, status string) error
	GetMaterialStatus(id uint) (string, error)
}

type materialStore struct {
	BaseStore
}

func (s *materialStore) CreateMaterial(material *models.Material) (uint, error) {
	if material == nil {
		return 0, errors.New(ErrMaterialCannotBeNil)
	}
	err := s.PerformDBTransaction(func(tx *gorm.DB) error {
		return tx.Create(material).Error
	})
	if err != nil {
		return 0, err
	}
	return material.ID, nil
}

func (s *materialStore) GetMaterialByID(id uint, UserUID string) (*models.Material, error) {
	log.Println("store material id", id)
	var material models.Material
	err := s.DB.Where("id = ? AND user_uid = ?", id, UserUID).Preload("Phrases").First(&material).Error
	return &material, err
}

func (s *materialStore) UpdateMaterial(id uint, material *models.Material) error {
	if material == nil {
		return errors.New(ErrMaterialCannotBeNil)
	}
	if id != material.ID {
		return errors.New(ErrMismatchedMaterialID)
	}
	return s.PerformDBTransaction(func(tx *gorm.DB) error {
		return tx.Model(&models.Material{}).Where("id = ?", id).Updates(material).Error
	})
}

func (s *materialStore) DeleteMaterial(id uint, UserUID string) error {
	return s.PerformDBTransaction(func(tx *gorm.DB) error {
		return tx.Where("id = ? AND user_uid = ?", id, UserUID).Delete(&models.Material{}).Error
	})
}

func (s *materialStore) GetAllMaterials(searchQuery string, UserUID string) ([]models.Material, error) {
	var materials []models.Material
	query := s.DB.Where("user_uid = ?", UserUID)
	if searchQuery != "" {
		query = query.Where("title LIKE ?", "%"+searchQuery+"%")
	}
	err := query.Find(&materials).Error
	return materials, err
}

func (s *materialStore) UpdateMaterialStatus(id uint, status string) error {
	if status != models.StatusProcessing && status != models.StatusCompleted && status != models.StatusFailed {
		return errors.New("invalid status")
	}
	return s.PerformDBTransaction(func(tx *gorm.DB) error {
		return tx.Model(&models.Material{}).Where("id = ?", id).Update("status", status).Error
	})
}

func (s *materialStore) GetMaterialStatus(id uint) (string, error) {
	var material models.Material
	err := s.DB.Select("status").Where("id = ?", id).First(&material).Error
	return material.Status, err
}
