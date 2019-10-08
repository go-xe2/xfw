package xerrors

import (
	"bytes"
	"github.com/gogf/gf/g/errors/gerror"
	"github.com/gogf/gf/g/util/gconv"
)

func New(msg ...interface{}) error {
	var buf bytes.Buffer
	for _, s := range msg {
		buf.WriteString(gconv.String(s))
	}
	return gerror.NewText(buf.String())
}
