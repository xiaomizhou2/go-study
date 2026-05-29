# 第一部分：语言核心差异

> Java 老兵的 Go 语言生存指南 —— 用比喻打通任督二脉 🧠

---

## 1. 并发模型：goroutine + channel vs Java 线程池 + Future

### 🎯 生动比喻

| Java 线程 | Go goroutine |
|-----------|-------------|
| 重型装甲车 🚛 | 轻量级绿色小精灵 🧚 |
| 启动一辆装甲车需要大量燃料（内存 ~1MB 栈空间） | 召唤一个小精灵几乎不花力气（~2KB 栈空间，可动态增长） |
| 你只能同时派出有限的装甲车（几百到几千） | 你可以同时召唤数万甚至数十万小精灵 |
| 装甲车之间通过复杂的无线电通信（共享内存 + 锁） | 小精灵们通过竹筒传话（channel）优雅地协作 |

**核心理念：** Go 的并发哲学是 *"Don't communicate by sharing memory; share memory by communicating."*
（不要通过共享内存来通信，而要通过通信来共享内存。）

---

### Java 代码：线程池 + Future

```java
import java.util.concurrent.*;
import java.util.ArrayList;
import java.util.List;

public class ConcurrentDemo {
    public static void main(String[] args) throws Exception {
        ExecutorService executor = Executors.newFixedThreadPool(10);
        List<Future<String>> futures = new ArrayList<>();

        // 提交 5 个任务
        for (int i = 0; i < 5; i++) {
            final int taskId = i;
            futures.add(executor.submit(() -> {
                Thread.sleep(1000); // 模拟耗时操作
                return "任务 " + taskId + " 完成";
            }));
        }

        // 收集结果
        for (Future<String> future : futures) {
            System.out.println(future.get()); // 阻塞等待
        }

        executor.shutdown(); // 别忘了关！
    }
}
```

### Go 代码：goroutine + channel

```go
package main

import (
    "fmt"
    "time"
)

func worker(id int, ch chan<- string) {
    time.Sleep(1 * time.Second) // 模拟耗时操作
    ch <- fmt.Sprintf("任务 %d 完成", id) // 通过 channel 发送结果
}

func main() {
    ch := make(chan string, 5) // 带缓冲的 channel

    // 启动 5 个 goroutine
    for i := 0; i < 5; i++ {
        go worker(i, ch) // go 关键字，一行搞定！
    }

    // 收集结果
    for i := 0; i < 5; i++ {
        fmt.Println(<-ch) // 从 channel 接收
    }
}
```

### 🔑 关键差异总结

| 维度 | Java | Go |
|------|------|-----|
| 启动并发 | `executor.submit(task)` | `go func()` 一行搞定 |
| 通信方式 | 共享变量 + 锁 / Future.get() | Channel（像管道一样自然） |
| 资源开销 | 线程 ~1MB 栈 | goroutine ~2KB 栈 |
| 数量上限 | 通常几百~几千 | 轻松数十万 |
| 错误处理 | try-catch 或 Future.get() 异常 | channel 传递 error |
| 资源清理 | 必须手动 shutdown | goroutine 完成后自动回收 |

