package utils

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/kokoc-hack/internal/logging"
)

func ErrGetRecordsFailed(name string, err error) error {
	logging.Log.Errorf("failed to get %s records: %+v", name, err)
	return fiber.NewError(http.StatusInternalServerError, err.Error())
}

func ErrCreateRecordsFailed(name string, err error) error {
	logging.Log.Errorf("failed to create %s records: %+v", name, err)
	return fiber.NewError(http.StatusBadRequest, err.Error())
}

func ErrUpdateRecordsFailed(name string, err error) error {
	logging.Log.Errorf("failed to update %s records: %+v", name, err)
	return fiber.NewError(http.StatusBadRequest, err.Error())
}

func ErrValidationError(field string, err error) error {
	logging.Log.Errorf("validation error in %s: %+v", field, err)
	return fiber.NewError(http.StatusUnprocessableEntity, err.Error())
}

func ErrCustomResponse(status int, msg string, err error) error {
	logging.Log.Errorf("%s: %+v", msg, err)
	return fiber.NewError(status, err.Error())
}
