package controller

import (
	"github.com/csby/gwin/config"
	"github.com/csby/gwsf/gtype"
)

type base struct {
	gtype.Base

	cfg *config.Config
}

func (s *base) createRootCatalog(doc gtype.Doc, names ...string) gtype.Catalog {
	root := doc.AddCatalog("API")

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
