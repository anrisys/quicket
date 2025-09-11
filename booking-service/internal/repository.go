package internal

type RepositoryInterface interface{}

type repo struct{}

func NewRepo() *repo {
	return &repo{}
}