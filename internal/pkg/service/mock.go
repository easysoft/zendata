package service

import (
	"context"
	"encoding/json"
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
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
			logUtils.PrintTo(err.Error())
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
				s.setEndPoint(pathItem.Get, &zendataDef, zendataDefPath, model.Get, &mp)
			}

			if pathItem.Delete != nil {
				s.setEndPoint(pathItem.Get, &zendataDef, zendataDefPath, model.Get, &mp)
			}

			if pathItem.Get != nil {
				s.setEndPoint(pathItem.Get, &zendataDef, zendataDefPath, model.Get, &mp)
			}

			if pathItem.Head != nil {
				s.setEndPoint(pathItem.Get, &zendataDef, zendataDefPath, model.Get, &mp)
			}

			if pathItem.Options != nil {
				s.setEndPoint(pathItem.Get, &zendataDef, zendataDefPath, model.Get, &mp)
			}

			if pathItem.Patch != nil {
				s.setEndPoint(pathItem.Get, &zendataDef, zendataDefPath, model.Get, &mp)
			}

			if pathItem.Post != nil {
				s.setEndPoint(pathItem.Get, &zendataDef, zendataDefPath, model.Get, &mp)
			}

			if pathItem.Put != nil {
				s.setEndPoint(pathItem.Get, &zendataDef, zendataDefPath, model.Get, &mp)
			}

			if pathItem.Trace != nil {
				s.setEndPoint(pathItem.Get, &zendataDef, zendataDefPath, model.Get, &mp)
			}

			mockDef.Paths[pathStr] = mp
		}

		s.saveFile(zendataDef, zendataDefPath)
		s.saveFile(mockDef, mockDefPath)
	}

	return
}

func (s *MockService) setEndPoint(operation *openapi3.Operation, zendataDef *model.DefData,
	zendataDefPath string, method model.HttpMethod, mp *map[string]map[string]map[string]*model.EndPoint) {

	codeToEndpointMap := s.createEndPoint(operation, zendataDef, zendataDefPath, method)
	(*mp)[method.String()] = codeToEndpointMap
}

func (s *MockService) createEndPoint(operation *openapi3.Operation, zendataDef *model.DefData,
	zendataDefPath string, method model.HttpMethod) (
	mockDef map[string]map[string]*model.EndPoint) {

	mockDef = map[string]map[string]*model.EndPoint{}

	for code, val := range operation.Responses {
		// map[string]*ResponseRef

		for mediaType, mediaItem := range val.Value.Content {
			// mediaType is like "responses => 501 => content => application/json"

			// zendata def
			fields := s.getZendataDefFromMedia(mediaItem)
			zendataDef.Fields = append(zendataDef.Fields, fields...) // maybe has no children from properties

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
		s.getFieldFromSchema("schema", &fields, schemaNode)
	}

	if exampleNode != nil {
		exampleField := s.getFieldFromExample("example", exampleNode)
		fields = append(fields, exampleField)
	}

	if examplesNode != nil {
		examplesFields := s.getFieldFromExamples("example", examplesNode)

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
		} else if schemaNode.Value.Type == string(consts.SchemaTypeObject) {
			endpoint.Type = consts.SchemaTypeObject
		} else if schemaNode.Value.Type != "" {
			endpoint.Type = consts.OpenApiSchemaType(schemaNode.Value.Type)
		}

	} else if exampleNode != nil {

	} else if examplesNode != nil {

	}

	return
}

func (s *MockService) getFieldFromSchema(name string, fields *[]model.DefField, schemaNodes ...*openapi3.SchemaRef) {
	for _, schemaNode := range schemaNodes {
		// properties based
		if len(schemaNode.Value.Properties) > 0 {
			for propName, prop := range schemaNode.Value.Properties {
				field := model.DefField{}
				field.Field = propName
				s.getRangeByType(prop.Value.Type, &field)

				*fields = append(*fields, field)
			}

		} else if schemaNode.Value.OneOf != nil {
			s.getFieldFromSchema(name+"-oneof", fields, schemaNode.Value.OneOf[0])

		} else if schemaNode.Value.AllOf != nil {
			s.getFieldFromSchema(name+"-allof", fields, schemaNode.Value.AllOf...)

		} else if schemaNode.Value.AnyOf != nil {
			arr := openapi3.SchemaRefs{schemaNode.Value.AnyOf[0]}
			if len(schemaNode.Value.AnyOf) > 1 {
				arr = append(arr, schemaNode.Value.AnyOf[len(schemaNode.Value.AnyOf)-1])
			}

			s.getFieldFromSchema(name+"-anyof", fields, arr...)

		}

		// example based
		if schemaNode.Value.Example != nil {
			exampleField := s.getFieldFromExample(name+"-example", schemaNode.Value.Example)

			*fields = append(*fields, exampleField)
		}

		// items based
		if schemaNode.Value.Items != nil {
			s.getFieldFromItems(name+"-items", fields, schemaNode.Value.Items)
		}
	}

	return
}

func (s *MockService) getFieldFromExample(name string, example interface{}) (field model.DefField) {
	bytes, _ := json.Marshal(example)

	field.Field = name
	field.Range = fmt.Sprintf("%s", bytes)

	return
}

func (s *MockService) getFieldFromExamples(name string, examples openapi3.Examples) (fields []model.DefField) {
	for key, val := range examples {
		bytes, _ := json.Marshal(val.Value.Value)

		field := model.DefField{}
		field.Field = name + "-" + key
		field.Range = fmt.Sprintf("%s", bytes)

		fields = append(fields, field)
	}

	return
}

func (s *MockService) getFieldFromItems(name string, fields *[]model.DefField, itemsDef *openapi3.SchemaRef) {
	s.getFieldFromSchema(name, fields, itemsDef)

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
		field.Range = "a-z"
		field.Loop = "6-8"

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

func (s *MockService) saveFile(obj interface{}, pth string) {
	fileUtils.MkDirIfNeeded(filepath.Dir(pth))

	bytes, err := yaml.Marshal(obj)
	str := string(bytes)
	if err != nil {
		str = err.Error()
	}

	fileUtils.WriteFile(pth, str)
}
