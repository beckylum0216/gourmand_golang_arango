package services

import (
	"context"
	"errors"

	"github.com/arangodb/go-driver/v2/arangodb"

	"gourmand.golang.arango/src/entities"
	"gourmand.golang.arango/src/interfaces"
)

const authors_collection = "authors"
const persons_authors_edges = "persons_authors"

type AuthorService struct {
	db       arangodb.Database
	_persons *PersonService
}

func NewAuthorService(db arangodb.Database) interfaces.IAuthor {
	personService := NewPersonService(db)
	return &AuthorService{db: db, _persons: personService.(*PersonService)}
}

func (s *AuthorService) CreateAuthor(ctx context.Context, 
	person *entities.Person, author *entities.Author) error {
	
	if author == nil {
		return errors.New("author is nil")
	}

	var personkey string

	if person.Id == "" {
		existing, err := s._persons.GetPerson(ctx, person.Id)
		if err != nil {
			return err
		}
		personkey = existing.Id
	} else {
		created, err := s._persons.CreatePerson(ctx, person)
		if err != nil {
			return err
		}
		personkey = created.Id
	}

	col, err := s.db.GetCollection(ctx, authors_collection, nil)
	if err != nil {
		return err
	}

	meta, err := col.CreateDocument(ctx, author)
	if err != nil {
		return err
	}

	author.Id = meta.Key

	edgeCol, err := s.db.GetCollection(ctx, persons_authors_edges, nil)
	if err != nil {
		return err
	}

	edge := map[string]interface{}{
		"_from": "persons/" + personkey,
		"_to":   "authors/" + author.Id,
	}

	if _, err := edgeCol.CreateDocument(ctx, edge); 
	err != nil {
		return err
	}

	
	return nil
}

func (s *AuthorService) GetAuthor(ctx context.Context, id string) (*entities.Author, error) {
	var author entities.Author

	col, err := s.db.GetCollection(ctx, collection, nil)
	if err != nil {
		return nil, err
	}

	meta, err := col.ReadDocument(ctx, id, &author)
	if err != nil {
		return nil, errors.New("author not found")
	}

	author.Id = meta.Key

	return &author, nil
}

func (s *AuthorService) GetAuthors(ctx context.Context) ([]*entities.Author, error) {
	query := "FOR a IN authors RETURN a"
	bindVars := map[string]interface{}{
		"@collection": authors_collection,
	}

	options := &arangodb.QueryOptions{
		BindVars: bindVars,
	}

	cursor, err := s.db.Query(ctx, query, options)
	if err != nil {
		return nil, err
	}

	var authors []*entities.Author
	for cursor.HasMore() {
		var author entities.Author
		_, err := cursor.ReadDocument(ctx, &author)
		if err != nil {
			return nil, err
		}
		authors = append(authors, &author)
	}

	return authors, nil
}

func (s *AuthorService) UpdateAuthor(ctx context.Context, id string, author *entities.Author) error {
	if author == nil {
		return errors.New("author is nil")
	}

	col, err := s.db.GetCollection(ctx, authors_collection, nil)
	if err != nil {
		return err
	}


	patch := map[string]interface{}{
		"source":        author.Source,
		"credit_string": author.CreditString,
	}

	_, err = col.UpdateDocument(ctx, id, patch)
	if err != nil {
		return errors.New("author not found")
	}


	return nil
}

func (s *AuthorService) DeleteAuthor(ctx context.Context, id string) error {
	col, err := s.db.GetCollection(ctx, authors_collection, nil)
	if err != nil {
		return err
	}

	_, err = col.DeleteDocument(ctx, id)
	if err != nil {
		return errors.New("author not found")
	}

	return nil
}
