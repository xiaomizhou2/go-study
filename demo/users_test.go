package demo

import "testing"

func TestUser_IsAdult(t *testing.T) {
	tests := []struct {
		name string
		age  int
		want bool
	}{
		{"成年", 25, true},
		{"刚成年", 18, true},
		{"未成年", 15, false},
		{"零岁", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUser(1, "测试", tt.age)
			if got := u.IsAdult(); got != tt.want {
				t.Errorf("IsAdult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_FormatInfo(t *testing.T) {
	user := NewUser(1, "张三", 25)
	got := user.FormatInfo()
	want := "[用户#1] 张三, 25岁 (成年)"
	if got != want {
		t.Errorf("FormatInfo() = %v, want %v", got, want)
	}
}

func TestUser_Rename(t *testing.T) {
	user := NewUser(1, "张三", 25)
	user.Rename("张三丰")
	if user.Name != "张三丰" {
		t.Errorf("Rename() 后 Name = %v, want 张三丰", user.Name)
	}
}
