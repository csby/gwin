package controller

import (
	"fmt"
	"github.com/csby/gwin/assist"
	"github.com/csby/gwin/config"
	"github.com/csby/gwin/model"
	"github.com/csby/gwsf/gtype"
	"strings"
)

func NewSvn(log gtype.Log, cfg *config.Config) *Svn {
	inst := &Svn{}
	inst.SetLog(log)
	inst.cfg = cfg

	return inst
}

type Svn struct {
	base
}

func (s *Svn) GetUsers(ctx gtype.Context, ps gtype.Params) {
	ad := &assist.MsAd{}
	if s.cfg != nil {
		ad.Host = s.cfg.Svn.Ad.Host
		ad.Port = s.cfg.Svn.Ad.Port
		ad.Base = s.cfg.Svn.Ad.Base
		ad.Account = s.cfg.Svn.Ad.Account
		ad.Password = s.cfg.Svn.Ad.Password
	}
	results, err := ad.GetAllUsers()
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(results)
}

func (s *Svn) GetUsersDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "存储库")
	function := catalog.AddFunction(method, uri, "获取用户列表")
	function.SetOutputDataExample([]*model.MsAdUser{
		{
			Id:   "",
			Name: "",
		},
	})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Svn) GetRepositories(ctx gtype.Context, ps gtype.Params) {
	svn := &assist.Svn{}
	results, err := svn.GetRepositories(false)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(results)
}

func (s *Svn) GetRepositoriesDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "存储库")
	function := catalog.AddFunction(method, uri, "获取存储库列表")
	function.SetOutputDataExample([]*model.SvnRepositoryItem{
		{
			Id:         gtype.NewGuid(),
			Repository: "test",
			Name:       "test",
			Path:       "/",
			Type:       model.SvnRepositoryItemKindRepository,
			Children:   []*model.SvnRepositoryItem{},
		},
	})
	function.AddOutputError(gtype.ErrInternal)
}

func (s *Svn) GetFolders(ctx gtype.Context, ps gtype.Params) {
	argument := &model.SvnRepository{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInput, "存储库名称(name)为空")
		return
	}
	if len(argument.Path) < 1 {
		ctx.Error(gtype.ErrInput, "路径(path)为空")
		return
	}

	svn := &assist.Svn{}
	results, err := svn.GetRepositoryFolders(argument.Name, argument.Path, false)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(&model.SvnRepositoryItem{
		Id:       argument.Id,
		Children: results,
	})
}

func (s *Svn) GetFoldersDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "存储库")
	function := catalog.AddFunction(method, uri, "获取文件夹列表")
	function.SetInputJsonExample(&model.SvnRepository{
		Name: "test",
		Path: "/",
	})
	function.SetOutputDataExample([]*model.SvnRepositoryItem{
		{
			Id: gtype.NewGuid(),
			Children: []*model.SvnRepositoryItem{
				{
					Id:         gtype.NewGuid(),
					Repository: "test",
					Name:       "trunk",
					Path:       "/trunk",
					Type:       model.SvnRepositoryItemKindFolder,
					Children:   []*model.SvnRepositoryItem{},
				},
			},
		},
	})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Svn) NewRepository(ctx gtype.Context, ps gtype.Params) {
	argument := &model.SvnRepositoryNew{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInput, "存储库名称(name)为空")
		return
	}

	svn := &assist.Svn{}

	rs, re := svn.GetRepositories(false)
	if re != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(re))
		return
	}
	for i := 0; i < len(rs); i++ {
		if strings.ToLower(argument.Name) == strings.ToLower(rs[i].Repository) {
			ctx.Error(gtype.ErrInput, fmt.Sprintf("存储库名称(%s)已存在", argument.Name))
			return
		}
	}

	err = svn.NewRepository(argument.Name)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(argument.Name)
}

func (s *Svn) NewRepositoryDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "存储库")
	function := catalog.AddFunction(method, uri, "新建存储库")
	function.SetNote("成功时返回存储库名称, 并创建3个文件夹: branches,tags,trunk")
	function.SetInputJsonExample(&model.SvnRepositoryNew{
		Name: "MyRepo",
	})
	function.SetOutputDataExample("MyRepo")
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Svn) GetPermissions(ctx gtype.Context, ps gtype.Params) {
	argument := &model.SvnRepository{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInput, "存储库名称(name)为空")
		return
	}
	if len(argument.Path) < 1 {
		ctx.Error(gtype.ErrInput, "路径(path)为空")
		return
	}

	svn := &assist.Svn{}
	results, err := svn.GetPermissions(argument.Name, argument.Path)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(results)
}

