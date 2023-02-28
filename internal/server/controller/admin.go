package controller

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	serverConfig "github.com/easysoft/zendata/internal/server/config"
	serverService "github.com/easysoft/zendata/internal/server/service"
	serverUtils "github.com/easysoft/zendata/internal/server/utils"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/kataras/iris/v12"
	"strconv"
)

type AdminCtrl struct {
	BaseCtrl

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

func (c *AdminCtrl) Handle(ctx iris.Context) {
	reqData := model.ReqData{}

	err := ctx.ReadJSON(&reqData)
	if err != nil {
		ctx.JSON(c.ErrResp(consts.CommErr, err.Error()))
		return
	}

	ret := model.RespData{Code: 1, Msg: "success"}

	switch reqData.Action {
	// common
	//case "getWorkDir":
	//	ret.WorkDir = vari.WorkDir

	// def
	//case "syncData":
	//	s.SyncService.SyncData()
	//case "listDef":
	//	ret.Data, ret.Total = s.DefService.List(reqData.Keywords, reqData.Page)

	case "getDef":
		ret.Data, ret.Res = c.DefService.Get(reqData.Id)
		ret.WorkDir = vari.WorkDir
	case "saveDef":
		def := serverUtils.ConvertDef(reqData.Data)
		c.DefService.Save(&def)
		ret.Data = def
	case "removeDef":
		err = c.DefService.Remove(reqData.Id)

	// field
	case "getDefFieldTree":
		ret.Data, err = c.FieldService.GetTree(uint(reqData.Id))
	case "getDefField":
		ret.Data, err = c.FieldService.Get(reqData.Id)
	case "createDefField":
		var field *model.ZdField
		field, err = c.FieldService.Create(0, uint(reqData.Id), "新字段", reqData.Mode)

		ret.Data, err = c.FieldService.GetTree(field.DefID)
		ret.Model = field
	case "saveDefField":
		field := serverUtils.ConvertField(reqData.Data)
		err = c.FieldService.Save(&field)
	case "removeDefField":
		var defId int
		defId, err = c.FieldService.Remove(reqData.Id)
		ret.Data, err = c.FieldService.GetTree(uint(defId))
	case "moveDefField":
		var defId uint
		defId, ret.Model, err = c.FieldService.Move(uint(reqData.Src), uint(reqData.Dist), reqData.Mode)
		ret.Data, err = c.FieldService.GetTree(defId)

	// preview
	case "previewDefData":
		ret.Data = c.PreviewService.PreviewDefData(uint(reqData.Id))
	case "previewFieldData":
		ret.Data = c.PreviewService.PreviewFieldData(uint(reqData.Id), reqData.Mode)

	// field or instances section
	case "listSection":
		ret.Data, err = c.SectionService.List(uint(reqData.Id), reqData.Mode)

	case "createSection":
		paramMap := serverUtils.ConvertParams(reqData.Data)
		ownerType, _ := paramMap["ownerType"]
		ownerId, _ := strconv.Atoi(paramMap["ownerId"])
		sectionId, _ := strconv.Atoi(paramMap["sectionId"])

		err = c.SectionService.Create(uint(ownerId), uint(sectionId), ownerType)
		ret.Data, err = c.SectionService.List(uint(ownerId), ownerType)
	case "updateSection":
		section := serverUtils.ConvertSection(reqData.Data)
		err = c.SectionService.Update(&section)

		ret.Data, err = c.SectionService.List(section.OwnerID, reqData.Mode)
	case "removeSection":
		var fieldId uint
		fieldId, err = c.SectionService.Remove(reqData.Id, reqData.Mode)
		ret.Data, err = c.SectionService.List(fieldId, reqData.Mode)

	// field or instances refer, be create when init its owner
	case "getRefer":
		var refer model.ZdRefer
		refer, err = c.ReferService.Get(uint(reqData.Id), reqData.Mode)
		ret.Data = refer
	case "updateRefer":
		refer := serverUtils.ConvertRefer(reqData.Data)
		err = c.ReferService.Update(&refer)
	case "listReferFileForSelection":
		ret.Data = c.ResService.ListReferFileForSelection(reqData.Mode)
	case "listReferSheetForSelection":
		ret.Data = c.ResService.ListReferSheetForSelection(reqData.Mode)

	case "listReferExcelColForSelection":
		ret.Data = c.ResService.ListReferExcelColForSelection(reqData.Mode)
	case "listReferResFieldForSelection":
		ret.Data = c.ResService.ListReferResFieldForSelection(reqData.Id, reqData.Mode)

	// resource
	case "listRanges":
		ret.Data, ret.Total = c.RangesService.List(reqData.Keywords, reqData.Page)
	case "getRanges":
		ret.Data, ret.Res = c.RangesService.Get(reqData.Id)
		ret.WorkDir = vari.WorkDir
	case "saveRanges":
		ranges := serverUtils.ConvertRanges(reqData.Data)
		ret.Data = c.RangesService.Save(&ranges)
	case "removeRanges":
		err = c.RangesService.Remove(reqData.Id)

	case "getResRangesItemTree":
		ret.Data = c.RangesService.GetItemTree(reqData.Id)
	case "getResRangesItem":
		ret.Data = c.RangesService.GetItem(reqData.Id)
	case "createResRangesItem":
		var rangesItem *model.ZdRangesItem
		rangesId := reqData.DomainId
		rangesItem, err = c.RangesService.CreateItem(rangesId, reqData.Id, reqData.Mode)

		ret.Data = c.RangesService.GetItemTree(rangesId)
		ret.Model = rangesItem
	case "saveRangesItem":
		rangesItem := serverUtils.ConvertRangesItem(reqData.Data)
		ret.Data = c.RangesService.SaveItem(&rangesItem)
	case "removeResRangesItem":
		err = c.RangesService.RemoveItem(reqData.Id, reqData.DomainId)
		ret.Data = c.RangesService.GetItemTree(reqData.DomainId)

	case "listInstances":
		ret.Data, ret.Total = c.InstancesService.List(reqData.Keywords, reqData.Page)
	case "getInstances":
		ret.Data, ret.Res = c.InstancesService.Get(reqData.Id)
		ret.WorkDir = vari.WorkDir
	case "saveInstances":
		ranges := serverUtils.ConvertInstances(reqData.Data)
		ret.Data = c.InstancesService.Save(&ranges)
	case "removeInstances":
		err = c.InstancesService.Remove(reqData.Id)
	case "getResInstancesItemTree":
		ret.Data = c.InstancesService.GetItemTree(uint(reqData.Id))
	case "getResInstancesItem":
		ret.Data = c.InstancesService.GetItem(reqData.Id)
	case "createResInstancesItem":
		var item *model.ZdInstancesItem
		ownerId := reqData.DomainId
		item, err = c.InstancesService.CreateItem(ownerId, reqData.Id, reqData.Mode)

		ret.Data = c.InstancesService.GetItemTree(uint(ownerId))
		ret.Model = item
	case "saveInstancesItem":
		rangesItem := serverUtils.ConvertInstancesItem(reqData.Data)
		ret.Data = c.InstancesService.SaveItem(&rangesItem)
	case "removeResInstancesItem":
		err = c.InstancesService.RemoveItem(reqData.Id)
		ret.Data = c.InstancesService.GetItemTree(uint(reqData.DomainId))

	case "listExcel":
		ret.Data, ret.Total = c.ExcelService.List(reqData.Keywords, reqData.Page)
	case "getExcel":
		ret.Data, ret.Res = c.ExcelService.Get(reqData.Id)
		ret.WorkDir = vari.WorkDir
	case "saveExcel":
		ranges := serverUtils.ConvertExcel(reqData.Data)
		ret.Data = c.ExcelService.Save(&ranges)
	case "removeExcel":
		err = c.ExcelService.Remove(reqData.Id)

	case "listText":
		ret.Data, ret.Total = c.TextService.List(reqData.Keywords, reqData.Page)
	case "getText":
		ret.Data, ret.Res = c.TextService.Get(reqData.Id)
		ret.WorkDir = vari.WorkDir
	case "saveText":
		ranges := serverUtils.ConvertText(reqData.Data)
		ret.Data = c.TextService.Save(&ranges)
	case "removeText":
		err = c.TextService.Remove(reqData.Id)

	case "listConfig":
		ret.Data, ret.Total = c.ConfigService.List(reqData.Keywords, reqData.Page)
	case "getConfig":
		ret.Data, ret.Res = c.ConfigService.Get(reqData.Id)
		ret.WorkDir = vari.WorkDir
	case "saveConfig":
		ranges := serverUtils.ConvertConfig(reqData.Data)
		ret.Data = c.ConfigService.Save(&ranges)
	case "removeConfig":
		err = c.ConfigService.Remove(reqData.Id)

	case "getResConfigItemTree":
		ret.Data = c.ConfigService.GConfigItemTree(reqData.Id)

	default:
		ret.Code = 0
		ret.Msg = "api not found"
	}

	if err != nil {
		ret.Code = 0
		ret.Msg = "api error: " + err.Error()
	}

	ctx.JSON(ret)
}
