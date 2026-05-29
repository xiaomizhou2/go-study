# 第二部分：从零构建 Go 项目

> 目标：用 Go 开发一个"用户管理" REST API 服务，每一步都和 Java 对比

---

## 项目最终结构预览 🗂️

```
user-service/
├── go.mod              # 依赖管理（对比 pom.xml）
├── go.sum              # 依赖校验和（类似 Maven 的校验机制）
├── main.go             # 入口文件（对比 Application.java）
├── config/             # 配置
│   └── config.go
├── model/              # 数据模型（对比 Entity/DTO）
│   └── user.go
├── handler/            # HTTP 处理器（对比 @RestController）
│   └── user_handler.go
├── service/            # 业务逻辑（对比 @Service）
│   └── user_service.go
├── repository/         # 数据访问（对比 Repository/Mapper）
│   └── user_repository.go
├── middleware/         # 中间件（对比 Filter/Interceptor）
│   └── logger.go
│   └── auth.go
├── util/               # 工具类（对比 Utils）
│   └── response.go
└── test/               # 测试（对比 src/test/java）
    └── user_test.go
```

---

## 步骤 1：环境准备 ✅

### Go 环境已就绪
```bash
$ go version
go version go1.24.6 darwin/arm64
```

### 理解 GOPATH vs Go Modules

| 概念 | 对应 Java | 说明 |
|------|----------|------|
| **GOPATH** | 类似 JAVA_HOME | Go 旧版依赖管理方式，已不推荐使用 |
| **Go Modules** | 类似 Maven/Gradle | Go 1.11+ 的现代依赖管理，**推荐使用** |

**关键区别：**

```bash
# Java (Maven)
~/.m2/repository/              # 本地 Maven 仓库
pom.xml                        # 项目依赖配置

# Go (Modules)
~/go/pkg/mod/                  # Go module 缓存（自动管理）
go.mod + go.sum                # 项目依赖配置
```

> 💡 **提示：** 现代 Go 项目不需要手动设置 GOPATH，Go Modules 会自动处理依赖。

---

## 步骤 2：项目初始化 🚀

### 2.1 创建项目目录

```bash
cd go-study
mkdir user-service
cd user-service
```

### 2.2 初始化 Go Module

```bash
# 对比 Java：mvn archetype:generate 或 Spring Initializr
go mod init github.com/example/user-service
```

生成的 `go.mod`：
```go
module github.com/example/user-service

go 1.22
```

**对比 Java：**
| Go | Java (Maven) |
|----|--------------|
| `go mod init` | 创建 `pom.xml` |
| `module name` | `<groupId>+<artifactId>` |
| `go 1.22` | `<maven.compiler.source>` |

### 2.3 创建基础目录结构

```bash
mkdir -p {model,handler,service,repository,middleware,util,config,test}
touch main.go
```

---

## 步骤 3：第一个 HTTP 端点 🌐

### 3.1 先看 Java 版本

```java
@RestController
@RequestMapping("/users")
public class UserController {

    private List<User> users = Arrays.asList(
        new User(1, "张三", 25),
        new User(2, "李四", 30)
    );

    @GetMapping
    public List<User> getAllUsers() {
        return users;
    }

    @GetMapping("/{id}")
    public User getUserById(@PathVariable Long id) {
        return users.stream()
                    .filter(u -> u.getId().equals(id))
                    .findFirst()
                    .orElse(null);
    }
}
```

### 3.2 Go 版本（使用标准库）

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

// User 结构体定义
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

var users = []User{
    {ID: 1, Name: "张三", Age: 25},
    {ID: 2, Name: "李四", Age: 30},
}

func main() {
    // 注册路由
    http.HandleFunc("/users", usersHandler)
    http.HandleFunc("/users/", userByIDHandler)

    fmt.Println("服务启动在 :8080")
    http.ListenAndServe(":8080", nil) // Java: server.start()
}

// GET /users
func usersHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users) // 自动序列化为 JSON
}

