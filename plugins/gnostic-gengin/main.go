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
	goTypeSrc := ""
	cTypeSrc := ""
	for _, model := range env.Request.Models {
		switch model.TypeUrl {
		case "openapi.v2.Document":
			documentv2 := &openapiv2.Document{}
			err = proto.Unmarshal(model.Value, documentv2)
			if err == nil {
				// Analyze the API document.
				callerSrc, goTypeSrc, cTypeSrc = v2doc2Gin(documentv2)
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
	file2.Data = []byte(goTypeSrc)
	env.Response.Files = append(env.Response.Files, file2)

	outputName3 := filepath.Join(
		filepath.Dir(env.Request.SourceName), "/ctype/metadatastruct.h")
	file3 := &plugins.File{}

	file3.Name = outputName3
	file3.Data = []byte(cTypeSrc)
	env.Response.Files = append(env.Response.Files, file3)

	env.RespondAndExitIfError(err)
	env.RespondAndExit()
}

const GIN_TEMPLATE = `
// gin handler file generated from openapi doc
// Created by plugin gengin of gnostic at {{.CreatedAt}}
// WARNING! All changes made in this file will be lost when rebuilding

package docs

import (
	"context"
	"github.com/gin-gonic/gin"
	"gitee.com/julytech/zlutils"
	"{{.ModName}}/netadapter/overhttp/netmodel"
	"{{.ModName}}/docs/optype"
	"{{.ModName}}/biz/bizutils"
	{{range .Imports}}{{.}}
	{{end}}
)

// optype name is the func name prefix with 'AT'
// func ExampleFunc(ctx context.Context, param *optype.ATExampleFunc) (result interface{}, err error) {
//	return result, err
// }

var FuncMap = map[string]interface{} {
	{{range .APIs}}"{{.BasePath}}{{.FullMethodPath}}" : {{.MethodPackage}}.{{.BareMethodName}},
	{{end}}
}

func RouterInit(r *gin.Engine) {
	{{range .APIs}}
	r.{{.HTTPMethod}}("{{.BasePath}}{{.FullMethodPath}}", func(c *gin.Context) {
		param := &optype.{{.TypeName}} {
		}
		{{.PreParams}}
		if err := bizutils.ValidateParam(param); err == nil {
			{{.CheckToken}}
			if result, err := {{.MethodPackage}}.{{.BareMethodName}}(context.TODO(), param); err == nil {
				c.JSON(200, netmodel.CallResult{
					BaseCallResult : netmodel.BaseCallResult {
						HTTPCode: 200,
					},
					Data: result,
				})
			} else {
				c.JSON(200, netmodel.CallResult{
					BaseCallResult : netmodel.BaseCallResult {
						HTTPCode: 200,
						Code: -1,
						ErrMsg: err.Error(),
					},
				})
			}
		} else {
			c.JSON(200, netmodel.CallResult{
				BaseCallResult : netmodel.BaseCallResult {
					HTTPCode: 200,
					Code: -2,
					ErrMsg: err.Error(),
				},
			})
		}
	})
	{{end}}
}
`

const checkTokenStatement = `
			if param.TokenInfo, err = bizutils.PreValidateToken(param.Token); err != nil {
				c.JSON(200, netmodel.CallResult{
					BaseCallResult : netmodel.BaseCallResult {
						HTTPCode: 200,
						Code: -3,
						ErrMsg: err.Error(),
					},
				})
				return
			}`

const GIN_TYPE_TEMPLATE = `
// type file generated from openapi doc
// Created by plugin gengin of gnostic at {{.CreatedAt}}
// WARNING! All changes made in this file will be lost when rebuilding

package optype

import (
	"time"
)

var _ time.Time

{{range .APITypes}}
type {{.TypeName}} struct {
{{range .Fields}}	{{.FieldName}}	{{.FieldType}}	{{.FieldRemark}}
{{end}}
}
{{end}}
`

type GinTemplateInfo struct {
	ModName        string
	CreatedAt      string
	Imports        []string
	APITypes       []*APIType
	APIs           []*GinPathTemplateInfo
	NetDefinitions []*NetDefinition
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
	TypeName   string
	CheckToken string
}

const GIN_DEFINITION4C_TEMPLATE = `
// c/c++ type file generated from openapi doc
// Created by plugin gengin of gnostic at {{.CreatedAt}}
// WARNING! All changes made in this file will be lost when rebuilding
#ifndef META_DATA_STRUCT_H_E72A284E_42B3_428B_B0DF_BC0EDAEF23B1
#define META_DATA_STRUCT_H_E72A284E_42B3_428B_B0DF_BC0EDAEF23B1

#include <qstring.h>
#include <qglobal.h>
#include <string.h>
{{range .NetDefinitions}}
// {{.Comment}}
struct {{.Name}} {
	{{.Name}}(): {{.MemberInits}} {}
	{{range .Fields}}{{.FieldType}}	{{.FieldName}};		// {{.FieldComment}}
	{{end}}
	XPACK(O({{.AllFields}}));
};
{{end}}
#endif
`

type NetDefinitionField struct {
	FieldType    string
	FieldName    string
	FieldComment string
}

type NetDefinition struct {
	Comment     string
	Name        string
	Fields      []*NetDefinitionField
	AllFields   string
	MemberInits string
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
	if isGoTypeArray(goType) {
		elemTypeConvFun := getTypeConvFun(goType[2:])
		if elemTypeConvFun == "" {
			elemTypeConvFun = "zlutils.Str2Str"
		}
		convFunc = elemTypeConvFun + "Slice"
	} else if goType == "string" {
	} else if goType == "int" {
		convFunc = "zlutils.Str2Int"
	} else if goType == "int64" {
		convFunc = "zlutils.Str2Int64"
	} else if goType == "float64" {
		convFunc = "zlutils.Str2Float64"
	} else if goType == "bool" {
		convFunc = "zlutils.String2Bool"
	} else if goType == "time.Time" {
		convFunc = "zlutils.Str2Time"
	}
	return convFunc
}

func openAPIType2Go(paramType, paramFormat string, item *openapiv2.PrimitivesItems) (goType string) {
	goType = "string"
	if paramType == "array" {
		goType = "[]" + openAPIType2Go(item.Type, item.Format, item.Items)
	} else if paramType == "string" {
		if paramFormat == "date-time" {
			goType = "time.Time"
		}
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

func isGoTypeArray(goType string) bool {
	return goType[:2] == "[]"
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

var standardFormat = map[string]string{
	"float":     "Y",
	"double":    "Y",
	"int32":     "Y",
	"int64":     "Y",
	"date":      "Y",
	"date-time": "Y",
	"password":  "Y",
	"byte":      "Y",
	"binary":    "Y",
}

func getValidate(isRequred, isArray bool, format string) string {
	var ls []string
	if isRequred {
		ls = append(ls, "required")
	}
	if format != "" {
		tmp := standardFormat[format]
		if tmp != "Y" {
			if tmp == "" {
				tmp = format
			}
			if !isRequred {
				ls = append(ls, "omitempty")
			}
			if isArray {
				ls = append(ls, "dive")
			}
			ls = append(ls, tmp)
		}
	}
	if len(ls) > 0 {
		return strings.Join(ls, ",")
	}
	return "-"
}

func getCType(value *openapiv2.Schema) (cTypeName string) {
	if value.XRef != "" {
		ls := strings.Split(value.XRef, "/")
		return convertStructName(ls[len(ls)-1])
	}
	cTypeName = "QString"
	switch value.Type.GetValue()[0] {
	case "integer":
		cTypeName = "quint32"
	case "number":
		cTypeName = "qreal"
	}
	return cTypeName
}

func convertStructName(name string) (outStr string) {
	// ls := strings.Split(name, ".")
	// return ls[len(ls)-1]
	outStr = strings.ReplaceAll(name, ".", "")
	outStr = strings.ToUpper(outStr[:1]) + outStr[1:]
	return outStr
}

func v2doc2Gin(doc *openapiv2.Document) (goSource, goTypeSrc, cTypeSrc string) {
	info := &GinTemplateInfo{
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	if tmp, err := template.New("caller").Parse(GIN_TEMPLATE); err == nil {
		for _, ext := range doc.VendorExtension {
			if strings.ToLower(ext.Name) == "x-wr-modname" {
				info.ModName = strings.Trim(ext.GetValue().GetYaml(), "\n")
				break
			}
		}
		if info.ModName == "" {
			info.ModName = "wrnetman"
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

			// paramList := make([]string, len(operation.Parameters))
			preParamList := make([]string, len(operation.Parameters))
			apiType := &APIType{
				TypeName: fmt.Sprintf("AT%s", methodInfo.BareMethodName),
			}
			methodInfo.TypeName = apiType.TypeName

			for index, param := range operation.Parameters {
				// paramList[index] = fmt.Sprintf("param%d", index+1)
				var paramName, paramType, paramFormat, getParamFuncName string
				var paramRequired bool
				var item *openapiv2.PrimitivesItems
				isHeader := false
				if nonBodyParam := param.GetParameter().GetNonBodyParameter(); nonBodyParam != nil {
					if subSchema := nonBodyParam.GetQueryParameterSubSchema(); subSchema != nil {
						paramName = subSchema.Name
						paramType = subSchema.Type
						paramFormat = subSchema.Format
						paramRequired = subSchema.Required
						getParamFuncName = "GetQuery"
						item = subSchema.Items
					} else if subSchema := nonBodyParam.GetFormDataParameterSubSchema(); subSchema != nil {
						paramName = subSchema.Name
						paramType = subSchema.Type
						paramFormat = subSchema.Format
						paramRequired = subSchema.Required
						getParamFuncName = "GetPostForm"
						item = subSchema.Items
					} else if subSchema := nonBodyParam.GetHeaderParameterSubSchema(); subSchema != nil {
						paramName = subSchema.Name
						paramType = subSchema.Type
						paramFormat = subSchema.Format
						paramRequired = subSchema.Required
						getParamFuncName = "GetHeader"
						item = subSchema.Items
						if strings.ToLower(paramName) == "token" {
							apiType.Fields = append(apiType.Fields, &FieldInfo{
								FieldName:   "TokenInfo",
								FieldType:   "interface{}",
								FieldRemark: "`json:\"-\"`",
							})
							methodInfo.CheckToken = checkTokenStatement
						}
						isHeader = true
					}
				} else if bodyParam := param.GetParameter().GetBodyParameter(); bodyParam != nil { // todo
					paramName = bodyParam.Name
					paramRequired = bodyParam.Required
					getParamFuncName = "GetQuery"
				} else {
					panic("parameter is not either body or nonbody")
				}
				paramGoType := openAPIType2Go(paramType, paramFormat, item)
				fieldInfo := &FieldInfo{
					FieldName:   UpperFirstLetter(paramName),
					FieldRemark: fmt.Sprintf("`json:\"%s,omitempty\" validate:\"%s\"`", paramName, getValidate(paramRequired, isGoTypeArray(paramGoType), paramFormat)),
				}
				if isHeader {
					fieldInfo.FieldType = paramGoType
					preParamList[index] = fmt.Sprintf(`
		param.%s = %s(c.%s("%s"))`, fieldInfo.FieldName, getTypeConvFun(paramGoType), getParamFuncName, paramName)
				} else if paramRequired || isGoTypeArray(paramGoType) {
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
			info.Imports = append(info.Imports, fmt.Sprintf(`"%s/netadapter/overhttp/%s"`, info.ModName, key))
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
	if tmp, err := template.New("ctype").Parse(GIN_DEFINITION4C_TEMPLATE); err == nil {
		for _, nSchema := range doc.Definitions.AdditionalProperties {
			if strings.Index(nSchema.Name, "Result") == -1 {
				nd := &NetDefinition{
					Name:    convertStructName(nSchema.Name),
					Comment: nSchema.GetValue().Description,
				}
				var allFields []string
				var memberInits []string
				for _, v := range nSchema.Value.Properties.AdditionalProperties {
					if value := v.GetValue(); value != nil {
						field := &NetDefinitionField{
							FieldName:    v.Name,
							FieldType:    getCType(value),
							FieldComment: value.Description,
						}
						if field.FieldType != "QString" {
							memberInits = append(memberInits, fmt.Sprintf("%s(0)", v.Name))
						}
						nd.Fields = append(nd.Fields, field)

						allFields = append(allFields, v.Name)
					}
				}
				nd.AllFields = strings.Join(allFields, ",")
				nd.MemberInits = strings.Join(memberInits, ",")
				info.NetDefinitions = append(info.NetDefinitions, nd)
			}
		}
		builder := &strings.Builder{}
		tmp.Execute(builder, info)
		cTypeSrc = builder.String()
	}
	return goSource, goTypeSrc, cTypeSrc
}
