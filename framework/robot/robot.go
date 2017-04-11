// Package robot 将动态插件所依赖的接口提取为一个独立的包（因为plugin机制不允许加载同主程序版本不一致的动态库）
package robot

// Robot 机器人接口
type Robot interface {
	OK() bool
	Do(name string) error
}
