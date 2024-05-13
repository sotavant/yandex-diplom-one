package middleware

import (
	"compress/gzip"
	"github.com/sotavant/yandex-diplom-one/internal"
	"io"
	"net/http"
	"strings"
)

const AcceptableEncoding = "gzip"

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func Gzip(h http.Handler) http.Handler {
	gzipFn := func(w http.ResponseWriter, r *http.Request) {
		ow := w
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, AcceptableEncoding)
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				internal.Logger.Infow("compressReaderError", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer func(cr *compressReader) {
				err = cr.Close()
				if err != nil {
					internal.Logger.Infow("compressReaderCloseError", "err", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}(cr)
		}

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")

		if supportGzip {
			cw := newCompressWriter(w)
			cw.Header().Set("Content-Encoding", "gzip")
			ow = cw
			defer func(cw *compressWriter) {
				err := cw.Close()
				if err != nil {
					internal.Logger.Infow("compressWriterCloseError", "err", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}(cw)
		}

		h.ServeHTTP(ow, r)
	}

	return http.HandlerFunc(gzipFn)
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < http.StatusMultipleChoices {
		c.w.Header().Set("Content-Encoding", "gzip")
	}

	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}

	return c.zr.Close()
}
