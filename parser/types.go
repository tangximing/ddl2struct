package parser

import (
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/types"
)

type goType string

const (
	goTypeInt8    goType = "int8"
	goTypeInt     goType = "int"
	goTypeInt64   goType = "int64"
	goTypeFloat64 goType = "float64"
	goTypeTime    goType = "time.Time"
	goTypeString  goType = "string"
)

var (
	typeToImport = map[goType]string{
		goTypeTime: "time",
	}
)

func getGoType(cType *types.FieldType) goType {
	evalType := cType.EvalType()
	switch evalType {
	case types.ETInt:
		if cType.Tp == mysql.TypeLonglong {
			return goTypeInt64
		}
		return goTypeInt
	case types.ETReal, types.ETDecimal:
		return goTypeFloat64
	case types.ETDatetime, types.ETTimestamp:
		return goTypeTime
	default:
		return goTypeString
	}
}

func getGoImport(t goType) string {
	return typeToImport[t]
}