### 💡 Go 进阶：用 sync.WaitGroup 等待完成

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    var wg sync.WaitGroup
    results := make(chan string, 5)

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()           // 函数结束时通知 WaitGroup
            time.Sleep(1 * time.Second)
            results <- fmt.Sprintf("任务 %d 完成", id)
        }(i) // 注意：把 i 作为参数传入，避免闭包陷阱！
    }

    // 等待所有 goroutine 完成后关闭 channel
    go func() {
        wg.Wait()
        close(results)
    }()

    // range 会自动在 channel 关闭时停止
    for result := range results {
        fmt.Println(result)
    }
}
```

> ⚠️ **闭包陷阱提醒：** 在 `for` 循环中启动 goroutine 时，如果直接使用循环变量 `i`，所有 goroutine 可能会看到同一个值。解决方案是将 `i` 作为参数传入（如上面的 `(id int)`）。

---

## 2. 错误处理：多返回值 error vs try-catch

### 🎯 生动比喻

| Java 异常 | Go error |
|-----------|---------|
| 像火灾警报器 🔥 —— 一旦触发，所有事情都停下来，层层上报 | 像快递单上的备注 📦 —— 每一步都检查有没有问题，有就当场处理 |
| 你永远不知道哪个角落藏着一个 unchecked RuntimeException | 错误就在返回值里，想忽略都难（编译器会提醒） |
| try-catch 像在楼下铺一张大网，等楼上抛东西下来 | Go 的 if err != nil 像每一步都检查脚下有没有坑 |

**核心理念：** Go 认为，错误是正常的业务逻辑的一部分，不应该用"异常"这种特殊机制来处理。
让错误处理显式化、透明化，代码读起来就能看到所有可能的失败路径。

---

### Java 代码：try-catch 异常

```java
import java.io.*;

public class ErrorDemo {
    public static void main(String[] args) {
        try {
            String content = readFile("config.txt");
            int port = parsePort(content);
            System.out.println("端口: " + port);
        } catch (FileNotFoundException e) {
            System.err.println("文件未找到: " + e.getMessage());
        } catch (NumberFormatException e) {
            System.err.println("端口格式错误: " + e.getMessage());
        } catch (IOException e) {
            System.err.println("IO 错误: " + e.getMessage());
        } finally {
            System.out.println("清理工作完成");
        }
    }

    static String readFile(String path) throws IOException {
        StringBuilder sb = new StringBuilder();
        try (BufferedReader reader = new BufferedReader(new FileReader(path))) {
            String line;
            while ((line = reader.readLine()) != null) {
                sb.append(line);
            }
        }
        return sb.toString();
    }

    static int parsePort(String content) throws NumberFormatException {
        return Integer.parseInt(content.trim());
    }
}
```

### Go 代码：多返回值 error

```go
package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    content, err := readFile("config.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "读文件失败: %v\n", err)
        return
    }

    port, err := parsePort(content)
    if err != nil {
        fmt.Fprintf(os.Stderr, "解析端口失败: %v\n", err)
        return
    }

    fmt.Printf("端口: %d\n", port)
}

func readFile(path string) (string, error) { // 多返回值！(result, error)
    data, err := os.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("读文件 %s 失败: %w", path, err) // 用 %w 包装错误
    }
    return string(data), nil // nil 表示没有错误
}

func parsePort(content string) (int, error) {
    port, err := strconv.Atoi(strings.TrimSpace(content))
    if err != nil {
        return 0, fmt.Errorf("端口格式错误: %w", err)
    }
    return port, nil
}
```

### 🔑 关键差异总结

| 维度 | Java | Go |
|------|------|-----|
| 错误表示 | Exception 类层次 | `error` 接口（只要实现 `Error() string`） |
| 错误传递 | throw → catch 层层上抛 | 返回值显式传递 |
| 强制处理 | checked exception 编译器强制 | `if err != nil` 惯用模式 |
| 未处理后果 | RuntimeException 可能被吞掉 | 未使用的返回值编译器会警告 |
| 错误包装 | `Exception.getCause()` | `fmt.Errorf("上下文: %w", err)` + `errors.Is/As` |
| 清理资源 | finally 块 | defer 语句（下一个主题会讲） |

### 💡 Go 错误处理最佳实践

```go
// ✅ 好的做法：每一步都检查错误
func processFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("打开文件失败: %w", err)
    }
    defer f.Close() // 离开函数时自动关闭文件

    // ... 处理文件
    return nil
}

// ❌ 坏的做法：忽略错误（Go 编译器不允许 _ 忽略 error）
data, _ := os.ReadFile("important.txt") // 危险！_ 吞掉了错误

