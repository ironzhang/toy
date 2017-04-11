// 为避免framework的修改影响动态插件，所以将动态插件所依赖的接口单独提取为一个包
package robot

// Robot 机器人接口
type Robot interface {
	OK() bool
	Do(name string) error
}
