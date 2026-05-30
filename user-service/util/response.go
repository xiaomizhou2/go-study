package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
// 对比 Java: class ApiResponse<T> { int code; String message; T data; }
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 返回成功响应
// 对比 Java: return ResponseEntity.ok(new ApiResponse<>(0, "success", data));
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Created 返回创建成功响应（201）
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "created",
		Data:    data,
	})
}

// Error 返回错误响应
// 对比 Java: throw new BusinessException(code, message);
func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}
