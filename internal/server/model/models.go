package model

var (
	CommonPrefix = "zd_"
	Models       = []interface{}{
		&ZdDef{},
		&ZdField{},
		&ZdSection{},
		&ZdRefer{},
		&ZdConfig{},
		&ZdRanges{},
		&ZdRangesItem{},
		&ZdInstances{},
		&ZdInstancesItem{},
		&ZdText{},
		&ZdExcel{},
	}
)
