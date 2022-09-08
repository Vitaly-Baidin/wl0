package httpErr

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime httpErr
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`         // user-level status message
	AppCode    int64  `json:"code,omitempty"` // application-specific httpErr code
	ErrorText  string `json:"-"`              // application-level httpErr message, for debugging (httpErr,omitempty)
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func Err404Render(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "not found",
		ErrorText:      err.Error(),
	}
}

func Err500Render(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "internal server error",
		ErrorText:      err.Error(),
	}
}
