package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/easysoft/zendata/res"
	"github.com/easysoft/zendata/src/action"
	"github.com/easysoft/zendata/src/gen"
	"github.com/easysoft/zendata/src/model"
	serverConfig "github.com/easysoft/zendata/src/server/config"
	serverService "github.com/easysoft/zendata/src/server/service"
	serverUtils "github.com/easysoft/zendata/src/server/utils"
	serverConst "github.com/easysoft/zendata/src/server/utils/const"
	"github.com/easysoft/zendata/src/service"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	configUtils "github.com/easysoft/zendata/src/utils/config"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/facebookgo/inject"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	configs     []string
	defaultFile string
	configFile  string
	//count       int
	fields string

	root   string
	input  string
	decode bool

	listData bool
	listRes  bool
	view     string
	md5      string

	example bool
	help    bool
	set     bool

	flagSet *flag.FlagSet
)

func main() {
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		cleanup()
		os.Exit(0)
	}()

	flagSet = flag.NewFlagSet("zd", flag.ContinueOnError)

	flagSet.StringVar(&defaultFile, "d", "", "")
	flagSet.StringVar(&defaultFile, "default", "", "")

	flagSet.StringVar(&configFile, "c", "", "")
	flagSet.StringVar(&configFile, "Config", "", "")

	flagSet.StringVar(&input, "i", "", "")
	flagSet.StringVar(&input, "input", "", "")

	flagSet.IntVar(&vari.Total, "n", -1, "")
	flagSet.IntVar(&vari.Total, "lines", -1, "")

	flagSet.StringVar(&fields, "F", "", "")
	flagSet.StringVar(&fields, "field", "", "")

	flagSet.StringVar(&vari.Out, "o", "", "")
	flagSet.StringVar(&vari.Out, "output", "", "")

	flagSet.BoolVar(&listData, "l", false, "")
	flagSet.BoolVar(&listData, "list", false, "")
	flagSet.BoolVar(&listRes, "L", false, "")

	flagSet.StringVar(&view, "v", "", "")
	flagSet.StringVar(&view, "view", "", "")

	flagSet.StringVar(&md5, "md5", "", "")

	flagSet.BoolVar(&vari.Human, "H", false, "")
	flagSet.BoolVar(&vari.Human, "human", false, "")

	flagSet.BoolVar(&decode, "D", false, "")
	flagSet.BoolVar(&decode, "decode", false, "")

	flagSet.StringVar(&vari.Ip, "b", "", "")
	flagSet.StringVar(&vari.Ip, "bind", "", "")
	flagSet.IntVar(&vari.Port, "p", 0, "")
	flagSet.IntVar(&vari.Port, "port", 0, "")
	flagSet.StringVar(&root, "R", "", "")
	flagSet.StringVar(&root, "root", "", "")

	flagSet.BoolVar(&vari.Trim, "T", false, "")
	flagSet.BoolVar(&vari.Trim, "trim", false, "")

	flagSet.BoolVar(&vari.Recursive, "r", false, "")
	flagSet.BoolVar(&vari.Recursive, "recursive", false, "")

	flagSet.BoolVar(&example, "e", false, "")
	flagSet.BoolVar(&example, "example", false, "")

	flagSet.BoolVar(&help, "h", false, "")
	flagSet.BoolVar(&help, "help", false, "")

	flagSet.BoolVar(&set, "S", false, "")
	flagSet.BoolVar(&set, "set", false, "")

	flagSet.StringVar(&vari.Table, "t", "", "")
	flagSet.StringVar(&vari.Table, "table", "", "")
	flagSet.StringVar(&vari.Server, "s", "mysql", "")
	flagSet.StringVar(&vari.Server, "server", "mysql", "")
	flagSet.StringVar(&vari.DBDsn, "dns", "", "")
	flagSet.BoolVar(&vari.DBClear, "clear", false, "")

	flagSet.StringVar(&vari.ProtoCls, "cls", "", "")

	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-help")
	}

	files, count := fileUtils.GetFilesFromParams(os.Args[1:])
	flagSet.Parse(os.Args[1+count:])
	if count == 0 {
		files = []string{defaultFile, configFile}
	}
	if vari.Port != 0 {
		vari.RunMode = constant.RunModeServer
	}

	configUtils.InitConfig(root)
	vari.DB, _ = configUtils.InitDB()
	defer vari.DB.Close()

	switch os.Args[1] {
	default:
		flagSet.SetOutput(ioutil.Discard)
		if err := flagSet.Parse(os.Args[1:]); err == nil {
			if example {
				logUtils.PrintExample()
				return
			} else if help {
				logUtils.PrintUsage()
				return
			} else if set {
				service.Set()
				return
			} else if listData {
				service.ListData()
				return
			} else if listRes {
				service.ListRes()
				return
			} else if view != "" {
				service.View(view)
				return
			} else if md5 != "" {
				service.AddMd5(md5)
				return
			} else if decode {
				gen.Decode(files, fields, input)
				return
			}

			if vari.Ip != "" || vari.Port != 0 {
				vari.RunMode = constant.RunModeServer
			} else if input != "" {
				vari.RunMode = constant.RunModeParse
			}

			toGen(files)
		} else {
			logUtils.PrintUsage()
		}
	}
}

