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

func ctxGetTemplateVars(r *http.Request) map[string]interface{} {
	if m, ok := ctxGet(r, ctxTemplateVars).(map[string]interface{}); ok {
		return m
	}
	return nil
}

func ctxSetTemplateVars(r *http.Request, vars map[string]interface{}) {
	ctxSet(r, ctxTemplateVars, vars)
}

func ctxAddTemplateVars(r *http.Request, data map[string]interface{}) {
	vars := ctxGetTemplateVars(r)
	if vars == nil {
		ctxSetTemplateVars(r, map[string]interface{}{})
		vars = ctxGetTemplateVars(r)
	}

	extendTemplateVars(vars, data)
}
