package util

import (
	"fmt"
	"time"
	"github.com/lazybeaver/xorshift"
)

type (
	UtilHandler struct {
	}
)

func (u *UtilHandler) RandString128() string {
	xor1024 := xorshift.NewXorShift1024Star(uint64(time.Now().Nanosecond()))
	code := ""
	for i := 0; i < 8; i++ {
		code = code + fmt.Sprintf("%x", xor1024.Next())
	}
	return code
}


func (u *UtilHandler) RandString64() string {
	xor1024 := xorshift.NewXorShift1024Star(uint64(time.Now().Nanosecond()))
	code := ""
	for i := 0; i < 4; i++ {
		code = code + fmt.Sprintf("%x", xor1024.Next())
	}
	return code
}
