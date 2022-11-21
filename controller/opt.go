package controller

import (
	"github.com/csby/gwin/config"
	"github.com/csby/gwin/model"
	"github.com/csby/gwsf/gtype"
)

func NewOpt(log gtype.Log, cfg *config.Config) *Opt {
	inst := &Opt{}
	inst.SetLog(log)
	inst.cfg = cfg

	return inst
}

type Opt struct {
	base
}

func (s *Opt) GetSetting(ctx gtype.Context, ps gtype.Params) {
	setting := &model.Opt{}
	setting.Dhcp.Enable = s.cfg.Dhcp.Enable
	setting.Dns.Enable = s.cfg.Dns.Enable
	setting.Dns.ZoomNames = s.cfg.Dns.ZoomNames
	if setting.Dns.ZoomNames == nil {
		setting.Dns.ZoomNames = make([]string, 0)
	}
	setting.Svn.Enable = s.cfg.Svn.Enable

	ctx.Success(setting)
}

func (s *Opt) GetSettingDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc)
	function := catalog.AddFunction(method, uri, "获取接口配置")
	function.SetNote("获取接口配置信息")
	function.SetOutputDataExample(&model.Opt{
		Dns: model.OptDns{
			ZoomNames: []string{
				"example.com",
			},
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Opt) createCatalog(doc gtype.Doc, names ...string) gtype.Catalog {
	root := s.createRootCatalog(doc, "服务设置")
	count := len(names)
	if count < 1 {
		return root
	}

	child := root
	for i := 0; i < count; i++ {
		name := names[i]
		child = child.AddChild(name)
	}

	return child
}

func (s *Opt) createRootCatalog(doc gtype.Doc, names ...string) gtype.Catalog {
	root := doc.AddCatalog("管理平台接口")

	count := len(names)
	if count < 1 {
		return root
	}

	child := root
	for i := 0; i < count; i++ {
		name := names[i]
		child = child.AddChild(name)
	}

	return child
}
