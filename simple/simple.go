package simple

import "errors"

type SimpleRepository struct {
	Error bool
}

type SimpleService struct {
	*SimpleRepository
}

func NewSimpleService(repository *SimpleRepository) (*SimpleService, error) {
	if repository.Error {
		return nil, errors.New("failed create service")
	} else {
		return &SimpleService{SimpleRepository: repository}, nil
	}
}

func NewSimpleRepository(isError bool) *SimpleRepository {
	return &SimpleRepository{Error: isError}
}
