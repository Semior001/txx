package proxy

import (
	"database/sql/driver"
)

// Handlers is a set of middlewares.
type Handlers struct {
	QueryContextHandler   QueryContextMiddleware
	PrepareContextHandler PrepareContextMiddleware
	BeginTxHandler        BeginTxMiddleware
	ExecContextHandler    ExecContextMiddleware
}

func (h *Handlers) fillStubs() {
	if h.QueryContextHandler == nil {
		h.QueryContextHandler = func(next driver.QueryerContext) driver.QueryerContext {
			return next
		}
	}

	if h.PrepareContextHandler == nil {
		h.PrepareContextHandler = func(next driver.ConnPrepareContext) driver.ConnPrepareContext {
			return next
		}
	}

	if h.BeginTxHandler == nil {
		h.BeginTxHandler = func(next driver.ConnBeginTx) driver.ConnBeginTx {
			return next
		}
	}

	if h.ExecContextHandler == nil {
		h.ExecContextHandler = func(next driver.ExecerContext) driver.ExecerContext {
			return next
		}
	}
}
