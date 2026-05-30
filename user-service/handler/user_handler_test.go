package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/user-service/middleware"
	"github.com/example/user-service/repository"
	"github.com/example/user-service/service"
	"github.com/example/user-service/util"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

// TestGetAllUsers 测试获取所有用户接口
// 对比 Java: @SpringBootTest + MockMvc.perform(get("/users"))
//
// Go 用 httptest.NewRecorder() 当作 ResponseWriter
// 再用 gin.Engine.ServeHTTP() 执行请求，完全不需要启动服务器
func TestGetAllUsers(t *testing.T) {
	router, _ := newTestRouter(t)

	// 创建测试请求
	// 对比 Java: MockMvcRequestBuilders.get("/api/v1/users")
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)

	// 录制响应（相当于 MockMvc 的 andReturn()）
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证状态码
	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %d, want %d", w.Code, http.StatusOK)
	}

	// 解析响应体
	var resp util.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	if resp.Code != 0 {
		t.Errorf("响应 code = %d, want 0", resp.Code)
	}
}

// TestGetByID 测试按 ID 查询用户
func TestGetByID(t *testing.T) {
	router, svc := newTestRouter(t)

	// 先创建一个用户
	user, _ := svc.Create("张三", 25)

	// 查询
	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %d, want %d", w.Code, http.StatusOK)
	}

	var resp util.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Errorf("响应 code = %d, want 0", resp.Code)
	}

	_ = user // 避免未使用变量警告
}

// TestGetByID_NotFound 测试查询不存在的用户
func TestGetByID_NotFound(t *testing.T) {
	router, _ := newTestRouter(t)

	req, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("状态码 = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// TestGetByID_InvalidID 测试无效 ID
func TestGetByID_InvalidID(t *testing.T) {
	router, _ := newTestRouter(t)

	req, _ := http.NewRequest("GET", "/api/v1/users/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestCreateUser 测试创建用户（需认证）
func TestCreateUser(t *testing.T) {
	router, _ := newTestRouter(t)

	body := `{"name":"赵六","age":22}`
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("状态码 = %d, want %d", w.Code, http.StatusCreated)
	}
}

// TestCreateUser_NoToken 测试无 Token 创建用户（应 401）
func TestCreateUser_NoToken(t *testing.T) {
	router, _ := newTestRouter(t)

	body := `{"name":"赵六","age":22}`
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("状态码 = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

// TestCreateUser_InvalidBody 测试无效请求体
func TestCreateUser_InvalidBody(t *testing.T) {
	router, _ := newTestRouter(t)

	body := `{"invalid"}`
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("状态码 = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestHealthCheck 测试健康检查
func TestHealthCheck(t *testing.T) {
	router, _ := newTestRouter(t)

	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("状态码 = %d, want %d", w.Code, http.StatusOK)
	}
}

// newTestRouter 创建测试路由
// 对比 Java: @BeforeEach + MockMvc 的 setup
//
// 返回 (router, service) 方便测试中操作数据
func newTestRouter(t *testing.T) (*gin.Engine, *service.UserService) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("打开内存数据库失败: %v", err)
	}
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
	svc := service.NewUserService(repo)
	h := NewUserHandler(svc)

	// 构建测试路由（只注册必要的中间件和路由）
	r := gin.New()
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())

	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "user-service"})
	})

	users := r.Group("/api/v1/users")
	{
		users.GET("", h.GetAll)
		users.GET("/:id", h.GetByID)

		auth := users.Group("")
		auth.Use(middleware.Auth())
		{
			auth.POST("", h.Create)
			auth.PUT("/:id", h.Update)
			auth.DELETE("/:id", h.Delete)
		}
	}

	return r, svc
}
