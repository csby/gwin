package assist

import (
	"encoding/json"
	"github.com/csby/gwin/model"
	"testing"
)

func TestSvn_GetRepositories(t *testing.T) {
	svn := &Svn{}
	items, err := svn.GetRepositories(false)
	if err != nil {
		t.Error(err)
		return
	}
	c := len(items)
	t.Log("count: ", c)
	for i := 0; i < c; i++ {
		item := items[i]
		t.Logf("%2d  %s  %s", i+1, item.Name, fmtItem(item))
	}
}

func TestSvn_GetRepositoryFolders(t *testing.T) {
	svn := &Svn{}
	items, err := svn.GetRepositoryFolders("test", "/", false)
	if err != nil {
		t.Error(err)
		return
	}
	c := len(items)
	t.Log("count: ", c)
	for i := 0; i < c; i++ {
		item := items[i]
		t.Logf("%2d  %s  %s", i+1, item.Name, fmtItem(item))
	}
}

func TestSvn_GetPermissions(t *testing.T) {
	svn := &Svn{}
	items, err := svn.GetPermissions("test", "/trunk")
	if err != nil {
		t.Error(err)
		return
	}
	c := len(items)
	t.Log("count: ", c)
	for i := 0; i < c; i++ {
		item := items[i]
		t.Logf("%2d  %s  %s", i+1, item.AccountName, fmtItem(item))
	}
}

func TestSvn_AddPermission(t *testing.T) {
	svn := &Svn{}
	repository := "prod"
	path := "/trunk"
	account := "S-1-5-21-1114322273-403004966-1807125474-1104"
	access := model.SvnPermissionReadOnly

	items, err := svn.GetPermissions(repository, path)
	if err != nil {
		t.Error(err)
		return
	}
	c := len(items)
	t.Log("count: ", c)
	for i := 0; i < c; i++ {
		item := items[i]
		t.Logf("%2d  %s  %s", i+1, item.AccountName, fmtItem(item))
	}

	err = svn.AddPermission(repository, path, account, access)
	if err != nil {
		t.Error("AddPermission fail: ", err)
		return
	}
	t.Log("AddPermission success")

	items, err = svn.GetPermissions(repository, path)
	if err != nil {
		t.Error(err)
		return
	}
	c = len(items)
	t.Log("count: ", c)
	for i := 0; i < c; i++ {
		item := items[i]
		t.Logf("%2d  %s  %s", i+1, item.AccountName, fmtItem(item))
	}

	access = model.SvnPermissionReadWrite
	err = svn.SetPermission(repository, path, account, access)
	if err != nil {
		t.Error("SetPermission fail: ", err)
		return
	}
	t.Log("SetPermission success")

	items, err = svn.GetPermissions(repository, path)
	if err != nil {
		t.Error(err)
		return
	}
	c = len(items)
	t.Log("count: ", c)
	for i := 0; i < c; i++ {
		item := items[i]
		t.Logf("%2d  %s  %s", i+1, item.AccountName, fmtItem(item))
	}

	err = svn.RemovePermission(repository, path, account)
	if err != nil {
		t.Error("RemovePermission fail: ", err)
		return
	}
	t.Log("RemovePermission success")

	items, err = svn.GetPermissions(repository, path)
	if err != nil {
		t.Error(err)
		return
	}
	c = len(items)
	t.Log("count: ", c)
	for i := 0; i < c; i++ {
		item := items[i]
		t.Logf("%2d  %s  %s", i+1, item.AccountName, fmtItem(item))
	}
}

func fmtItem(v interface{}) string {
	if v == nil {
		return ""
	}

	d, e := json.MarshalIndent(v, "", "  ")
	if e != nil {
		return ""
	}

	return string(d)
}
