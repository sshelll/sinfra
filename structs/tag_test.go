package structs

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type base struct {
	ID string `json:"id" gorm:"primary_key"`
}

func TestExtractStructTags(t *testing.T) {

	s := struct {
		b        *base
		Name     string `json:"name,omitempty" gorm:"column:name"`
		Age      int    `json:"age,omitempty" gorm:"column:age"`
		Birthday int64  `json:"birthday,omitempty" gorm:"column:birthday"`
	}{}

	tagList := ExtractStructTags(reflect.TypeOf(s), "json", func(s string) string {
		return strings.Split(s, ",")[0]
	})

	assert.Equal(t, 3, len(tagList))
	assert.Equal(t, "name", tagList[0])
	assert.Equal(t, "age", tagList[1])
	assert.Equal(t, "birthday", tagList[2])

}

func TestExtractStructTagsWithEmbeddedField(t *testing.T) {

	s := struct {
		*base
		Name     string `json:"name,omitempty"`
		Age      int    `json:"age,omitempty"`
		Birthday int64  `json:"birthday,omitempty"`
	}{}

	tagList := ExtractStructTags(reflect.TypeOf(s), "json", func(s string) string {
		return strings.Split(s, ",")[0]
	})

	assert.Equal(t, 4, len(tagList))
	assert.Equal(t, "id", tagList[0])
	assert.Equal(t, "name", tagList[1])
	assert.Equal(t, "age", tagList[2])
	assert.Equal(t, "birthday", tagList[3])

}

func TestExtractGormColumnName(t *testing.T) {

	s := struct {
		*base
		Name     string `json:"name,omitempty" gorm:"column:name"`
		Age      int    `json:"age,omitempty" gorm:"column:age"`
		Birthday int64  `json:"birthday,omitempty" gorm:"column:birthday"`
	}{}

	colNames := ExtractGormColumnName(reflect.TypeOf(s))

	assert.Equal(t, 4, len(colNames))
	assert.Equal(t, "id", colNames[0])
	assert.Equal(t, "name", colNames[1])
	assert.Equal(t, "age", colNames[2])
	assert.Equal(t, "birthday", colNames[3])

}

func TestExtractJsonKey(t *testing.T) {

	s := struct {
		*base
		Name     string `json:"name,omitempty" gorm:"column:name"`
		Age      int    `json:"age,omitempty" gorm:"column:age"`
		Birthday int64  `json:"birthday,omitempty" gorm:"column:birthday"`
	}{}

	jsonKeys := ExtractJsonKey(reflect.TypeOf(s))
	assert.Equal(t, 4, len(jsonKeys))
	assert.Equal(t, "id", jsonKeys[0])
	assert.Equal(t, "name", jsonKeys[1])
	assert.Equal(t, "age", jsonKeys[2])
	assert.Equal(t, "birthday", jsonKeys[3])

}
