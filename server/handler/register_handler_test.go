package handler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"basic_server/db"
	"basic_server/internal/config"
	"basic_server/internal/repository"
	"basic_server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
	"gorm.io/gorm"
)

const dbTableNameToClean = "users"

var cleaner = dbcleaner.New()
var connection *gorm.DB

type TestRegisterUserSuite struct {
	suite.Suite
}

func (suite *TestRegisterUserSuite) SetupSuite() {
	if err := godotenv.Load("../../.env.testing"); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	cleaner.SetEngine(engine.NewMySQLEngine(dsn))
}

func (suite *TestRegisterUserSuite) SetupTest() {
	connection = db.InitDB(config.NewConfig().DB)
	cleaner.Acquire(dbTableNameToClean)
}

func (suite *TestRegisterUserSuite) TearDownTest() {
	cleaner.Clean(dbTableNameToClean)
	storage, err := connection.DB()
	if err != nil {
		log.Fatalln(err)
	}
	if err := storage.Close(); err != nil {
		log.Print(err)
	}
}

func (suite *TestRegisterUserSuite) TestRegisterUser() {
	gin.SetMode(gin.TestMode)

	server := gin.New()
	server.POST(
		"/users",
		NewRegisterHandler().RegisterUser(service.NewUserService(repository.NewUserRepository(connection))),
	)

	recorder := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/users",
		bytes.NewBuffer([]byte(`{"email":"test@test.com","password":"test"}`)),
	)
	req.Header.Add("Content-Type", "application/json")

	server.ServeHTTP(recorder, req)

	suite.Equal(http.StatusOK, recorder.Code)
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(TestRegisterUserSuite))
}
