package request

// Функция извлекает из контекста информацию о требуемом ключе
func (r *Request) contextFetchStringValue(key ContextKey) string {
	if ivalue := r.ctx.Value(key); ivalue != nil {
		res, ok := ivalue.(string)
		if ok {
			return res
		}
	}
	return ""
}

// RequestID - функция возвращает request_id из контекста запроса
func (r *Request) RequestID() string {
	return r.contextFetchStringValue(ContextKeyRequestID)
}

// SessionUsername - функция возвращает session_username за контекста запроса
func (r *Request) SessionUsername() string {
	return r.contextFetchStringValue(ContextKeySessionUsername)
}

// SessionToken - функция возвращает session_token за контекста запроса
func (r *Request) SessionToken() string {
	return r.contextFetchStringValue(ContextKeySessionToken)
}
