package httpx_ext

import (
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
)

const (
	ContentTypeHtml = "text/html"
	ContentTypePng  = "image/png"
	ContentTypeJPEG = "image/jpeg"
	ContentTypeGif  = "image/gif"
	ContentTypeIcon = "image/x-icon"
	ContentTypeCss  = "text/css"
	ContentTypeJs   = "application/x-javascript"
)

// WriteJson writes v as json string into w with code.
func WriteHtml(w http.ResponseWriter, code int, v []byte) {
	w.Header().Set(httpx.ContentType, ContentTypeHtml)
	w.WriteHeader(code)

	if n, err := w.Write(v); err != nil {
		// http.ErrHandlerTimeout has been handled by http.TimeoutHandler,
		// so it's ignored here.
		if err != http.ErrHandlerTimeout {
			logx.Errorf("write response failed, error: %s", err)
		}
	} else if n < len(v) {
		logx.Errorf("actual bytes: %d, written bytes: %d", len(v), n)
	}
}

// WriteJson writes v as json string into w with code.
func Write(w http.ResponseWriter, code int, ContentType string, v []byte) {
	w.Header().Set(httpx.ContentType, ContentType)
	w.WriteHeader(code)

	if n, err := w.Write(v); err != nil {
		// http.ErrHandlerTimeout has been handled by http.TimeoutHandler,
		// so it's ignored here.
		if err != http.ErrHandlerTimeout {
			logx.Errorf("write response failed, error: %s", err)
		}
	} else if n < len(v) {
		logx.Errorf("actual bytes: %d, written bytes: %d", len(v), n)
	}
}
