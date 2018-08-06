package server

import (
	"context"
	"net/http"
)

type ctxKey int

const (
	ctxTemplateVars ctxKey = iota
)

func ctxSet(r *http.Request, key ctxKey, value interface{}) *http.Request {
	ctx := context.WithValue(r.Context(), key, value)
	return r.WithContext(ctx)
}

func ctxGet(r *http.Request, key ctxKey) interface{} {
	return r.Context().Value(key)
}

func ctxSetTemplateVars(r *http.Request, d map[string]interface{}) *http.Request {
	return ctxSet(r, ctxTemplateVars, &d)
}

func ctxGetTemplateVars(r *http.Request) *map[string]interface{} {
	if a, ok := ctxGet(r, ctxTemplateVars).(*map[string]interface{}); ok {
		return a
	}
	return nil
}

func ctxAppendTemplateVars(r *http.Request, data map[string]interface{}) *http.Request {
	m := ctxGetTemplateVars(r)
	if m == nil {
		data := map[string]interface{}{}
		r = ctxSetTemplateVars(r, data)
		m = &data
	}

	extendTemplateVars(data, m)
	return r
}