// ✅ 哨兵错误（Sentinel Errors）
var ErrNotFound = errors.New("资源未找到")

func findUser(id int) (*User, error) {
    // ...
    return nil, ErrNotFound
}

// 调用方可以这样判断：
user, err := findUser(42)
if errors.Is(err, ErrNotFound) {
    // 处理未找到的情况
}
```

---

## 3. 接口：隐式实现（鸭子类型）vs 显式 implements

### 🎯 生动比喻

| Java 接口 | Go 接口 |
|-----------|--------|
| 像入职手续 📋 —— 你必须先签合同、填表格，公司才承认你能干活 | 像看能力不看证书 👀 —— 只要你会游泳，你就是鸭子，不需要贴标签 |
| 必须写 `implements Comparable` | 只要你写了 `Compare()` 方法，你就自动满足了 |
| 改接口要改所有实现类 | 接口是按需定义的，谁用谁定义 |
| 接口通常很大（很多方法） | Go 鼓励小接口（1-2 个方法），"接口越小，抽象越强" |

**核心理念：** Go 的接口是隐式的——你不需要声明"我实现了这个接口"，
只要你的类型拥有接口要求的所有方法，编译器就认为你实现了它。这就是所谓的"鸭子类型"：
"如果它走起来像鸭子，叫起来像鸭子，那它就是鸭子。" 🦆

---

### Java 代码：显式 implements

```java
// 定义接口
public interface Notifier {
    void send(String message);
}

// 显式声明实现
public class EmailNotifier implements Notifier {
    @Override
    public void send(String message) {
        System.out.println("发送邮件: " + message);
    }
}

public class SmsNotifier implements Notifier {
    @Override
    public void send(String message) {
        System.out.println("发送短信: " + message);
    }
}

// 使用
public class NotificationService {
    private final Notifier notifier;

    public NotificationService(Notifier notifier) {
        this.notifier = notifier;
    }

    public void notify(String message) {
        notifier.send(message);
    }
}

// 如果想让一个新类支持通知，必须 implements Notifier
public class WechatNotifier implements Notifier {  // 必须显式声明！
    @Override
    public void send(String message) {
        System.out.println("发送微信: " + message);
    }
}
```

### Go 代码：隐式实现（鸭子类型）

```go
package main

import "fmt"

// 定义接口 —— 由使用方定义，不是实现方！
type Notifier interface {
    Send(message string)
}

// EmailNotifier —— 注意：没有 implements 声明！
type EmailNotifier struct {
    SMTP string
}

func (e EmailNotifier) Send(message string) { // 只要有 Send 方法就够了
    fmt.Printf("发送邮件(%s): %s\n", e.SMTP, message)
}

// SmsNotifier —— 同样没有 implements
type SmsNotifier struct {
    Phone string
}

func (s SmsNotifier) Send(message string) {
    fmt.Printf("发送短信(%s): %s\n", s.Phone, message)
}

// 使用接口
func Notify(n Notifier, message string) {
    n.Send(message) // 只要传入的类型有 Send 方法，就能用！
}

func main() {
    email := EmailNotifier{SMTP: "smtp.example.com"}
    sms := SmsNotifier{Phone: "13800138000"}

    Notify(email, "你好，世界！")
    Notify(sms, "验证码 123456")

    // 甚至可以这样 —— 第三方库的类型也能用，只要方法签名匹配！
    // 不需要修改原始类型，不需要 implements
}
```

### 🔑 关键差异总结

| 维度 | Java | Go |
|------|------|-----|
| 声明方式 | `class A implements B` | 隐式，无需声明 |
| 接口定义者 | 通常是接口作者定义 | 使用方定义（谁用谁定义） |
| 接口大小 | 通常较大（5-10个方法不罕见） | 鼓励小接口（1-3个方法） |
| 跨包实现 | 可以但要 import 接口 | 天然支持，零耦合 |
| 空接口 | `Object` 是万能基类 | `interface{}` / `any` 可以接受任何值 |
| 类型断言 | `instanceof` | `v, ok := i.(T)` |

### 💡 Go 接口设计原则

```go
// ✅ Go 风格：小接口，在使用方定义
// io.Reader 只有一个方法 —— Go 标准库中最经典的接口
type Reader interface {
    Read(p []byte) (n int, err error)
}

