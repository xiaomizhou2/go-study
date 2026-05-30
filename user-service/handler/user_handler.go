package handler

import (
	"net/http"
	"strconv"

	"github.com/example/user-service/model"
	"github.com/example/user-service/service"
	"github.com/example/user-service/util"
	"github.com/gin-gonic/gin"
)

// UserHandler 用户 HTTP 处理器
// 对比 Java: @RestController @RequestMapping("/users") public class UserController
//
// Go 没有注解，handler 就是普通 struct + 方法，在 main.go 里手动注册路由
type UserHandler struct {
	svc *service.UserService
}

// NewUserHandler 构造函数
// 对比 Java: @Autowired public UserController(UserService svc)
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GetAll 获取所有用户
// 对比 Java: @GetMapping public List<User> getAll()
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.svc.GetAll()
	if err != nil {
		util.Error(c, http.StatusInternalServerError, 500, "获取用户列表失败")
		return
	}
	util.Success(c, users)
}

// GetByID 按 ID 获取用户
// 对比 Java: @GetMapping("/{id}") public User getById(@PathVariable int id)
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.Error(c, http.StatusBadRequest, 400, "无效的 ID")
		return
	}

	user, err := h.svc.GetByID(id)
	if err != nil {
		util.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	util.Success(c, user)
}

// Create 创建用户
// 对比 Java: @PostMapping public User create(@RequestBody @Valid User user)
func (h *UserHandler) Create(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, 400, "请求参数错误: "+err.Error())
		return
	}

	user, err := h.svc.Create(req.Name, req.Age)
	if err != nil {
		util.Error(c, http.StatusInternalServerError, 500, "创建用户失败")
		return
	}

	util.Created(c, user)
}

// Update 更新用户
// 对比 Java: @PutMapping("/{id}") public User update(@PathVariable int id, @RequestBody User updates)
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.Error(c, http.StatusBadRequest, 400, "无效的 ID")
		return
	}

	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Error(c, http.StatusBadRequest, 400, "请求参数错误: "+err.Error())
		return
	}

	user, err := h.svc.Update(id, req.Name, req.Age)
	if err != nil {
		util.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	util.Success(c, user)
}

// Delete 删除用户
// 对比 Java: @DeleteMapping("/{id}") public void delete(@PathVariable int id)
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.Error(c, http.StatusBadRequest, 400, "无效的 ID")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		util.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
