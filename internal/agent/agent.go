package agent

import (
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/action"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen"
	serverUtils "github.com/easysoft/zendata/internal/server/utils"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"net/http"
	"time"
)

var (
	configs           []string
	defaultFile       string
	configFile        string
	defaultDefContent []byte
	configDefContent  []byte

	fields string
	input  string
	decode bool
)

func DataHandler(writer http.ResponseWriter, req *http.Request) {
	logUtils.HttpWriter = writer

	if req.Method == http.MethodGet {
		defaultFile, configFile, fields, vari.GenVars.Total,
			vari.Format, vari.Trim, vari.Table, decode, input, vari.Out = serverUtils.ParseGenParams(req)
	} else if req.Method == http.MethodPost {
		defaultDefContent, configDefContent, fields, vari.GenVars.Total,
			vari.Format, vari.Trim, vari.Table, decode, input, vari.Out = serverUtils.ParseGenParamsToByte(req)
	}

	if decode {
		files := []string{defaultFile, configFile}
		gen.Decode(files, fields, input)
		return
	}

	if defaultDefContent != nil || configDefContent != nil {
		vari.RunMode = constant.RunModeServerRequest
		logUtils.PrintToWithoutNewLine(i118Utils.I118Prt.Sprintf("server_request", req.Method, req.URL))

		genData()
		// Avoid variable affecting the results of request.
		defaultDefContent = nil
		configDefContent = nil

	} else if defaultFile != "" || configFile != "" {
		vari.RunMode = constant.RunModeServerRequest
		logUtils.PrintToWithoutNewLine(i118Utils.I118Prt.Sprintf("server_request", req.Method, req.URL))

		genData()

		// Avoid variable affecting the results of request.
		defaultFile = ""
		configFile = ""
	}

	//writer.WriteHeader(http.StatusAccepted)
}

func genData() {
	tmStart := time.Now()
	if vari.Verbose {
		logUtils.PrintTo(fmt.Sprintf("Start at %s.", tmStart.Format("2006-01-02 15:04:05")))
	}

	vari.Format = constant.FormatJson
	if defaultFile != "" || configFile != "" {
		files := []string{defaultFile, configFile}
		action.Generate(files, fields, vari.Format, vari.Table)
	} else {
		contents := [][]byte{defaultDefContent, configDefContent}
		action.GenerateByContent(contents, fields, vari.Format, vari.Table)
	}

	tmEnd := time.Now()
	if vari.Verbose {
		logUtils.PrintTo(fmt.Sprintf("End at %s.", tmEnd.Format("2006-01-02 15:04:05")))

		dur := tmEnd.Unix() - tmStart.Unix()
		logUtils.PrintTo(fmt.Sprintf("Duriation %d sec.", dur))
	}
}