// io.Writer 也只有一个方法
type Writer interface {
    Write(p []byte) (n int, err error)
}

// 组合接口！
type ReadWriter interface {
    Reader
    Writer
}

// ❌ Java 风格（不要在 Go 中这样做）：
// 在包里定义一个巨大接口，然后让所有类型去实现它
type UserService interface {
    CreateUser()
    GetUser()
    UpdateUser()
    DeleteUser()
    ListUsers()
    // ... 20 个方法 ...
}
```

---

## 4. 对象模型：struct 组合 vs class 继承

### 🎯 生动比喻

| Java 继承 | Go 组合 |
|-----------|--------|
| 像家谱 👨‍👩‍👧‍👦 —— 子承父业，血缘关系绑定一生 | 像乐高积木 🧱 —— 需要什么功能就拼什么模块 |
| 你爸是医生，你就自动有行医资格（is-a 关系） | 你和医生团队签了合同，就有行医资格（has-a 关系） |
| 多重继承是禁忌（钻石问题 💎） | 组合随便拼，没有冲突 |
| 改父类影响所有子类（脆弱基类问题） | 改一个组件不影响其他组件 |

**核心理念：** Go 没有 class，没有 extends，没有继承。
Go 用 struct（结构体）+ 组合（embedding）来组织代码。
设计模式上更偏爱"组合优于继承"。

---

### Java 代码：class 继承

```java
// 基类
public class Animal {
    protected String name;

    public Animal(String name) {
        this.name = name;
    }

    public void eat() {
        System.out.println(name + " 在吃东西");
    }
}

// 继承
public class Dog extends Animal {
    private String breed;

    public Dog(String name, String breed) {
        super(name);  // 必须调用父类构造器
        this.breed = breed;
    }

    public void bark() {
        System.out.println(name + " 汪汪汪！");
    }
}

// 多层继承（容易变得复杂）
public class GuideDog extends Dog {
    private String owner;

    public GuideDog(String name, String breed, String owner) {
        super(name, breed);  // 一层一层往上穿
        this.owner = owner;
    }

    public void guide() {
        System.out.println(name + " 正在引导 " + owner);
    }
}
```

### Go 代码：struct 组合

```go
package main

import "fmt"

// 基础结构体
type Animal struct {
    Name string
}

func (a Animal) Eat() {
    fmt.Printf("%s 在吃东西\n", a.Name)
}

// 组合：嵌入 Animal（不是继承！）
type Dog struct {
    Animal       // 嵌入（embedding），Dog "拥有" Animal 的所有字段和方法
    Breed string
}

func (d Dog) Bark() {
    fmt.Printf("%s 汪汪汪！\n", d.Name) // 可以直接访问 Animal.Name
}

// 继续组合
type GuideDog struct {
    Dog           // 嵌入 Dog（Dog 里又嵌入了 Animal）
    Owner string
}

func (g GuideDog) Guide() {
    fmt.Printf("%s 正在引导 %s\n", g.Name, g.Owner)
}

