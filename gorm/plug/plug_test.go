package plug

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type PlugTestSuite struct {
	suite.Suite
	db   *gorm.DB
	stmt *gorm.Statement
}

func TestPlugTestSuite(t *testing.T) {
	suite.Run(t, new(PlugTestSuite))
}

func (s *PlugTestSuite) SetupSuite() {
	dsn := "root:123456@tcp(localhost:3309)/"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{DryRun: true})
	s.Require().NoError(err)
	s.db = db
}

func (s *PlugTestSuite) BeforeTest(_, _ string) {
	db := s.db
	timePlug := &TimePlug{}
	timePlug.Finalize(db)
}

func (s *PlugTestSuite) AfterTest(_, _ string) {
	// log stmt
	stmt := s.stmt
	s.T().Log(s.db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...))
}

type TestModel struct {
	CreatedAtX time.Time `gorm:"column:created_at_x"`
	CreatedAtY time.Time `gorm:"column:created_at_y" sinfra_gorm_time_plug:""`
}

func (*TestModel) TableName() string {
	return "test"
}

func (s *PlugTestSuite) TestTimePlugByField() {
	err := s.db.Use(NewTimePlug().WithFields("CreatedAtX"))
	s.Nil(err)
	stmt := s.db.Create(&TestModel{}).Statement
	s.stmt = stmt
}

func (s *PlugTestSuite) TestTimePlugByTag() {
	err := s.db.Use(NewTimePlug().WithTag("sinfra_gorm_time_plug"))
	s.Nil(err)
	stmt := s.db.Create(&TestModel{}).Statement
	s.stmt = stmt
}

func (s *PlugTestSuite) TestTimePlugWithCustomizedTimer() {
	err := s.db.Use(NewTimePlug().
		WithFields("CreatedAtX", "CreatedAtY").
		WithTimeFn(func() time.Time {
			return time.Date(1949, 10, 1, 0, 0, 0, 0, time.UTC)
		}),
	)
	s.Nil(err)
	stmt := s.db.Create(&TestModel{}).Statement
	s.stmt = stmt
}
