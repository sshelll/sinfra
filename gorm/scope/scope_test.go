package scope

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type PaginateTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func TestPaginateTestSuite(t *testing.T) {
	suite.Run(t, new(PaginateTestSuite))
}

func (s *PaginateTestSuite) SetupSuite() {
	dsn := "root:123456@tcp(localhost:3309)/"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DryRun: true})
	s.Require().NoError(err)
	// db.Exec("DROP DATABASE IF EXISTS scope_test")
	// db.Exec("CREATE DATABASE scope_test")
	// db.Exec("USE scope_test")
	// db.AutoMigrate(&TestModel{})
	s.db = db
}

type TestModel struct {
	gorm.Model
}

func (*TestModel) TableName() string {
	return "test"
}

func (s *PaginateTestSuite) TestPaginate() {
	stmt := s.db.
		Scopes(Paginate(1, 10)).
		Where("id > ?", 1).
		Find(&TestModel{}).
		Statement
	s.T().Log(s.db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...))
}

func (s *PaginateTestSuite) TestPaginateByCol() {
	stmt := s.db.
		Scopes(PaginateByCol("id", 1, 10, false)).
		Find(&TestModel{}).
		Statement
	s.T().Log(s.db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...))
}

func (s *PaginateTestSuite) TestChore() {
	s.db.Callback().Create().Before("gorm:create").Register("scope:chore", func(db *gorm.DB) {
		for _, v := range db.Statement.Schema.Fields {
			fv, ok := v.ValueOf(s.db.Statement.Context, db.Statement.ReflectValue)
			s.T().Log(v.Name, fv, ok)
			if v.Name == "CreatedAt" {
				rv := v.ReflectValueOf(s.db.Statement.Context, db.Statement.ReflectValue)
				rv.Set(reflect.ValueOf(time.Now().Add(-time.Hour * 24)))
			}
		}
	})
	stmt := s.db.Create(&TestModel{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
		},
	}).Statement
	s.T().Log(s.db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...))
}
