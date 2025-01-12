package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func Decode(c *gin.Context, object interface{}) error {
	if !strings.Contains(c.Request.Header.Get("Content-type"), "application/json") {
		return errors.New("only application/json is supported")
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1048576)

	if err := c.BindJSON(&object); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return errors.New("request body contains badly-formed JSON (at position " + cast.ToString(syntaxError.Offset) + ")")
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("request body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			return errors.New("request body contains an invalid value for the " + unmarshalTypeError.Field + " field (at position " + cast.ToString(unmarshalTypeError.Offset) + ")")
		case errors.Is(err, io.EOF):
			return errors.New("request body must not be empty")
		case err.Error() == "http: request body too large":
			return errors.New("request body must not be larger than 1MB")
		default:
			return err
		}
	}
	return nil
}
