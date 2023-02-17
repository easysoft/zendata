package service

import (
	"context"
	"encoding/json"
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	shellUtils "github.com/easysoft/zendata/pkg/utils/shell"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

type MockService struct {
	FileService *FileService `inject:""`
}

func (s *MockService) GenMockDef(input string) (
	name, mockDefPath, zendataDefPath string, err error) { // return the last ones for client spec uploading

	var files []string
	fileUtils.GetFilesByExtInDir(input, ".yaml,.json", &files)

	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}

	for _, f := range files {
		if filepath.Ext(f) == ".json" {
			f = s.convertPostmanSpec(f)
		}

		doc3, err := loader.LoadFromFile(f)
		if err != nil {
			logUtils.PrintTo(fmt.Sprintf("skip file %s which is not a vaild openapi3, swagger and postman spec.", f))
			continue
		}

		fileName := filepath.Base(f)
		dir := filepath.Dir(f)
		if vari.GlobalVars.Output != "" {
			dir = vari.GlobalVars.Output
		}
		mockDefPath, zendataDefPath = s.getFilePaths(fileName, dir)

		zendataDef := domain.DefData{}
		zendataDef.ClsInfo.Title = doc3.Info.Title

		mockDef := domain.MockData{}
		mockDef.Title = doc3.Info.Title
		name = mockDef.Title

		if mockDef.Paths == nil {
			mockDef.Paths = map[string]map[string]map[string]map[string]*domain.EndPoint{}
		}

		for pathStr, pathItem := range doc3.Paths {
			mp := map[string]map[string]map[string]*domain.EndPoint{}

			if pathItem.Connect != nil {
				s.setEndPoint(pathItem.Connect, &zendataDef, zendataDefPath, domain.Connect, &mp)
			}

			if pathItem.Delete != nil {
				s.setEndPoint(pathItem.Delete, &zendataDef, zendataDefPath, domain.Delete, &mp)
			}

			if pathItem.Get != nil {
				s.setEndPoint(pathItem.Get, &zendataDef, zendataDefPath, domain.Get, &mp)
			}

			if pathItem.Head != nil {
				s.setEndPoint(pathItem.Head, &zendataDef, zendataDefPath, domain.Head, &mp)
			}

			if pathItem.Options != nil {
				s.setEndPoint(pathItem.Options, &zendataDef, zendataDefPath, domain.Options, &mp)
			}

			if pathItem.Patch != nil {
				s.setEndPoint(pathItem.Patch, &zendataDef, zendataDefPath, domain.Patch, &mp)
			}

			if pathItem.Post != nil {
				s.setEndPoint(pathItem.Post, &zendataDef, zendataDefPath, domain.Post, &mp)
			}

			if pathItem.Put != nil {
				s.setEndPoint(pathItem.Put, &zendataDef, zendataDefPath, domain.Put, &mp)
			}

			if pathItem.Trace != nil {
				s.setEndPoint(pathItem.Trace, &zendataDef, zendataDefPath, domain.Trace, &mp)
			}

			mockDef.Paths[pathStr] = mp
		}

		s.saveFile(mockDef, mockDefPath)
		s.saveFile(zendataDef, zendataDefPath)
	}

	return
}

func (s *MockService) convertPostmanSpec(input string) (ret string) {
	// npm i postman-to-openapi -g

	ret = input

	content, _ := os.ReadFile(input)

	if strings.Contains(string(content), "_postman_id") {
		ret = ret + ".yaml"
		cmd := fmt.Sprintf("p2o %s -f %s", input, ret)

		shellUtils.ExecInDir(cmd, filepath.Dir(ret))
	}

	return
}

func (s *MockService) setEndPoint(operation *openapi3.Operation, zendataDef *domain.DefData,
	zendataDefPath string, method domain.HttpMethod, mp *map[string]map[string]map[string]*domain.EndPoint) {

	codeToEndpointMap := s.createEndPoint(operation, zendataDef, zendataDefPath, method)
	(*mp)[method.String()] = codeToEndpointMap
}

