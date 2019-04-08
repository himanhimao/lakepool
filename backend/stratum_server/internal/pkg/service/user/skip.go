package user

import "github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/cellnet/proto"

type SkipUserService struct {
}

func NewSkipUserService() *SkipUserService {
	return &SkipUserService{}
}

func (g *SkipUserService) Login(username string, password string) proto.AuthorizeResult {
	return true
}

func (g *SkipUserService) Init() error {
	return nil
}
