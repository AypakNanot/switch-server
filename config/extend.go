package config

var ExtConfig Extend

// Extend 扩展配置
//  extend:
//    demo:
//      name: demo-name
// 使用方法： config.ExtConfig......即可！！
type Extend struct {
	AMap AMap   // 这里配置对应配置文件的结构即可

	// Runtime 配置（用于内存优化）
	Runtime RuntimeConfig `yaml:"runtime" json:"runtime"`

	// Application 扩展配置（用于功能开关）
	ApplicationEx ApplicationExConfig `yaml:"applicationEx" json:"applicationEx"`
}

type AMap struct {
	Key string
}

// RuntimeConfig 运行时内存调优配置
type RuntimeConfig struct {
	// GOMAXPROCS 设置（0 = 自动检测）
	GoMaxProcs int `yaml:"gomaxprocs" json:"gomaxprocs"`

	// GOGC 设置（100 = 默认）
	GOGC int `yaml:"gogc" json:"gogc"`

	// 软内存限制（MB，0 = 不限制）
	MemoryLimit int `yaml:"memoryLimit" json:"memoryLimit"`

	// 最大线程数（0 = 不限制）
	MaxThreads int `yaml:"maxThreads" json:"maxThreads"`
}

// ApplicationExConfig 应用程序扩展配置
type ApplicationExConfig struct {
	// 是否启用前端静态文件（默认: true）
	EnableFrontend bool `yaml:"enableFrontend" json:"enableFrontend"`

	// 中间件开关
	EnableMiddleware MiddlewareConfig `yaml:"enableMiddleware" json:"enableMiddleware"`
}

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	// Sentinel 限流中间件（默认: true）
	Sentinel bool `yaml:"sentinel" json:"sentinel"`

	// RequestID 中间件（默认: true）
	RequestID bool `yaml:"requestID" json:"requestID"`

	// Metrics 中间件（默认: false）
	Metrics bool `yaml:"metrics" json:"metrics"`
}