func (s *MockService) createEndPoint(operation *openapi3.Operation, zendataDef *domain.DefData,
	zendataDefPath string, method domain.HttpMethod) (
	mockDef map[string]map[string]*domain.EndPoint) {

	mockDef = map[string]map[string]*domain.EndPoint{}

	for code, val := range operation.Responses {
		// map[string]*ResponseRef

		for mediaType, mediaItem := range val.Value.Content {
			if mediaItem == nil {
				continue
			}
			// mediaType is like "responses => 501 => content => application/json"

			// zendata def
			fields := s.genZendataDefFromMedia(mediaItem)
			zendataDef.Fields = append(zendataDef.Fields, fields...) // maybe has no children from properties

			// mock def
			endpoint := s.genMockDefFromMedia(mediaItem, fields)
			endpoint.Method = method
			endpoint.MediaType = mediaType
			endpoint.Config = filepath.Base(zendataDefPath) // set a relative path
			endpoint.Lines = 10

			if mockDef[code] == nil {
				mockDef[code] = map[string]*domain.EndPoint{}
			}
			mockDef[code][mediaType] = &endpoint
		}
	}

	return
}

func (s *MockService) genZendataDefFromMedia(item *openapi3.MediaType) (fields []domain.DefField) {
	schemaNode := item.Schema
	//exampleNode := item.Example
	//examplesNode := item.Examples
	//encodingNode := item.Encoding

	if schemaNode != nil {
		s.getFieldFromSchema(&fields, schemaNode)
	}

	return
}

func (s *MockService) genMockDefFromMedia(item *openapi3.MediaType, fields []domain.DefField) (endpoint domain.EndPoint) {
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
		} else {
			endpoint.Type = consts.SchemaTypeObject
		}
	}

	endpoint.Samples = map[string]string{}
	if schemaNode != nil && schemaNode.Value.Example != nil { // from schema's example
		s.getExample(schemaNode.Value.Example, &endpoint.Samples)
	}

	if exampleNode != nil {
		s.getExample(exampleNode, &endpoint.Samples)
	}
	if examplesNode != nil {
		s.getExamples(examplesNode, &endpoint.Samples)
	}

	return
}

func (s *MockService) getFieldFromSchema(fields *[]domain.DefField, schemaNodes ...*openapi3.SchemaRef) {
	propsMap := map[string]bool{}

	for _, schemaNode := range schemaNodes {
		if len(schemaNode.Value.Properties) > 0 { // properties based
			for propName, prop := range schemaNode.Value.Properties {
				if propsMap[propName] {
					continue
				}

				field := domain.DefField{}
				field.Field = propName // s.getSchemaNameFromRef(schemaNode.Ref) + "~" + propName
				propsMap[propName] = true

				if prop.Ref == "" { // leaf property
					s.getRangeByTypeFormat(consts.OpenApiDataType(prop.Value.Type),
						prop.Value.Enum, prop.Value.Default, prop.Value.Min, prop.Value.Max, &field)
				} else { // refer to another schema
					s.getFieldFromSchema(&field.Fields, prop)
				}

				*fields = append(*fields, field)
			}

		} else if schemaNode.Value.OneOf != nil {
			s.getFieldFromSchema(fields, schemaNode.Value.OneOf[0])

		} else if schemaNode.Value.AllOf != nil {
			s.getFieldFromSchema(fields, schemaNode.Value.AllOf...)

		} else if schemaNode.Value.AnyOf != nil {
			arr := openapi3.SchemaRefs{schemaNode.Value.AnyOf[0]}
			if len(schemaNode.Value.AnyOf) > 1 {
				arr = append(arr, schemaNode.Value.AnyOf[len(schemaNode.Value.AnyOf)-1])
			}

			s.getFieldFromSchema(fields, arr...)

		}

		// items based
		if schemaNode.Value.Items != nil {
			s.getFieldFromItems(fields, schemaNode.Value.Items)
		}
	}

	return
}

func (s *MockService) getSchemaNameFromRef(ref string) (ret string) {
	arr := strings.Split(ref, "/")
	ret = arr[len(arr)-1]

	return
}

func (s *MockService) getFieldFromItems(fields *[]domain.DefField, itemsDef *openapi3.SchemaRef) {
	s.getFieldFromSchema(fields, itemsDef)

	return
}

