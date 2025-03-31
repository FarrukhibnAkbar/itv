package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcErrWithMessage(w http.ResponseWriter, r *http.Request, err error, args ...interface{}) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok || st.Code() == codes.Internal {
		// logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, response{
			Error:     true,
			Message:   st.Message(),
			Status:    http.StatusInternalServerError,
			Timestamp: time.Now().Format(time.RFC3339),
			Path:      r.URL.Path,
		})
		return err

	} else if st.Code() == codes.NotFound {
		// logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, response{
			Error:     true,
			Message:   st.Message(),
			Status:    http.StatusNotFound,
			Timestamp: time.Now().Format(time.RFC3339),
			Path:      r.URL.Path,
		})
		return err

	} else if st.Code() == codes.Unavailable {
		// logger.Error(err)
		w.WriteHeader(http.StatusBadGateway)
		writeJSON(w, response{
			Error:     true,
			Message:   st.Message(),
			Status:    http.StatusBadGateway,
			Timestamp: time.Now().Format(time.RFC3339),
			Path:      r.URL.Path,
		})
		return err

	} else if st.Code() == codes.AlreadyExists {
		// logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, response{
			Error:     true,
			Message:   st.Message(),
			Status:    http.StatusBadRequest,
			Timestamp: time.Now().Format(time.RFC3339),
			Path:      r.URL.Path,
		})
		return err

	} else if st.Code() == codes.InvalidArgument {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, response{
			Error:     true,
			Message:   st.Message(),
			Status:    http.StatusBadRequest,
			Timestamp: time.Now().Format(time.RFC3339),
			Path:      r.URL.Path,
		})
		return err

	} else if st.Code() == codes.DataLoss {
		// logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, response{
			Error:     true,
			Message:   st.Message(),
			Status:    http.StatusBadRequest,
			Timestamp: time.Now().Format(time.RFC3339),
			Path:      r.URL.Path,
		})
		return err

	} else if st.Code() == codes.PermissionDenied {
		// logger.Error(err)
		w.WriteHeader(http.StatusForbidden)
		writeJSON(w, response{
			Error:     true,
			Message:   st.Message(),
			Status:    http.StatusForbidden,
			Timestamp: time.Now().Format(time.RFC3339),
			Path:      r.URL.Path,
		})
		return err

	} else if strings.Contains("User blocked in user service", st.Message()) {
		// logger.Error(err)
		w.WriteHeader(http.StatusForbidden)
		writeJSON(w, response{
			Error:     true,
			Message:   st.Message(),
			Status:    http.StatusForbidden,
			Timestamp: time.Now().Format(time.RFC3339),
			Path:      r.URL.Path,
		})
		return err
	} else if st.Code() == codes.Unauthenticated {
		// logger.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		writeJSON(w, response{
			Error:     true,
			Message:   st.Message(),
			Status:    http.StatusUnauthorized,
			Timestamp: time.Now().Format(time.RFC3339),
			Path:      r.URL.Path,
		})
		return err

	}
	// logger.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
	writeJSON(w, response{
		Error:     true,
		Message:   st.Message(),
		Status:    http.StatusInternalServerError,
		Timestamp: time.Now().Format(time.RFC3339),
		Path:      r.URL.Path,
	})

	return err
}

// ValidationError ...
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ResponseForExec for insert, update, delete
type ResponseForExec struct {
	Info string `json:"info"`
	ID   string `json:"id"`
}

func HandleBodyParseError(err error) ValidationError {
	var validationErrors ValidationError

	var ute *json.UnmarshalTypeError
	if errors.As(err, &ute) {
		validationErrors = ValidationError{
			Field:   ute.Field,
			Message: fmt.Sprintf("Type mismatch: expected %s but got %s", ute.Type.String(), ute.Value),
		}
	} else if syntaxErr, ok := err.(*json.SyntaxError); ok {
		validationErrors = ValidationError{
			Field:   "JSON",
			Message: fmt.Sprintf("Syntax error at offset %d", syntaxErr.Offset),
		}
	} else {
		validationErrors = ValidationError{
			Field:   "Unknown",
			Message: err.Error(),
		}
	}

	return validationErrors
}

type response struct {
	Timestamp string      `json:"timestamp"`
	Status    int         `json:"status"`
	Message   interface{} `json:"message"`
	Error     bool        `json:"error"`
	Path      string      `json:"path"`
	Data      interface{} `json:"data"`
}

type Response struct {
	Timestamp string      `json:"timestamp"`
	Status    int         `json:"status"`
	Message   interface{} `json:"message"`
	Error     bool        `json:"error"`
	Path      string      `json:"path"`
	Data      interface{} `json:"data"`
}

func BodyParser(r *http.Request, body interface{}) error {
	return json.NewDecoder(r.Body).Decode(&body)
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	bytes, _ := json.MarshalIndent(data, "", "  ")

	w.Header().Set("Content-Type", "Application/json")
	w.Write(bytes)
}

func WriteJSONWithSuccess(w http.ResponseWriter, r *http.Request, data interface{}) {
	if data != nil {
		v := reflect.ValueOf(data)
		if (v.Kind() == reflect.Slice || v.Kind() == reflect.Array) && v.Len() == 0 {
			data = []interface{}{}
		}
	}

	data = response{
		Error:     false,
		Status:    http.StatusOK,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
		Path:      r.URL.Path,
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("WriteJSONWithSuccess err:", err)
	}

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// ParsePagination parses page and limit value from request query
func ParsePagination(r *http.Request) (page int, limit int, err error) {
	q := r.URL.Query()

	pageStr := q.Get("page")

	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return 0, 0, err
		}
		if page == 0 {
			page = 1
		}
	}

	limitStr := q.Get("limit")

	if limitStr == "" {
		limit = 10
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, err
		}
	}

	return page, limit, nil
}
