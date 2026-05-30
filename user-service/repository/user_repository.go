package repository

import (
	"database/sql"
	"fmt"

	"github.com/example/user-service/model"
)

// UserRepository 用户数据仓库（数据库版）
// 对比 Java: @Repository public class UserRepositoryImpl implements UserRepository
//
// 之前是内存版（sync.RWMutex + map），现在换成 database/sql
// Go 的 database/sql = Java 的 JDBC，是最底层的数据库操作接口
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository 构造函数，接收数据库连接
// 对比 Java: @Autowired public UserRepositoryImpl(DataSource dataSource)
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindAll 查询所有用户
// 对比 Java JDBC:
//   PreparedStatement ps = conn.prepareStatement("SELECT ...");
//   ResultSet rs = ps.executeQuery();
//   while (rs.next()) { ... }
func (r *UserRepository) FindAll() ([]*model.User, error) {
	rows, err := r.db.Query("SELECT id, name, age FROM users ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %w", err)
	}
	// defer rows.Close() —— 离开房间自动关灯，无论从哪个门出去
	// 这是 Go 处理资源释放的惯用模式，对比 Java 的 try-with-resources
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		var u model.User
		// Scan 相当于 Java 的 rs.getInt("id"), rs.getString("name")
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			return nil, fmt.Errorf("读取用户数据失败: %w", err)
		}
		users = append(users, &u)
	}

	// 检查遍历过程中是否有错误（别忘了这步！）
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历用户数据失败: %w", err)
	}

	return users, nil
}

// FindByID 按 ID 查询用户
// 对比 Java: rs = ps.executeQuery(); if (rs.next()) { ... }
func (r *UserRepository) FindByID(id int) (*model.User, error) {
	var u model.User
	// QueryRow 查询单行，相当于 PreparedStatement + 只取第一条结果
	// 参数用 ? 占位符（SQLite 语法），PostgreSQL 用 $1, $2
	err := r.db.QueryRow("SELECT id, name, age FROM users WHERE id = ?", id).Scan(&u.ID, &u.Name, &u.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			// sql.ErrNoRows 相当于 Java 的 EmptyResultDataAccessException
			return nil, fmt.Errorf("用户 %d 不存在", id)
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &u, nil
}

// Create 创建用户
// 对比 Java JDBC:
//   ps = conn.prepareStatement(sql, Statement.RETURN_GENERATED_KEYS);
//   ps.executeUpdate();
//   ResultSet keys = ps.getGeneratedKeys();
func (r *UserRepository) Create(user *model.User) error {
	// Result 返回自增 ID（相当于 JDBC 的 getGeneratedKeys）
	result, err := r.db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", user.Name, user.Age)
	if err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取自增ID失败: %w", err)
	}
	user.ID = int(id)

	return nil
}

// Update 更新用户
func (r *UserRepository) Update(id int, name string, age int) (*model.User, error) {
	// 先查后改（也可以用 UPDATE ... RETURNING，但 SQLite 不支持）
	var u model.User
	err := r.db.QueryRow("SELECT id, name, age FROM users WHERE id = ?", id).Scan(&u.ID, &u.Name, &u.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("用户 %d 不存在", id)
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	if name != "" {
		u.Name = name
	}
	if age > 0 {
		u.Age = age
	}

	_, err = r.db.Exec("UPDATE users SET name = ?, age = ? WHERE id = ?", u.Name, u.Age, u.ID)
	if err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	return &u, nil
}

// Delete 删除用户
func (r *UserRepository) Delete(id int) error {
	// result.RowsAffected() 相当于 JDBC 的 getUpdateCount()
	result, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("用户 %d 不存在", id)
	}

	return nil
}
