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

type ServicePersonTest struct {
	suite.Suite
	ctx       context.Context
	container *tcarangodb.Container
	db        arangodb.Database
	service   *PersonService
}

func (s *ServicePersonTest) SetupSuite() {
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

	s.service = NewPersonService(s.db).(*PersonService)
}

func (s *ServicePersonTest) TearDownSuite() {
	if s.container != nil {
		_ = s.container.Terminate(s.ctx)
	}
}

func (s *ServicePersonTest) TestCreatePerson_Success() {
	person := &entities.Person{
		FirstName: "John",
		LastName:  "Doe",
	}

	created, err := s.service.CreatePerson(s.ctx, person)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), created.Id)

	saved, err := s.service.GetPerson(s.ctx, created.Id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "John", saved.FirstName)
	assert.Equal(s.T(), "Doe", saved.LastName)
}

func (s *ServicePersonTest) TestGetPerson_NotFound() {
	_, err := s.service.GetPerson(s.ctx, "does-not-exist")
	assert.Error(s.T(), err)
}

func (s *ServicePersonTest) TestUpdatePerson_Success() {
	person := &entities.Person{FirstName: "Jane", LastName: "Smith"}
	created, err := s.service.CreatePerson(s.ctx, person)
	assert.NoError(s.T(), err)

	created.FirstName = "Janet"
	err = s.service.UpdatePerson(s.ctx, created.Id, created)
	assert.NoError(s.T(), err)

	updated, err := s.service.GetPerson(s.ctx, created.Id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Janet", updated.FirstName)
}

func (s *ServicePersonTest) TestDeletePerson_Success() {
	person := &entities.Person{FirstName: "Temp", LastName: "Person"}
	created, err := s.service.CreatePerson(s.ctx, person)
	assert.NoError(s.T(), err)

	err = s.service.DeletePerson(s.ctx, created.Id)
	assert.NoError(s.T(), err)

	_, err = s.service.GetPerson(s.ctx, created.Id)
	assert.Error(s.T(), err)
}

func TestServicePersonSuite(t *testing.T) {
	suite.Run(t, new(ServicePersonTest))
}