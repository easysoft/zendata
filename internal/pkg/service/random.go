package service

import commonUtils "github.com/easysoft/zendata/pkg/utils/common"

type RandomService struct {
}

func (s *RandomService) RandomValues(values []interface{}) (ret []interface{}) {
	length := len(values)

	for i := 0; i < length; i++ {
		num := commonUtils.RandNum(length * 10000)
		ret = append(ret, values[num%len(values)])
	}

	return
}
func (s *RandomService) RandomStrValues(values []string) (ret []string) {
	length := len(values)

	for i := 0; i < length; i++ {
		num := commonUtils.RandNum(length * 10000)
		ret = append(ret, values[num%len(values)])
	}

	return
}
