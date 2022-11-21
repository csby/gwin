package controller

import (
	"fmt"
	"github.com/csby/gwin/assist"
	"github.com/csby/gwin/config"
	"github.com/csby/gwin/model"
	"github.com/csby/gwsf/gtype"
	"net"
	"strings"
)

func NewDhcp(log gtype.Log, cfg *config.Config) *Dhcp {
	inst := &Dhcp{}
	inst.SetLog(log)
	inst.cfg = cfg

	return inst
}

type Dhcp struct {
	base
}

func (s *Dhcp) GetFilters(ctx gtype.Context, ps gtype.Params) {
	dhcp := &assist.Dhcp{}
	results, err := dhcp.GetFilters()
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(results)
}

func (s *Dhcp) GetFiltersDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "筛选器")
	function := catalog.AddFunction(method, uri, "获取筛选器列表")
	function.SetNote("获取IPv4筛选器列表")
	function.SetOutputDataExample([]*model.DhcpFilter{
		{
			Allow:   true,
			Address: "00-1C-23-20-AF-4A",
			Comment: "描述信息",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Dhcp) AddFilter(ctx gtype.Context, ps gtype.Params) {
	argument := &model.DhcpFilter{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}

	if len(argument.Address) < 1 {
		ctx.Error(gtype.ErrInput, "MAC地址为空")
		return
	}
	_, err = net.ParseMAC(argument.Address)
	if err != nil {
		ctx.Error(gtype.ErrInput, fmt.Sprintf("MAC地址(%s)无效", argument.Address))
		return
	}
	argument.Address = strings.ToUpper(strings.ReplaceAll(argument.Address, ":", "-"))

	dhcp := &assist.Dhcp{}
	err = dhcp.AddFilter(argument)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(nil)
}

func (s *Dhcp) AddFilterDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "筛选器")
	function := catalog.AddFunction(method, uri, "添加筛选器")
	function.SetNote("添加IPv4筛选器到允许或拒绝列表")
	function.SetInputJsonExample(&model.DhcpFilter{
		Allow:   true,
		Address: "00-1C-23-20-AF-4A",
		Comment: "描述信息",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Dhcp) DelFilter(ctx gtype.Context, ps gtype.Params) {
	argument := &model.DhcpFilterDelete{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}

	if len(argument.Address) < 1 {
		ctx.Error(gtype.ErrInput, "MAC地址为空")
		return
	}
	_, err = net.ParseMAC(argument.Address)
	if err != nil {
		ctx.Error(gtype.ErrInput, fmt.Sprintf("MAC地址(%s)无效", argument.Address))
		return
	}
	argument.Address = strings.ToUpper(strings.ReplaceAll(argument.Address, ":", "-"))

	dhcp := &assist.Dhcp{}
	err = dhcp.DeleteFilter(argument.Address)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(nil)
}

func (s *Dhcp) DelFilterDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "筛选器")
	function := catalog.AddFunction(method, uri, "删除筛选器")
	function.SetNote("从IPv4筛选器允许或拒绝列表中删除指定的筛选器")
	function.SetInputJsonExample(&model.DhcpFilterDelete{
		Address: "00-1C-23-20-AF-4A",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Dhcp) ModFilter(ctx gtype.Context, ps gtype.Params) {
	argument := &model.DhcpFilterModify{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}

	if len(argument.Address) < 1 {
		ctx.Error(gtype.ErrInput, "原MAC地址为空")
		return
	}
	if len(argument.Filter.Address) < 1 {
		ctx.Error(gtype.ErrInput, "新MAC地址为空")
		return
	}
	_, err = net.ParseMAC(argument.Filter.Address)
	if err != nil {
		ctx.Error(gtype.ErrInput, fmt.Sprintf("MAC地址(%s)无效", argument.Filter.Address))
		return
	}

	oldAddr := strings.ToUpper(strings.ReplaceAll(argument.Address, ":", "-"))
	newAddr := strings.ToUpper(strings.ReplaceAll(argument.Filter.Address, ":", "-"))
	argument.Filter.Address = newAddr

	dhcp := &assist.Dhcp{}
	if oldAddr == newAddr {
		err = dhcp.DeleteFilter(oldAddr)
		if err != nil {
			ctx.Error(gtype.ErrInternal.SetDetail(err))
			return
		}
		err = dhcp.AddFilter(&argument.Filter)
		if err != nil {
			ctx.Error(gtype.ErrInternal.SetDetail(err))
			return
		}
	} else {
		err = dhcp.AddFilter(&argument.Filter)
		if err != nil {
			ctx.Error(gtype.ErrInternal.SetDetail(err))
			return
		}
		err = dhcp.DeleteFilter(argument.Address)
		if err != nil {
			ctx.Error(gtype.ErrInternal.SetDetail(err))
			return
		}
	}

	ctx.Success(nil)
}

func (s *Dhcp) ModFilterDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "筛选器")
	function := catalog.AddFunction(method, uri, "修改筛选器")
	function.SetNote("修改IPv4筛选器允许或拒绝列表中已存在的筛选器")
	function.SetInputJsonExample(&model.DhcpFilterModify{
		Address: "00-1C-23-20-AF-4A",
		Filter: model.DhcpFilter{
			Allow:   true,
			Address: "00-1C-23-20-AF-4B",
			Comment: "描述信息",
		},
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Dhcp) GetLeases(ctx gtype.Context, ps gtype.Params) {
	dhcp := &assist.Dhcp{}
	results, err := dhcp.GetLeases()
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(results)
}

func (s *Dhcp) GetLeasesDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "地址租用")
	function := catalog.AddFunction(method, uri, "获取地址租用列表")
	function.SetNote("获取IPv4地址租用列表")
	function.SetOutputDataExample([]*model.DhcpLease{
		{
			IpV4:    "192.168.1.103",
			Address: "00-1C-23-20-AF-4A",
			Comment: "描述信息",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Dhcp) createCatalog(doc gtype.Doc, names ...string) gtype.Catalog {
	root := s.createRootCatalog(doc, "DHCP")
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
