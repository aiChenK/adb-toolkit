package request

type CommandForm struct {
	Ip          string `form:"ip"`
	Port        string `form:"port"`                  //默认5555
	Op          string `form:"op" binding:"required"` //执行命令，setProxy|delProxy|clear|stop|free
	ProxyAddr   string `form:"proxyAddr"`             //代理地址 setProxy可用
	PackageName string `form:"packageName"`           //包名 clear|stop可用
	Cmd         string `form:"cmd"`
}
