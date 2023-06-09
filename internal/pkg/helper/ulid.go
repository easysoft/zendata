package helper

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

func GenerateUlid(format string) (ret string) {
	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

	ret = ulid.MustNew(ulid.Timestamp(t), entropy).String()

	return
}
