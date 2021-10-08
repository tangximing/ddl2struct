package parser

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pingcap/parser/mysql"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/types"
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

func (p *Parser) ToStructs() (data map[string][]byte, err error) {
	data = make(map[string][]byte)
	for _, table := range p.ts {
		data[table.name], err = format.Source([]byte(table.toStruct()))
		if err != nil {
			err = errors.Wrapf(err, "format.Source")
			return
		}
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
		t.columns = append(t.columns, column{
			name:  col.Name.Name.String(),
			cType: getColumnType(col.Tp),
		})
	}

	p.ts = append(p.ts, t)
	return
}

func getColumnType(cType *types.FieldType) string {
	evalType := cType.EvalType()
	switch evalType {
	case types.ETInt:
		if cType.Tp == mysql.TypeLonglong {
			return "int64"
		}
		return "int"
	case types.ETReal, types.ETDecimal:
		return "float64"
	case types.ETDatetime, types.ETTimestamp:
		return "time.Time"
	default:
		return "string"
	}
}

type table struct {
	name    string
	columns columns
}

func (t *table) toStruct() string {
	return fmt.Sprintf("type %s struct { %s }", strcase.ToCamel(t.name), t.columns.toStructFields())
}

type columns []column

func (c columns) toStructFields() string {
	fields := make([]string, 0)
	for _, column := range c {
		fields = append(fields, column.toStructField())
	}

	return strings.Join(fields, "\n")
}

type column struct {
	name  string
	cType string
}

func (c column) toStructField() string {
	tag := fmt.Sprintf("`json:\"%s\" gorm:\"Column:%s\"`", strcase.ToSnake(c.name), c.name)
	return fmt.Sprintf("%s %s %s", strcase.ToCamel(c.name), c.cType, tag)
}
