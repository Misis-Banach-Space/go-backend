package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/kokoc-hack/internal/logging"
)

func errGetRecordsFailed(name string, err error) error {
	logging.Log.Errorf("failed to get %s records: %+v", name, err)
	return fiber.NewError(http.StatusInternalServerError, err.Error())
}

func errCreateRecordsFailed(name string, err error) error {
	logging.Log.Errorf("failed to create %s records: %+v", name, err)
	return fiber.NewError(http.StatusBadRequest, err.Error())
}

func errUpdateRecordsFailed(name string, err error) error {
	logging.Log.Errorf("failed to update %s records: %+v", name, err)
	return fiber.NewError(http.StatusBadRequest, err.Error())
}

func errValidationError(field string, err error) error {
	logging.Log.Errorf("validation error in %s: %+v", field, err)
	return fiber.NewError(http.StatusUnprocessableEntity, err.Error())
}

func errCustomResponse(status int, msg string, err error) error {
	logging.Log.Errorf("%s: %+v", msg, err)
	return fiber.NewError(status, err.Error())
}
