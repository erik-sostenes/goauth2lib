package response

import (
	"encoding/json"
	"mime"
	"net/http"
)

// Response represents the http request body
type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// JSON encodes the body of the http request in json format
func JSON(w http.ResponseWriter, method int, body any) error {
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(method)

	return json.NewEncoder(w).Encode(body)
}

// Bind parses the request data to an object
// checks the content type
func Bind(w http.ResponseWriter, r *http.Request, body any) (err error) {
	content := r.Header.Get("Content-type")
	if content == "" {
		return JSON(w, http.StatusUnprocessableEntity, Response{
			Message: "missing content type",
		})
	}

	mediaType, _, err := mime.ParseMediaType(content)
	if err != nil {
		return JSON(w, http.StatusUnprocessableEntity, Response{
			Message: err.Error(),
		})
	}

	switch mediaType {
	case "application/json; charset=utf-8":
		err = json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			return JSON(w, http.StatusUnprocessableEntity, Response{
				Message: "the format of the body of the request is malformed",
			})
		}
	}

	return
}
