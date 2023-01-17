package service

import (
	"context"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/getkin/kin-openapi/openapi3"
	"log"
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
	for code, val := range operation.Responses {
		// map[string]*ResponseRef

		for mediaType, val := range val.Value.Content {
			schema := val.Schema
			example := val.Example
			examples := val.Examples
			encoding := val.Encoding

			log.Print(schema, example, examples, encoding)

			endpoint := model.EndPoint{}

			if ret[code] == nil {
				ret[code] = map[string]*model.EndPoint{}
			}
			ret[code][mediaType] = &endpoint
		}
	}

	return
}
