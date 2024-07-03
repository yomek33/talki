package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/logger"
	"github.com/yomek33/talki/internal/models"
	"github.com/yomek33/talki/internal/services"
)

type MaterialHandler interface {
	CreateMaterial(c echo.Context) error
	GetMaterialByID(c echo.Context) error
	UpdateMaterial(c echo.Context) error
	DeleteMaterial(c echo.Context) error
	GetAllMaterials(c echo.Context) error
	CheckMaterialStatus(c echo.Context) error
}

type materialHandler struct {
	services.MaterialService
	services.PhraseService
}

func NewMaterialHandler(materialService services.MaterialService, phraseService services.PhraseService) MaterialHandler {
	return &materialHandler{
		MaterialService: materialService,
		PhraseService:   phraseService,
	}
}

func (h *materialHandler) CreateMaterial(c echo.Context) error {
	var material models.Material
	if err := bindAndValidateMaterial(c, &material); err != nil {
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}

	UserUID, err := getUserUIDFromContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}

	material.UserUID = UserUID
	material.Status = "processing"

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	id, err := h.MaterialService.CreateMaterial(&material)
	if err != nil {
		logger.Errorf("Error creating material: %v, UserUID: %v", err, UserUID)
		return respondWithError(c, http.StatusInternalServerError, ErrFailedCreateMaterial)
	}

	material.ID = id
	go h.processMaterialAsync(ctx, material.ID, UserUID)

	logger.Info("Material created successfully")
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Material created successfully",
		"id":      material.ID,
	})
}

func (h *materialHandler) GetMaterialByID(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidMaterialID)
	}

	UserUID, err := getUserUIDFromContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}

	material, err := h.MaterialService.GetMaterialByID(id, UserUID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, ErrMaterialNotFound)
	}

	logger.Infof("Retrieved material MaterialID;%v", id)
	return c.JSON(http.StatusOK, material)
}

func (h *materialHandler) UpdateMaterial(c echo.Context) error {
	UserUID, err := getUserUIDFromContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}

	materialID, err := parseUintParam(c, "id")
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidMaterialID)
	}

	material, err := h.MaterialService.GetMaterialByID(materialID, UserUID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, ErrMaterialNotFound)
	}

	if err := bindAndValidateMaterial(c, material); err != nil {
		return respondWithError(c, http.StatusBadRequest, err.Error())
	}

	if material.UserUID != UserUID {
		return respondWithError(c, http.StatusForbidden, ErrForbiddenModify)
	}

	if err := h.MaterialService.UpdateMaterial(materialID, material); err != nil {
		logger.Errorf("Failed to update material: %v, MaterialID: %v, UserUID: %v", err, materialID, UserUID)
		return respondWithError(c, http.StatusInternalServerError, ErrFailedUpdateMaterial)
	}

	logger.Infof("Updated material, MaterialID: %v, UserUID: %v", materialID, UserUID)
	return c.JSON(http.StatusOK, material)
}

func (h *materialHandler) DeleteMaterial(c echo.Context) error {
	materialID, err := parseUintParam(c, "id")
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidID)
	}

	UserUID, err := getUserUIDFromContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}

	if err := h.MaterialService.DeleteMaterial(materialID, UserUID); err != nil {
		logger.Errorf("Failed to delete material: %v, MaterialID: %v, UserUID: %v", err, materialID, UserUID)
		return respondWithError(c, http.StatusInternalServerError, ErrFailedDeleteMaterial)
	}

	logger.Infof("Deleted material, MaterialID: %v, UserUID: %v", materialID, UserUID)
	return c.NoContent(http.StatusNoContent)
}

func (h *materialHandler) GetAllMaterials(c echo.Context) error {
	searchQuery := c.QueryParam("search")

	UserUID, err := getUserUIDFromContext(c)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, ErrInvalidUserToken)
	}

	materials, err := h.MaterialService.GetAllMaterials(searchQuery, UserUID)
	if err != nil {
		logger.Errorf("Failed to retrieve materials: %v, UserUID: %v", err, UserUID)
		return respondWithError(c, http.StatusInternalServerError, ErrFailedRetrieveMaterials)
	}

	logger.Infof("Retrieved materials, MaterialCount: %v, UserUID: %v", len(materials), UserUID)
	return c.JSON(http.StatusOK, materials)
}

func (h *materialHandler) CheckMaterialStatus(c echo.Context) error {
	materialID, err := parseUintParam(c, "id")
	if err != nil {
		return respondWithError(c, http.StatusBadRequest, ErrInvalidMaterialID)
	}

	status, err := h.MaterialService.GetMaterialStatus(materialID)
	if err != nil {
		logger.Errorf("Failed to get material status: %v, MaterialID: %v", err, materialID)
		return respondWithError(c, http.StatusInternalServerError, err.Error())
	}

	logger.Infof("Checked material status, MaterialID: %v, Status: %v", materialID, status)
	return c.JSON(http.StatusOK, map[string]string{"status": status})
}

func (h *materialHandler) processMaterialAsync(ctx context.Context, materialID uint, userUID string) {
	h.MaterialService.UpdateMaterialStatus(materialID, "processing")

	phrases, err := h.PhraseService.GeneratePhrases(ctx, materialID, userUID)
	if err != nil {
		logger.Errorf("Failed to generate phrases: %v, MaterialID: %v, UserUID: %v", err, materialID, userUID)
		h.MaterialService.UpdateMaterialStatus(materialID, "failed")
		return
	}

	if err = h.PhraseService.StorePhrases(materialID, phrases); err != nil {
		logger.Errorf("Failed to store phrases: %v, MaterialID: %v, UserUID: %v", err, materialID, userUID)
		h.MaterialService.UpdateMaterialStatus(materialID, "failed")
		return
	}

	logger.Infof("Phrases generated and stored successfully, MaterialID: %v, UserUID: %v", materialID, userUID)
	h.MaterialService.UpdateMaterialStatus(materialID, "completed")
}
