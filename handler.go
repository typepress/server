package server

import (
	"github.com/typepress/core"
	"github.com/typepress/i18n"
	"github.com/typepress/static"
)

func init() {
	core.Handler(
		static.Handler,
		i18n.Translate("zh-cn"),
	)
}
