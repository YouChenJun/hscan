package libs

import "fmt"

const (
	VERSION = "1.0.1"
	DESC    = "进攻性全自动信息搜集"
	BINARY  = "hscan"
	AUTHOR  = "ChenDark"
)

var LOGDIR = fmt.Sprintf("/tmp/%s-log/", BINARY)
