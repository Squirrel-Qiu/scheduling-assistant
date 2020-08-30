package internal

import (
	"schedule/dbb"
	"schedule/snowid"
	"schedule/wechatid"
)

type Implement struct {
	DB           dbb.DBApi
	OpenidGetter wechatid.OpenidGetter
	RotaidGetter snowid.RotaidGetter
}
