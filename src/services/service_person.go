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

func (s *PersonService) GetPersonDetails(ctx context.Context, key string) (*entities.PersonWithDetails, error) {
	// 1. Fetch person
	colPersons, err := s.db.GetCollection(ctx, "persons", nil)
	if err != nil {
		return nil, err
	}

	var person entities.Person
	id := "persons/" + key

	meta, err := colPersons.ReadDocument(ctx, id, &person)

	if err != nil {
		return nil, errors.New("person not found")
	}

	person.Id = meta.Key

	// 2. Fetch user via edge
	var user entities.User
	queryUser := `
	FOR e IN persons_users
		FILTER e._from == @from
		FOR u IN users
			FILTER u._id == e._to
			RETURN u`
	bindVars := map[string]interface{}{
		"from": "persons/" + key,
	}

	cursor, err := s.db.Query(ctx, queryUser, &arangodb.QueryOptions{BindVars: bindVars})
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for cursor.HasMore() {
		meta, err := cursor.ReadDocument(ctx, &user)

		user.Id = meta.Key

		if err != nil {
			return nil, err
		}
	}

	// 3. Fetch author via edge
	var author entities.Author
	queryAuthor := `
        FOR e IN persons_authors
            FILTER e._from == @from
            FOR a IN authors
                FILTER a._id == e._to
                RETURN a
    `
	cursor, err = s.db.Query(ctx, queryAuthor, &arangodb.QueryOptions{BindVars: bindVars})

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	if cursor.HasMore() {
		meta, err := cursor.ReadDocument(ctx, &author)
		author.Id = meta.Key

		if err != nil {
			return nil, err
		}
	}

	var u *entities.User = nil
	if user.Id != "" {
		u = &user
	}

	var a *entities.Author = nil
	if author.Id != "" {
		a = &author
	}

	// 4. Return combined result
	return &entities.PersonWithDetails{
		Person: person,
		User:   u,
		Author: a,
	}, nil
}

func (s *PersonService) GetPeopleWithDetails(ctx context.Context) ([]*entities.PersonWithDetails, error) {
	query := `
        FOR p IN persons
            LET user = FIRST(
                FOR e IN persons_users
                    FILTER e._from == p._id
                    RETURN DOCUMENT(e._to)
            )
            LET author = FIRST(
                FOR e IN persons_authors
                    FILTER e._from == p._id
                    RETURN DOCUMENT(e._to)
            )
            RETURN {
                person: p,
                user: user,
                author: author
            }
    `

	opts := &arangodb.QueryOptions{
		BindVars: map[string]interface{}{},
	}

	cursor, err := s.db.Query(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	var results []*entities.PersonWithDetails

	for cursor.HasMore() {
		var item entities.PersonWithDetails
		meta, err := cursor.ReadDocument(ctx, &item)

		item.Person.Id = meta.Key

		if err != nil {
			return nil, err
		}
		results = append(results, &item)
	}

	return results, nil
}
