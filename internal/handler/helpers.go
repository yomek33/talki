package handler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/logger"
	"github.com/yomek33/talki/internal/models"
)

func respondWithError(c echo.Context, code int, message string) error {
	return c.JSON(code, map[string]string{"error": message})
}

func getUserUIDFromContext(c echo.Context) (string, error) {
	user, ok := c.Get("userUID").(string)
	if !ok || user == "" {
		logger.Errorf("Invalid user token")
		return "", errors.New(ErrInvalidUserToken)
	}
	return user, nil
}

func bindAndValidateMaterial(c echo.Context, material *models.Material) error {
	if err := c.Bind(material); err != nil {
		logger.Errorf("Error binding material: %v", err)
		return errors.New(ErrInvalidMaterialData)
	}
	if err := validateMaterial(material); err != nil {
		return err
	}
	return nil
}

func validateMaterial(material *models.Material) error {
	validate := validator.New()
	if err := validate.Struct(material); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error in field '%s': %s", strings.ToLower(err.Field()), err.Tag())
			errorMessages = append(errorMessages, errorMessage)

		}
		logger.Errorf("Error validating material: %v", errors.New(strings.Join(errorMessages, ", ")))
		return errors.New(strings.Join(errorMessages, ", "))
	}
	return nil
}

func parseUintParam(c echo.Context, paramName string) (uint, error) {
	param := c.Param(paramName)
	value, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		logger.Errorf("Error parsing uint param: %v", err)
		return 0, err
	}
	return uint(value), err
}
