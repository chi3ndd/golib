package rest

import (
	"github.com/labstack/echo/v4"
)

type jsonHandler struct {
	ctx     echo.Context
	code    int
	message string
	data    interface{}
}

func JSON(c echo.Context) JSONInterface {
	// Success
	return &jsonHandler{ctx: c}
}

func (h *jsonHandler) Code(code int) JSONInterface {
	// Success
	h.code = code
	return h
}

func (h *jsonHandler) Message(message string) JSONInterface {
	// Success
	h.message = message
	return h
}

func (h *jsonHandler) Data(data interface{}) JSONInterface {
	// Success
	h.data = data
	return h
}

func (h *jsonHandler) Log(data interface{}) JSONInterface {
	if h.code < StatusBadRequest {
		logger.Infof("code %d: %v", h.code, data)
	} else {
		logger.Errorf("code %d: %v", h.code, data)
	}
	// Success
	return h
}

func (h *jsonHandler) Go() error {
	// Success
	return h.ctx.JSON(h.code, &response{
		Status:  StatusText(h.code),
		Message: h.message,
		Data:    h.data,
	})
}

func (h *jsonHandler) ResponseOK(data interface{}) error {
	// Success
	return h.Code(StatusOK).Data(data).Go()
}

func (h *jsonHandler) ResponseError(code int, err error, message string) error {
	if message == "" && err != nil {
		message = err.Error()
	}
	// Success
	return h.Code(code).Message(message).Log(err).Go()
}
