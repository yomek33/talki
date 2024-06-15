package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/stores"
)

type MaterialService interface {
	CreateMaterial(material *models.Material) (uint, error)
	GetMaterialByID(id uint, UserUID string) (*models.Material, error)
	UpdateMaterial(id uint, material *models.Material) error
	DeleteMaterial(id uint, UserUID string) error
	GetAllMaterials(searchQuery string, UserUID string) ([]models.Material, error)
	UpdateMaterialStatus(id uint, status string) error
	GetMaterialStatus(id uint) (string, error)
}

type materialService struct {
	store stores.MaterialStore
	mu    sync.Mutex
}

var (
	ErrMaterialNil          = errors.New("material cannot be nil")
	ErrMismatchedMaterialID = errors.New("mismatched material ID")
)

func (s *materialService) CreateMaterial(material *models.Material) (uint, error) {
	if material == nil {
		return 0, errors.New("material cannot be nil")
	}
	return s.store.CreateMaterial(material)
}

func (s *materialService) GetMaterialByID(id uint, UserUID string) (*models.Material, error) {
	material, err := s.store.GetMaterialByID(id, UserUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get material by ID: %w", err)
	}
	return material, nil
}

func (s *materialService) UpdateMaterial(id uint, material *models.Material) error {
	if material == nil {
		return ErrMaterialNil
	}
	if id != material.ID {
		return ErrMismatchedMaterialID
	}
	return s.store.UpdateMaterial(id, material)
}

func (s *materialService) DeleteMaterial(id uint, UserUID string) error {
	return s.store.DeleteMaterial(id, UserUID)
}

func (s *materialService) GetAllMaterials(searchQuery string, UserUID string) ([]models.Material, error) {
	return s.store.GetAllMaterials(searchQuery, UserUID)
}

func (s *materialService) UpdateMaterialStatus(id uint, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.store.UpdateMaterialStatus(id, status)
}

func (s *materialService) GetMaterialStatus(id uint) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.store.GetMaterialStatus(id)
}
