package helper

import (
	"fmt"
	"math/rand"
	"strings"
)

func GenerateMac() (ret string) {
	buf := make([]byte, 6)

	rand.Read(buf)

	buf[0] |= 2
	ret = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])

	ret = strings.ToUpper(ret)

	return
}
