package autolink_pages

import (
	"context"
	"embed"
	"html/template"
	"sort"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

//go:embed templates
var templates embed.FS

type fileInfoByNameLength []Page

func (a fileInfoByNameLength) Len() int           { return len(a) }
func (a fileInfoByNameLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a fileInfoByNameLength) Less(i, j int) bool { return len(a[i].Name()) > len(a[j].Name()) }

func init() {
	Listen(AfterWrite, UpdatePagesList)
	Listen(AfterDelete, UpdatePagesList)

	RegisterWidget(AFTER_VIEW_WIDGET, 1, backlinksSection)
	RegisterTemplate(templates, "templates")
}

var autolinkPages []Page

func UpdatePagesList(_ Page) (err error) {
	ps := []Page{}
	EachPage(context.Background(), func(p Page) {
		ps = append(ps, p)
	})
	sort.Sort(fileInfoByNameLength(ps))
	autolinkPages = ps
	return
}

func countTodos(p Page) (total int, done int) {
	tasks := FindAllInAST[*east.TaskCheckBox](p.AST())
	for _, v := range tasks {
		total++
		if v.IsChecked {
			done++
		}
	}

	return
}

func backlinksSection(p Page) template.HTML {
	if p.Name() == INDEX {
		return ""
	}

	pages := []Page{}

	EachPage(context.Background(), func(a Page) {
		// a page shouldn't mention itself
		if a.Name() == p.Name() {
			return
		}

		if containLinkTo(a.AST(), p) {
			pages = append(pages, a)
		}
	})

	return Partial("backlinks", Locals{"pages": pages})
}

func containLinkTo(n ast.Node, p Page) bool {
	if n.Kind() == KindPageLink {
		t, _ := n.(*PageLink)
		if t.page.FileName() == p.FileName() {
			return true
		}
	}

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if containLinkTo(c, p) {
			return true
		}

		if c == n.LastChild() {
			break
		}
	}

	return false
}
