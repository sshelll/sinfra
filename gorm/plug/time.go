package plug

import (
	"errors"
	"reflect"
	"time"

	"gorm.io/gorm"
)

type TimePlug struct {
	tag    *string
	fields map[string]struct{}
	timeFn func() time.Time
}

func NewTimePlug() *TimePlug {
	return &TimePlug{
		timeFn: time.Now,
		fields: make(map[string]struct{}),
	}
}

func (p *TimePlug) WithTag(tag string) *TimePlug {
	p.tag = &tag
	return p
}

func (p *TimePlug) WithFields(fields ...string) *TimePlug {
	if p.fields == nil {
		p.fields = make(map[string]struct{})
	}
	for _, v := range fields {
		p.fields[v] = struct{}{}
	}
	return p
}

func (p *TimePlug) WithTimeFn(fn func() time.Time) *TimePlug {
	p.timeFn = fn
	return p
}

func (p *TimePlug) Name() string {
	return "sinfra:time_plug"
}

func (p *TimePlug) Initialize(db *gorm.DB) error {
	if p.tag == nil && len(p.fields) == 0 {
		return errors.New("tag or field must be set")
	}
	callback := p.setTimeByField
	if p.tag != nil {
		callback = p.setTimeByTag
	}
	return db.Callback().Create().Before("gorm:create").Register(p.Name(), callback)
}

func (p *TimePlug) Finalize(db *gorm.DB) error {
	delete(db.Config.Plugins, p.Name())
	return db.Callback().Create().Before("gorm:create").Remove(p.Name())
}

func (p *TimePlug) setTimeByTag(db *gorm.DB) {
	for _, v := range db.Statement.Schema.Fields {
		_, ok := v.StructField.Tag.Lookup(*p.tag)
		if ok {
			rv := v.ReflectValueOf(db.Statement.Context, db.Statement.ReflectValue)
			rv.Set(reflect.ValueOf(p.timeFn()))
		}
	}
}

func (p *TimePlug) setTimeByField(db *gorm.DB) {
	for _, v := range db.Statement.Schema.Fields {
		_, ok := p.fields[v.Name]
		if ok {
			rv := v.ReflectValueOf(db.Statement.Context, db.Statement.ReflectValue)
			rv.Set(reflect.ValueOf(p.timeFn()))
		}
	}
}
