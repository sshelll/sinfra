package plug

import (
	"context"
	"errors"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// CryptoPlug is a plug for encrypt / decrypt data.
// NOTE: The encryptFn and decryptFn must be set, and the target column must be []byte type.
type CryptoPlug struct {
	BasePlug
	encryptFn func(context.Context, []byte) ([]byte, error)
	decryptFn func(context.Context, []byte) ([]byte, error)
}

func NewCryptoPlug(name ...string) *CryptoPlug {
	p := &CryptoPlug{
		BasePlug: *NewBasePlug(name...),
	}
	return p
}

func (p *CryptoPlug) WithTag(tag string) *CryptoPlug {
	p.BasePlug.WithTag(tag)
	return p
}

func (p *CryptoPlug) WithFields(fields ...string) *CryptoPlug {
	p.BasePlug.WithFields(fields...)
	return p
}

func (p *CryptoPlug) WithEncryptFn(fn func(context.Context, []byte) ([]byte, error)) *CryptoPlug {
	p.encryptFn = fn
	return p
}

func (p *CryptoPlug) WithDecryptFn(fn func(context.Context, []byte) ([]byte, error)) *CryptoPlug {
	p.decryptFn = fn
	return p
}

func (p *CryptoPlug) Name() string {
	if p.name != nil {
		return "sinfra:crypto_plug:" + *p.name
	}
	return "sinfra:crypto_plug"
}

func (p *CryptoPlug) Initialize(db *gorm.DB) error {
	if p.tag == nil && len(p.fields) == 0 {
		return errors.New("tag or field must be set")
	}

	if p.encryptFn == nil && p.decryptFn == nil {
		return errors.New("encryptFn and decryptFn must be set")
	}

	encryptCallback := func(db *gorm.DB) {
		p.WalkFields(db, func(db *gorm.DB, f *schema.Field) {
			v, zero := f.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
			if zero {
				return
			}
			bytes, ok := v.([]byte)
			if !ok {
				return
			}
			encrypted, err := p.encryptFn(db.Statement.Context, bytes)
			if err != nil {
				return
			}
			rv := f.ReflectValueOf(db.Statement.Context, db.Statement.ReflectValue)
			rv.Set(reflect.ValueOf(encrypted))
		})
	}

	decryptCallback := func(db *gorm.DB) {
		p.WalkFields(db, func(db *gorm.DB, f *schema.Field) {
			v, ok := f.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
			if !ok {
				return
			}
			bytes, ok := v.([]byte)
			if !ok {
				return
			}
			decrypted, err := p.decryptFn(db.Statement.Context, bytes)
			if err != nil {
				return
			}
			rv := f.ReflectValueOf(db.Statement.Context, db.Statement.ReflectValue)
			rv.Set(reflect.ValueOf(decrypted))
		})
	}

	if p.encryptFn != nil {
		if err := db.Callback().Create().Before("gorm:create").
			Register(p.Name(), encryptCallback); err != nil {
			return err
		}
	}

	if p.decryptFn != nil {
		return db.Callback().Query().After("gorm:query").Register(p.Name(), decryptCallback)
	}

	return nil
}

func (p *CryptoPlug) Finalize(db *gorm.DB) error {
	delete(db.Config.Plugins, p.Name())
	if err := db.Callback().Create().Before("gorm:create").
		Remove(p.Name()); err != nil {
		return err
	}
	return db.Callback().Query().After("gorm:query").Remove(p.Name())
}