func (s *MockService) getRangeByTypeFormat(typ consts.OpenApiDataType,
	enums []interface{}, defaultVal interface{},
	min, max *float64,
	field *domain.DefField) {
	if enums != nil {
		field.Range = s.getRangeFromEnum(enums)
		return
	}

	if consts.OpenApiDataTypeInteger == typ {
		start, end := s.getStartEnd(1, 99, min, max, typ)
		field.Range = fmt.Sprintf("%d-%d", start, end)

	} else if consts.OpenApiDataTypeLong == typ {
		start, end := s.getStartEnd(9223372036854775801, 9223372036854775807, min, max, typ)
		field.Range = fmt.Sprintf("%d-%d", start, end)

	} else if consts.OpenApiDataTypeFloat == typ {
		start, end := s.getStartEnd(1.01, 99, min, max, typ)
		field.Range = fmt.Sprintf("%f-%f", start, end)

	} else if consts.OpenApiDataTypeDouble == typ {
		start, end := s.getStartEnd(1.000000000000009, 99, min, max, typ)
		field.Range = fmt.Sprintf("%f-%f", start, end)

	} else if consts.OpenApiDataTypeString == typ {
		field.Range = "a-z"
		field.Loop = "6-8"

	} else if consts.OpenApiDataTypeByte == typ {
		start, end := s.getStartEnd('a', 'z', min, max, typ)
		field.Range = fmt.Sprintf("%c-%c", start, end)

	} else if consts.OpenApiDataTypeBinary == typ {
		field.Format = "binary"

	} else if consts.OpenApiDataTypeBoolean == typ {
		field.Range = "[true,false]"

	} else if consts.OpenApiDataTypeDate == typ {
		field.Range = "20230101 000000-20230101 235959:60"
		field.Type = "timestamp"
		field.Format = "YY/MM/DD"

	} else if consts.OpenApiDataTypeDateTime == typ {
		field.Range = "20230101 000000-20230101 235959:60"
		field.Type = "timestamp"
		field.Format = "YY/MM/DD hh:mm:ss"

	} else if consts.OpenApiDataTypePassword == typ {
		field.Format = "password(8)"

	}

	if defaultVal != nil && (consts.OpenApiDataTypeInteger == typ || consts.OpenApiDataTypeLong == typ ||
		consts.OpenApiDataTypeFloat == typ || consts.OpenApiDataTypeDouble == typ ||
		consts.OpenApiDataTypeString == typ || consts.OpenApiDataTypeByte == typ) {

		field.Range = fmt.Sprintf("%v, ", defaultVal) + field.Range
	}
}

func (s *MockService) getRangeFromEnum(enums []interface{}) (ret string) {
	var arr []string
	for _, e := range enums {
		arr = append(arr, fmt.Sprintf("%v", e))
	}

	ret = fmt.Sprintf("[%s]", strings.Join(arr, ","))

	return
}

func (s *MockService) getFilePaths(name string, dir string) (mockPath, zendataPath string) {
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

func (s *MockService) getStartEnd(startDefault, endDefault interface{}, min, max *float64, typ consts.OpenApiDataType) (startRet, endRet interface{}) {
	startRet = startDefault
	endRet = endDefault

	if min != nil {
		if typ == consts.OpenApiDataTypeInteger {
			startRet = int(*min)
		} else if typ == consts.OpenApiDataTypeLong {
			startRet = int64(*min)
		} else if typ == consts.OpenApiDataTypeFloat {
			startRet = float32(*min)
		} else if typ == consts.OpenApiDataTypeDouble {
			startRet = *min
		} else if typ == consts.OpenApiDataTypeByte {
			startRet = int(*min)
		}
	}

	if max != nil {
		if typ == consts.OpenApiDataTypeInteger {
			endRet = int(*max)
		} else if typ == consts.OpenApiDataTypeLong {
			endRet = int64(*max)
		} else if typ == consts.OpenApiDataTypeFloat {
			endRet = float32(*max)
		} else if typ == consts.OpenApiDataTypeDouble {
			endRet = *max
		} else if typ == consts.OpenApiDataTypeByte {
			endRet = int(*max)
		}
	}

	return
}

func (s *MockService) getExample(exampleNode interface{}, sample *map[string]string) {
	bytes, _ := json.Marshal(exampleNode)
	(*sample)["example"] = string(bytes)
}

func (s *MockService) getExamples(exampleNodes openapi3.Examples, sample *map[string]string) {
	for key, val := range exampleNodes {
		bytes, _ := json.Marshal(val.Value.Value)
		(*sample)[key] = string(bytes)
	}
}
