package gen

import (
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"strings"
	"sync"
	"time"
)

var (
	c = SafeCounter{v: make(map[string]int)}
)

func SyncGenCacheAndReturnFirstPart(fileContents [][]byte, fieldsToExport *[]string) (
	rows [][]string, colIsNumArr []bool, err error) {

	FixTotalNum()
	genResData(fieldsToExport)

	threadCount := getTheadCount(vari.GlobalVars.Total, consts.MaxNumbForAsync)

	var waitGroup sync.WaitGroup
	waitGroup.Add(threadCount)

	vari.GlobalVars.DefData = LoadDataContentDef(fileContents, fieldsToExport)

	if err = CheckParams(); err != nil {
		return
	}

	for i := 0; i < threadCount; i++ {
		index := i
		c.Inc("count")

		recordNumb := getTheadRecordNum(vari.GlobalVars.Total, consts.MaxNumbForAsync, threadCount, index)

		go func() {
			defer func() {
				c.Dec("count")
				waitGroup.Done()
			}()

			fmt.Printf("Thread %d to create %d records.\n", index, recordNumb)

			topLevelFieldNameToValuesMap := genFieldsData(fieldsToExport, &colIsNumArr, recordNumb)
			twoDimArr := genDataTwoDimArr(topLevelFieldNameToValuesMap, fieldsToExport, recordNumb)
			rowsPart := populateRowsFromTwoDimArr(twoDimArr, vari.GlobalVars.Recursive, true, recordNumb)

			fmt.Printf("Thread %d actual create %d records.\n", index, len(rowsPart))

			cacheKey, _, _, _, _ := ParseCache()
			cacheKey = getBatchKey(cacheKey, index)

			ClearCache(cacheKey)
			baseKey := removeBatchNumInKey(cacheKey)
			if index == 0 {
				CreateCacheIsNumTable(baseKey, *fieldsToExport)
				rows = rowsPart
			}
			CreateCacheDataTable(cacheKey, *fieldsToExport)

			CreateCacheData(cacheKey, *fieldsToExport, rowsPart)
			CreateCacheIsNum(baseKey, *fieldsToExport, colIsNumArr)
		}()
	}

	go func() {
		for {
			fmt.Printf("=thread count= %d\n", c.Value("count"))
			time.Sleep(time.Millisecond * 50)
		}
	}()

	waitGroup.Wait()

	return
}

func getTheadCount(total, size int) (threadCount int) {
	threadCount = total / size
	if total%size > 0 {
		threadCount++
	}

	return
}
func getTheadRecordNum(total, size, theadCount, index int) (num int) {
	if index == theadCount-1 {
		num = total - ((theadCount - 1) * size)
	} else {
		num = size
	}

	return
}

func getBatchKey(cacheKey string, index int) string {
	return fmt.Sprintf("%s_%d", cacheKey, index)
}

func removeBatchNumInKey(cacheKey string) string {
	return cacheKey[:strings.LastIndex(cacheKey, "_")]
}

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	c.v[key]++
	c.mu.Unlock()
}
func (c *SafeCounter) Dec(key string) {
	c.mu.Lock()
	c.v[key]--
	c.mu.Unlock()
}
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v[key]
}
