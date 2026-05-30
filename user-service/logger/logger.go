package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log 全局 Logger 实例
// 对比 Java: private static final Logger log = LoggerFactory.getLogger(App.class);
//
// Go 用包级变量，所有代码通过 logger.Log.Info(...) 访问
var Log *zap.Logger

// Config 日志配置
// 对比 Java: logging.level.root=INFO, logging.pattern.console=%d{...}
type Config struct {
	Level  string `mapstructure:"level"`  // debug, info, error
	Format string `mapstructure:"format"` // json, text
}

// Init 初始化日志
// 对比 Java: logback-spring.xml 里配 Appender、Encoder、Level
//
// Go 的 zap 初始化只需几行代码：
//   json 格式 → 适合生产（方便 ELK/Grafana Loki 采集）
//   text 格式 → 适合开发（人类可读，彩色输出）
func Init(cfg Config) error {
	var zapCfg zap.Config

	if cfg.Format == "json" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		// 开发模式用更可读的编码器
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// 设置日志级别（对比 Java: logging.level.root=INFO）
	switch cfg.Level {
	case "debug":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "error":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	var err error
	Log, err = zapCfg.Build()
	if err != nil {
		return err
	}

	return nil
}

// Sync 刷新缓冲区（程序退出前调用）
// 对比 Java: LogManager.shutdown()
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
