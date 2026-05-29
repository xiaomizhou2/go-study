package model

// User 用户模型
// 对比 Java: 相当于 User Entity 或 POJO
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name" binding:"required,min=2,max=50"`
    Age  int    `json:"age" binding:"required,gte=0,lte=150"`
}

// 字段标签说明：
// `json:"id"`                      -> JSON 序列化时字段名为 "id"
// `binding:"required"`             -> 请求时必填（类似 @NotNull）
// `binding:"min=2,max=50"`         -> 字符串长度限制（类似 @Size）
// `binding:"gte=0,lte=150"`        -> 数值范围（类似 @Min/@Max）
//
// 对比 Java (Bean Validation):
// @NotNull
// @Size(min = 2, max = 50)
// @Min(0) @Max(150)
// private String name;

// 字段标签说明：
// `json:"id"`     -> JSON 序列化时字段名为 "id"
// `json:"name"`   -> JSON 序列化时字段名为 "name"
// `json:"age"`    -> JSON 序列化时字段名为 "age"
//
// 对比 Java:
// @JsonProperty("id")
// private int id;