func main() {
    dog := Dog{
        Animal: Animal{Name: "旺财"},
        Breed:  "柴犬",
    }
    dog.Eat()  // 通过嵌入 "继承" 了 Animal 的方法
    dog.Bark()

    guide := GuideDog{
        Dog:   Dog{Animal: Animal{Name: "导导"}, Breed: "拉布拉多"},
        Owner: "张三",
    }
    guide.Eat()   // Animal 的方法
    guide.Bark()  // Dog 的方法
    guide.Guide() // 自己的方法
}
```

### 🔑 关键差异总结

| 维度 | Java | Go |
|------|------|-----|
| 定义数据 | class（数据和方法的容器） | struct（纯数据） + method（独立绑定） |
| 复用机制 | extends 继承 | embedding 组合 |
| 多态 | 通过继承 + 方法重写 | 通过接口（隐式实现） |
| 构造 | new + 构造函数 | 工厂函数 `NewXxx()` |
| 方法绑定 | 方法属于类 | 方法可以绑定到任何类型（包括基本类型别名） |
| nil 安全 | NullPointerException | 可以检查 nil（但也要小心 nil 指针） |

### 💡 Go 的方法可以绑定到任何类型

```go
// Go 甚至可以给基本类型起别名后绑定方法！Java 做不到
type Celsius float64

func (c Celsius) ToFahrenheit() float64 {
    return float64(c)*9/5 + 32
}

func main() {
    temp := Celsius(100)
    fmt.Printf("%.1f°C = %.1f°F\n", temp, temp.ToFahrenheit())
    // 输出: 100.0°C = 212.0°F
}
```

---

## 5. 包与模块：go mod vs Maven/Gradle

### 🎯 生动比喻

| Maven/Gradle | go mod |
|-------------|--------|
| 像一个庞大的仓库管理员 🏢 —— pom.xml/build.gradle 是你的超级购物清单 | 像一个简约的便利店老板 🏪 —— go.mod 只有你真正需要的东西 |
| 依赖传递复杂（mvn dependency:tree） | 直接记录最小版本（Minimal Version Selection） |
| 配置文件动辄上百行 | go.mod 通常十几行搞定 |
| 仓库中心（Maven Central） | 去中心化（直接引用 Git 仓库） |

---

### Java：Maven pom.xml

```xml
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>com.example</groupId>
    <artifactId>user-service</artifactId>
    <version>1.0.0</version>
    <packaging>jar</packaging>

    <parent>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-parent</artifactId>
        <version>3.2.0</version>
    </parent>

    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>
        <dependency>
            <groupId>org.postgresql</groupId>
            <artifactId>postgresql</artifactId>
        </dependency>
    </dependencies>

    <build>
        <plugins>
            <plugin>
                <groupId>org.springframework.boot</groupId>
                <artifactId>spring-boot-maven-plugin</artifactId>
            </plugin>
        </plugins>
    </build>
</project>
```

### Go：go.mod

```go
module github.com/example/user-service

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/lib/pq v1.10.9
)
```

**就这么简单！** 🎉

### 常用命令对比

| 操作 | Maven | Go |
|------|-------|-----|
| 初始化项目 | 创建 pom.xml | `go mod init <module-name>` |
| 添加依赖 | 编辑 pom.xml + `mvn install` | `go get github.com/gin-gonic/gin` |
| 下载依赖 | `mvn install` | `go mod download` |
| 更新依赖 | 修改版本号 + `mvn install` | `go get -u <package>` |
| 查看依赖树 | `mvn dependency:tree` | `go mod graph` |
| 整理依赖 | - | `go mod tidy`（删除无用、添加缺失） |
| 编译 | `mvn compile` | `go build` |
| 运行测试 | `mvn test` | `go test ./...` |
| 打包 | `mvn package` | `go build`（直接产出二进制！不需要容器） |

---

## 6. 访问控制：首字母大小写导出 vs public/private

### 🎯 生动比喻

| Java 访问控制 | Go 访问控制 |
|-------------|-----------|
| 像门禁卡系统 🔐 —— 有 public、protected、private、default 四个等级 | 像一扇简单的门 🚪 —— 大写字母 = 开门（公开），小写字母 = 关门（私有） |
| 四种修饰符，需要仔细区分作用范围 | 只有两种：导出（大写）或未导出（小写） |
| 作用范围是类级别 | 作用范围是包级别（同一个包内都能访问小写成员） |

---

### Java 代码：四种访问级别

```java
public class User {
    public String name;        // 全世界都能访问
    protected int age;         // 同包 + 子类能访问
    String email;              // 同包能访问（default）
    private String password;   // 只有本类能访问

