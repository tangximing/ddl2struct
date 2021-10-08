package parser

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	driver "github.com/pingcap/tidb/types/parser_driver"
	"github.com/pkg/errors"
)

type Parser struct {
	p *parser.Parser

	ts  []table
	err error
}

func New() *Parser {
	return &Parser{
		p: parser.New(),
	}
}

func (p *Parser) Parse(sql string) (err error) {
	nodes, _, err := p.p.Parse(sql, "", "")
	if err != nil {
		err = errors.Wrap(err, "p.p.Parse")
		return
	}

	for _, node := range nodes {
		node.Accept(p)
		if p.err != nil {
			err = errors.Wrapf(err, "node.Accept")
		}
	}

	return
}

func (p *Parser) ToStructs() (data map[string]*goStruct) {
	data = make(map[string]*goStruct)
	for _, table := range p.ts {
		goStruct := table.toGoStruct()
		data[goStruct.name] = goStruct
	}

	return
}

func (p *Parser) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	switch n := n.(type) {
	case *ast.CreateTableStmt:
		p.err = p.parseCreateTableStmt(n)
	}
	return n, true
}

func (p *Parser) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}

func (p *Parser) parseCreateTableStmt(stmt *ast.CreateTableStmt) (err error) {
	t := table{}
	t.name = stmt.Table.Name.String()
	for _, col := range stmt.Cols {
		c := column{
			name:  col.Name.Name.String(),
			cType: col.Tp,
		}
		for _, option := range col.Options {
			if option.Tp != ast.ColumnOptionComment {
				continue
			}

			exprNode, ok := option.Expr.(*driver.ValueExpr)
			if !ok {
				continue
			}

			c.comment = exprNode.GetDatumString()
			break
		}

		t.columns = append(t.columns, c)
	}

	p.ts = append(p.ts, t)
	return
}