// GET /users/1
func userByIDHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // 解析路径参数
    id := 0
    fmt.Sscanf(r.URL.Path[len("/users/"):], "%d", &id)

    for _, u := range users {
        if u.ID == id {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(u)
            return
        }
    }

    http.NotFound(w, r)
}
```

### 3.3 对比总结

| Java | Go |
|------|-----|
| `@RestController` | 普通函数 + `http.HandleFunc` |
| `@GetMapping` | `if r.Method != http.MethodGet` |
| `@PathVariable` | 手动解析 `r.URL.Path` |
| `@RequestParam` | `r.URL.Query().Get("name")` |
| `@RequestBody` | `json.NewDecoder(r.Body).Decode(&obj)` |
| 自动 JSON 序列化 | `json.NewEncoder(w).Encode(obj)` |
| 返回对象 | `w.Write([]byte("..."))` |
| 自动异常处理 | 需要手动返回错误码 |

> 💡 **注意：** 标准库比较底层，实际项目我们通常用 Gin 框架（下一节会讲）。

---

## 步骤 4：使用 Gin 框架（推荐）🔥

### 4.1 为什么选 Gin？

| 特性 | 标准库 | Gin |
|------|--------|-----|
| 路由 | 手动解析 | 自动匹配 |
| 参数绑定 | 手动解析 | 自动绑定 |
| 中间件 | 需要自己封装 | 内置丰富 |
| JSON | 手动编码 | 自动序列化 |
| 生态 | 零依赖 | 插件丰富 |

### 4.2 安装 Gin

```bash
go get -u github.com/gin-gonic/gin
```

对比 Java：
```xml
<!-- pom.xml -->
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-web</artifactId>
</dependency>
```

### 4.3 用 Gin 重写

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type User struct {
    ID   int    `json:"id" binding:"required"`
    Name string `json:"name" binding:"required"`
    Age  int    `json:"age" binding:"gte=0,lte=150"`
}

var users = []User{
    {ID: 1, Name: "张三", Age: 25},
    {ID: 2, Name: "李四", Age: 30},
}

func main() {
    r := gin.Default() // 默认带 Logger 和 Recovery 中间件

    // 路由组（类似 @RequestMapping）
    api := r.Group("/api/v1")
    {
        api.GET("/users", getAllUsers)
        api.GET("/users/:id", getUserByID)
        api.POST("/users", createUser)
        api.PUT("/users/:id", updateUser)
        api.DELETE("/users/:id", deleteUser)
    }

    r.Run(":8080") // Java: server.start()
}

// GET /api/v1/users
func getAllUsers(c *gin.Context) {
    c.JSON(http.StatusOK, users) // 自动返回 JSON
}

// GET /api/v1/users/:id
func getUserByID(c *gin.Context) {
    // 路径参数（类似 @PathVariable）
    id := c.Param("id")

    for _, u := range users {
        if fmt.Sprintf("%d", u.ID) == id {
            c.JSON(http.StatusOK, u)
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
}

// POST /api/v1/users
func createUser(c *gin.Context) {
    var newUser User

    // 绑定请求体（类似 @RequestBody）
    if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    users = append(users, newUser)
    c.JSON(http.StatusCreated, newUser)
}

// PUT /api/v1/users/:id
func updateUser(c *gin.Context) {
    id := c.Param("id")
    var updates User

    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    for i, u := range users {
        if fmt.Sprintf("%d", u.ID) == id {
            if updates.Name != "" {
                users[i].Name = updates.Name
            }
            if updates.Age > 0 {
                users[i].Age = updates.Age
            }
            c.JSON(http.StatusOK, users[i])
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
}

// DELETE /api/v1/users/:id
func deleteUser(c *gin.Context) {
    id := c.Param("id")

    for i, u := range users {
        if fmt.Sprintf("%d", u.ID) == id {
            users = append(users[:i], users[i+1:]...)
            c.Status(http.StatusNoContent)
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
}
```

### 4.4 Gin vs Spring Boot 对比

| 功能 | Spring Boot | Gin |
|------|-------------|-----|
| 路由定义 | `@GetMapping("/users/:id")` | `r.GET("/users/:id", handler)` |
| 路径参数 | `@PathVariable Long id` | `c.Param("id")` |
| 查询参数 | `@RequestParam String name` | `c.Query("name")` |
| 请求体 | `@RequestBody User user` | `c.ShouldBindJSON(&user)` |
| 响应 JSON | `return user` | `c.JSON(200, user)` |
| 错误处理 | `@ExceptionHandler` | `c.Error()` 或自定义中间件 |
| 路由组 | `@RequestMapping("/api/v1")` | `api := r.Group("/api/v1")` |

---

## 步骤 5：路由与中间件 🔧

### 5.1 什么是中间件？

**比喻：**
- Java Filter = 商场的安检员 🛂 —— 所有进门的人都要检查
- Go Middleware = 工厂流水线 🏭 —— 每个环节处理一部分工作

### 5.2 Java Filter 示例

