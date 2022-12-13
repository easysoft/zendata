package main

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zendata/cmd/test/comm"
	"github.com/easysoft/zendata/cmd/test/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"strings"
)

func main() {
	tableName := "food"
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/歌手歌名歌词 %d.json"

	tableName = model.Place{}.TableName()
	db := comm.GetDB()
	err := db.Exec(fmt.Sprintf(comm.TruncateTable, tableName)).Error
	if err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		content := fileUtils.ReadFileBuf(fmt.Sprintf(filePath, i+1))

		data := make([]model.SongData, 0)
		json.Unmarshal(content, &data)

		for _, song := range data {
			po := model.Song{
				Name:   song.Name,
				Singer: song.Singer,
				Lyric:  strings.Join(song.Lyric, "\n"),
			}

			if po.Name != "" {
				db.Save(&po)
			}
		}
	}

}
