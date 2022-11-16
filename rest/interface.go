package rest

import (
	"io"
	"os"

	"github.com/chi3ndd/golib/log"
)

type (
	JSONInterface interface {
		Code(code int) JSONInterface
		Message(message string) JSONInterface
		Data(data interface{}) JSONInterface
		Log(data interface{}) JSONInterface
		Go() error
		ResponseOK(data interface{}) error
		ResponseError(code int, err error, message string) error
	}

	StreamInterface interface {
		Code(code int) StreamInterface
		ContentType(contentType string) StreamInterface
		Body(data io.Reader) StreamInterface
		Go() error
	}

	AttachmentInterface interface {
		Name(name string) AttachmentInterface
		Path(path string) AttachmentInterface
		Go() error
	}
)

var logger log.Logger

func init() {
	logger, _ = log.New(Module, log.DebugLevel, true, os.Stdout)
}
