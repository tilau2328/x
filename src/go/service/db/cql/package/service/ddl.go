package service

import (
	. "context"
)

type (
	DDLServiceOptions struct {
		KeySpaceProvider
		TableProvider
	}
	DDLService struct {
		DDLServiceOptions
	}
)

var _ DDL = &DDLService{}

func NewDDL(opts DDLServiceOptions) *DDLService {
	return &DDLService{DDLServiceOptions: opts}
}

func (s *DDLService) ListKeySpaces(ctx Context, space KeySpace) ([]KeySpace, error) {
	return s.KeySpaceProvider.List(ctx, space)
}

func (s *DDLService) GetKeySpace(ctx Context, key KeySpaceKey) (KeySpace, error) {
	//TODO implement me
	panic("implement me")
}

func (s *DDLService) CreateKeySpace(ctx Context, space KeySpace) error {
	//TODO implement me
	panic("implement me")
}

func (s *DDLService) AlterKeySpace(ctx Context, key KeySpaceKey, patches []Patch) error {
	//TODO implement me
	panic("implement me")
}

func (s *DDLService) DropKeySpace(ctx Context, key KeySpaceKey) error {
	//TODO implement me
	panic("implement me")
}

func (s *DDLService) ListTables(ctx Context, table Table) ([]Table, error) {
	//TODO implement me
	panic("implement me")
}

func (s *DDLService) GetTable(ctx Context, key TableKey) (Table, error) {
	//TODO implement me
	panic("implement me")
}

func (s *DDLService) CreateTable(ctx Context, table Table) error {
	//TODO implement me
	panic("implement me")
}

func (s *DDLService) AlterTable(ctx Context, key TableKey, patches []Patch) error {
	//TODO implement me
	panic("implement me")
}

func (s *DDLService) DropTable(ctx Context, key TableKey) error {
	//TODO implement me
	panic("implement me")
}
