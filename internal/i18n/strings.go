package i18n

func T(lang Lang, key string) string {
	if lang == LangZH {
		if val, ok := zh[key]; ok {
			return val
		}
	}
	if val, ok := en[key]; ok {
		return val
	}
	return key
}

var en = map[string]string{
	"root.short":              "Composable AI pipeline for logs/text.",
	"root.long":               "aip is a CLI tool for building Unix-style AI pipelines over logs and text.",
	"cmd.summary.short":       "Single-pass: whole input → whole output.",
	"cmd.map.short":           "Chunked processing with streaming output.",
	"cmd.watch.short":         "Windowed processing for long-running streams.",
	"cmd.norm.short":          "Normalize raw logs into signatures + meta.",
	"cmd.reduce.short":        "Aggregate records by key (top-k, time range, samples).",
	"cmd.cluster.short":       "Approximate clustering for signatures.",
	"cmd.sample.short":        "Sample raw records from top-k sig/cluster.",
	"cmd.diagnose.short":      "Opinionated pipeline for log diagnosis.",
	"cmd.cache.short":         "Cache management.",
	"cmd.config.short":        "Configuration management.",
	"cmd.config.show.short":   "Show current config.",
	"cmd.config.show.merged":  "Show merged config (file + env).",
	"cmd.config.path.short":   "Show config file path.",
	"cmd.config.get.short":    "Get a config value by key.",
	"cmd.config.set.short":    "Set a config value by key.",
	"cmd.config.wizard.short": "Interactive config wizard.",
	"cmd.version.short":       "Show version information.",
	"err.not_implemented":     "not implemented yet",
	"err.command_disabled":    "command disabled",
	"err.config_missing":      "config not found",
	"err.config_key":          "unknown config key",
	"msg.config_saved":        "config saved",
}

var zh = map[string]string{
	"root.short":              "面向日志/文本的可组合 AI 管道工具。",
	"root.long":               "aip 是用于构建 Unix 风格日志/文本 AI 管道的命令行工具。",
	"cmd.summary.short":       "一次性处理：整体输入 → 整体输出。",
	"cmd.map.short":           "分块处理：流式输出结果。",
	"cmd.watch.short":         "长流输入：窗口化处理。",
	"cmd.norm.short":          "归一化：raw → sig + meta。",
	"cmd.reduce.short":        "聚合：按 key 统计 top-k、时间范围、样本。",
	"cmd.cluster.short":       "近似聚类：签名聚类。",
	"cmd.sample.short":        "回查样本：对 top-k sig/cluster 抽样。",
	"cmd.diagnose.short":      "封装流水线：面向日志诊断。",
	"cmd.cache.short":         "缓存管理。",
	"cmd.config.short":        "配置管理。",
	"cmd.config.show.short":   "显示当前配置。",
	"cmd.config.show.merged":  "显示合并配置（文件+环境变量）。",
	"cmd.config.path.short":   "显示配置文件路径。",
	"cmd.config.get.short":    "按 key 读取配置值。",
	"cmd.config.set.short":    "按 key 设置配置值。",
	"cmd.config.wizard.short": "交互式配置向导。",
	"cmd.version.short":       "版本信息。",
	"err.not_implemented":     "尚未实现",
	"err.command_disabled":    "命令被禁用",
	"err.config_missing":      "配置不存在",
	"err.config_key":          "未知配置 key",
	"msg.config_saved":        "配置已保存",
}
