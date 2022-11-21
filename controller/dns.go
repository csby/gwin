package controller

import (
	"github.com/csby/gwin/assist"
	"github.com/csby/gwin/config"
	"github.com/csby/gwin/model"
	"github.com/csby/gwsf/gtype"
)

func NewDns(log gtype.Log, cfg *config.Config) *Dns {
	inst := &Dns{}
	inst.SetLog(log)
	inst.cfg = cfg

	return inst
}

type Dns struct {
	base
}

func (s *Dns) GetRecords(ctx gtype.Context, ps gtype.Params) {
	argument := &model.DnsRecordArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.ZoneName) < 1 {
		ctx.Error(gtype.ErrInput, "域名(zoneName)为空")
		return
	}
	dns := &assist.Dns{
		ZoneName: argument.ZoneName,
	}
	results, err := dns.GetRecords()
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(results)
}

func (s *Dns) GetRecordsDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "域名解析")
	function := catalog.AddFunction(method, uri, "获取记录列表")
	function.SetNote("获取解析A记录列表")
	function.SetInputJsonExample(&model.DnsRecordArgument{
		ZoneName: "example.com",
	})
	function.SetOutputDataExample([]*model.DnsRecord{
		{
			Name: "server-a",
			Data: "192.168.1.11",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Dns) AddRecord(ctx gtype.Context, ps gtype.Params) {
	argument := &model.DnsRecordModify{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.ZoneName) < 1 {
		ctx.Error(gtype.ErrInput, "域名(zoneName)为空")
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInput, "记录名称(name)为空")
		return
	}
	if len(argument.Data) < 1 {
		ctx.Error(gtype.ErrInput, "记录数据(data)为空")
		return
	}
	dns := &assist.Dns{
		ZoneName: argument.ZoneName,
	}
	err = dns.AddRecord(argument.Name, argument.Data)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(nil)
}

func (s *Dns) AddRecordDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "域名解析")
	function := catalog.AddFunction(method, uri, "添加记录")
	function.SetNote("添加解析A记录")
	function.SetInputJsonExample(&model.DnsRecordModify{
		DnsRecordArgument: model.DnsRecordArgument{
			ZoneName: "example.com",
		},
		DnsRecord: model.DnsRecord{
			Name: "server-a",
			Data: "192.168.1.11",
		},
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Dns) DelRecord(ctx gtype.Context, ps gtype.Params) {
	argument := &model.DnsRecordModify{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.ZoneName) < 1 {
		ctx.Error(gtype.ErrInput, "域名(zoneName)为空")
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInput, "记录名称(name)为空")
		return
	}
	if len(argument.Data) < 1 {
		ctx.Error(gtype.ErrInput, "记录数据(data)为空")
		return
	}
	dns := &assist.Dns{
		ZoneName: argument.ZoneName,
	}
	err = dns.DeleteRecord(argument.Name, argument.Data)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(nil)
}

func (s *Dns) DelRecordDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "域名解析")
	function := catalog.AddFunction(method, uri, "删除记录")
	function.SetNote("删除解析A记录")
	function.SetInputJsonExample(&model.DnsRecordModify{
		DnsRecordArgument: model.DnsRecordArgument{
			ZoneName: "example.com",
		},
		DnsRecord: model.DnsRecord{
			Name: "server-a",
			Data: "192.168.1.11",
		},
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Dns) createCatalog(doc gtype.Doc, names ...string) gtype.Catalog {
	root := s.createRootCatalog(doc, "DNS")
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
