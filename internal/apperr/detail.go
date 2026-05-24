package apperr

import "errors"

// RootCause возвращает текст самой внутренней ошибки в цепочке.
func RootCause(err error) string {
	if err == nil {
		return ""
	}
	root := err
	for {
		next := errors.Unwrap(root)
		if next == nil {
			break
		}
		root = next
	}
	return root.Error()
}

// PayloadWithDetail добавляет поле detail в payload для ответа API.
func PayloadWithDetail(payload any, err error) map[string]any {
	out := map[string]any{}
	if m, ok := payload.(map[string]any); ok {
		for k, v := range m {
			out[k] = v
		}
	}
	if err != nil {
		out["detail"] = RootCause(err)
	}
	return out
}
