package service

import "github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/cellnet/proto"

type UserService interface {
	Init() error
	Login(username string, password string) proto.AuthorizeResult
}
