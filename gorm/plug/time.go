package plug

import (
	"errors"
	"reflect"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type TimePlug struct {
	BasePlug
	timeFn func() time.Time
}

// NewTimePlug creates a new time plug, name is optional.
// If you want to call db.Use() twice, you must set a unique name,
// otherwise the second call will cause a error by gorm.
func NewTimePlug(name ...string) *TimePlug {
	p := &TimePlug{
		timeFn:   time.Now,
		BasePlug: *NewBasePlug(name...),
	}
	return p
}

func (p *TimePlug) WithTag(tag string) *TimePlug {
	p.BasePlug.WithTag(tag)
	return p
}

func (p *TimePlug) WithFields(fields ...string) *TimePlug {
	p.BasePlug.WithFields(fields...)
	return p
}

func (p *TimePlug) WithTimeFn(fn func() time.Time) *TimePlug {
	p.timeFn = fn
	return p
}

func (p *TimePlug) Name() string {
	if p.name != nil {
		return "sinfra:time_plug:" + *p.name
	}
	return "sinfra:time_plug"
}

func (p *TimePlug) Initialize(db *gorm.DB) error {
	if p.tag == nil && len(p.fields) == 0 {
		return errors.New("tag or field must be set")
	}
	callback := func(db *gorm.DB) {
		p.WalkFields(db, func(db *gorm.DB, f *schema.Field) {
			rv := f.ReflectValueOf(db.Statement.Context, db.Statement.ReflectValue)
			rv.Set(reflect.ValueOf(p.timeFn()))
		})
	}
	return db.Callback().Create().Before("gorm:create").Register(p.Name(), callback)
}

func (p *TimePlug) Finalize(db *gorm.DB) error {
	delete(db.Config.Plugins, p.Name())
	return db.Callback().Create().Before("gorm:create").Remove(p.Name())
}
