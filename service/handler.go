package main

import (
	"fmt"
	"github.com/csby/gwsf/gopt"
	"github.com/csby/gwsf/gtype"
	"net/http"
)

func NewHandler(log gtype.Log) gtype.Handler {
	instance := &Handler{}
	instance.SetLog(log)

	instance.apiController = &Controller{}
	instance.apiController.SetLog(log)

	return instance
}

type Handler struct {
	gtype.Base

	apiController *Controller
}

func (s *Handler) InitRouting(router gtype.Router) {
	router.POST(apiPath.Uri("/dhcp/filter/list"), nil,
		s.apiController.GetDhcpFilters, s.apiController.GetDhcpFiltersDoc)
	router.POST(apiPath.Uri("/dhcp/filter/add"), nil,
		s.apiController.AddDhcpFilter, s.apiController.AddDhcpFilterDoc)
	router.POST(apiPath.Uri("/dhcp/filter/del"), nil,
		s.apiController.DelDhcpFilter, s.apiController.DelDhcpFilterDoc)
	router.POST(apiPath.Uri("/dhcp/filter/mod"), nil,
		s.apiController.ModDhcpFilter, s.apiController.ModDhcpFilterDoc)
}

func (s *Handler) BeforeRouting(ctx gtype.Context) {
	method := ctx.Method()

	// enable across access
	if method == "OPTIONS" {
		ctx.Response().Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "content-type,token")
		ctx.SetHandled(true)
		return
	}

	// default to opt site
	if method == "GET" {
		path := ctx.Path()
		if "/" == path || "" == path || gopt.WebPath == path {
			redirectUrl := fmt.Sprintf("%s://%s%s/", ctx.Schema(), ctx.Host(), gopt.WebPath)
			http.Redirect(ctx.Response(), ctx.Request(), redirectUrl, http.StatusMovedPermanently)
			ctx.SetHandled(true)
			return
		}
	}
}

func (s *Handler) AfterRouting(ctx gtype.Context) {

}

func (s *Handler) ExtendOptApi(router gtype.Router, path *gtype.Path, preHandle gtype.HttpHandle, wsc gtype.SocketChannelCollection) {

}
