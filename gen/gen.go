package gen

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/pkg/errors"

	"github.com/tangximing/ddl2struct/parser"
	"github.com/tangximing/ddl2struct/util"
)

func Generate(sqlFile string, dir string, packageName string) (err error) {
	sqlFilePath, err := filepath.Abs(sqlFile)
	if err != nil {
		err = errors.Wrapf(err, "filepath.Abs")
		return
	}
	if !util.IsExist(sqlFile) {
		err = fmt.Errorf("sql file(%s) does not exist", sqlFile)
		return
	}

	dirPath, err := filepath.Abs(dir)
	if err != nil {
		err = errors.Wrapf(err, "filepath.Abs")
		return
	}
	if !util.IsExist(dirPath) {
		err = os.Mkdir(dirPath, 0755)
		if err != nil {
			err = errors.Wrapf(err, "os.Mkdir dir(%s)", dirPath)
			return
		}
	}

	if err = generate(sqlFilePath, dirPath, packageName); err != nil {
		err = errors.Wrapf(err, "generate")
		return
	}

	return
}

func generate(absSqlFile string, absDir string, packageName string) (err error) {
	sql, err := ioutil.ReadFile(absSqlFile)
	if err != nil {
		err = errors.Wrapf(err, "ioutil.ReadFile")
		return
	}

	p := parser.New()
	if err = p.Parse(string(sql)); err != nil {
		err = errors.Wrapf(err, "p.Parse")
		return
	}

	structs, err := p.ToStructs()
	if err != nil {
		err = errors.Wrapf(err, "p.ToStructs")
		return
	}
	for tableName, structData := range structs {
		goData := fmt.Sprintf("package %s\n%s", packageName, string(structData))
		goDataBytes, e := format.Source([]byte(goData))
		if e != nil {
			err = errors.Wrapf(e, "format.Source")
			return
		}

		goFileName := fmt.Sprintf("%s.go", strcase.ToSnake(tableName))
		goFile := filepath.Join(absDir, goFileName)
		e = ioutil.WriteFile(goFile, goDataBytes, 0644)
		if e != nil {
			err = errors.Wrapf(e, "ioutil.WriteFile")
			return
		}
	}

	return
}
