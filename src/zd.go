package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/easysoft/zendata/src/action"
	"github.com/easysoft/zendata/src/gen"
	"github.com/easysoft/zendata/src/model"
	serverConfig "github.com/easysoft/zendata/src/server/config"
	serverRepo "github.com/easysoft/zendata/src/server/repo"
	serverService "github.com/easysoft/zendata/src/server/service"
	serverUtils "github.com/easysoft/zendata/src/server/utils"
	"github.com/easysoft/zendata/src/service"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	configUtils "github.com/easysoft/zendata/src/utils/config"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

var (
	defaultFile string
	configFile  string
	//count       int
	fields      string

	root string
	input  string
	output string
	table  string
	format = constant.FormatText
	decode bool

	article string

	listRes bool
	viewRes string
	viewDetail string
	md5 string

	example bool
	help   bool
	set   bool

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

	flagSet.StringVar(&output, "o", "", "")
	flagSet.StringVar(&output, "output", "", "")

	flagSet.StringVar(&table, "t", "", "")
	flagSet.StringVar(&table, "table", "", "")

	flagSet.BoolVar(&listRes, "l", false, "")
	flagSet.BoolVar(&listRes, "list", false, "")

	flagSet.StringVar(&viewRes, "v", "", "")
	flagSet.StringVar(&viewRes, "view", "", "")

	flagSet.StringVar(&md5, "md5", "", "")

	flagSet.BoolVar(&vari.Human, "H", false, "")
	flagSet.BoolVar(&vari.Human, "human", false, "")

	flagSet.BoolVar(&decode, "D", false, "")
	flagSet.BoolVar(&decode, "decode", false, "")

	flagSet.StringVar(&article, "a", "", "")
	flagSet.StringVar(&article, "article", "", "")

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

	flagSet.BoolVar(&set, "s", false, "")
    flagSet.BoolVar(&set, "set", false, "")

	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-help")
	}

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
            } else if listRes {
				service.ListRes()
				return
			} else if viewRes != "" {
				service.ViewRes(viewRes)
				return
			} else if md5 != "" {
				service.AddMd5(md5)
				return
			} else if decode {
				gen.Decode(defaultFile, configFile, fields, input, output)
				return
			} else if article != "" {
				service.ConvertArticle(article, output)
				return
			}

			if vari.Ip != "" || vari.Port != 0 {
				vari.RunMode = constant.RunModeServer
			} else if input != "" {
				vari.RunMode = constant.RunModeParse
			}

			toGen()
		} else {
			logUtils.PrintUsage()
		}
	}
}

func toGen() {
	if vari.RunMode == constant.RunModeServer {
		if root != "" {
			if fileUtils.IsAbosutePath(root) {
				vari.WorkDir = root
			} else {
				vari.WorkDir = vari.WorkDir + root
			}
			vari.WorkDir = fileUtils.AddSepIfNeeded(vari.WorkDir)
		}
		constant.SqliteData = strings.Replace(constant.SqliteData, "file:", "file:" + vari.WorkDir, 1)

		StartServer()
	} else if vari.RunMode == constant.RunModeServerRequest {
		format = constant.FormatJson
		action.Generate(defaultFile, configFile, fields, format, table)
	} else if vari.RunMode == constant.RunModeParse {
		action.ParseSql(input, output)
	} else if vari.RunMode == constant.RunModeGen {
		if vari.Human {
			vari.WithHead = true
		}

		if output != "" {
			fileUtils.MkDirIfNeeded(filepath.Dir(output))
			fileUtils.RemoveExist(output)

			ext := strings.ToLower(path.Ext(output))
			if len(ext) > 1 {
				ext = strings.TrimLeft(ext,".")
			}
			if stringUtils.InArray(ext, constant.Formats) {
				format = ext
			}

			if format == constant.FormatExcel {
				logUtils.FilePath = output
			} else {
				logUtils.FileWriter, _ = os.OpenFile(output, os.O_RDWR | os.O_CREATE, 0777)
				defer logUtils.FileWriter.Close()
			}
		}

		if format == constant.FormatSql && table == "" {
			logUtils.PrintErrMsg(i118Utils.I118Prt.Sprintf("miss_table_name"))
			return
		}

		action.Generate(defaultFile, configFile, fields, format, table)
	}
}

func StartServer() {
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
	err := Init()
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("start_server_fail", port), color.FgRed)
	}
}