    public User(String name, int age, String email, String password) {
        this.name = name;
        this.age = age;
        this.email = email;
        this.password = password;
    }

    public String getDisplayName() {  // public 方法
        return name + " (" + age + ")";
    }

    private void validatePassword() {  // private 方法
        // ...
    }
}
```

### Go 代码：首字母大小写

```go
package model

// User —— 大写开头，导出（相当于 public）
type User struct {
    Name     string  // 大写开头，导出字段
    Age      int     // 大写开头，导出字段
    email    string  // 小写开头，未导出（包内私有）
    password string  // 小写开头，未导出（包内私有）
}

// NewUser —— 大写开头，导出的工厂函数（构造器）
func NewUser(name string, age int, email, password string) *User {
    return &User{
        Name:     name,
        Age:      age,
        email:    email,
        password: password,
    }
}

// DisplayName —— 大写开头，导出方法
func (u *User) DisplayName() string {
    return fmt.Sprintf("%s (%d)", u.Name, u.Age)
}

// validatePassword —— 小写开头，未导出（包内私有）
func (u *User) validatePassword() bool {
    return len(u.password) >= 8
}
```

### 🔑 关键差异总结

| 维度 | Java | Go |
|------|------|-----|
| 可见性级别 | 4种（public/protected/default/private） | 2种（导出/未导出） |
| 判断方式 | 修饰符关键字 | 首字母大小写 |
| 作用范围 | 类级别 | **包级别**（这是最大的区别！） |
| getter/setter | 习惯用 getXxx()/setXxx() | Go 习惯直接访问字段，需要时用方法封装 |

> 💡 **包级别的含义：** 在 Go 中，同一个包下的所有文件都可以访问小写字母开头的标识符。
> 这和 Java 的 private（类级别）不同。Go 没有类级别的私有，只有包级别的私有。

---

## 7. 其他关键差异：defer、指针、slice

### 7.1 defer：离开房间时自动关灯 💡

**比喻：** defer 就像你在家门口贴了一张便签 "离开时请关灯"，
无论你从哪个门出去（正常返回、提前返回、甚至 panic），灯都会被关上。

```go
func readFile(path string) (string, error) {
    f, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer f.Close() // ← 无论函数怎么退出，f.Close() 都会被调用

    data, err := io.ReadAll(f)
    if err != nil {
        return "", err // 这里返回时，f.Close() 会自动执行
    }

    return string(data), nil // 这里返回时，f.Close() 也会自动执行
}

// 多个 defer 按 LIFO（后进先出）顺序执行
func example() {
    fmt.Println("开始")
    defer fmt.Println("第一个 defer")  // 最后执行
    defer fmt.Println("第二个 defer")  // 先执行
    fmt.Println("结束")
    // 输出顺序：开始 → 结束 → 第二个 defer → 第一个 defer
}
```

**Java 对比：** defer ≈ finally，但更优雅——紧跟资源分配写清理代码，不用跑到 finally 块里。

---

### 7.2 指针：直接操作 vs 间接引用 📍

**比喻：**
- Java 的引用像门牌号 —— 你拿着门牌号找到房子，但门牌号本身不能做加减运算
- Go 的指针像 GPS 坐标 —— 你不仅知道位置，还可以做偏移计算（但日常开发中很少需要）

**好消息：** Go 的指针比 C 简单多了！没有指针运算，只是"地址传递"。

```go
// 值传递 vs 指针传递
type User struct {
    Name string
    Age  int
}

// 值接收者：操作的是副本（就像复印一份文件再修改）
func (u User) BirthdayValue() {
    u.Age++ // 修改的是副本，原始数据不变！
}

