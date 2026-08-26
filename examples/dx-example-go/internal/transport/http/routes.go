package http

import httpx "github.com/datakaveri/dx-common-go/platform/http"

const URNs httpx.URNSpace = "example"

func Routes(handler *Handler) httpx.RouteSet {
	return httpx.Routes("/widgets",
		httpx.POST("/", httpx.Handle(handler.Create, httpx.WithURNs(URNs)), httpx.OpID("createWidget")),
		httpx.GET("/{id}", httpx.Handle(handler.Get, httpx.WithURNs(URNs)), httpx.OpID("getWidget")),
	)
}
