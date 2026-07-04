package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/arangodb/go-driver/v2/arangodb"
	"gourmand.golang.arango/src/entities"
	"gourmand.golang.arango/src/interfaces"
)

const persons_collection = "persons"
const users_collection = "users"
const persons_users_edges = "persons_users"

type UserService struct {
	db       arangodb.Database
	_persons *PersonService
}

func NewUserService(db arangodb.Database) interfaces.IUser {
	personService := NewPersonService(db)
	return &UserService{db: db, _persons: personService.(*PersonService)}
}

func (s *UserService) CreateUser(ctx context.Context,
	person *entities.Person, user *entities.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	if person.Id != "" {
		existing, err := s._persons.GetPerson(ctx, person.Id)
		if err != nil {
			return err
		}
		person.Id = existing.Id
	} else {
		created, err := s._persons.CreatePerson(ctx, person)
		if err != nil {
			return err
		}
		person.Id = created.Id
	}

	col, err := s.db.GetCollection(ctx, users_collection, nil)
	if err != nil {
		return err
	}

	meta, err := col.CreateDocument(ctx, user)

	if err != nil {
		return err
	}

	user.Id = meta.Key

	edgeCol, err := s.db.GetCollection(ctx, persons_users_edges, nil)
	if err != nil {
		return err
	}

	if person.Id == "" || user.Id == "" {
		return fmt.Errorf("cannot create edge: personKey=%q userKey=%q", person.Id, user.Id)
	}

	edge := map[string]interface{}{
		"_from": "persons/" + person.Id,
		"_to":   "users/" + user.Id,
	}

	_, err = edgeCol.CreateDocument(ctx, edge)

	if err != nil {
		return errors.New("failed to create edge between person and user: " + err.Error())
	}

	return nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	col, err := s.db.GetCollection(ctx, users_collection, nil)
	if err != nil {
		return nil, err
	}

	meta, err := col.ReadDocument(ctx, id, &user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Id = meta.Key

	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user *entities.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	var existingUser entities.User
	col, err := s.db.GetCollection(ctx, users_collection, nil)
	if err != nil {
		return err
	}

	meta, err := col.ReadDocument(ctx, id, &existingUser)
	if err != nil {
		return errors.New("user not found")
	}

	existingUser.Id = meta.Key

	patch := map[string]interface{}{
		"authentication_id": user.Id,
		"active":            user.Active,
	}

	_, err = col.UpdateDocument(ctx, id, patch)
	if err != nil {
		return errors.New("failed to update user: " + err.Error())
	}

	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	col, err := s.db.GetCollection(ctx, users_collection, nil)
	if err != nil {
		return err
	}

	_, err = col.DeleteDocument(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}
	return nil
}
