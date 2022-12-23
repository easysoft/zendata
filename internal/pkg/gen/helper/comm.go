package genHelper

import (
	"github.com/easysoft/zendata/pkg/utils/vari"
)

func GetRandFieldSection(pth string) (key int) {
	max := 0

	for k, v := range vari.GlobalVars.RandFieldSectionShortKeysToPathMap {
		if pth == v {
			key = k
			return
		}

		if k > max {
			max = k
		}
	}

	if key == 0 {
		key = max + 1
	}

	return
}
