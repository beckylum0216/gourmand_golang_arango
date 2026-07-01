package services

import (
	"context"
	"errors"

	"github.com/arangodb/go-driver/v2/arangodb"

	"gourmand.golang.arango/src/entities"
	"gourmand.golang.arango/src/interfaces"
)

const collection = "persons"

type PersonService struct {
	db arangodb.Database
}

func NewPersonService(db arangodb.Database) interfaces.IPerson {
	return &PersonService{db: db}
}

func (s *PersonService) CreatePerson(ctx context.Context, person *entities.Person) (*entities.Person, error) {
	if person == nil {
		return nil, errors.New("person is nil")
	}

	col, err := s.db.GetCollection(ctx, collection, nil)
	if err != nil {
		return nil, err
	}

	meta, err := col.CreateDocument(ctx, person)
	if err != nil {
		return nil, err
	}

	person.Id = meta.Key
	return person, nil
}

func (s *PersonService) GetPerson(ctx context.Context, id string) (*entities.Person, error) {
	var person entities.Person

	col, err := s.db.GetCollection(ctx, collection, nil)
	if err != nil {
		return nil, err
	}

	meta, err := col.ReadDocument(ctx, id, &person)
	if err != nil {
		return nil, errors.New("person not found")
	}

	person.Id = meta.Key

	return &person, nil
}

func (s *PersonService) GetPersons(ctx context.Context) ([]*entities.Person, error) {
	query := `FOR p IN @@collection RETURN p`

	bindVars := map[string]interface{}{
		"@collection": collection,
	}

	options := &arangodb.QueryOptions{
		BindVars: bindVars,
	}

	cursor, err := s.db.Query(ctx, query, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var persons []*entities.Person
	for cursor.HasMore() {
		var person entities.Person
		_, err := cursor.ReadDocument(ctx, &person)

		if err != nil {
			return nil, err
		}
		persons = append(persons, &person)
	}

	return persons, nil
}

func (s *PersonService) UpdatePerson(ctx context.Context, id string, person *entities.Person) error {
	if person == nil {
		return errors.New("person is nil")
	}

	col, err := s.db.GetCollection(ctx, collection, nil)
	if err != nil {
		return err
	}

	patch := map[string]interface{}{
		"firstname": person.FirstName,
		"lastname":  person.LastName,
	}

	_, err = col.UpdateDocument(ctx, person.Id, patch)
	if err != nil {
		return errors.New("person not found")
	}

	return nil
}

func (s *PersonService) DeletePerson(ctx context.Context, id string) error {
	col, err := s.db.GetCollection(ctx, collection, nil)
	if err != nil {
		return err
	}

	_, err = col.DeleteDocument(ctx, id)

	return err
}
