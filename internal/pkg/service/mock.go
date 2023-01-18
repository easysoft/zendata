package service

import (
	"context"
	"encoding/json"
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/getkin/kin-openapi/openapi3"
	uuid "github.com/satori/go.uuid"
	"log"
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

		data := model.MockData{}
		if data.Paths == nil {
			data.Paths = map[string]map[string]map[string]map[string]*model.EndPoint{}
		}

		for pathStr, pathItem := range doc3.Paths {
			mp := map[string]map[string]map[string]*model.EndPoint{}

			if pathItem.Connect != nil {

			}

			if pathItem.Delete != nil {

			}

			if pathItem.Get != nil {
				codeToEndpointMap := s.createEndPoint(pathItem.Get)
				mp[model.Connect.String()] = codeToEndpointMap
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

			data.Paths[pathStr] = mp
		}

		log.Print(doc3)
	}

	return
}

func (s *MockService) createEndPoint(operation *openapi3.Operation) (ret map[string]map[string]*model.EndPoint) {
	ret = map[string]map[string]*model.EndPoint{}

	for code, val := range operation.Responses {
		// map[string]*ResponseRef

		for mediaType, val := range val.Value.Content {
			schemaNode := val.Schema
			exampleNode := val.Example
			examplesNode := val.Examples
			encodingNode := val.Encoding

			log.Print(schemaNode, exampleNode, examplesNode, encodingNode)

			// gen zendata def
			defName := ""
			if schemaNode != nil && schemaNode.Ref != "" {
				defName = strings.TrimPrefix(strings.ReplaceAll(schemaNode.Ref, "/", "-"), "#/")
			} else {
				defName = fmt.Sprintf("field-%s", strings.ReplaceAll(uuid.NewV4().String(), "-", ""))
			}
			log.Print(defName)

			def := model.DefData{}
			for propName, prop := range schemaNode.Value.Properties {
				field := model.DefField{}
				field.Field = propName
				s.getRangeByType(prop.Value.Type, &field)

				if exampleNode != nil {
					s.getExample(exampleNode, &field)
				} else if examplesNode != nil {
					s.getExamples(examplesNode, &field)
				}

				def.Fields = append(def.Fields, field)
			}

			// gen mock def
			endpoint := model.EndPoint{}
			if ret[code] == nil {
				ret[code] = map[string]*model.EndPoint{}
			}
			ret[code][mediaType] = &endpoint
		}
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

func (s *MockService) getExample(node interface{}, field *model.DefField) {
	bytes, _ := json.Marshal(node)
	field.Range = fmt.Sprintf("`%s`", bytes)
}

func (s *MockService) getExamples(examples openapi3.Examples, field *model.DefField) {
	for _, val := range examples {
		bytes, _ := json.Marshal(val.Value.Value)
		field.Range = fmt.Sprintf("`%s`", bytes)
	}
}
