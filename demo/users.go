package demo

import "fmt"

// 对比 Java: public class User { private int id; private String name; ... }
//
// Go struct 语法要点：
//   1. type User struct 开头（不是 class）
//   2. 字段用换行分隔，不用逗号
//   3. 首字母大写 = public，小写 = private
type User struct {
	ID   int
	Name string
	Age  int
}

// NewUser 构造函数
// 对比 Java: public User(int id, String name, int age) { this.id = id; ... }
//
// Go 没有 new 关键字的构造函数，用 NewXxx 工厂函数代替
func NewUser(id int, name string, age int) *User {
	return &User{
		ID:   id,
		Name: name,
		Age:  age,
	}
}

// IsAdult 判断是否成年
// 对比 Java: public boolean isAdult() { return this.age >= 18; }
//
// (u *User) 是指针接收者，相当于 Java 的 this
func (u *User) IsAdult() bool {
	return u.Age >= 18
}

// FormatInfo 格式化输出用户信息
// 对比 Java: public String formatInfo() { return String.format(...); }
func (u *User) FormatInfo() string {
	adult := ""
	if u.IsAdult() {
		adult = " (成年)"
	} else {
		adult = " (未成年)"
	}
	return fmt.Sprintf("[用户#%d] %s, %d岁%s", u.ID, u.Name, u.Age, adult)
}

// Rename 重命名
// 对比 Java: public void rename(String newName) { this.name = newName; }
//
// 必须用指针接收者 (*User)，因为要修改字段值
// 如果用值接收者 (u User)，修改的是副本，原对象不会变
func (u *User) Rename(newName string) {
	u.Name = newName
}
