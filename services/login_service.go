package services

import "github.com/MrWestbury/terraxen/backend"

type LoginService struct {
	backend.MongoBackend
}

func NewLoginService(options Options) *ModuleService {
	svc := &ModuleService{
		backend.MongoBackend{
			ConnectionString: options.ConnectionString(),
			Database:         options.Database,
		},
	}

	return svc
}
