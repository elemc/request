package request

import "context"

// ContextKey - тип строки для ключей контекста
type ContextKey string

// ключи контекста
const (
	ContextKeyRequestID       ContextKey = "request_id"
	ContextKeySessionUsername ContextKey = "session_username"
	ContextKeySessionToken    ContextKey = "session_token"
)

// Context - обобщенная структура данных и обертка над стандартным контекстом
type Context struct {
	ctx context.Context
}

// NewContext - функция возвращает указатель на контекст запроса
func NewContext(ctx context.Context) *Context {
	return &Context{
		ctx: ctx,
	}
}

// Context - возвращает типовой контекст
func (ctx *Context) Context() context.Context {
	return ctx.ctx
}

// Функция извлекает из контекста информацию о требуемом ключе
func (ctx *Context) contextFetchStringValue(key ContextKey) string {
	if ctx.ctx == nil {
		return ""
	}
	if value := ctx.ctx.Value(key); value != nil {
		res, ok := value.(string)
		if ok {
			return res
		}
	}
	return ""
}

// RequestID - функция возвращает request_id из контекста запроса
func (ctx *Context) RequestID() string {
	return ctx.contextFetchStringValue(ContextKeyRequestID)
}

// SessionUsername - функция возвращает session_username за контекста запроса
func (ctx *Context) SessionUsername() string {
	return ctx.contextFetchStringValue(ContextKeySessionUsername)
}

// SessionToken - функция возвращает session_token за контекста запроса
func (ctx *Context) SessionToken() string {
	return ctx.contextFetchStringValue(ContextKeySessionToken)
}
