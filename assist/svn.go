package assist

import (
	"bytes"
	"fmt"
	"github.com/csby/gwin/model"
	"io"
	"strconv"
	"strings"
)

// Svn
// https://www.visualsvn.com/support/topic/00088/
type Svn struct {
	base
}

func (s *Svn) GetRepositories(folder bool) ([]*model.SvnRepositoryItem, error) {
	output, err := s.runCmd("Get-SvnRepository", "|", "Select", "Name,Revisions,URL")
	if err != nil {
		return nil, err
	}

	valid := false
	items := make([]*model.SvnRepositoryItem, 0)
	reader := &bytes.Buffer{}
	reader.Write(output)
	for {
		l, e := reader.ReadString('\n')
		if e == io.EOF {
			break
		}
		if len(l) < 5 {
			continue
		}
		if l[0] == '-' {
			valid = true
			continue
		}
		if !valid {
			continue
		}

		item := s.getRepository(l)
		if item != nil {
			items = append(items, item)
		}
	}

	if folder {
		c := len(items)
		for i := 0; i < c; i++ {
			item := items[i]
			err = s.getRepositoryFolders(item, folder)
			if err != nil {
				return nil, err
			}
		}
	}

	return items, nil
}

func (s *Svn) NewRepository(repository string) error {
	if len(repository) < 1 {
		return fmt.Errorf("repository is empty")
	}

	_, err := s.runCmd("New-SvnRepository", repository)
	if err != nil {
		return err
	}

	_, err = s.runCmd("New-SvnRepositoryItem", repository, "-Path", "/branches,/tags,/trunk", "-Type", "Folder")
	if err != nil {
		return err
	}

	return nil
}

func (s *Svn) GetRepositoryFolders(repository, path string, recursive bool) ([]*model.SvnRepositoryItem, error) {
	parent := &model.SvnRepositoryItem{
		Repository: repository,
		Path:       path,
	}
	err := s.getRepositoryFolders(parent, recursive)
	if err != nil {
		return nil, err
	}

	return parent.Children, nil
}

func (s *Svn) GetPermissions(repository, path string) ([]*model.SvnPermission, error) {
	if len(repository) < 1 {
		return nil, fmt.Errorf("repository is empty")
	}
	if len(path) < 1 {
		return nil, fmt.Errorf("path is empty")
	}
	output, err := s.runCmd("Select-SvnAccessRule", repository, "-Path", path, "|", "Select", "Path,Access,AccountId,AccountName")
	if err != nil {
		return nil, err
	}

	valid := false
	items := make([]*model.SvnPermission, 0)
	reader := &bytes.Buffer{}
	reader.Write(output)
	for {
		l, e := reader.ReadString('\n')
		if e == io.EOF {
			break
		}
		if len(l) < 5 {
			continue
		}
		if l[0] == '-' {
			valid = true
			continue
		}
		if !valid {
			continue
		}

		item := s.getPermission(path, l)
		if item != nil {
			items = append(items, item)
		}
	}

	return items, nil
}

func (s *Svn) GetUserPermissions(accountId string) ([]*model.SvnPermissionUser, error) {
	if len(accountId) < 1 {
		return nil, fmt.Errorf("account id is empty")
	}
	output, err := s.runCmd("Get-SvnAccessRule", "-AccountId", accountId, "|", "Select", "Repository,Path,Access")
	if err != nil {
		return nil, err
	}

	valid := false
	items := make([]*model.SvnPermissionUser, 0)
	reader := &bytes.Buffer{}
	reader.Write(output)
	for {
		l, e := reader.ReadString('\n')
		if e == io.EOF {
			break
		}
		if len(l) < 5 {
			continue
		}
		if l[0] == '-' {
			valid = true
			continue
		}
		if !valid {
			continue
		}

		item := s.getUserPermission(l)
		if item != nil {
			items = append(items, item)
		}
	}

	return items, nil
}

func (s *Svn) AddPermission(repository, path, accountId string, accessLevel int) error {
	if len(repository) < 1 {
		return fmt.Errorf("repository is empty")
	}
	if len(path) < 1 {
		return fmt.Errorf("path is empty")
	}
	if len(accountId) < 1 {
		return fmt.Errorf("accountId is empty")
	}
	level := "NoAccess"
	if accessLevel == model.SvnPermissionReadOnly {
		level = "ReadOnly"
	} else if accessLevel == model.SvnPermissionReadWrite {
		level = "ReadWrite"
	}

	_, err := s.runCmd("Add-SvnAccessRule",
		repository, "-Path", path, "-AccountId", accountId, "-Access", level)
	if err != nil {
		return err
	}

	return nil
}

func (s *Svn) SetPermission(repository, path, accountId string, accessLevel int) error {
	if len(repository) < 1 {
		return fmt.Errorf("repository is empty")
	}
	if len(path) < 1 {
		return fmt.Errorf("path is empty")
	}
	if len(accountId) < 1 {
		return fmt.Errorf("accountId is empty")
	}
	level := "NoAccess"
	if accessLevel == model.SvnPermissionReadOnly {
		level = "ReadOnly"
	} else if accessLevel == model.SvnPermissionReadWrite {
		level = "ReadWrite"
	}

	_, err := s.runCmd("Set-SvnAccessRule",
		repository, "-Path", path, "-AccountId", accountId, "-Access", level)
	if err != nil {
		return err
	}

	return nil
}

