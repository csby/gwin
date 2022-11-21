package main

import (
	"github.com/csby/gwin/controller"
	"github.com/csby/gwsf/gtype"
)

type Controller struct {
	opt  *controller.Opt
	dhcp *controller.Dhcp
	dns  *controller.Dns
	svn  *controller.Svn
}

func (s *Controller) Init(h *Handler) {
	s.opt = controller.NewOpt(log, cfg)
	s.dhcp = controller.NewDhcp(log, cfg)
	s.dns = controller.NewDns(log, cfg)
	s.svn = controller.NewSvn(log, cfg)
}

func (s *Controller) InitRouting(router gtype.Router, path *gtype.Path) {
	// DHCP
	if cfg.Dhcp.Enable {
		router.POST(path.Uri("/dhcp/filter/list"), nil,
			s.dhcp.GetFilters, s.dhcp.GetFiltersDoc)
		router.POST(path.Uri("/dhcp/filter/add"), nil,
			s.dhcp.AddFilter, s.dhcp.AddFilterDoc)
		router.POST(path.Uri("/dhcp/filter/del"), nil,
			s.dhcp.DelFilter, s.dhcp.DelFilterDoc)
		router.POST(path.Uri("/dhcp/filter/mod"), nil,
			s.dhcp.ModFilter, s.dhcp.ModFilterDoc)

		router.POST(path.Uri("/dhcp/lease/list"), nil,
			s.dhcp.GetLeases, s.dhcp.GetLeasesDoc)
	}

	// DNS
	if cfg.Dns.Enable {
		router.POST(path.Uri("/dns/record/list"), nil,
			s.dns.GetRecords, s.dns.GetRecordsDoc)
		router.POST(path.Uri("/dns/record/add"), nil,
			s.dns.AddRecord, s.dns.AddRecordDoc)
		router.POST(path.Uri("/dns/record/del"), nil,
			s.dns.DelRecord, s.dns.DelRecordDoc)
	}

	// SVN
	if cfg.Svn.Enable {
		router.POST(path.Uri("/svn/user/all/list"), nil,
			s.svn.GetUsers, s.svn.GetUsersDoc)
		router.POST(path.Uri("/svn/repository/new"), nil,
			s.svn.NewRepository, s.svn.NewRepositoryDoc)
		router.POST(path.Uri("/svn/repository/list"), nil,
			s.svn.GetRepositories, s.svn.GetRepositoriesDoc)
		router.POST(path.Uri("/svn/folder/list"), nil,
			s.svn.GetFolders, s.svn.GetFoldersDoc)
		router.POST(path.Uri("/svn/permission/list"), nil,
			s.svn.GetPermissions, s.svn.GetPermissionsDoc)
		router.POST(path.Uri("/svn/user/permission/list"), nil,
			s.svn.GetUserPermissions, s.svn.GetUserPermissionsDoc)
		router.POST(path.Uri("/svn/permission/add"), nil,
			s.svn.AddPermission, s.svn.AddPermissionDoc)
		router.POST(path.Uri("/svn/permission/mod"), nil,
			s.svn.ModPermission, s.svn.ModPermissionDoc)
		router.POST(path.Uri("/svn/permission/del"), nil,
			s.svn.DelPermission, s.svn.DelPermissionDoc)
	}
}
