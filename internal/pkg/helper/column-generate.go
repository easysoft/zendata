package helper

import (
	"fmt"
	"math"
	"strconv"
)

func GenBit() (ret, note string) {
	ret = "0,1:R"
	note = "bit"
	return
}
func GenTinyint(hasSign bool) (ret, note string) {
	if hasSign {
		ret = "-128-127:R"
	} else {
		ret = "0-255:R"
	}
	note = "tinyint 2^8"
	return
}
func GenSmallint(hasSign bool) (ret, note string) {
	if hasSign {
		ret = "-32768-32767:R"
	} else {
		ret = "0-65535:R"
	}
	note = "smallint 2^16"
	return
}
func GenMediumint(hasSign bool) (ret, note string) {
	if hasSign {
		ret = "-8388608-8388607:R"
	} else {
		ret = "0-16777215:R"
	}

	note = "mediumint 2^24"
	return
}
func GenInt(hasSign bool) (ret, note string) {
	if hasSign {
		ret = "-2147483648-2147483647:R"
	} else {
		ret = "0-4294967295:R"
	}

	note = "int 2^32"
	return
}
func GenBigint(hasSign bool) (ret, note string) {
	//MaxInt64  = 1<<63 - 1           // 9223372036854775807
	//MinInt64  = -1 << 63            // -9223372036854775808

	if hasSign {
		ret = fmt.Sprintf("%d-%d:R", math.MinInt32, math.MaxInt32) // compatible for 32 bits system
	} else {
		ret = fmt.Sprintf("0-%d:R", math.MaxInt32) // compatible for 32 bits system
	}

	note = "bigint 2^64"
	return
}

func GenFloat(hasSign bool) (ret, note string) {
	if hasSign {
		ret = "-99-99.999:R"
	} else {
		ret = "0-99.999:R"
	}

	note = "float"
	return
}
func GenDouble(hasSign bool) (ret, note string) {
	if hasSign {
		ret = "-99-99.999999:R"
	} else {
		ret = "0-99.999999:R"
	}

	note = "double"
	return
}
func GenDecimal(hasSign bool) (ret, note string) {
	if hasSign {
		ret = "-99-99.99:R"
	} else {
		ret = "0-99.99:R"
	}

	note = "decimal"
	return
}

func GenChar(param string) (ret string, loop string) {
	ret = `a-z`

	loopInt, _ := strconv.Atoi(param)
	if loopInt > 0 {
		loop = fmt.Sprintf("%d", loopInt)
	}

	return
}

func GenBin() (from, format string) {
	format = "binary()"
	return
}

// date time
func GenDate() (rang, typ, format string) {
	rang = `"(-1M)-(+1w):86400"`
	typ = "timestamp"
	format = `"YY/MM/DD"`

	return
}
func GenTime() (rang, typ, format string) {
	rang = `"(-1M)-(+1w):60"`
	typ = "timestamp"
	format = `"hh:mm:ss"`

	return
}
func GenYear() (rang, typ, format string) {
	rang = `"(-6Y)-(+6Y):31536000"`
	typ = "timestamp"
	format = `"YYYY"`

	return
}

func GenDatetime() (rang, typ, format string) {
	rang = `"(-1M)-(+1w):60"`
	typ = "timestamp"
	format = `"YY/MM/DD hh:mm:ss"`

	return
}
func GenTimestamp() (rang, typ, format string) {
	rang = `"(-1M)-(+1w):60"`
	typ = "timestamp"
	format = `""`

	return
}
