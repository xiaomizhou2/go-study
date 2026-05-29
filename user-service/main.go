package main

import (
    "fmt"
    "log"
    "net/http"
    "strconv"

    "github.com/example/user-service/model"
    "github.com/gin-gonic/gin"
)

// 内存数据存储（模拟数据库）
var users = []model.User{
    {ID: 1, Name: "张三", Age: 25},
    {ID: 2, Name: "李四", Age: 30},
    {ID: 3, Name: "王五", Age: 28},
}

// 自增 ID 计数器（模拟数据库自增）
var nextID = 4

func main() {
    // 创建 Gin 路由器
    // gin.Default() 包含 Logger 和 Recovery 中间件
    // 对比 Java: new SpringApplicationBuilder().sources(AppConfig.class).run()
    r := gin.Default()

    // 路由组 - v1 API
    // 对比 Java: @RequestMapping("/api/v1")
    api := r.Group("/api/v1")
    {
        // 用户相关路由
        usersGroup := api.Group("/users")
        {
            usersGroup.GET("", getAllUsers)         // GET /api/v1/users
            usersGroup.GET("/:id", getUserByID)      // GET /api/v1/users/:id
            usersGroup.POST("", createUser)          // POST /api/v1/users
            usersGroup.PUT("/:id", updateUser)       // PUT /api/v1/users/:id
            usersGroup.DELETE("/:id", deleteUser)    // DELETE /api/v1/users/:id
        }

        // 健康检查
        api.GET("/health", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "status": "ok",
                "service": "user-service",
            })
        })
    }

    // 启动服务
    // 对比 Java: SpringApplication.run(App.class, args)
    port := ":8080"
    fmt.Printf("🚀 服务启动在 http://localhost%s\n", port)
    fmt.Printf("📍 可用的端点:\n")
    fmt.Printf("   GET    /api/v1/users       - 获取所有用户\n")
    fmt.Printf("   GET    /api/v1/users/:id   - 按 ID 获取用户\n")
    fmt.Printf("   POST   /api/v1/users       - 创建新用户\n")
    fmt.Printf("   PUT    /api/v1/users/:id   - 更新用户\n")
    fmt.Printf("   DELETE /api/v1/users/:id   - 删除用户\n")
    fmt.Printf("   GET    /api/v1/health      - 健康检查\n")

    if err := r.Run(port); err != nil {
        log.Fatalf("服务启动失败: %v", err)
    }
}

// getAllUsers 获取所有用户
// 对比 Java: @GetMapping public List<User> getAllUsers()
func getAllUsers(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "code":    0,
        "message": "success",
        "data":    users,
    })
}

// getUserByID 按 ID 获取用户
// 对比 Java: @GetMapping("/{id}") public User getUserById(@PathVariable Long id)
func getUserByID(c *gin.Context) {
    // 获取路径参数
    // 对比 Java: @PathVariable Long id
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    400,
            "message": "无效的 ID",
            "error":   err.Error(),
        })
        return
    }

    // 查找用户
    for _, user := range users {
        if user.ID == id {
            c.JSON(http.StatusOK, gin.H{
                "code":    0,
                "message": "success",
                "data":    user,
            })
            return
        }
    }

    // 未找到
    c.JSON(http.StatusNotFound, gin.H{
        "code":    404,
        "message": "用户不存在",
    })
}

// createUser 创建新用户
// 对比 Java: @PostMapping public User createUser(@RequestBody @Valid User user)
func createUser(c *gin.Context) {
    var newUser model.User

    // 绑定请求体到结构体，并自动验证
    // 对比 Java: @RequestBody @Valid User user
    if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    400,
            "message": "请求参数错误",
            "error":   err.Error(),
        })
        return
    }

    // 分配 ID
    newUser.ID = nextID
    nextID++

    // 添加到列表
    users = append(users, newUser)

    // 返回创建的用户
    c.JSON(http.StatusCreated, gin.H{
        "code":    0,
        "message": "用户创建成功",
        "data":    newUser,
    })
}

// updateUser 更新用户
// 对比 Java: @PutMapping("/{id}") public User updateUser(@PathVariable Long id, @RequestBody User updates)
func updateUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    400,
            "message": "无效的 ID",
        })
        return
    }

    var updates model.User
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    400,
            "message": "请求参数错误",
            "error":   err.Error(),
        })
        return
    }

    // 查找并更新用户
    for i, user := range users {
        if user.ID == id {
            // 只更新提供的字段（部分更新）
            if updates.Name != "" {
                users[i].Name = updates.Name
            }
            if updates.Age > 0 {
                users[i].Age = updates.Age
            }

            c.JSON(http.StatusOK, gin.H{
                "code":    0,
                "message": "用户更新成功",
                "data":    users[i],
            })
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{
        "code":    404,
        "message": "用户不存在",
    })
}

// deleteUser 删除用户
// 对比 Java: @DeleteMapping("/{id}") public void deleteUser(@PathVariable Long id)
func deleteUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    400,
            "message": "无效的 ID",
        })
        return
    }

    // 查找并删除用户
    for i, user := range users {
        if user.ID == id {
            // 删除元素
            users = append(users[:i], users[i+1:]...)

            c.Status(http.StatusNoContent) // 204 No Content
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{
        "code":    404,
        "message": "用户不存在",
    })
}
