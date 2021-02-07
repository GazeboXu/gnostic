package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/golang/protobuf/proto"
	discovery_v1 "github.com/googleapis/gnostic/discovery"
	openapiv2 "github.com/googleapis/gnostic/openapiv2"
	openapiv3 "github.com/googleapis/gnostic/openapiv3"
	plugins "github.com/googleapis/gnostic/plugins"
)

// This is the main function for the plugin.
func main() {
	env, err := plugins.NewEnvironment()
	env.RespondAndExitIfError(err)

	fileContent := ""
	for _, model := range env.Request.Models {
		switch model.TypeUrl {
		case "openapi.v2.Document":
			documentv2 := &openapiv2.Document{}
			err = proto.Unmarshal(model.Value, documentv2)
			if err == nil {
				// Analyze the API document.
				fileContent = v2doc2Gin(documentv2)
			}
		case "openapi.v3.Document":
			documentv3 := &openapiv3.Document{}
			err = proto.Unmarshal(model.Value, documentv3)
			if err == nil {
				// Analyze the API document.
			}
		case "discovery.v1.Document":
			discoveryDocument := &discovery_v1.Document{}
			err = proto.Unmarshal(model.Value, discoveryDocument)
			if err == nil {
				// Analyze the API document.
			}
		default:
			// log.Printf("unsupported document type %s", model.TypeUrl)
		}
	}

	outputName1 := filepath.Join(
		filepath.Dir(env.Request.SourceName), "gin_autogen.go")
	file := &plugins.File{}

	file.Name = outputName1
	file.Data = []byte(fileContent)
	env.Response.Files = append(env.Response.Files, file)
	env.RespondAndExitIfError(err)
	env.RespondAndExit()
}

const GIN_TEMPLATE = `
// gin handler file generated from openapi doc
// Created by plugin gengin of gnostic at {{.CreatedAt}}
// WARNING! All changes made in this file will be lost when building

package docs

import (
	"github.com/gin-gonic/gin"
	"wrnetman/wrutils"
	"wrnetman/netadapter/overhttp/netmodel"
	{{range .Imports}}{{.}}
	{{end}}
)

func RouterInit(r *gin.Engine) {
	{{range .APIs}}{{.}}{{end}}
}
`
const GIN_PATH_TEMPLATE = `
	r.{{.HTTPMethod}}("{{.BasePath}}{{.FullMethodPath}}", func(c *gin.Context) {
		c.JSON(200, netmodel.CallResult{
			Data: {{.MethodPackage}}.{{.BareMethodName}}({{.Params}}),
		})
	})`

type GinTemplateInfo struct {
	CreatedAt string
	Imports   []string
	APIs      []string
}

type GinPathTemplateInfo struct {
	HTTPMethod     string
	BasePath       string
	FullMethodPath string
	MethodPackage  string
	BareMethodName string
	Params         string
}

func splitFullMethodPath(fullMethodPath string) (methodPackage, bareMethod string) {
	strList := strings.Split(strings.Trim(fullMethodPath, "/"), "/")
	listLen := len(strList)
	if listLen > 1 {
		methodPackage = strList[listLen-2]
	}
	bareMethod = strList[listLen-1]
	return methodPackage, bareMethod
}

func getTypeConvFun(paramType, paramFormat string) (convFunc string) {
	if paramType == "string" {
	} else if paramType == "integer" {
		convFunc = "wrutils.String2Int"
		if paramFormat == "int64" {
			convFunc = "wrutils.String2Int64"
		}
	} else if paramType == "number" {
		convFunc = "wrutils.String2Float"
	} else if paramType == "boolean" {
		convFunc = "wrutils.String2Bool"
	}
	return convFunc
}

var methodList = []string{"GET", "PUT", "POST", "DELETE", "OPTIONS", "HEAD", "PATCH"}
var supportedMethodMap = map[string]bool{
	"GET":  true,
	"POST": true,
}

func v2doc2Gin(doc *openapiv2.Document) (goSource string) {
	if tmp, err := template.New("test").Parse(GIN_TEMPLATE); err == nil {
		subTmp, _ := template.New("subtest").Parse(GIN_PATH_TEMPLATE)
		builder := &strings.Builder{}
		info := &GinTemplateInfo{
			CreatedAt: time.Now().Format(time.RFC3339),
		}
		importMap := make(map[string]string)
		for _, pathItem := range doc.Paths.Path {
			methodInfo := &GinPathTemplateInfo{
				BasePath:       "/v1",
				FullMethodPath: pathItem.Name,
			}
			methodInfo.MethodPackage, methodInfo.BareMethodName = splitFullMethodPath(pathItem.Name)
			methodInfo.MethodPackage = "http" + methodInfo.MethodPackage
			importMap[methodInfo.MethodPackage] = methodInfo.MethodPackage
			var operation *openapiv2.Operation
			var index int
			for index, operation = range []*openapiv2.Operation{
				pathItem.Value.Get,
				pathItem.Value.Put,
				pathItem.Value.Post,
				pathItem.Value.Delete,
				pathItem.Value.Options,
				pathItem.Value.Head,
				pathItem.Value.Patch,
			} {
				if operation != nil {
					methodInfo.HTTPMethod = methodList[index]
					break
				}
			}
			if !supportedMethodMap[methodInfo.HTTPMethod] {
				msg := fmt.Sprintf("unsupported method:%s!\n", methodInfo.HTTPMethod)
				log.Panicln(msg)
				methodInfo.FullMethodPath = msg
				break
			}

			var paramList []string
			for _, param := range operation.Parameters {
				nonBodyParam := param.GetParameter().GetNonBodyParameter()
				if methodInfo.HTTPMethod == "GET" {
					querySchema := nonBodyParam.GetQueryParameterSubSchema()
					paramList = append(paramList, fmt.Sprintf(`%s(c.Query("%s"))`, getTypeConvFun(querySchema.Type, querySchema.Format), querySchema.Name))
				} else if methodInfo.HTTPMethod == "POST" {
					formSchema := nonBodyParam.GetFormDataParameterSubSchema()
					paramList = append(paramList, fmt.Sprintf(`%s(c.PostForm("%s"))`, getTypeConvFun(formSchema.Type, formSchema.Format), formSchema.Name))
				}
			}
			methodInfo.Params = strings.Join(paramList, ", ")
			methodBuilder := &strings.Builder{}
			subTmp.Execute(methodBuilder, methodInfo)

			info.APIs = append(info.APIs, methodBuilder.String())
		}
		for key := range importMap {
			info.Imports = append(info.Imports, fmt.Sprintf(`"wrnetman/netadapter/overhttp/%s"`, key))
		}

		tmp.Execute(builder, info)
		goSource = builder.String()
	}
	return goSource
}
