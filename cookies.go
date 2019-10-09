package request

// GetCookieValue - returning cookie value as Value
// if doesn't exists then returning empty string
func (r *Request) GetCookieValue(name string) Value {
	cookie, err := r.r.Cookie(name)
	if err != nil {
		return ""
	}
	return Value(cookie.Value)
}