func toGen(files []string) {
	tmStart := time.Now()
	if vari.Verbose {
		logUtils.PrintTo(fmt.Sprintf("Start at %s.", tmStart.Format("2006-01-02 15:04:05")))
	}

	if vari.RunMode == constant.RunModeParse {
		ext := filepath.Ext(input)
		if ext == ".sql" {
			action.ParseSql(input, vari.Out)
		} else if ext == ".txt" {
			action.ParseArticle(input, vari.Out)
		}

	} else if vari.RunMode == constant.RunModeServer {
		vari.AgentLogDir = vari.ZdPath + serverConst.AgentLogDir + constant.PthSep
		err := fileUtils.MkDirIfNeeded(vari.AgentLogDir)
		if err != nil {
			logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("perm_deny", vari.AgentLogDir), color.FgRed)
			os.Exit(1)
		}

		startServer() // will init its own db

	} else if vari.RunMode == constant.RunModeServerRequest {
		//  use the files from post data
		files := []string{defaultFile, configFile}

		vari.Format = constant.FormatJson
		action.Generate(files, fields, vari.Format, vari.Table)

	} else if vari.RunMode == constant.RunModeGen {
		if vari.Human {
			vari.WithHead = true
		}

		if vari.Out != "" {
			fileUtils.MkDirIfNeeded(filepath.Dir(vari.Out))
			fileUtils.RemoveExist(vari.Out)

			ext := strings.ToLower(filepath.Ext(vari.Out))
			if len(ext) > 1 {
				ext = strings.TrimLeft(ext, ".")
			}
			if stringUtils.InArray(ext, constant.Formats) {
				vari.Format = ext
			}

			if vari.Format == constant.FormatExcel {
				logUtils.FilePath = vari.Out
			} else {
				logUtils.FileWriter, _ = os.OpenFile(vari.Out, os.O_RDWR|os.O_CREATE, 0777)
				defer logUtils.FileWriter.Close()
			}
		}
		if vari.DBDsn != "" {
			vari.Format = constant.FormatSql
		}

		if vari.Format == constant.FormatSql && vari.Table == "" {
			logUtils.PrintErrMsg(i118Utils.I118Prt.Sprintf("miss_table_name"))
			return
		}

		action.Generate(files, fields, vari.Format, vari.Table)
	}

	tmEnd := time.Now()
	if vari.Verbose {
		logUtils.PrintTo(fmt.Sprintf("End at %s.", tmEnd.Format("2006-01-02 15:04:05")))

		dur := tmEnd.Unix() - tmStart.Unix()
		logUtils.PrintTo(fmt.Sprintf("Duriation %d sec.", dur))
	}
}

