package browser

type IBrowser interface {
	BaseArgs() string					// 基础参数

	ProgramDir() string                // 程序目录
	Execable() string                  // 可执行文件路径
	Profile() string                   // 配置名
	ProfileDir() string                // 配置文件路径
	Incognito() string                 // 匿名/无痕/隐私模式
	ProxyServer(address string) string // 代理接口

	Tpl() string // 默认的命令模板
}

func appendFuncMap(browser IBrowser) {

}