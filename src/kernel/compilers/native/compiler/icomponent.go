package compiler

import (
	"go/ast"
	"go/token"
)

type IComponent interface {
	Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error

	Prepare() error

	Compile() error
}
