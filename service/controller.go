package main

import (
	"fmt"
	"github.com/csby/gwsf/gtype"
	"net"
	"strings"
)

type Controller struct {
	gtype.Base
}

func (s *Controller) GetDhcpFilters(ctx gtype.Context, ps gtype.Params) {
	dhcp := &Dhcp{}
	results, err := dhcp.GetFilters()
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(results)
}

func (s *Controller) GetDhcpFiltersDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "筛选器")
	function := catalog.AddFunction(method, uri, "获取筛选器列表")
	function.SetNote("获取IPv4筛选器列表")
	function.SetOutputDataExample([]*Filter{
		{
			Allow:   true,
			Address: "00-1C-23-20-AF-4A",
			Comment: "描述信息",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Controller) AddDhcpFilter(ctx gtype.Context, ps gtype.Params) {
	argument := &Filter{}
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

	dhcp := &Dhcp{}
	err = dhcp.AddFilter(argument)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(nil)
}

func (s *Controller) AddDhcpFilterDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "筛选器")
	function := catalog.AddFunction(method, uri, "添加筛选器")
	function.SetNote("添加IPv4筛选器到允许或拒绝列表")
	function.SetInputJsonExample(&Filter{
		Allow:   true,
		Address: "00-1C-23-20-AF-4A",
		Comment: "描述信息",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Controller) DelDhcpFilter(ctx gtype.Context, ps gtype.Params) {
	argument := &FilterDelete{}
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

	dhcp := &Dhcp{}
	err = dhcp.DeleteFilter(argument.Address)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(nil)
}

func (s *Controller) DelDhcpFilterDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "筛选器")
	function := catalog.AddFunction(method, uri, "删除筛选器")
	function.SetNote("从IPv4筛选器允许或拒绝列表中删除指定的筛选器")
	function.SetInputJsonExample(&FilterDelete{
		Address: "00-1C-23-20-AF-4A",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Controller) ModDhcpFilter(ctx gtype.Context, ps gtype.Params) {
	argument := &FilterModify{}
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

	dhcp := &Dhcp{}
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

func (s *Controller) ModDhcpFilterDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "筛选器")
	function := catalog.AddFunction(method, uri, "修改筛选器")
	function.SetNote("修改IPv4筛选器允许或拒绝列表中已存在的筛选器")
	function.SetInputJsonExample(&FilterModify{
		Address: "00-1C-23-20-AF-4A",
		Filter: Filter{
			Allow:   true,
			Address: "00-1C-23-20-AF-4B",
			Comment: "描述信息",
		},
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrInput)
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Controller) createCatalog(doc gtype.Doc, names ...string) gtype.Catalog {

	root := doc.AddCatalog("DHCP")

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
