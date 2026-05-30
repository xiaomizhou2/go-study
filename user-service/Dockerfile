# ============================================================
# 阶段1：构建（使用完整的 Go 镜像编译）
# 对比 Java: FROM maven:3.9 AS builder → mvn package
# ============================================================
FROM golang:1.25-alpine AS builder

# 使用国内 Go 代理加速依赖下载
# 对比 Java: settings.xml 里配阿里云 Maven 镜像
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /build

# 先复制依赖文件，利用 Docker 缓存层
# 只要 go.mod/go.sum 没变，这一层就不会重新下载依赖
# 对比 Java: COPY pom.xml → mvn dependency:resolve
COPY go.mod go.sum ./
RUN go mod download

# 再复制源码并编译
COPY . .

# CGO_ENABLED=0: 静态编译，不依赖 C 库
# 这样跑在 alpine（精简 Linux）上也不会缺动态库
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o user-service .

# ============================================================
# 阶段2：运行（使用最小镜像，只拷贝编译产物）
# 对比 Java: FROM eclipse-temurin:17-jre-alpine → COPY *.jar
# ============================================================
FROM alpine:latest

# 安装 CA 证书（HTTPS 请求需要）
RUN apk --no-cache add ca-certificates

WORKDIR /app

# 从构建阶段只拷贝二进制文件和配置文件
# 对比 Java: COPY --from=builder /app/target/*.jar app.jar
COPY --from=builder /build/user-service .
COPY --from=builder /build/config.yaml .

EXPOSE 8080

CMD ["./user-service"]
