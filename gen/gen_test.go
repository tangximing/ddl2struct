package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	_testDir = "../test"
)

func TestGenerate(t *testing.T) {
	sqlFile := _testDir + "/table.sql"
	outDir := _testDir + "/model"
	packageName := "model"
	err := Generate(sqlFile, outDir, packageName)
	assert.Nil(t, err)
}
