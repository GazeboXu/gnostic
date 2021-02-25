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

	callerSrc := ""
	typeSrc := ""
	for _, model := range env.Request.Models {
		switch model.TypeUrl {
		case "openapi.v2.Document":
			documentv2 := &openapiv2.Document{}
			err = proto.Unmarshal(model.Value, documentv2)
			if err == nil {
				// Analyze the API document.
				callerSrc, typeSrc = v2doc2Gin(documentv2)
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
	file.Data = []byte(callerSrc)
	env.Response.Files = append(env.Response.Files, file)

	outputName2 := filepath.Join(
		filepath.Dir(env.Request.SourceName), "/optype/optype_autogen.go")
	file2 := &plugins.File{}

	file2.Name = outputName2
	file2.Data = []byte(typeSrc)
	env.Response.Files = append(env.Response.Files, file2)

	env.RespondAndExitIfError(err)
	env.RespondAndExit()
}

const GIN_TEMPLATE = `
// gin handler file generated from openapi doc
// Created by plugin gengin of gnostic at {{.CreatedAt}}
// WARNING! All changes made in this file will be lost when building

package docs

import (
	"context"
	"github.com/gin-gonic/gin"
	"wrnetman/wrutils"
	"wrnetman/netadapter/overhttp/netmodel"
	"wrnetman/docs/optype"
	{{range .Imports}}{{.}}
	{{end}}
)

// optype name is the func name prefix with 'AT'
// func ExampleFunc(ctx context.Context, param *optype.ATExampleFunc) (result interface{}, err error) {
//	return result, err
// }

func RouterInit(r *gin.Engine) {
	{{range .APIs}}
	r.{{.HTTPMethod}}("{{.BasePath}}{{.FullMethodPath}}", func(c *gin.Context) {
		param := &optype.{{.TypeName}}{}
		{{.PreParams}}
		if result, err := {{.MethodPackage}}.{{.BareMethodName}}(context.TODO(), param); err == nil {
			c.JSON(200, netmodel.CallResult{
				HTTPCode: 200,
				Data: result,
			})
		} else {
			c.JSON(200, netmodel.CallResult{
				HTTPCode: 200,
				Code: -1,
				ErrMsg: err.Error(),
			})
		}
	})
	{{end}}
}
`

const GIN_TYPE_TEMPLATE = `
// type file generated from openapi doc
// Created by plugin gengin of gnostic at {{.CreatedAt}}
// WARNING! All changes made in this file will be lost when building

package optype

{{range .APITypes}}
type {{.TypeName}} struct {
{{range .Fields}}	{{.FieldName}}	{{.FieldType}}	{{.FieldRemark}}
{{end}}
}
{{end}}
`

type GinTemplateInfo struct {
	CreatedAt string
	Imports   []string
	APITypes  []*APIType
	APIs      []*GinPathTemplateInfo
}

type APIType struct {
	TypeName string
	Fields   []*FieldInfo
}

type FieldInfo struct {
	FieldName   string
	FieldType   string
	FieldRemark string
}

type GinPathTemplateInfo struct {
	HTTPMethod     string
	BasePath       string
	FullMethodPath string
	MethodPackage  string
	BareMethodName string
	PreParams      string
	// Params         string
	TypeName string
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

func getTypeConvFun(goType string) (convFunc string) {
	if goType == "string" {
	} else if goType == "int" {
		convFunc = "wrutils.String2Int"
	} else if goType == "int64" {
		convFunc = "wrutils.String2Int64"
	} else if goType == "float64" {
		convFunc = "wrutils.String2Float"
	} else if goType == "bool" {
		convFunc = "wrutils.String2Bool"
	}
	return convFunc
}

func openAPIType2Go(paramType, paramFormat string) (goType string) {
	goType = "string"
	if paramType == "string" {
	} else if paramType == "integer" {
		goType = "int"
		if paramFormat == "int64" {
			goType = "int64"
		}
	} else if paramType == "number" {
		goType = "float64"
	} else if paramType == "boolean" {
		goType = "bool"
	}
	return goType
}

var methodList = []string{"GET", "PUT", "POST", "DELETE", "OPTIONS", "HEAD", "PATCH"}
var supportedMethodMap = map[string]bool{
	"GET":  true,
	"POST": true,
}

func UpperFirstLetter(str string) string {
	if len(str) > 0 {
		str = strings.ToUpper(str[:1]) + str[1:]
	}
	return str
}

func v2doc2Gin(doc *openapiv2.Document) (goSource, goTypeSrc string) {
	info := &GinTemplateInfo{
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	if tmp, err := template.New("caller").Parse(GIN_TEMPLATE); err == nil {
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

			// paramList := make([]string, len(operation.Parameters))
			preParamList := make([]string, len(operation.Parameters))
			apiType := &APIType{
				TypeName: fmt.Sprintf("AT%s", methodInfo.BareMethodName),
			}
			methodInfo.TypeName = apiType.TypeName

			// apiType.Fields = append(apiType.Fields, &FieldInfo{
			// 	FieldName:   "MapData",
			// 	FieldType:   "map[string]interface{}",
			// 	FieldRemark: "`json:\"-\"`",
			// })
			for index, param := range operation.Parameters {
				// paramList[index] = fmt.Sprintf("param%d", index+1)
				var paramName, paramType, paramFormat, getParamFuncName string
				var paramRequired bool
				if nonBodyParam := param.GetParameter().GetNonBodyParameter(); nonBodyParam != nil {
					if methodInfo.HTTPMethod == "GET" {
						subSchema := nonBodyParam.GetQueryParameterSubSchema()
						paramName = subSchema.Name
						paramType = subSchema.Type
						paramFormat = subSchema.Format
						paramRequired = subSchema.Required
						getParamFuncName = "GetQuery"
					} else if methodInfo.HTTPMethod == "POST" {
						subSchema := nonBodyParam.GetFormDataParameterSubSchema()
						paramName = subSchema.Name
						paramType = subSchema.Type
						paramFormat = subSchema.Format
						paramRequired = subSchema.Required
						getParamFuncName = "GetPostForm"
					}
				} else if bodyParam := param.GetParameter().GetBodyParameter(); bodyParam != nil { // todo
					paramName = bodyParam.Name
					paramRequired = bodyParam.Required
					getParamFuncName = "GetQuery"
				} else {
					panic("parameter is not either body or nonbody")
				}
				paramGoType := openAPIType2Go(paramType, paramFormat)
				fieldInfo := &FieldInfo{
					FieldName:   UpperFirstLetter(paramName),
					FieldRemark: fmt.Sprintf("`json:\"%s,omitempty\"`", paramName),
				}
				if paramRequired {
					fieldInfo.FieldType = paramGoType
					preParamList[index] = fmt.Sprintf(`
		if strValue, isExist := c.%s("%s"); isExist {
			param.%s = %s(strValue)
		} else {
		}`, getParamFuncName, paramName, fieldInfo.FieldName, getTypeConvFun(paramGoType))
				} else {
					fieldInfo.FieldType = fmt.Sprintf("*%s", paramGoType)
					preParamList[index] = fmt.Sprintf(`
		if strValue, isExist := c.%s("%s"); isExist {
			tmpValue := %s(strValue)
			param.%s = &tmpValue
		}`, getParamFuncName, paramName, getTypeConvFun(paramGoType), fieldInfo.FieldName)
				}
				apiType.Fields = append(apiType.Fields, fieldInfo)
			}
			methodInfo.PreParams = strings.Join(preParamList, "\n")
			// methodInfo.Params = strings.Join(paramList, ", ")
			info.APIs = append(info.APIs, methodInfo)
			info.APITypes = append(info.APITypes, apiType)
		}
		for key := range importMap {
			info.Imports = append(info.Imports, fmt.Sprintf(`"wrnetman/netadapter/overhttp/%s"`, key))
		}

		builder := &strings.Builder{}
		tmp.Execute(builder, info)
		goSource = builder.String()
	}
	if tmp, err := template.New("type").Parse(GIN_TYPE_TEMPLATE); err == nil {
		builder := &strings.Builder{}
		tmp.Execute(builder, info)
		goTypeSrc = builder.String()
	}
	return goSource, goTypeSrc
}
