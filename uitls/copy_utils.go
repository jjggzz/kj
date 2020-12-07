package uitls

import (
	"bytes"
	"encoding/gob"
)

// 传入两个结构体的地址,将src的可导出字段的值拷贝到dst的可导出字段上
// 只有相同字段名的值才会被拷贝,多余的值会被忽略,过长的字段会被截断(int64->int)
func Copy(src interface{}, dst interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
