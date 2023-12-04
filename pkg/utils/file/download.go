package fileUtils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/fatih/color"
)

func Download(url string, dst string) (err error) {
	fmt.Printf("DownloadToFile From: %s to %s.\n", url, dst)

	MkDirIfNeeded(filepath.Dir(dst))

	var data []byte
	data, err = HTTPDownload(url)
	if err == nil {
		logUtils.PrintTo(i118Utils.Sprintf("file_downloaded", url))

		err = WriteDownloadFile(dst, data)
		if err == nil {
			logUtils.Info(i118Utils.Sprintf("file_download_saved", url, dst))
		}
	}

	return
}

func HTTPDownload(uri string) ([]byte, error) {
	res, err := http.Get(uri)
	if err != nil {
		logUtils.Infof(color.RedString("download file failed, error: %s.", err.Error()))
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logUtils.Infof(color.RedString("read downloaded file failed, error: %s.", err.Error()))
	}
	return d, err
}

func WriteDownloadFile(dst string, d []byte) error {
	err := ioutil.WriteFile(dst, d, 0444)
	if err != nil {
		logUtils.Infof(color.RedString("write download file failed, error: %s.", err.Error()))
	}
	return err
}
