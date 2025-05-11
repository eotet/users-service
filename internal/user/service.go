package user

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(user User) (User, error) {
	return s.repo.CreateUser(user)
}

func (s *Service) GetUserByID(id uint32) (User, error) {
	return s.repo.GetUserByID(id)
}

func (s *Service) GetAllUsers() ([]User, error) {
	return s.repo.GetAllUsers()
}

func (s *Service) UpdateUserByID(id uint32, user UpdateUserRequest) (User, error) {
	return s.repo.UpdateUserByID(id, user)
}

func (s *Service) DeleteUserByID(id uint32) error {
	return s.repo.DeleteUserByID(id)
}
