package apperr

// ValidationFields — ошибки валидации по полям для ответа 422.
func ValidationFields(fields map[string]string) *AppError {
	return &AppError{
		Code:    CodeValidation,
		Message: "validation failed",
		Payload: fields,
	}
}

// AsFieldErrors извлекает map полей из Payload.
func AsFieldErrors(payload any) map[string]string {
	switch p := payload.(type) {
	case map[string]string:
		return p
	case map[string]any:
		out := make(map[string]string, len(p))
		for k, v := range p {
			if s, ok := v.(string); ok {
				out[k] = s
			}
		}
		return out
	default:
		return nil
	}
}
