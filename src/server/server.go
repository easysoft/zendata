package server

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zendata/src/model"
	serverConfig "github.com/easysoft/zendata/src/server/config"
	serverService "github.com/easysoft/zendata/src/server/service"
	serverUtils "github.com/easysoft/zendata/src/server/utils"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/facebookgo/inject"
	"github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

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

func InitServer(config *serverConfig.Config, gormDb *gorm.DB) (server *Server, err error) {
	var g inject.Graph

	server = &Server{}

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

	return
}

func (s *Server) Admin(writer http.ResponseWriter, req *http.Request) {
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