// 指针接收者：操作的是原件（直接在原件上改）
func (u *User) BirthdayPointer() {
    u.Age++ // 修改的是原始数据
}

func main() {
    user := User{Name: "张三", Age: 25}

    user.BirthdayValue()
    fmt.Println(user.Age) // 还是 25！副本上的修改不生效

    user.BirthdayPointer()
    fmt.Println(user.Age) // 26！指针修改生效了

    // 取地址用 &，取值用 *
    ptr := &user           // ptr 是 *User 类型
    fmt.Println(ptr.Name)  // Go 自动解引用，不需要 ->（不同于 C）
    fmt.Println((*ptr).Age) // 也可以显式解引用
}
```

**Java 对比：** Java 中对象天然是引用传递（虽然技术上是值传递引用），
Go 需要显式用 `*` 表示指针，但 Go 会自动帮你解引用（不像 C 要用 `->`）。

---

### 7.3 slice：Go 的动态数组，比 Java ArrayList 更底层 🔪

**比喻：**
- Java ArrayList 像一个全自动收纳箱 —— 自动扩容，接口友好
- Go slice 像一把可以伸缩的刀 🔪 —— 底层是数组，但可以灵活切割和追加

```go
func main() {
    // 创建 slice（三种方式）
    nums := []int{1, 2, 3, 4, 5}          // 字面量
    nums = make([]int, 0, 10)              // make(类型, 长度, 容量)
    nums = make([]int, 5)                  // 长度5，容量5

    // 追加（类似 ArrayList.add）
    nums = append(nums, 6, 7, 8)

    // 切片（子数组，类似 subList 但更灵活）
    sub := nums[1:4] // 索引 1 到 3（不包含 4）

    // ⚠️ 重要：slice 是引用类型！修改子切片会影响原切片
    sub[0] = 100
    fmt.Println(nums[1]) // 100！被修改了

    // 安全复制
    copySlice := make([]int, len(sub))
    copy(copySlice, sub)

    // 删除元素（Go 没有内置删除，用切片技巧）
    // 删除索引 2 的元素
    nums = append(nums[:2], nums[3:]...)

    // 遍历
    for i, v := range nums {
        fmt.Printf("索引 %d: 值 %d\n", i, v)
    }
}
```

**Java 对比：**

```java
// Java ArrayList
List<Integer> nums = new ArrayList<>(Arrays.asList(1, 2, 3, 4, 5));
nums.add(6);                           // 添加
List<Integer> sub = nums.subList(1,4); // 子列表
nums.remove(2);                        // 删除（直接有方法！）
for (int i = 0; i < nums.size(); i++) {
    System.out.println(nums.get(i));
}
```

| 操作 | Java ArrayList | Go slice |
|------|---------------|----------|
| 创建 | `new ArrayList<>()` | `make([]T, 0)` 或 `[]T{}` |
| 添加 | `list.add(x)` | `slice = append(slice, x)` ⚠️ 注意重新赋值 |
| 访问 | `list.get(i)` | `slice[i]` |
| 删除 | `list.remove(i)` | `slice = append(slice[:i], slice[i+1:]...)` |
| 长度 | `list.size()` | `len(slice)` |
| 容量 | 无直接对应 | `cap(slice)` |
| 子集 | `subList(from, to)` | `slice[from:to]` |

---

> 🎓 **恭喜你完成了第一部分的学习！**
>
> 现在你已经了解了 Go 和 Java 在以下方面的核心差异：
> 1. 并发模型（goroutine vs 线程）
> 2. 错误处理（error vs exception）
> 3. 接口（隐式 vs 显式）
> 4. 对象模型（组合 vs 继承）
> 5. 包管理（go mod vs Maven）
> 6. 访问控制（大小写 vs 修饰符）
> 7. defer / 指针 / slice
>
> 准备好进入第二部分了吗？🚀