func DataHandler(writer http.ResponseWriter, req *http.Request) {
	logUtils.HttpWriter = writer

	defaultFile, configFile, fields, vari.Total,
		format, table, decode, input, output = serverUtils.ParseGenParams(req)

	if decode {
		gen.Decode(defaultFile, configFile, fields, input, output)
	} else if defaultFile != "" || configFile != "" {
		vari.RunMode = constant.RunModeServerRequest
		logUtils.PrintToWithoutNewLine(i118Utils.I118Prt.Sprintf("server_request", req.Method, req.URL))

		toGen()
	}
}

// for admin server
type Server struct {
	config        *serverConfig.Config

	defService *serverService.DefService
	fieldService *serverService.FieldService
	sectionService *serverService.SectionService
	referService *serverService.ReferService
	resService *serverService.ResService
	syncService *serverService.SyncService

	rangesService *serverService.RangesService
	instancesService *serverService.InstancesService
	textService *serverService.TextService
	excelService *serverService.ExcelService
	configService *serverService.ConfigService
}

func Init() (err error) {
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

	defService := serverService.NewDefService(defRepo, fieldRepo, referRepo)
	fieldService := serverService.NewFieldService(defRepo, fieldRepo, referRepo, defService)
	sectionService := serverService.NewSectionService(fieldRepo, sectionRepo)
	rangesService := serverService.NewRangesService(rangesRepo)
	instancesService := serverService.NewInstancesService(instancesRepo, referRepo)
	textService := serverService.NewTextService(textRepo)
	excelService := serverService.NewExcelService(excelRepo)
	configService := serverService.NewConfigService(configRepo)
	resService := serverService.NewResService(rangesRepo, instancesRepo,
		configRepo, excelRepo, textRepo, defRepo)

	referService := serverService.NewReferService(fieldRepo, referRepo)
	syncService := serverService.NewSyncService(defService,
		fieldService, rangesService, instancesService, configService, excelService, textService,
		referService, resService)

	server := NewServer(config, defService, fieldService, sectionService, referService,
		rangesService, instancesService, textService, excelService, configService, resService, syncService)
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

	mux.HandleFunc("/", DataHandler)
	mux.HandleFunc("/admin", s.admin)

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

	ret := model.ResData{ Code: 1, Msg: "success"}
	switch reqData.Action {
	// def
	case "syncData":
		s.syncService.SyncData(reqData.Mode)
	case "listDef":
		ret.Data, ret.Total = s.defService.List(reqData.Keywords, reqData.Page)
	case "getDef":
		var def model.ZdDef
		def, ret.Res = s.defService.Get(reqData.Id)

		ret.Data = def
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
	case "listReferForSelection":
		ret.Data = s.resService.ListReferForSelection(reqData.Mode)
	case "listReferFieldForSelection":
		ret.Data = s.resService.ListReferFieldForSelection(reqData.Id, reqData.Mode)

	// resource
	case "listRanges":
		ret.Data, ret.Total = s.rangesService.List(reqData.Keywords, reqData.Page)
	case "getRanges":
		ret.Data, ret.Res = s.rangesService.Get(reqData.Id)
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
	case "saveExcel":
		ranges := serverUtils.ConvertExcel(reqData.Data)
		ret.Data = s.excelService.Save(&ranges)
	case "removeExcel":
		err = s.excelService.Remove(reqData.Id)

	case "listText":
		ret.Data, ret.Total = s.textService.List(reqData.Keywords, reqData.Page)
	case "getText":
		ret.Data, ret.Res = s.textService.Get(reqData.Id)
	case "saveText":
		ranges := serverUtils.ConvertText(reqData.Data)
		ret.Data = s.textService.Save(&ranges)
	case "removeText":
		err = s.textService.Remove(reqData.Id)

	case "listConfig":
		ret.Data, ret.Total = s.configService.List(reqData.Keywords, reqData.Page)
	case "getConfig":
		ret.Data, ret.Res = s.configService.Get(reqData.Id)
	case "saveConfig":
		ranges := serverUtils.ConvertConfig(reqData.Data)
		ret.Data = s.configService.Save(&ranges)
	case "removeConfig":
		err = s.configService.Remove(reqData.Id)
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
	resService *serverService.ResService, syncService *serverService.SyncService) *Server {
	return &Server{
		config:        config,
		defService: defService,
		fieldService: fieldServer,
		sectionService: sectionService,
		referService: referService,
		rangesService: rangesService,
		instancesService: instancesService,
		textService: textService,
		excelService: excelService,
		configService: configService,
		resService: resService,
		syncService: syncService,
	}
}

func init() {
	cleanup()
	configUtils.InitConfig()
}

func cleanup() {
	color.Unset()
}
