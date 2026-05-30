package service

import (
	"database/sql"
	"testing"

	"github.com/example/user-service/model"
	"github.com/example/user-service/repository"
	_ "modernc.org/sqlite"
)

// TestCreateAndGet 表驱动测试 —— Go 最经典的测试模式
// 对比 Java: @ParameterizedTest + @CsvSource
//
// 比喻：像考试的选择题表格，一行就是一个测试用例
// 每行有：输入 → 期望输出
// 全部跑一遍，错了哪行一目了然
func TestCreateAndGet(t *testing.T) {
	// 用 SQLite 内存数据库做测试
	// 对比 Java: @DataJpaTest 用 H2 内存数据库
	svc := newTestService(t)

	tests := []struct {
		name      string
		inputName string
		inputAge  int
		wantName  string
		wantAge   int
		wantErr   bool
	}{
		{"正常创建", "张三", 25, "张三", 25, false},
		{"年龄边界", "婴儿", 0, "婴儿", 0, false},
		{"大龄用户", "老王", 150, "老王", 150, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.Create(tt.inputName, tt.inputAge)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if user.Name != tt.wantName {
				t.Errorf("Create() name = %v, want %v", user.Name, tt.wantName)
			}
			if user.Age != tt.wantAge {
				t.Errorf("Create() age = %v, want %v", user.Age, tt.wantAge)
			}
			// 创建后应该能查到
			if user.ID <= 0 {
				t.Error("Create() 应该返回有效的 ID")
			}
		})
	}
}

// TestGetByID_NotFound 测试查询不存在的用户
func TestGetByID_NotFound(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.GetByID(999)
	if err == nil {
		t.Error("GetByID(999) 应该返回错误，但没有")
	}
}

// TestGetAll 测试获取所有用户
func TestGetAll(t *testing.T) {
	svc := newTestService(t)

	// 先创建 3 个用户
	svc.Create("张三", 25)
	svc.Create("李四", 30)
	svc.Create("王五", 28)

	users, err := svc.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}
	if len(users) != 3 {
		t.Errorf("GetAll() 返回 %d 个用户，期望 3 个", len(users))
	}
}

// TestUpdate 测试更新用户
func TestUpdate(t *testing.T) {
	svc := newTestService(t)

	// 先创建
	user, _ := svc.Create("张三", 25)

	// 再更新
	updated, err := svc.Update(user.ID, "张三丰", 99)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Name != "张三丰" {
		t.Errorf("Update() name = %v, want 张三丰", updated.Name)
	}
	if updated.Age != 99 {
		t.Errorf("Update() age = %v, want 99", updated.Age)
	}

	// 部分更新：只改名字
	partial, err := svc.Update(user.ID, "张无忌", 0)
	if err != nil {
		t.Fatalf("Update(partial) error = %v", err)
	}
	if partial.Name != "张无忌" {
		t.Errorf("Update(partial) name = %v, want 张无忌", partial.Name)
	}
	if partial.Age != 99 {
		t.Errorf("Update(partial) age 不应该改变，got %d, want 99", partial.Age)
	}
}

// TestDelete 测试删除用户
func TestDelete(t *testing.T) {
	svc := newTestService(t)

	user, _ := svc.Create("张三", 25)

	// 删除
	err := svc.Delete(user.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// 再查应该找不到
	_, err = svc.GetByID(user.ID)
	if err == nil {
		t.Error("删除后再查询应该返回错误")
	}
}

// TestDelete_NotFound 测试删除不存在的用户
func TestDelete_NotFound(t *testing.T) {
	svc := newTestService(t)

	err := svc.Delete(999)
	if err == nil {
		t.Error("Delete(999) 应该返回错误")
	}
}

// newTestService 创建一个用内存数据库的测试用 service
// 对比 Java: @BeforeEach void setUp() { ... }
//
// 每个测试函数都调用这个，拿到独立的、干净的 service 实例
func newTestService(t *testing.T) *UserService {
	t.Helper() // 标记为测试辅助函数，报错时显示调用方的行号

	// SQLite 内存数据库，每个连接都是独立的
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("打开内存数据库失败: %v", err)
	}

	// 建表
	_, err = db.Exec(`
		CREATE TABLE users (
			id   INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			age  INTEGER NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("建表失败: %v", err)
	}

	repo := repository.NewUserRepository(db)
	return NewUserService(repo)
}

// TestModel_User 测试 model 的 JSON tag
// 演示子测试的写法
func TestModel_User(t *testing.T) {
	user := model.User{ID: 1, Name: "张三", Age: 25}

	if user.Name != "张三" {
		t.Errorf("User.Name = %v, want 张三", user.Name)
	}
}
