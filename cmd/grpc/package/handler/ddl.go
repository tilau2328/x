package handler

import (
	"github.com/tilau2328/cql/package/domain/provider"
	"grpc/package/model"
)

type (
	DDLOptions struct {
		provider.DDL
	}
	DDL struct {
		model.UnimplementedDDLServer
		DDLOptions
	}
)

var _ model.DDLServer = &DDL{}

func NewDDL(opts DDLOptions) *DDL {
	return &DDL{
		DDLOptions: opts,
	}
}