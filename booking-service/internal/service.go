package internal

type ServiceInterface interface {
	FindByID(id uint) error
}

type srv struct {
	repo RepositoryInterface
}

func Newsrv(repo RepositoryInterface) *srv {
	return &srv{
		repo: repo,
	}
}

func (s *srv) FindByID(id uint) error {
	return nil
}