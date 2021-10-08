package parser

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pingcap/parser/types"
)

const (
	enumFlagPrefix    = "枚举（"
	enumFlagSuffix    = "）"
	enumItemsSplitter = "，"
	enumItemSplitter  = "："
	commentSplitter   = "。"
)

type goStruct struct {
	name    string
	fields  []*goField
	imports map[string]struct{}
}

func (s *goStruct) ToGo() string {
	bs := new(bytes.Buffer)
	if len(s.imports) != 0 {
		bs.WriteString("import (")
		bs.WriteString("\n")
		for goImport := range s.imports {
			bs.WriteString(fmt.Sprintf("\"%s\"", goImport))
			bs.WriteString("\n")
		}
		bs.WriteString(")")
	}
	bs.WriteString("\n")

	for _, field := range s.fields {
		if field.goEnum == nil {
			continue
		}

		bs.WriteString(field.goEnum.ToGo())
		bs.WriteString("\n")
	}

	bs.WriteString(fmt.Sprintf("type %s struct {", s.name))
	bs.WriteString("\n")
	for _, field := range s.fields {
		tag := fmt.Sprintf("`gorm:\"column:%s\" json:\"%s\"`", field.gormTag, strcase.ToSnake(field.name))
		if field.comment != "" {
			bs.WriteString(fmt.Sprintf("// %s", field.comment))
			bs.WriteString("\n")
		}
		if field.goEnum == nil {
			bs.WriteString(fmt.Sprintf("%s %s %s", field.name, field.goType, tag))
		} else {
			bs.WriteString(fmt.Sprintf("%s %s %s", field.name, field.goEnum.GetName(), tag))
		}
		bs.WriteString("\n")
	}
	bs.WriteString("}")
	bs.WriteString("\n")

	return bs.String()
}

type goField struct {
	name    string
	gormTag string
	comment string
	goEnum  *goEnum
	goType  goType
}

type goEnum struct {
	namePrefix string
	goType     goType
	items      []*goEnumItem
}

type goEnumItem struct {
	key     string
	comment string
}

func (e *goEnum) ToGo() string {
	enumName := e.GetName()

	bs := new(bytes.Buffer)
	bs.WriteString(fmt.Sprintf("type %s %s", enumName, e.goType))
	bs.WriteString("\n")
	bs.WriteString("const (")
	bs.WriteString("\n")
	for _, item := range e.items {
		bs.WriteString(fmt.Sprintf("// %s - %s", item.key, item.comment))
		bs.WriteString("\n")
		bs.WriteString(fmt.Sprintf("%s%s %s = %s", e.namePrefix, strcase.ToCamel(item.key), enumName, item.key))
		bs.WriteString("\n")
	}
	bs.WriteString(")")

	return bs.String()
}

func (e *goEnum) GetName() string {
	return fmt.Sprintf("%sEnum", e.namePrefix)
}

type table struct {
	name    string
	columns columns
}

func (t *table) toGoStruct() *goStruct {
	s := &goStruct{
		name:    strcase.ToCamel(t.name),
		fields:  t.columns.toGoFields(),
		imports: make(map[string]struct{}),
	}
	for _, field := range s.fields {
		goImport := getGoImport(field.goType)
		if goImport != "" {
			s.imports[goImport] = struct{}{}
		}
	}

	return s
}

type columns []column

func (c columns) toGoFields() []*goField {
	fields := make([]*goField, 0)
	for _, column := range c {
		fields = append(fields, column.toGoField())
	}

	return fields
}

type column struct {
	name    string
	cType   *types.FieldType
	comment string
}

func (c column) toGoField() *goField {
	f := &goField{
		name:    strcase.ToCamel(c.name),
		gormTag: c.name,
		goType:  getGoType(c.cType),
		goEnum:  c.parseEnum(),
		comment: c.parseComment(),
	}

	return f
}

func (c column) parseEnum() *goEnum {
	enumIndex := strings.LastIndex(c.comment, enumFlagPrefix)
	if enumIndex == -1 {
		return nil
	}

	e := &goEnum{
		namePrefix: strcase.ToCamel(c.name),
	}

	enumStr := c.comment[enumIndex:]
	enumStr = strings.TrimPrefix(enumStr, enumFlagPrefix)
	enumStr = strings.TrimSuffix(enumStr, enumFlagSuffix)
	enumItems := strings.Split(enumStr, enumItemsSplitter)
	for _, enumItem := range enumItems {
		enumItemDetail := strings.Split(enumItem, enumItemSplitter)
		if len(enumItemDetail) != 2 {
			return nil
		}
		enumItemKey := enumItemDetail[0]
		enumItemComment := enumItemDetail[1]

		enumGoType := goTypeString
		_, err := strconv.Atoi(enumItemKey)
		if err == nil {
			enumGoType = goTypeInt8
		}
		if e.goType == "" {
			e.goType = enumGoType
		}
		if e.goType != enumGoType {
			return nil
		}

		e.items = append(e.items, &goEnumItem{
			key:     enumItemKey,
			comment: enumItemComment,
		})
	}

	return e
}

func (c column) parseComment() string {
	enumIndex := strings.LastIndex(c.comment, enumFlagPrefix)
	if enumIndex == -1 {
		return c.comment
	}

	comment := strings.TrimSpace(c.comment[:enumIndex])
	comment = strings.TrimSuffix(comment, commentSplitter)
	return comment
}
