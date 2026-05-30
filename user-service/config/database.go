package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// InitDB 根据配置初始化数据库连接
// 对比 Java: @Bean public DataSource dataSource(DatabaseConfig cfg) { ... }
//
// 换数据库只需改 config.yaml：
//   PostgreSQL: driver: "postgres", dsn: "postgres://user:pass@localhost/db"
//   MySQL:      driver: "mysql",    dsn: "user:pass@tcp(localhost:3306)/db"
//   SQLite:     driver: "sqlite",   dsn: "./user.db"
func InitDB(cfg DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	log.Println("数据库连接成功")

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("建表失败: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	sql := `
	CREATE TABLE IF NOT EXISTS users (
		id   INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age  INTEGER NOT NULL
	);`

	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	log.Println("数据表就绪")
	return nil
}
