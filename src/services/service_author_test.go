package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	tcarangodb "github.com/testcontainers/testcontainers-go/modules/arangodb"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/connection"

	"gourmand.golang.arango/src/entities"
)

type ServiceAuthorTest struct {
	suite.Suite
	ctx       context.Context
	container *tcarangodb.Container
	db        arangodb.Database
	service   *AuthorService
}

func (s *ServiceAuthorTest) SetupSuite() {
	s.ctx = context.Background()

	const password = "t3stc0ntain3rs!"
	container, err := tcarangodb.Run(s.ctx, "arangodb:3.11.5",
		tcarangodb.WithRootPassword(password))
	assert.NoError(s.T(), err)
	s.container = container

	endpoint, err := container.HTTPEndpoint(s.ctx)
	assert.NoError(s.T(), err)

	conn := connection.NewHttp2Connection(
		connection.DefaultHTTP2ConfigurationWrapper(
			connection.NewRoundRobinEndpoints([]string{endpoint}), true,
		),
	)
	auth := connection.NewBasicAuth("root", password)
	err = conn.SetAuthentication(auth)
	assert.NoError(s.T(), err)

	client := arangodb.NewClient(conn)

	db, err := client.CreateDatabase(s.ctx, "gourmand_test", nil)
	assert.NoError(s.T(), err)
	s.db = db

	_, err = db.CreateCollectionV2(s.ctx, "persons", nil)
	assert.NoError(s.T(), err)
	_, err = db.CreateCollectionV2(s.ctx, "authors", nil)
	assert.NoError(s.T(), err)
	_, err = db.CreateCollectionV2(s.ctx, 
		"person_is_author", nil)

	assert.NoError(s.T(), err)

	s.service = NewAuthorService(s.db).(*AuthorService)
}

func (s *ServiceAuthorTest) TearDownSuite() {
	if s.container != nil {
		_ = s.container.Terminate(s.ctx)
	}
}

func (s *ServiceAuthorTest) TestCreateAuthor_NewPerson_Success() {
	person := &entities.Person{FirstName: "Jane", LastName: "Doe"}
	author := &entities.Author{Source: "blog", CreditString: "Jane D."}

	err := s.service.CreateAuthor(s.ctx, person, author)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), author.Id)

	saved, err := s.service.GetAuthor(s.ctx, author.Id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "blog", saved.Source)
	assert.Equal(s.T(), "Jane D.", saved.CreditString)
}

func (s *ServiceAuthorTest) TestCreateAuthor_NilAuthor_Errors() {
	person := &entities.Person{FirstName: "X", LastName: "Y"}
	err := s.service.CreateAuthor(s.ctx, person, nil)
	assert.Error(s.T(), err)
}

func TestServiceAuthorSuite(t *testing.T) {
	suite.Run(t, new(ServiceAuthorTest))
}