```java
@Component
@Order(1)
public class AuthFilter implements Filter {
    @Override
    public void doFilter(ServletRequest req, ServletResponse res, FilterChain chain) {
        HttpServletRequest request = (HttpServletRequest) req;
        String token = request.getHeader("Authorization");

        if (token == null || !token.startsWith("Bearer ")) {
            ((HttpServletResponse) res).setStatus(401);
            return;
        }

        chain.doFilter(req, res); // 继续下一个 Filter
    }
}
```

### 5.3 Go 中间件示例

```go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
}

// 认证中间件
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")

        if token == "" || !strings.HasPrefix(token, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
            c.Abort() // 停止后续处理
            return
        }

        // 提取用户信息（简化版）
        userID := strings.TrimPrefix(token, "Bearer ")
        c.Set("userID", userID) // 存储到上下文

        c.Next() // 继续处理下一个中间件
    }
}

// 日志中间件
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path

        c.Next()

        latency := time.Since(start)
        status := c.Writer.Status()

        log.Printf("[%s] %s %d %v",
            c.Request.Method,
            path,
            status,
            latency,
        )
    }
}

// CORS 中间件
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}
```

### 5.4 使用中间件

```go
func main() {
    r := gin.New() // 不使用默认中间件

    // 全局中间件
    r.Use(middleware.Logger())
    r.Use(middleware.CORS())
    r.Use(gin.Recovery()) // panic 恢复

    // 需要认证的路由
    auth := r.Group("/api")
    auth.Use(middleware.Auth())
    {
        auth.GET("/profile", getProfile)
        auth.PUT("/profile", updateProfile)
    }

    // 公开路由
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    r.Run(":8080")
}
```

### 5.5 中间件对比

| Java | Go (Gin) |
|------|----------|
| `@Component` + `implements Filter` | 普通函数返回 `gin.HandlerFunc` |
| `FilterChain.doFilter()` | `c.Next()` 继续或 `c.Abort()` 停止 |
| `@Order` 定义顺序 | `r.Use()` 顺序就是执行顺序 |
| `web.xml` 或 `@Bean` 注册 | 直接调用 `r.Use()` 或 `group.Use()` |

---

## 步骤 6：数据库操作 🗄️

### 6.1 两个选择：database/sql vs GORM

| 方面 | database/sql | GORM |
|------|---------------|------|
| 定位 | 标准库，底层 | ORM 框架，高层 |
| 类比 | JDBC + MyBatis | Hibernate / JPA |
| 学习曲线 | 需要手写 SQL | 自动生成 SQL |
| 性能 | 更高 | 略低（可接受） |
| 灵活性 | 完全控制 | 适合 CRUD |
| 适用场景 | 复杂查询、性能敏感 | 快速开发、标准 CRUD |

### 6.2 使用 database/sql（像 JDBC）

```go
package repository

import (
    "database/sql"
    "fmt"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

// 查询单个用户（类似 MyBatis 的 mapped statement）
func (r *UserRepository) FindByID(id int) (*User, error) {
    var user User
    query := `SELECT id, name, age FROM users WHERE id = $1`

    err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Age)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("用户不存在")
        }
        return nil, err
    }

    return &user, nil
}

// 查询多个用户（像 ResultSet 遍历）
func (r *UserRepository) FindAll() ([]User, error) {
    query := `SELECT id, name, age FROM users ORDER BY id`
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close() // ← 记得关闭！

    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    return users, rows.Err() // 检查遍历过程是否有错误
}

// 插入用户
func (r *UserRepository) Create(user *User) error {
    query := `INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id`
    return r.db.QueryRow(query, user.Name, user.Age).Scan(&user.ID)
}

// 更新用户
func (r *UserRepository) Update(user *User) error {
    query := `UPDATE users SET name = $1, age = $2 WHERE id = $3`
    result, err := r.db.Exec(query, user.Name, user.Age, user.ID)
    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("更新失败：用户不存在")
    }

    return nil
}

// 删除用户
func (r *UserRepository) Delete(id int) error {
    query := `DELETE FROM users WHERE id = $1`
    result, err := r.db.Exec(query, id)
    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("删除失败：用户不存在")
    }

    return nil
}

// 事务处理（类似 @Transactional）
func (r *UserRepository) TransferMoney(fromID, toID int, amount float64) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }

    defer func() {
        if err != nil {
            tx.Rollback() // 回滚
        }
    }()

    // 扣钱
    _, err = tx.Exec(`UPDATE accounts SET balance = balance - $1 WHERE id = $2`, amount, fromID)
    if err != nil {
        return err
    }

    // 加钱
    _, err = tx.Exec(`UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount, toID)
    if err != nil {
        return err
    }

    return tx.Commit() // 提交
}
```

### 6.3 使用 GORM（像 JPA）

```go
package repository

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model        // 自动添加 ID、CreatedAt、UpdatedAt、DeletedAt
    Name  string      `gorm:"type:varchar(100);not null"`
    Age   int         `gorm:"not null"`
}

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    // 自动迁移（类似 JPA 的 auto-ddl）
    db.AutoMigrate(&User{})
    return &UserRepository{db: db}
}