func startServer() {
	if vari.Ip == "" {
		vari.Ip = commonUtils.GetIp()
	}
	if vari.Port == 0 {
		vari.Port = constant.DefaultPort
	}

	port := strconv.Itoa(vari.Port)
	logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("start_server",
		vari.Ip, port, vari.Ip, port, vari.Ip, port), color.FgCyan)

	// start admin server
	err := InitServer()
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("start_server_fail", port), color.FgRed)
	}
}

func DataHandler(writer http.ResponseWriter, req *http.Request) {
	logUtils.HttpWriter = writer

	defaultFile, configFile, fields, vari.Total,
		vari.Format, vari.Trim, vari.Table, decode, input, vari.Out = serverUtils.ParseGenParams(req)

	if decode {
		files := []string{defaultFile, configFile}
		gen.Decode(files, fields, input)
	} else if defaultFile != "" || configFile != "" {
		vari.RunMode = constant.RunModeServerRequest
		logUtils.PrintToWithoutNewLine(i118Utils.I118Prt.Sprintf("server_request", req.Method, req.URL))

		toGen(nil)
	}
}

// Server for admin server
type Server struct {
	Config *serverConfig.Config `inject:""`

	DefService     *serverService.DefService     `inject:""`
	FieldService   *serverService.FieldService   `inject:""`
	PreviewService *serverService.PreviewService `inject:""`
	SectionService *serverService.SectionService `inject:""`
	ReferService   *serverService.ReferService   `inject:""`
	ResService     *serverService.ResService     `inject:""`
	SyncService    *serverService.SyncService    `inject:""`

	RangesService    *serverService.RangesService    `inject:""`
	InstancesService *serverService.InstancesService `inject:""`
	TextService      *serverService.TextService      `inject:""`
	ExcelService     *serverService.ExcelService     `inject:""`
	ConfigService    *serverService.ConfigService    `inject:""`
}

func InitServer() (err error) {
	var g inject.Graph

	config := serverConfig.NewConfig()
	gormDb, err := serverConfig.NewGormDB(config)
	defer gormDb.Close()

	server := NewServer()

	// inject objects
	if err := g.Provide(
		&inject.Object{Value: config},
		&inject.Object{Value: gormDb},
		&inject.Object{Value: server},
	); err != nil {
		logUtils.PrintErrMsg(fmt.Sprintf("provide usecase objects to the Graph: %v", err))
	}
	err = g.Populate()
	if err != nil {
		logUtils.PrintErrMsg(fmt.Sprintf("populate the incomplete Objects: %v", err))
	}

	server.Run()

	return
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Config.ServerPort),
		Handler: s.Handler(),
	}

	httpServer.ListenAndServe()
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer( // client static
		&assetfs.AssetFS{Asset: res.Asset, AssetDir: res.AssetDir, AssetInfo: res.AssetInfo, Prefix: "ui/dist"}))
	mux.HandleFunc("/admin", s.admin)    // data admin
	mux.HandleFunc("/data", DataHandler) // data gen

	return mux
}

