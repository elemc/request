package request

// ContextKey - тип строки для ключей контекста
type ContextKey string

// ключи контекста
const (
	ContextKeyRequestID       ContextKey = "request_id"
	ContextKeySessionUsername ContextKey = "session_username"
	ContextKeySessionToken    ContextKey = "session_token"
)
