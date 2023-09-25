package rtl

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func init() {
	MarkDownRenderer.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(addDirAuto(0), 0),
		),
	)
}

type addDirAuto int

func (t addDirAuto) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	tags := []ast.Node{}

	ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		kind := node.Kind()
		if kind == ast.KindParagraph ||
			kind == ast.KindHeading ||
			kind == ast.KindList ||
			kind == ast.KindBlockquote {
			tags = append(tags, node)
		}

		return ast.WalkContinue, nil
	})

	for _, t := range tags {
		t.SetAttributeString("dir", []byte("auto"))
	}
}