func (s *Server) admin(writer http.ResponseWriter, req *http.Request) {
	serverUtils.SetupCORS(&writer, req)

	bytes, err := ioutil.ReadAll(req.Body)
	if len(bytes) == 0 {
		return
	}

	reqData := model.ReqData{}
	err = serverUtils.ParserJsonReq(bytes, &reqData)
	if err != nil {
		serverUtils.OutputErr(err, writer)
		return
	}

	ret := model.ResData{Code: 1, Msg: "success"}
	switch reqData.Action {
	// common
	case "getWorkDir":
		ret.WorkDir = vari.ZdPath

	// def
	case "syncData":
		s.SyncService.SyncData(reqData.Mode)
	case "listDef":
		ret.Data, ret.Total = s.DefService.List(reqData.Keywords, reqData.Page)
	case "getDef":
		ret.Data, ret.Res = s.DefService.Get(reqData.Id)
		ret.WorkDir = vari.ZdPath
	case "saveDef":
		def := serverUtils.ConvertDef(reqData.Data)
		s.DefService.Save(&def)
		ret.Data = def
	case "removeDef":
		err = s.DefService.Remove(reqData.Id)

	// field
	case "getDefFieldTree":
		ret.Data, err = s.FieldService.GetTree(uint(reqData.Id))
	case "getDefField":
		ret.Data, err = s.FieldService.Get(reqData.Id)
	case "createDefField":
		var field *model.ZdField
		field, err = s.FieldService.Create(0, uint(reqData.Id), "新字段", reqData.Mode)

		ret.Data, err = s.FieldService.GetTree(field.DefID)
		ret.Model = field
	case "saveDefField":
		field := serverUtils.ConvertField(reqData.Data)
		err = s.FieldService.Save(&field)
	case "removeDefField":
		var defId int
		defId, err = s.FieldService.Remove(reqData.Id)
		ret.Data, err = s.FieldService.GetTree(uint(defId))
	case "moveDefField":
		var defId uint
		defId, ret.Model, err = s.FieldService.Move(uint(reqData.Src), uint(reqData.Dist), reqData.Mode)
		ret.Data, err = s.FieldService.GetTree(defId)

	// preview
	case "previewDefData":
		ret.Data = s.PreviewService.PreviewDefData(uint(reqData.Id))
	case "previewFieldData":
		ret.Data = s.PreviewService.PreviewFieldData(uint(reqData.Id), reqData.Mode)

	// field or instances section
	case "listSection":
		ret.Data, err = s.SectionService.List(uint(reqData.Id), reqData.Mode)

	case "createSection":
		paramMap := serverUtils.ConvertParams(reqData.Data)
		ownerType, _ := paramMap["ownerType"]
		ownerId, _ := strconv.Atoi(paramMap["ownerId"])
		sectionId, _ := strconv.Atoi(paramMap["sectionId"])

		err = s.SectionService.Create(uint(ownerId), uint(sectionId), ownerType)
		ret.Data, err = s.SectionService.List(uint(ownerId), ownerType)
	case "updateSection":
		section := serverUtils.ConvertSection(reqData.Data)
		err = s.SectionService.Update(&section)

		ret.Data, err = s.SectionService.List(section.OwnerID, reqData.Mode)
	case "removeSection":
		var fieldId uint
		fieldId, err = s.SectionService.Remove(reqData.Id, reqData.Mode)
		ret.Data, err = s.SectionService.List(fieldId, reqData.Mode)

	// field or instances refer, be create when init its owner
	case "getRefer":
		var refer model.ZdRefer
		refer, err = s.ReferService.Get(uint(reqData.Id), reqData.Mode)
		ret.Data = refer
	case "updateRefer":
		refer := serverUtils.ConvertRefer(reqData.Data)
		err = s.ReferService.Update(&refer)
	case "listReferFileForSelection":
		ret.Data = s.ResService.ListReferFileForSelection(reqData.Mode)
	case "listReferSheetForSelection":
		ret.Data = s.ResService.ListReferSheetForSelection(reqData.Mode)

	case "listReferExcelColForSelection":
		ret.Data = s.ResService.ListReferExcelColForSelection(reqData.Mode)
	case "listReferResFieldForSelection":
		ret.Data = s.ResService.ListReferResFieldForSelection(reqData.Id, reqData.Mode)

	// resource
	case "listRanges":
		ret.Data, ret.Total = s.RangesService.List(reqData.Keywords, reqData.Page)
	case "getRanges":
		ret.Data, ret.Res = s.RangesService.Get(reqData.Id)
		ret.WorkDir = vari.ZdPath
	case "saveRanges":
		ranges := serverUtils.ConvertRanges(reqData.Data)
		ret.Data = s.RangesService.Save(&ranges)
	case "removeRanges":
		err = s.RangesService.Remove(reqData.Id)

	case "getResRangesItemTree":
		ret.Data = s.RangesService.GetItemTree(reqData.Id)
	case "getResRangesItem":
		ret.Data = s.RangesService.GetItem(reqData.Id)
	case "createResRangesItem":
		var rangesItem *model.ZdRangesItem
		rangesId := reqData.DomainId
		rangesItem, err = s.RangesService.CreateItem(rangesId, reqData.Id, reqData.Mode)

		ret.Data = s.RangesService.GetItemTree(rangesId)
		ret.Model = rangesItem
	case "saveRangesItem":
		rangesItem := serverUtils.ConvertRangesItem(reqData.Data)
		ret.Data = s.RangesService.SaveItem(&rangesItem)
	case "removeResRangesItem":
		err = s.RangesService.RemoveItem(reqData.Id, reqData.DomainId)
		ret.Data = s.RangesService.GetItemTree(reqData.DomainId)

	case "listInstances":
		ret.Data, ret.Total = s.InstancesService.List(reqData.Keywords, reqData.Page)
	case "getInstances":
		ret.Data, ret.Res = s.InstancesService.Get(reqData.Id)
		ret.WorkDir = vari.ZdPath
	case "saveInstances":
		ranges := serverUtils.ConvertInstances(reqData.Data)
		ret.Data = s.InstancesService.Save(&ranges)
	case "removeInstances":
		err = s.InstancesService.Remove(reqData.Id)
	case "getResInstancesItemTree":
		ret.Data = s.InstancesService.GetItemTree(uint(reqData.Id))
	case "getResInstancesItem":
		ret.Data = s.InstancesService.GetItem(reqData.Id)
	case "createResInstancesItem":
		var item *model.ZdInstancesItem
		ownerId := reqData.DomainId
		item, err = s.InstancesService.CreateItem(ownerId, reqData.Id, reqData.Mode)

		ret.Data = s.InstancesService.GetItemTree(uint(ownerId))
		ret.Model = item
	case "saveInstancesItem":
		rangesItem := serverUtils.ConvertInstancesItem(reqData.Data)
		ret.Data = s.InstancesService.SaveItem(&rangesItem)
	case "removeResInstancesItem":
		err = s.InstancesService.RemoveItem(reqData.Id)
		ret.Data = s.InstancesService.GetItemTree(uint(reqData.DomainId))

	case "listExcel":
		ret.Data, ret.Total = s.ExcelService.List(reqData.Keywords, reqData.Page)
	case "getExcel":
		ret.Data, ret.Res = s.ExcelService.Get(reqData.Id)
		ret.WorkDir = vari.ZdPath
	case "saveExcel":
		ranges := serverUtils.ConvertExcel(reqData.Data)
		ret.Data = s.ExcelService.Save(&ranges)
	case "removeExcel":
		err = s.ExcelService.Remove(reqData.Id)

	case "listText":
		ret.Data, ret.Total = s.TextService.List(reqData.Keywords, reqData.Page)
	case "getText":
		ret.Data, ret.Res = s.TextService.Get(reqData.Id)
		ret.WorkDir = vari.ZdPath
	case "saveText":
		ranges := serverUtils.ConvertText(reqData.Data)
		ret.Data = s.TextService.Save(&ranges)
	case "removeText":
		err = s.TextService.Remove(reqData.Id)

	case "listConfig":
		ret.Data, ret.Total = s.ConfigService.List(reqData.Keywords, reqData.Page)
	case "getConfig":
		ret.Data, ret.Res = s.ConfigService.Get(reqData.Id)
		ret.WorkDir = vari.ZdPath
	case "saveConfig":
		ranges := serverUtils.ConvertConfig(reqData.Data)
		ret.Data = s.ConfigService.Save(&ranges)
	case "removeConfig":
		err = s.ConfigService.Remove(reqData.Id)

	case "getResConfigItemTree":
		ret.Data = s.ConfigService.GConfigItemTree(reqData.Id)

	default:
		ret.Code = 0
		ret.Msg = "api not found"
	}
	if err != nil {
		ret.Code = 0
		ret.Msg = "api error: " + err.Error()
	}

	bytes, _ = json.Marshal(ret)
	io.WriteString(writer, string(bytes))
}

func NewServer() *Server {
	return &Server{}
}

func init() {
	cleanup()
}

func cleanup() {
	color.Unset()
}
