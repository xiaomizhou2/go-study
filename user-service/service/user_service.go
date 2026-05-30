package service

import (
	"github.com/example/user-service/model"
	"github.com/example/user-service/repository"
)

// UserService 用户业务逻辑层
// 对比 Java: @Service public class UserServiceImpl implements UserService
//
// Go 通过 struct 持有依赖（而非 @Autowired），构造函数注入
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService 构造函数，接收依赖
// 对比 Java: @Autowired public UserServiceImpl(UserRepository repo)
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetAll 获取所有用户
func (s *UserService) GetAll() ([]*model.User, error) {
	return s.repo.FindAll()
}

// GetByID 按 ID 获取用户
func (s *UserService) GetByID(id int) (*model.User, error) {
	return s.repo.FindByID(id)
}

// Create 创建用户（可以在这里加业务校验）
func (s *UserService) Create(name string, age int) (*model.User, error) {
	user := &model.User{
		Name: name,
		Age:  age,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Update 更新用户
func (s *UserService) Update(id int, name string, age int) (*model.User, error) {
	return s.repo.Update(id, name, age)
}

// Delete 删除用户
func (s *UserService) Delete(id int) error {
	return s.repo.Delete(id)
}
