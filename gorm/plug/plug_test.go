package plug

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/sshelll/sinfra/util"
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
	plugs := []Plug{
		&TimePlug{},
	}
	for _, p := range plugs {
		p.Finalize(db)
	}
}

func (s *PlugTestSuite) AfterTest(_, _ string) {
	// log stmt
	stmt := s.stmt
	s.T().Log(s.db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...))
}

type TestModel struct {
	CreatedAtX time.Time `gorm:"column:created_at_x"`
	CreatedAtY time.Time `gorm:"column:created_at_y" sinfra_gorm_time_plug:""`
	CreatedAtZ time.Time `gorm:"column:created_at_z" sinfra_gorm_time_plug:""`
	CryptoData []byte    `gorm:"column:crypto_data" sinfra_gorm_crypto_plug:""`
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

func (s *PlugTestSuite) TestCryptoPlug() {
	plug := NewCryptoPlug().WithFields("CryptoData").
		WithEncryptFn(func(ctx context.Context, b []byte) ([]byte, error) {
			cipher := ctx.Value("cipher")
			if cipher == nil {
				panic("cipher is nil")
			}
			encrypted, err := util.AesEncrypt(b, cipher.([]byte))
			if err != nil {
				return nil, err
			}
			encoded := base64.StdEncoding.EncodeToString(encrypted)
			return []byte(encoded), nil
		})
	err := s.db.Use(plug)
	ctx := context.WithValue(context.Background(), "cipher", []byte("0123456789abcdef"))
	db := s.db.WithContext(ctx)
	s.Nil(err)
	stmt := db.Create(&TestModel{
		CryptoData: []byte("hello"),
	}).Statement
	s.stmt = stmt
}
