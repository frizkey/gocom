package gocom

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt"
	"github.com/jarcoal/httpmock"
	"github.com/maxatome/go-testdeep/td"
	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"hash/fnv"
	"time"
)

type CommonTestSuite struct {
	suite.Suite
	MockSql      sqlmock.Sqlmock
	GormInstance *gorm.DB
	td           *td.T
}

func (suite *CommonTestSuite) AfterTest(suiteName, testName string) {
	//we make sure that all expectations were met
	if err := suite.MockSql.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("%s there were unfulfilled expectations: %s", testName, err)
	}
	httpmock.DeactivateAndReset()
}

func (suite *CommonTestSuite) SetupTest() {
	dbMock, sMock, _ := sqlmock.New()
	gormInstance, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: dbMock,
	}), &gorm.Config{})
	suite.MockSql = sMock
	suite.GormInstance = gormInstance
	suite.td = td.NewT(suite.T())
	httpmock.ActivateNonDefault(DefaultRestyClient().GetHttpClient())
}

func (suite *CommonTestSuite) SQLCheckCommit(fn func(sqlmock.Sqlmock)) *CommonTestSuite {
	mockSql := suite.MockSql
	mockSql.ExpectBegin()
	fn(mockSql)
	mockSql.ExpectCommit()
	return suite
}

func (suite *CommonTestSuite) SQLCheck(fn func(sqlmock.Sqlmock)) *CommonTestSuite {
	mockSql := suite.MockSql
	fn(mockSql)
	return suite
}

func (suite *CommonTestSuite) SQLCheckRollback(fn func(sqlmock.Sqlmock)) *CommonTestSuite {
	mockSql := suite.MockSql
	mockSql.ExpectBegin()
	fn(mockSql)
	mockSql.ExpectRollback()
	return suite
}

func (suite *CommonTestSuite) AssertTD() *td.T {
	return suite.td
}

func (suite *CommonTestSuite) UseRealConn() {
	suite.GormInstance = DB()
}

var DefaultErrTest = eris.New("ERROR")

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func GenToken(sid int, user struct {
	ID    string
	Name  string
	Email string
	Layer string
	Grade string
	Roles []struct {
		ID string
	}
}, duration string) (string, string) {

	dur, _ := time.ParseDuration(duration)
	exp := time.Now().Add(dur)

	role := ""

	for _, v := range user.Roles {

		role += "," + v.ID
	}

	if role != "" {
		role = role[1:]
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid":   sid,
		"uid":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"layer": user.Layer,
		"grade": user.Grade,
		"role":  role,
		"iat":   time.Now().Unix(),
		"exp":   exp.Unix(),
	})

	tokenString, _ := token.SignedString([]byte("SuperAppKey2021"))

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid":       sid,
		"tokenHash": hash(tokenString),
	})

	refreshTokenString, _ := refreshToken.SignedString([]byte("SuperAppKey2021"))

	return tokenString, refreshTokenString
}

func hash(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}