// CRUD 操作（非常简洁！）
func (r *UserRepository) FindByID(id uint) (*User, error) {
    var user User
    result := r.db.First(&user, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

func (r *UserRepository) FindAll() ([]User, error) {
    var users []User
    result := r.db.Find(&users)
    return users, result.Error
}

func (r *UserRepository) Create(user *User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *User) error {
    return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
    return r.db.Delete(&User{}, id).Error
}

// 复杂查询（链式调用）
func (r *UserRepository) FindByAgeRange(min, max int) ([]User, error) {
    var users []User
    result := r.db.Where("age >= ? AND age <= ?", min, max).
               Order("age DESC").
               Limit(10).
               Find(&users)
    return users, result.Error
}

// 关联查询（类似 JPA 的 @OneToMany）
type Order struct {
    gorm.Model
    UserID uint
    User   User      `gorm:"foreignKey:UserID"`
    Items  []OrderItem
}

func (r *OrderRepository) FindWithUser(id uint) (*Order, error) {
    var order Order
    result := r.db.Preload("User").Preload("Items").First(&order, id)
    return &order, result.Error
}
```

### 6.4 数据库连接池

```go
// Java (Spring Boot)
// application.properties
// spring.datasource.hikari.maximum-pool-size=10
// spring.datasource.hikari.minimum-idle=5

// Go (database/sql)
func NewDB(dataSource string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dataSource)
    if err != nil {
        return nil, err
    }

    // 连接池配置
    db.SetMaxOpenConns(25)        // 最大连接数（像 HikariCP maximum-pool-size）
    db.SetMaxIdleConns(10)        // 最大空闲连接（像 minimum-idle）
    db.SetConnMaxLifetime(5 * time.Minute) // 连接最大存活时间

    // 验证连接
    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}
```

---

## 步骤 7：配置管理 ⚙️

### 7.1 Java application.yml vs Go Viper

**Java (application.yml)：**
```yaml
server:
  port: 8080

spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/userdb
    username: admin
    password: secret

logging:
  level:
    root: INFO
    com.example: DEBUG
```

### 7.2 安装 Viper

```bash
go get github.com/spf13/viper
```

### 7.3 配置实现

```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Log      LogConfig
}

type ServerConfig struct {
    Port            string
    ReadTimeout     int
    WriteTimeout    int
    ShutdownTimeout int
}

type DatabaseConfig struct {
    Driver          string
    Host            string
    Port            int
    User            string
    Password        string
    DBName          string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime int
}

type LogConfig struct {
    Level  string
    Format string // json 或 text
}

// 加载配置
func Load(configPath string) (*Config, error) {
    v := viper.New()
    v.SetConfigFile(configPath)

    // 也可以支持环境变量
    v.AutomaticEnv()
    v.SetEnvPrefix("APP") // APP_SERVER_PORT 会被绑定到 server.port

    if err := v.ReadInConfig(); err != nil {
        return nil, err
    }

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}
```

### 7.4 config.yaml

```yaml
server:
  port: "8080"
  read_timeout: 30
  write_timeout: 30
  shutdown_timeout: 10

database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  user: "admin"
  password: "secret"
  db_name: "userdb"
  max_open_conns: 25
  max_idle_conns: 10
  conn_max_lifetime: 300

log:
  level: "info"
  format: "json"
```

---

## 步骤 8：错误处理与日志 📝

### 8.1 统一错误响应

```go
package util

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

// 错误码常量（像 Java 的枚举）
const (
    ErrCodeBadRequest    = 400
    ErrCodeUnauthorized   = 401
    ErrCodeNotFound       = 404
    ErrCodeInternalError  = 500
)

// 响应错误
func ResponseError(c *gin.Context, code int, message string, err error) {
    resp := ErrorResponse{
        Code:    code,
        Message: message,
    }

    if err != nil && gin.Mode() == gin.DebugMode {
        resp.Detail = err.Error()
    }

    c.JSON(code, resp)
}

