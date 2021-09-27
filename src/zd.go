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
	serverRepo "github.com/easysoft/zendata/src/server/repo"
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
	flagSet.StringVar(&configFile, "config", "", "")

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
				gen.Decode(defaultFile, configFile, fields, input)
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
		gen.Decode(defaultFile, configFile, fields, input)
	} else if defaultFile != "" || configFile != "" {
		vari.RunMode = constant.RunModeServerRequest
		logUtils.PrintToWithoutNewLine(i118Utils.I118Prt.Sprintf("server_request", req.Method, req.URL))

		toGen(nil)
	}
}

// for admin server
type Server struct {
	config *serverConfig.Config

	defService     *serverService.DefService
	fieldService   *serverService.FieldService
	previewService *serverService.PreviewService
	sectionService *serverService.SectionService
	referService   *serverService.ReferService
	resService     *serverService.ResService
	syncService    *serverService.SyncService

	rangesService    *serverService.RangesService
	instancesService *serverService.InstancesService
	textService      *serverService.TextService
	excelService     *serverService.ExcelService
	configService    *serverService.ConfigService
}

func InitServer() (err error) {
	config := serverConfig.NewConfig()
	gormDb, err := serverConfig.NewGormDB(config)
	defer gormDb.Close()

	defRepo := serverRepo.NewDefRepo(gormDb)
	fieldRepo := serverRepo.NewFieldRepo(gormDb)
	sectionRepo := serverRepo.NewSectionRepo(gormDb)
	referRepo := serverRepo.NewReferRepo(gormDb)
	rangesRepo := serverRepo.NewRangesRepo(gormDb)
	instancesRepo := serverRepo.NewInstancesRepo(gormDb)
	textRepo := serverRepo.NewTextRepo(gormDb)
	excelRepo := serverRepo.NewExcelRepo(gormDb)
	configRepo := serverRepo.NewConfigRepo(gormDb)

	defService := serverService.NewDefService(defRepo, fieldRepo, sectionRepo, referRepo)
	fieldService := serverService.NewFieldService(defRepo, fieldRepo, referRepo, defService)

	referService := serverService.NewReferService(fieldRepo, referRepo, defService)
	rangesService := serverService.NewRangesService(rangesRepo, sectionRepo)
	instancesService := serverService.NewInstancesService(instancesRepo, referRepo, sectionRepo)
	textService := serverService.NewTextService(textRepo)
	excelService := serverService.NewExcelService(excelRepo)
	configService := serverService.NewConfigService(configRepo, sectionRepo)
	sectionService := serverService.NewSectionService(
		fieldRepo, configRepo, rangesRepo, instancesRepo, sectionRepo, defService,
		instancesService, rangesService, configService)
	resService := serverService.NewResService(rangesRepo, instancesRepo,
		configRepo, excelRepo, textRepo, defRepo)

	syncService := serverService.NewSyncService(defService,
		fieldService, rangesService, instancesService, configService, excelService, textService,
		referService, resService)
	previewService := serverService.NewPreviewService(defRepo, fieldRepo, referRepo, instancesRepo)

	server := NewServer(config, defService, fieldService, sectionService, referService,
		rangesService, instancesService, textService, excelService, configService, resService,
		syncService, previewService)
	server.Run()

	return
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.ServerPort),
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
		s.syncService.SyncData(reqData.Mode)
	case "listDef":
		ret.Data, ret.Total = s.defService.List(reqData.Keywords, reqData.Page)
	case "getDef":
		ret.Data, ret.Res = s.defService.Get(reqData.Id)
		ret.WorkDir = vari.ZdPath
	case "saveDef":
		def := serverUtils.ConvertDef(reqData.Data)
		s.defService.Save(&def)
		ret.Data = def
	case "removeDef":
		err = s.defService.Remove(reqData.Id)

	// field
	case "getDefFieldTree":
		ret.Data, err = s.fieldService.GetTree(uint(reqData.Id))
	case "getDefField":
		ret.Data, err = s.fieldService.Get(reqData.Id)
	case "createDefField":
		var field *model.ZdField
		field, err = s.fieldService.Create(0, uint(reqData.Id), "新字段", reqData.Mode)

		ret.Data, err = s.fieldService.GetTree(field.DefID)
		ret.Model = field
	case "saveDefField":
		field := serverUtils.ConvertField(reqData.Data)
		err = s.fieldService.Save(&field)
	case "removeDefField":
		var defId int
		defId, err = s.fieldService.Remove(reqData.Id)
		ret.Data, err = s.fieldService.GetTree(uint(defId))
	case "moveDefField":
		var defId uint
		defId, ret.Model, err = s.fieldService.Move(uint(reqData.Src), uint(reqData.Dist), reqData.Mode)
		ret.Data, err = s.fieldService.GetTree(defId)

	// preview
	case "previewDefData":
		ret.Data = s.previewService.PreviewDefData(uint(reqData.Id))
	case "previewFieldData":
		ret.Data = s.previewService.PreviewFieldData(uint(reqData.Id), reqData.Mode)

	// field or instances section
	case "listSection":
		ret.Data, err = s.sectionService.List(uint(reqData.Id), reqData.Mode)

	case "createSection":
		paramMap := serverUtils.ConvertParams(reqData.Data)
		ownerType, _ := paramMap["ownerType"]
		ownerId, _ := strconv.Atoi(paramMap["ownerId"])
		sectionId, _ := strconv.Atoi(paramMap["sectionId"])

		err = s.sectionService.Create(uint(ownerId), uint(sectionId), ownerType)
		ret.Data, err = s.sectionService.List(uint(ownerId), ownerType)
	case "updateSection":
		section := serverUtils.ConvertSection(reqData.Data)
		err = s.sectionService.Update(&section)

		ret.Data, err = s.sectionService.List(section.OwnerID, reqData.Mode)
	case "removeSection":
		var fieldId uint
		fieldId, err = s.sectionService.Remove(reqData.Id, reqData.Mode)
		ret.Data, err = s.sectionService.List(fieldId, reqData.Mode)

	// field or instances refer, be create when init its owner
	case "getRefer":
		var refer model.ZdRefer
		refer, err = s.referService.Get(uint(reqData.Id), reqData.Mode)
		ret.Data = refer
	case "updateRefer":
		refer := serverUtils.ConvertRefer(reqData.Data)
		err = s.referService.Update(&refer)
	case "listReferFileForSelection":
		ret.Data = s.resService.ListReferFileForSelection(reqData.Mode)
	case "listReferSheetForSelection":
		ret.Data = s.resService.ListReferSheetForSelection(reqData.Mode)

	case "listReferExcelColForSelection":
		ret.Data = s.resService.ListReferExcelColForSelection(reqData.Mode)
	case "listReferResFieldForSelection":
		ret.Data = s.resService.ListReferResFieldForSelection(reqData.Id, reqData.Mode)

	// resource
	case "listRanges":
		ret.Data, ret.Total = s.rangesService.List(reqData.Keywords, reqData.Page)
	case "getRanges":
		ret.Data, ret.Res = s.rangesService.Get(reqData.Id)
		ret.WorkDir = vari.ZdPath
	case "saveRanges":
		ranges := serverUtils.ConvertRanges(reqData.Data)
		ret.Data = s.rangesService.Save(&ranges)
	case "removeRanges":
		err = s.rangesService.Remove(reqData.Id)

	case "getResRangesItemTree":
		ret.Data = s.rangesService.GetItemTree(reqData.Id)
	case "getResRangesItem":
		ret.Data = s.rangesService.GetItem(reqData.Id)
	case "createResRangesItem":
		var rangesItem *model.ZdRangesItem
		rangesId := reqData.DomainId
		rangesItem, err = s.rangesService.CreateItem(rangesId, reqData.Id, reqData.Mode)

		ret.Data = s.rangesService.GetItemTree(rangesId)
		ret.Model = rangesItem
	case "saveRangesItem":
		rangesItem := serverUtils.ConvertRangesItem(reqData.Data)
		ret.Data = s.rangesService.SaveItem(&rangesItem)
	case "removeResRangesItem":
		err = s.rangesService.RemoveItem(reqData.Id, reqData.DomainId)
		ret.Data = s.rangesService.GetItemTree(reqData.DomainId)

	case "listInstances":
		ret.Data, ret.Total = s.instancesService.List(reqData.Keywords, reqData.Page)
	case "getInstances":
		ret.Data, ret.Res = s.instancesService.Get(reqData.Id)
		ret.WorkDir = vari.ZdPath
	case "saveInstances":
		ranges := serverUtils.ConvertInstances(reqData.Data)
		ret.Data = s.instancesService.Save(&ranges)
	case "removeInstances":
		err = s.instancesService.Remove(reqData.Id)
	case "getResInstancesItemTree":
		ret.Data = s.instancesService.GetItemTree(uint(reqData.Id))
	case "getResInstancesItem":
		ret.Data = s.instancesService.GetItem(reqData.Id)
	case "createResInstancesItem":
		var item *model.ZdInstancesItem
		ownerId := reqData.DomainId
		item, err = s.instancesService.CreateItem(ownerId, reqData.Id, reqData.Mode)

		ret.Data = s.instancesService.GetItemTree(uint(ownerId))
		ret.Model = item
	case "saveInstancesItem":
		rangesItem := serverUtils.ConvertInstancesItem(reqData.Data)
		ret.Data = s.instancesService.SaveItem(&rangesItem)
	case "removeResInstancesItem":
		err = s.instancesService.RemoveItem(reqData.Id)
		ret.Data = s.instancesService.GetItemTree(uint(reqData.DomainId))

	case "listExcel":
		ret.Data, ret.Total = s.excelService.List(reqData.Keywords, reqData.Page)
	case "getExcel":
		ret.Data, ret.Res = s.excelService.Get(reqData.Id)
		ret.WorkDir = vari.ZdPath
	case "saveExcel":
		ranges := serverUtils.ConvertExcel(reqData.Data)
		ret.Data = s.excelService.Save(&ranges)
	case "removeExcel":
		err = s.excelService.Remove(reqData.Id)

	case "listText":
		ret.Data, ret.Total = s.textService.List(reqData.Keywords, reqData.Page)
	case "getText":
		ret.Data, ret.Res = s.textService.Get(reqData.Id)
		ret.WorkDir = vari.ZdPath
	case "saveText":
		ranges := serverUtils.ConvertText(reqData.Data)
		ret.Data = s.textService.Save(&ranges)
	case "removeText":
		err = s.textService.Remove(reqData.Id)

	case "listConfig":
		ret.Data, ret.Total = s.configService.List(reqData.Keywords, reqData.Page)
	case "getConfig":
		ret.Data, ret.Res = s.configService.Get(reqData.Id)
		ret.WorkDir = vari.ZdPath
	case "saveConfig":
		ranges := serverUtils.ConvertConfig(reqData.Data)
		ret.Data = s.configService.Save(&ranges)
	case "removeConfig":
		err = s.configService.Remove(reqData.Id)

	case "getResConfigItemTree":
		ret.Data = s.configService.GConfigItemTree(reqData.Id)

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

func NewServer(config *serverConfig.Config, defService *serverService.DefService,
	fieldServer *serverService.FieldService, sectionService *serverService.SectionService,
	referService *serverService.ReferService, rangesService *serverService.RangesService,
	instancesService *serverService.InstancesService, textService *serverService.TextService,
	excelService *serverService.ExcelService, configService *serverService.ConfigService,
	resService *serverService.ResService, syncService *serverService.SyncService,
	previewService *serverService.PreviewService) *Server {
	return &Server{
		config:           config,
		defService:       defService,
		fieldService:     fieldServer,
		sectionService:   sectionService,
		referService:     referService,
		rangesService:    rangesService,
		instancesService: instancesService,
		textService:      textService,
		excelService:     excelService,
		configService:    configService,
		resService:       resService,
		syncService:      syncService,
		previewService:   previewService,
	}
}

func init() {
	cleanup()
}

func cleanup() {
	color.Unset()
}
