package request

import "fmt"

// ErrorResponse - структура ответа при возникновении
// ошибки
type ErrorResponse struct {
	HasError  bool   `json:"has_error"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

// Error - интерфейсный метод для интерфейса error
func (er *ErrorResponse) Error() string {
	return er.Message
}

// ErrorJSON - возвращает указатель на структуру ErrorResponse с указанным сообщением
func (r *Request) ErrorJSON(msg string, args ...interface{}) *ErrorResponse {
	return &ErrorResponse{
		HasError:  true,
		Message:   fmt.Sprintf(msg, args...),
		RequestID: r.Context().RequestID(),
	}
}