func (s *Svn) RemovePermission(repository, path, accountId string) error {
	if len(repository) < 1 {
		return fmt.Errorf("repository is empty")
	}
	if len(path) < 1 {
		return fmt.Errorf("path is empty")
	}
	if len(accountId) < 1 {
		return fmt.Errorf("accountId is empty")
	}

	_, err := s.runCmd("Remove-SvnAccessRule",
		repository, "-Path", path, "-AccountId", accountId, "-Confirm:$false")
	if err != nil {
		return err
	}

	return nil
}

func (s *Svn) getRepository(line string) *model.SvnRepositoryItem {
	// Name Revisions URL
	// ---- --------- ---
	// dev          0 https://svn.example.com/svn/dev/
	if len(line) < 5 {
		return nil
	}

	fields := s.getFields(line, " ")
	if len(fields) < 3 {
		return nil
	}
	r, e := strconv.Atoi(fields[1])
	if e != nil {
		return nil
	}

	item := &model.SvnRepositoryItem{
		Repository: fields[0],
		Name:       fields[0],
		Path:       "/",
		Type:       model.SvnRepositoryItemKindRepository,
		Url:        fields[2],
		Revisions:  r,
		Children:   make([]*model.SvnRepositoryItem, 0),
	}
	item.Id = s.uniqueId(item.Repository, item.Path)

	return item
}

func (s *Svn) getRepositoryFolders(parent *model.SvnRepositoryItem, recursive bool) error {
	if parent == nil {
		return fmt.Errorf("parent is nil")
	}
	if len(parent.Repository) < 1 {
		return fmt.Errorf("repository is empty")
	}
	if len(parent.Path) < 1 {
		return fmt.Errorf("path is empty")
	}
	if parent.Children == nil {
		parent.Children = make([]*model.SvnRepositoryItem, 0)
	}

	output, err := s.runCmd("Get-SvnRepositoryItem",
		parent.Repository, parent.Path, "-Type", "Folder", "|", "Select", "Repository,Name,Path,Url")
	if err != nil {
		return err
	}

	valid := false
	reader := &bytes.Buffer{}
	reader.Write(output)
	for {
		l, e := reader.ReadString('\n')
		if e == io.EOF {
			break
		}
		if len(l) < 5 {
			continue
		}
		if l[0] == '-' {
			valid = true
			continue
		}
		if !valid {
			continue
		}
		item := s.getRepositoryFolder(parent.Repository, l)
		if item != nil {
			parent.Children = append(parent.Children, item)
		}
	}

	if recursive {
		c := len(parent.Children)
		for i := 0; i < c; i++ {
			item := parent.Children[i]
			err = s.getRepositoryFolders(item, recursive)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Svn) getRepositoryFolder(repository, line string) *model.SvnRepositoryItem {
	// Repository Name     Path      Url
	// ---------- ----     ----      ---
	// test       trunk    /trunk    https://svn.example.com/svn/test/trunk
	if len(line) < 7 {
		return nil
	}

	fields := s.getFields(line, " ")
	if len(fields) != 4 {
		return nil
	}
	if strings.ToLower(repository) != strings.ToLower(fields[0]) {
		return nil
	}

	item := &model.SvnRepositoryItem{
		Repository: fields[0],
		Name:       fields[1],
		Path:       fields[2],
		Type:       model.SvnRepositoryItemKindFolder,
		Url:        fields[3],
		Children:   make([]*model.SvnRepositoryItem, 0),
	}
	item.Id = s.uniqueId(item.Repository, item.Path)

	return item
}

func (s *Svn) getPermission(path, line string) *model.SvnPermission {
	// Path    Access     AccountId                                     AccountName
	// ----    ------     ---------                                     -----------
	// /       ReadWrite  S-1-5-32-545                                  BUILTIN\Users
	// /trunk  NoAccess   S-1-5-21-1114322273-403004966-1807125474-1104 EXAMPLE\dev
	if len(line) < 7 {
		return nil
	}

	fields := s.getFields(line, " ")
	if len(fields) != 4 {
		return nil
	}
	item := &model.SvnPermission{
		AccountId:   fields[2],
		AccountName: fields[3],
		AccessLevel: model.SvnPermissionNoAccess,
	}

	level := strings.ToLower(fields[1])
	if level == "readonly" {
		item.AccessLevel = model.SvnPermissionReadOnly
	} else if level == "readwrite" {
		item.AccessLevel = model.SvnPermissionReadWrite
	}

	if strings.ToLower(path) != strings.ToLower(fields[0]) {
		item.Inherited = true
	}

	return item
}

func (s *Svn) getUserPermission(line string) *model.SvnPermissionUser {
	// Repository Path      Access
	// ---------- ----      ------
	// *          /         ReadWrite
	// prod       /trunk    ReadOnly
	// test       /tags     NoAccess
	if len(line) < 7 {
		return nil
	}

	fields := s.getFields(line, " ")
	if len(fields) != 3 {
		return nil
	}
	item := &model.SvnPermissionUser{
		Repository:  fields[0],
		Path:        fields[1],
		AccessLevel: model.SvnPermissionNoAccess,
	}

	level := strings.ToLower(fields[2])
	if level == "readonly" {
		item.AccessLevel = model.SvnPermissionReadOnly
	} else if level == "readwrite" {
		item.AccessLevel = model.SvnPermissionReadWrite
	}

	return item
}

func (s *Svn) runCmd(arg ...string) ([]byte, error) {
	return s.runShell(arg...)
}
