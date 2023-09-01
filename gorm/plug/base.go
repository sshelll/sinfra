package plug

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Plug interface {
	gorm.Plugin
	Finalize(*gorm.DB) error
}

type BasePlug struct {
	name   *string
	tag    *string
	fields map[string]struct{}
}

func NewBasePlug(name ...string) *BasePlug {
	p := &BasePlug{
		fields: make(map[string]struct{}),
	}
	if len(name) > 0 {
		p.name = &name[0]
	}
	return p
}

func (p *BasePlug) WithTag(tag string) *BasePlug {
	p.tag = &tag
	return p
}

func (p *BasePlug) WithFields(fields ...string) *BasePlug {
	if p.fields == nil {
		p.fields = make(map[string]struct{})
	}
	for _, v := range fields {
		p.fields[v] = struct{}{}
	}
	return p
}

// WalkFields walks all fields of the model, and call callback.
// This is a useful method for building a plug, see time.go for example.
func (p *BasePlug) WalkFields(db *gorm.DB, callback func(*gorm.DB, *schema.Field)) {
	for _, v := range db.Statement.Schema.Fields {
		var ok1, ok2 bool
		if p.tag != nil {
			_, ok1 = v.StructField.Tag.Lookup(*p.tag)
		}
		_, ok2 = p.fields[v.Name]
		if ok1 || ok2 {
			callback(db, v)
		}
	}
}

func (p *BasePlug) WalkFieldsByTag(db *gorm.DB, callback func(*gorm.DB, *schema.Field)) {
	for _, v := range db.Statement.Schema.Fields {
		_, ok := v.StructField.Tag.Lookup(*p.tag)
		if ok {
			callback(db, v)
		}
	}
}

func (p *BasePlug) WalkFieldsByName(db *gorm.DB, callback func(*gorm.DB, *schema.Field)) {
	for _, v := range db.Statement.Schema.Fields {
		if _, ok := p.fields[v.Name]; ok {
			callback(db, v)
		}
	}
}
