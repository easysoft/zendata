package service

import (
	"context"
	"encoding/json"
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"strings"
)

type MockService struct {
	FileService *FileService `inject:""`
}

func (s *MockService) GenMockDef(input string) (err error) {
	var files []string
	fileUtils.GetFilesByExtInDir(input, ".yaml", &files)

	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}

	for _, f := range files {
		doc3, err := loader.LoadFromFile(f)
		if err != nil {
			continue
		}

		name := filepath.Base(f)
		dir := filepath.Dir(f)
		if vari.GlobalVars.Output != "" {
			dir = vari.GlobalVars.Output
		}
		zendataDefPath, mockDefPath := s.getFilePaths(name, dir)

		zendataDef := model.DefData{}
		zendataDef.ClsInfo.Title = doc3.Info.Title

		mockDef := model.MockData{}
		mockDef.Title = doc3.Info.Title

		if mockDef.Paths == nil {
			mockDef.Paths = map[string]map[string]map[string]map[string]*model.EndPoint{}
		}

		for pathStr, pathItem := range doc3.Paths {
			mp := map[string]map[string]map[string]*model.EndPoint{}

			if pathItem.Connect != nil {

			}

			if pathItem.Delete != nil {

			}

			if pathItem.Get != nil {
				codeToEndpointMap := s.createEndPoint(pathItem.Get, &zendataDef, zendataDefPath, model.Get)
				mp[model.Get.String()] = codeToEndpointMap
			}

			if pathItem.Head != nil {

			}

			if pathItem.Options != nil {

			}

			if pathItem.Patch != nil {

			}

			if pathItem.Post != nil {

			}

			if pathItem.Put != nil {

			}

			if pathItem.Trace != nil {

			}

			mockDef.Paths[pathStr] = mp
		}

		bytesZd, err := yaml.Marshal(zendataDef)
		s.saveFile(bytesZd, zendataDefPath)

		bytesMock, err := yaml.Marshal(mockDef)
		s.saveFile(bytesMock, mockDefPath)
	}

	return
}

func (s *MockService) createEndPoint(operation *openapi3.Operation, zendataDef *model.DefData,
	zendataDefPath string, method model.HttpMethod) (
	mockDef map[string]map[string]*model.EndPoint) {

	mockDef = map[string]map[string]*model.EndPoint{}

	for code, val := range operation.Responses {
		// map[string]*ResponseRef

		for mediaType, mediaItem := range val.Value.Content {
			// mediaType is responses => 501 => content => application/json:

			// zendata def
			fields := s.getZendataDefFromMedia(mediaItem)
			zendataDef.Fields = append(zendataDef.Fields, fields...)

			// mock def
			endpoint := s.getMockDefFromMedia(mediaItem, fields)
			endpoint.Method = method
			endpoint.MediaType = mediaType
			endpoint.Config = zendataDefPath
			endpoint.Lines = 10

			if mockDef[code] == nil {
				mockDef[code] = map[string]*model.EndPoint{}
			}
			mockDef[code][mediaType] = &endpoint
		}
	}

	return
}

func (s *MockService) getZendataDefFromMedia(item *openapi3.MediaType) (fields []model.DefField) {
	schemaNode := item.Schema
	exampleNode := item.Example
	examplesNode := item.Examples
	//encodingNode := item.Encoding

	if schemaNode != nil {
		schemaField := s.getFieldFromSchema(schemaNode)
		fields = append(fields, schemaField)
	}

	if exampleNode != nil {
		exampleField := s.getFieldFromExample(exampleNode)
		fields = append(fields, exampleField)
	}

	if examplesNode != nil {
		examplesFields := s.getFieldFromExamples(examplesNode)

		for _, field := range examplesFields {
			fields = append(fields, field)
		}
	}

	return
}

func (s *MockService) getMockDefFromMedia(item *openapi3.MediaType, fields []model.DefField) (endpoint model.EndPoint) {
	var fieldNames []string
	for _, f := range fields {
		fieldNames = append(fieldNames, f.Field)
	}
	endpoint.Fields = strings.Join(fieldNames, ",")

	schemaNode := item.Schema
	exampleNode := item.Example
	examplesNode := item.Examples
	//encodingNode := item.Encoding

	if schemaNode != nil {
		if schemaNode.Value.Type == string(consts.SchemaTypeArray) {
			endpoint.Type = consts.SchemaTypeArray
		} else {
			endpoint.Type = consts.OpenApiSchemaType(schemaNode.Value.Type)
		}

	} else if exampleNode != nil {

	} else if examplesNode != nil {

	}

	return
}

func (s *MockService) getFieldFromSchema(schemaNode *openapi3.SchemaRef) (ret model.DefField) {
	ret.Field = "schema"

	for propName, prop := range schemaNode.Value.Properties {
		field := model.DefField{}
		field.Field = propName
		s.getRangeByType(prop.Value.Type, &field)

		ret.Fields = append(ret.Fields, field)
	}

	return
}

func (s *MockService) getFieldFromExample(example interface{}) (field model.DefField) {
	bytes, _ := json.Marshal(example)

	field.Field = "example"
	field.Range = fmt.Sprintf("%s", bytes)

	return
}

func (s *MockService) getFieldFromExamples(examples openapi3.Examples) (fields []model.DefField) {
	for key, val := range examples {
		bytes, _ := json.Marshal(val.Value.Value)

		field := model.DefField{}
		field.Field = "examples_" + key
		field.Range = fmt.Sprintf("%s", bytes)

		fields = append(fields, field)
	}

	return
}

func (s *MockService) getRangeByType(typ string, field *model.DefField) {
	if string(consts.Integer) == typ {
		field.Range = "1-99"

	} else if string(consts.Long) == typ {
		field.Range = "1-99"

	} else if string(consts.Float) == typ {
		field.Range = "1.01-99"

	} else if string(consts.Double) == typ {
		field.Range = "1.01-99"

	} else if string(consts.String) == typ {
		field.Range = "1.01-99"

	} else if string(consts.Byte) == typ {
		field.Range = "a-z"
		field.Loop = "6"

	} else if string(consts.Binary) == typ {
		field.Format = "binary"

	} else if string(consts.Boolean) == typ {
		field.Range = "[true,false]"

	} else if string(consts.Date) == typ {
		field.Range = "20210101 000000-20210101 230000:60"
		field.Type = "timestamp"
		field.Format = "YY/MM/DD"

	} else if string(consts.DateTime) == typ {
		field.Range = "20210101 000000-20210101 230000:60"
		field.Type = "timestamp"
		field.Format = "YY/MM/DD hh:mm:ss"

	} else if string(consts.Password) == typ {
		field.Format = "password(8)"

	}
}

func (s *MockService) getFilePaths(name string, dir string) (zendataPath, mockPath string) {
	ext := filepath.Ext(name)

	zendataPath = strings.ReplaceAll(name, ext, "-zd"+ext)
	zendataPath = filepath.Join(dir, zendataPath)

	mockPath = strings.ReplaceAll(name, ext, "-mock"+ext)
	mockPath = filepath.Join(dir, mockPath)

	return
}

func (s *MockService) saveFile(bytes []byte, pth string) {
	fileUtils.MkDirIfNeeded(filepath.Dir(pth))

	str := string(bytes)

	fileUtils.WriteFile(pth, str)
}
