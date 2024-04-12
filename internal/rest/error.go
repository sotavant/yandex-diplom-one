package rest

import (
	"github.com/go-chi/render"
	"net/http"
)

type errResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func errorRender(code int, err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: code,
		StatusText:     http.StatusText(code),
		ErrorText:      err.Error(),
	}
}
