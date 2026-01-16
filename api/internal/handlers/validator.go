package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// Validate valida uma struct
func Validate(data interface{}) error {
	return validate.Struct(data)
}

// RespondWithValidationError envia resposta de erro de validação
func RespondWithValidationError(w http.ResponseWriter, err error) {
	fields := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			fields[e.Field()] = translateError(e)
		}
	}

	RespondWithError(w, http.StatusBadRequest, "validation_error", "Os dados fornecidos são inválidos", fields)
}

// RespondWithError envia resposta de erro
func RespondWithError(w http.ResponseWriter, statusCode int, errorType, message string, fields map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   errorType,
		Message: message,
		Fields:  fields,
	})
}

// RespondWithJSON envia resposta de sucesso
func RespondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func translateError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "campo obrigatório"
	case "email":
		return "deve ser um e-mail válido"
	case "min":
		return "deve ter no mínimo " + e.Param()
	case "max":
		return "deve ter no máximo " + e.Param()
	case "gt":
		return "deve ser maior que " + e.Param()
	case "gte":
		return "deve ser maior ou igual a " + e.Param()
	case "lt":
		return "deve ser menor que " + e.Param()
	case "lte":
		return "deve ser menor ou igual a " + e.Param()
	case "uuid":
		return "deve ser um UUID válido"
	case "url":
		return "deve ser uma URL válida"
	case "oneof":
		return "deve ser um dos valores: " + e.Param()
	default:
		return "valor inválido"
	}
}