// 响应成功
func ResponseSuccess(c *gin.Context, data interface{}) {
    c.JSON(200, gin.H{
        "code":    0,
        "message": "success",
        "data":    data,
    })
}
```

### 8.2 结构化日志

```bash
# 安装 zap（uber 出品的高性能日志库）
go get go.uber.org/zap
```

```go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// 初始化日志
func Init(level string, format string) error {
    var config zap.Config

    if format == "json" {
        config = zap.NewProductionConfig()
    } else {
        config = zap.NewDevelopmentConfig()
    }

    // 设置日志级别
    switch level {
    case "debug":
        config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
    case "info":
        config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
    case "error":
        config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
    }

    log, err := config.Build()
    if err != nil {
        return err
    }

    Log = log
    return nil
}

// 使用示例
func example() {
    Log.Debug("调试信息",
        zap.String("user", "张三"),
        zap.Int("age", 25),
    )

    Log.Info("用户登录",
        zap.String("ip", "127.0.0.1"),
        zap.String("method", "POST"),
    )

    Log.Error("数据库错误",
        zap.String("query", "SELECT ..."),
        zap.Error(err),
    )
}
```

---

## 步骤 9：单元测试 🧪

### 9.1 Go 测试惯例

| Java | Go |
|------|-----|
| `src/test/java/` | 文件名加 `_test.go` |
| `@Test` | 函数名以 `Test` 开头 |
| `assertEquals()` | 直接比较或 `assert.Equal()` |
| JUnit | `testing` 包 |
| Mockito | `gomock` 或手写 fake |

### 9.2 表驱动测试（Go 特色！）

```go
package service

import (
    "testing"
)

// 测试数据
func TestValidateAge(t *testing.T) {
    tests := []struct {
        name    string
        age     int
        wantErr bool
    }{
        {"正常年龄", 25, false},
        {"未成年", 17, true},
        {"负数", -1, true},
        {"超大年龄", 200, true},
        {"边界值-18", 18, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateAge(tt.age)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateAge(%d) error = %v, wantErr %v",
                    tt.age, err, tt.wantErr)
            }
        })
    }
}
```

### 9.3 HTTP 测试

```go
func TestGetAllUsers(t *testing.T) {
    // 设置 Gin 测试路由
    r := setupTestRouter()

    // 创建测试请求
    req, _ := http.NewRequest("GET", "/api/v1/users", nil)
    w := httptest.NewRecorder()

    // 执行请求
    r.ServeHTTP(w, req)

    // 验证响应
    if w.Code != http.StatusOK {
        t.Errorf("状态码错误: got %v want %v", w.Code, http.StatusOK)
    }

    var users []User
    json.Unmarshal(w.Body.Bytes(), &users)

    if len(users) == 0 {
        t.Error("应该返回用户列表")
    }
}
```

---

## 步骤 10：构建与部署 🚢

### 10.1 交叉编译

```bash
# Java: 打包成 JAR，目标平台有 JVM 就能跑
mvn package

# Go: 编译成平台相关的二进制文件
GOOS=linux GOARCH=amd64 go build -o user-service-linux

# 常用组合
GOOS=linux GOARCH=amd64 go build   # Linux 64位
GOOS=darwin GOARCH=amd64 go build   # macOS Intel
GOOS=darwin GOARCH=arm64 go build   # macOS Apple Silicon
GOOS=windows GOARCH=amd64 go build  # Windows
```

### 10.2 最小 Docker 镜像

```dockerfile
# 构建阶段（使用完整 Go 镜像编译）
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service

# 运行阶段（使用最小基础镜像）
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/user-service .
EXPOSE 8080

CMD ["./user-service"]
```

对比 Java：
```dockerfile
# Java 需要 JRE（~200MB）
FROM openjdk:17-jre-slim
COPY target/user-service.jar app.jar
ENTRYPOINT ["java", "-jar", "app.jar"]
```

| 特性 | Go | Java |
|------|-----|------|
| 基础镜像大小 | ~5MB (alpine) | ~200MB (JRE) |
| 最终镜像大小 | ~10-20MB | ~200MB+ |
| 启动速度 | 毫秒级 | 秒级 |
| 运行时依赖 | 无需 JVM | 需要 JRE |

---

## 完整项目骨架 🦴

> 详见项目根目录的 `user-service/` 文件夹

---

## 继续学习路线 🗺️

1. **深入 Go 标准库**：`net/http`、`database/sql`、`context`
2. **进阶并发模式**：`errgroup`、`semaphore`、`worker pool`
3. **微服务实战**：gRPC、服务发现、链路追踪
4. **性能优化**：pprof、trace、benchstat
5. **工程实践**：CI/CD、代码生成、Go 生成式编程

---

🎓 **恭喜完成第二部分！现在准备好进入实战挑战了吗？**