func (s *Svn) GetPermissionsDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "存储库")
	function := catalog.AddFunction(method, uri, "获取项目访问权限列表")
	function.SetInputJsonExample(&model.SvnRepository{
		Name: "test",
		Path: "/",
	})
	function.SetOutputDataExample([]*model.SvnPermission{
		{
			AccountId:   "",
			AccountName: "",
			AccessLevel: model.SvnPermissionReadWrite,
			Inherited:   true,
		},
	})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Svn) GetUserPermissions(ctx gtype.Context, ps gtype.Params) {
	argument := &model.SvnPermissionUserArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.AccountId) < 1 {
		ctx.Error(gtype.ErrInput, "账号ID(accountId)为空")
		return
	}

	svn := &assist.Svn{}
	results, err := svn.GetUserPermissions(argument.AccountId)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(results)
}

func (s *Svn) GetUserPermissionsDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "存储库")
	function := catalog.AddFunction(method, uri, "获取用户访问权限列表")
	function.SetInputJsonExample(&model.SvnPermissionUserArgument{
		AccountId: "S-1-5-32-545",
	})
	function.SetOutputDataExample([]*model.SvnPermissionUser{
		{
			Repository:  "*",
			Path:        "/",
			AccessLevel: model.SvnPermissionReadWrite,
		},
	})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Svn) AddPermission(ctx gtype.Context, ps gtype.Params) {
	argument := &model.SvnPermissionArgumentEdit{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Repository) < 1 {
		ctx.Error(gtype.ErrInput, "存储库名称(repository)为空")
		return
	}
	if len(argument.Path) < 1 {
		ctx.Error(gtype.ErrInput, "路径(path)为空")
		return
	}
	if len(argument.AccountId) < 1 {
		ctx.Error(gtype.ErrInput, "账号ID(accountId)为空")
		return
	}

	svn := &assist.Svn{}
	err = svn.AddPermission(argument.Repository, argument.Path, argument.AccountId, argument.AccessLevel)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(nil)
}

func (s *Svn) AddPermissionDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "存储库")
	function := catalog.AddFunction(method, uri, "添加访问权限")
	function.SetInputJsonExample(&model.SvnPermissionArgumentEdit{})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Svn) ModPermission(ctx gtype.Context, ps gtype.Params) {
	argument := &model.SvnPermissionArgumentEdit{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Repository) < 1 {
		ctx.Error(gtype.ErrInput, "存储库名称(repository)为空")
		return
	}
	if len(argument.Path) < 1 {
		ctx.Error(gtype.ErrInput, "路径(path)为空")
		return
	}
	if len(argument.AccountId) < 1 {
		ctx.Error(gtype.ErrInput, "账号ID(accountId)为空")
		return
	}

	svn := &assist.Svn{}
	err = svn.SetPermission(argument.Repository, argument.Path, argument.AccountId, argument.AccessLevel)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(nil)
}

func (s *Svn) ModPermissionDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "存储库")
	function := catalog.AddFunction(method, uri, "修改访问权限")
	function.SetInputJsonExample(&model.SvnPermissionArgumentEdit{})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Svn) DelPermission(ctx gtype.Context, ps gtype.Params) {
	argument := &model.SvnPermissionArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Repository) < 1 {
		ctx.Error(gtype.ErrInput, "存储库名称(repository)为空")
		return
	}
	if len(argument.Path) < 1 {
		ctx.Error(gtype.ErrInput, "路径(path)为空")
		return
	}
	if len(argument.AccountId) < 1 {
		ctx.Error(gtype.ErrInput, "账号ID(accountId)为空")
		return
	}

	svn := &assist.Svn{}
	err = svn.RemovePermission(argument.Repository, argument.Path, argument.AccountId)
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(nil)
}

func (s *Svn) DelPermissionDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, "存储库")
	function := catalog.AddFunction(method, uri, "删除访问权限")
	function.SetInputJsonExample(&model.SvnPermissionArgument{})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Svn) createCatalog(doc gtype.Doc, names ...string) gtype.Catalog {
	root := s.createRootCatalog(doc, "SVN")
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
