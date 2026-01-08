package parsers

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"devops.aishu.cn/AISHUDevOps/DIP/_git/agent-operator-integration/server/infra/errors"
	"devops.aishu.cn/AISHUDevOps/DIP/_git/agent-operator-integration/server/interfaces"
	"devops.aishu.cn/AISHUDevOps/DIP/_git/agent-operator-integration/server/interfaces/model"
	"devops.aishu.cn/AISHUDevOps/DIP/_git/agent-operator-integration/server/utils"
	o11y "devops.aishu.cn/AISHUDevOps/DIP/_git/mdl-go-lib/observability"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-python/gpython/ast"
	"github.com/go-python/gpython/parser"
)

// pythonFunctionParser Python 函数解析器
type pythonFunctionParser struct {
	Logger    interfaces.Logger
	Validator interfaces.Validator
}

func (p *pythonFunctionParser) Type() interfaces.MetadataType {
	return interfaces.MetadataTypeFunc
}

func (p *pythonFunctionParser) validate(ctx context.Context, inputValue any) (input *interfaces.FunctionInput, err error) {
	input, ok := inputValue.(*interfaces.FunctionInput)
	if !ok {
		err = errors.DefaultHTTPError(ctx, http.StatusBadRequest, "input value is not *interfaces.FunctionInput")
		return
	}
	if input == nil {
		err = errors.DefaultHTTPError(ctx, http.StatusBadRequest, "input content is empty")
		return
	}
	// Code 校验
	if input.Code == "" {
		err = errors.DefaultHTTPError(ctx, http.StatusBadRequest, "python function code is empty")
		return
	}
	// 校验参数定义
	err = p.Validator.ValidatorStruct(ctx, input)
	if err != nil {
		return
	}
	return
}

// 检查是否包含入口函数handler
func checkRegexpHandler(ctx context.Context, code string) (err error) {
	// 使用正则表达式检查是否包含 handler 函数定义
	pattern := `def\s+handler\s*\(`
	matched, err := regexp.MatchString(pattern, code)
	if err != nil {
		return errors.DefaultHTTPError(ctx, http.StatusInternalServerError, fmt.Sprintf("check handler regexp failed: %v", err))
	}
	if !matched {
		// 必须包含 handler 函数定义
		return errors.NewHTTPError(ctx, http.StatusBadRequest, errors.ErrExtFunctionNoHandlerFound, "python function must have a handler(event) function")
	}
	return nil
}

// func checAstkHandler(ctx context.Context, code string) (err error) {
// 	// 解析Python代码
// 	mod, err := parser.ParseString(code, py.ExecMode)
// 	if err != nil {
// 		err = errors.DefaultHTTPError(ctx, http.StatusBadRequest, fmt.Sprintf("parse python code failed: %v", err))
// 		return
// 	}
// 	// 检查是否包含入口函数handler
// 	var hasHandler bool
// 	ast.Walk(mod, func(node ast.Ast) bool {
// 		n, ok := node.(*ast.FunctionDef)
// 		if ok && n.Name == "handler" {
// 			hasHandler = true
// 		}
// 		return true
// 	})
// 	if !hasHandler {
// 		err = errors.DefaultHTTPError(ctx, http.StatusBadRequest, "python function must have a handler function")
// 	}
// 	return
// }

// Parse 解析 Python 函数
func (p *pythonFunctionParser) Parse(ctx context.Context, inputValue any) (metadatas []interfaces.Metadata, err error) {
	// 记录可观测性
	ctx, _ = o11y.StartInternalSpan(ctx)
	defer o11y.EndSpan(ctx, err)
	input, err := p.validate(ctx, inputValue)
	if err != nil {
		return nil, err
	}
	err = checkRegexpHandler(ctx, input.Code)
	if err != nil {
		return nil, err
	}
	pathItem := convertToPathItemContent(input)
	desc := pathItem.Description
	if desc == "" {
		desc = pathItem.Summary
	}
	metadatas = make([]interfaces.Metadata, 0)
	metadataDB := &model.FunctionMetadataDB{
		ScriptType:   string(input.ScriptType),
		Code:         input.Code,
		Dependencies: utils.ObjectToJSON(input.Dependencies),
		Summary:      pathItem.Summary,
		Description:  desc,
		Path:         pathItem.Path,
		ServerURL:    pathItem.ServerURL,
		Method:       pathItem.Method,
		APISpec:      pathItem.APISpec.ToJSON(),
	}
	metadatas = append(metadatas, metadataDB)
	return
}

// GetAllContent 获取所有内容
func (p *pythonFunctionParser) GetAllContent(ctx context.Context, inputValue any) (content any, err error) {
	input, err := p.validate(ctx, inputValue)
	if err != nil {
		return nil, err
	}
	// 解析Python代码
	mod, err := parser.ParseString(input.Code, "exec")
	if err != nil {
		err = errors.DefaultHTTPError(ctx, http.StatusBadRequest, fmt.Sprintf("解析Python代码失败: %v", err))
		return
	} // 检查是否包含入口函数handler
	var hasHandler bool
	ast.Walk(mod, func(node ast.Ast) bool {
		n, ok := node.(*ast.FunctionDef)
		if ok && n.Name == "handler" {
			hasHandler = true
		}
		return true
	})
	if !hasHandler {
		err = errors.DefaultHTTPError(ctx, http.StatusBadRequest, "python function must have a handler function")
		return
	}
	content = convertToPathItemContent(input)
	return
}

// 将input\output转换成 PathItemContent
func convertToPathItemContent(input *interfaces.FunctionInput) (result *interfaces.PathItemContent) {
	result = &interfaces.PathItemContent{
		Path:        interfaces.AOIFuncExecPath,
		Method:      http.MethodPost,
		ServerURL:   interfaces.AOIServerURL,
		Summary:     input.Name,
		Description: input.Description,
		APISpec:     &interfaces.APISpec{},
	}
	// 添加超时时间参数
	result.APISpec.Parameters = createParameter()
	// 根据处理输入参数创建请求体
	result.APISpec.RequestBody = createRequestBody(input.Inputs)
	// 处理输出参数
	result.APISpec.Responses = createResponseBody(input.Outputs)
	return
}

/*
// 构造请求体
type FunctionExecuteProxyReq struct {
	Version string         `uri:"version" validate:"required,uuid4"`
	Timeout int64          `query:"timeout"` // 毫秒
	Event   map[string]any `json:"event"`
}
*/

// 构造Parameter参数
func createParameter() []*interfaces.Parameter {
	parameters := make([]*interfaces.Parameter, 0)
	// 超时时间
	timeoutParam := &interfaces.Parameter{
		Name:        "timeout",
		In:          "query",
		Description: "函数执行超时时间，单位毫秒",
		Required:    false,
		Schema: openapi3.NewSchemaRef("", &openapi3.Schema{
			Type:        &openapi3.Types{openapi3.TypeNumber},
			Description: "函数执行超时时间，单位毫秒",
		}),
	}
	parameters = append(parameters, timeoutParam)
	return parameters
}

// 构造请求体结构
func createRequestBody(inputs []interfaces.ParameterDef) *interfaces.RequestBody {
	// 创建schema定义
	requestBodySchema := openapi3.NewObjectSchema()
	if len(inputs) > 0 {
		for _, input := range inputs {
			if input.Description == "" {
				input.Description = input.Name
			}
			propertySchema := &openapi3.Schema{
				Type:        mapTypeToOpenAPI(input.Type),
				Description: input.Description,
			}
			// 设置默认值
			if input.Default != nil {
				propertySchema.Default = input.Default
			}
			// 设置枚举值
			if len(input.Enum) > 0 {
				propertySchema.Enum = input.Enum
			}
			// 设置示例值
			if input.Example != nil {
				propertySchema.Example = input.Example
			}
			// 如果是数组类型，设置items
			if input.Type == "array" {
				propertySchema.Items = openapi3.NewSchemaRef("", &openapi3.Schema{
					Type: &openapi3.Types{openapi3.TypeString},
				})
			}
			requestBodySchema.Properties[input.Name] = openapi3.NewSchemaRef("", propertySchema)
			// 设置必填字段
			if input.Required {
				requestBodySchema.Required = append(requestBodySchema.Required, input.Name)
			}
		}
	}
	// 创建请求体
	requestBody := &interfaces.RequestBody{
		Description: "函数输入参数",
		Content:     openapi3.NewContentWithJSONSchema(requestBodySchema),
		Required:    true,
	}
	return requestBody
}

// 处理输出参数
// 根据处理输出参数创建响应体
func createResponseBody(outputs []interfaces.ParameterDef) []*interfaces.Response {
	// 创建schema定义
	responseSchema := openapi3.NewObjectSchema()
	responseSchema.Properties["stdout"] = openapi3.NewSchemaRef("", &openapi3.Schema{
		Type:        &openapi3.Types{openapi3.TypeString},
		Description: "标准输出流内容",
	})
	responseSchema.Properties["stderr"] = openapi3.NewSchemaRef("", &openapi3.Schema{
		Type:        &openapi3.Types{openapi3.TypeString},
		Description: "标准错误流内容",
	})

	resultSchema := &openapi3.Schema{
		Type:        &openapi3.Types{openapi3.TypeObject},
		Description: "Handler 函数返回的业务结果: any or null",
		Properties:  make(openapi3.Schemas),
	}
	for _, output := range outputs {
		propertySchema := &openapi3.Schema{
			Type:        mapTypeToOpenAPI(output.Type),
			Description: output.Description,
		}
		// 设置默认值
		if output.Default != nil {
			propertySchema.Default = output.Default
		}
		// 设置枚举值
		if len(output.Enum) > 0 {
			propertySchema.Enum = output.Enum
		}
		// 设置示例值
		if output.Example != nil {
			propertySchema.Example = output.Example
		}
		// 如果是数组类型，设置items
		if output.Type == "array" {
			propertySchema.Items = openapi3.NewSchemaRef("", &openapi3.Schema{
				Type: &openapi3.Types{openapi3.TypeString},
			})
		}
		resultSchema.Properties[output.Name] = openapi3.NewSchemaRef("", propertySchema)
		// 设置必填字段
		if output.Required {
			resultSchema.Required = append(resultSchema.Required, output.Name)
		}
	}
	responseSchema.Properties["result"] = openapi3.NewSchemaRef("", resultSchema)
	// 添加指标
	metricsSchema := &openapi3.Schema{
		Type:        &openapi3.Types{openapi3.TypeObject},
		Description: "指标",
		Properties:  make(openapi3.Schemas),
	}
	metricsSchema.Properties["duration_ms"] = openapi3.NewSchemaRef("", &openapi3.Schema{
		Type:        &openapi3.Types{openapi3.TypeNumber},
		Description: "执行总耗时（毫秒）",
	})
	metricsSchema.Properties["memory_peak_mb"] = openapi3.NewSchemaRef("", &openapi3.Schema{
		Type:        &openapi3.Types{openapi3.TypeNumber},
		Description: "峰值内存占用（MB）",
	})
	metricsSchema.Properties["cpu_time_ms"] = openapi3.NewSchemaRef("", &openapi3.Schema{
		Type:        &openapi3.Types{openapi3.TypeNumber},
		Description: "CPU 时间（毫秒）",
	})
	responseSchema.Properties["metrics"] = openapi3.NewSchemaRef("", metricsSchema)
	// 添加错误响应体
	errSchema := &openapi3.Schema{
		Type:        &openapi3.Types{openapi3.TypeObject},
		Description: "失败详情",
		Properties:  map[string]*openapi3.SchemaRef{},
	}
	errSchema.Properties["code"] = openapi3.NewSchemaRef("", &openapi3.Schema{
		Type:        &openapi3.Types{openapi3.TypeString},
		Description: "错误码",
	})
	errSchema.Properties["description"] = openapi3.NewSchemaRef("", &openapi3.Schema{
		Type:        &openapi3.Types{openapi3.TypeString},
		Description: "错误描述",
	})
	errSchema.Properties["detail"] = openapi3.NewSchemaRef("", &openapi3.Schema{
		Type:        &openapi3.Types{openapi3.TypeObject},
		Description: "错误详情",
	})
	errSchema.Properties["solution"] = openapi3.NewSchemaRef("", &openapi3.Schema{
		Type:        &openapi3.Types{openapi3.TypeString},
		Description: "错误解决办法",
	})
	errSchema.Properties["link"] = openapi3.NewSchemaRef("", &openapi3.Schema{
		Type:        &openapi3.Types{openapi3.TypeString},
		Description: "错误链接",
	})
	// 创建响应体
	responseBody := []*interfaces.Response{
		{
			StatusCode:  "200",
			Description: "成功",
			Content:     openapi3.NewContentWithJSONSchema(responseSchema),
		},
		{
			StatusCode:  "400",
			Description: "参数校验失败",
			Content:     openapi3.NewContentWithJSONSchema(errSchema),
		},
		{
			StatusCode:  "404",
			Description: "资源不存在",
			Content:     openapi3.NewContentWithJSONSchema(errSchema),
		},
		{
			StatusCode:  "500",
			Description: "函数执行失败",
			Content:     openapi3.NewContentWithJSONSchema(errSchema),
		},
	}
	return responseBody
}

// mapTypeToOpenAPI 将参数类型映射到OpenAPI类型
func mapTypeToOpenAPI(paramType string) *openapi3.Types {
	switch strings.ToLower(paramType) {
	case "string":
		return &openapi3.Types{openapi3.TypeString}
	case "int", "integer", "number":
		return &openapi3.Types{openapi3.TypeNumber}
	case "float", "double":
		return &openapi3.Types{openapi3.TypeNumber}
	case "bool", "boolean":
		return &openapi3.Types{openapi3.TypeBoolean}
	case "array":
		return &openapi3.Types{openapi3.TypeArray}
	case "object":
		return &openapi3.Types{openapi3.TypeObject}
	default:
		return &openapi3.Types{openapi3.TypeString}
	}
}

// createRequestBody 创建请求体结构
// func (fi *FunctionInput) createRequestBody() *RequestBody {
// 	if len(fi.Inputs) == 0 {
// 		return nil
// 	}

// 	// 创建schema定义
// 	schema := &openapi3.Schema{
// 		Type:       "object",
// 		Properties: make(map[string]*openapi3.SchemaRef),
// 		Required:   []string{},
// 	}

// 	// 将Inputs转换为schema属性
// 	for _, input := range fi.Inputs {
// 		propertySchema := &openapi3.Schema{
// 			Type:        fi.mapTypeToOpenAPI(input.Type),
// 			Description: input.Description,
// 		}

// 		// 设置默认值
// 		if input.Default != nil {
// 			propertySchema.Default = input.Default
// 		}

// 		// 设置枚举值
// 		if len(input.Enum) > 0 {
// 			propertySchema.Enum = make([]interface{}, len(input.Enum))
// 			for i, enumVal := range input.Enum {
// 				propertySchema.Enum[i] = enumVal
// 			}
// 		}

// 		// 设置示例值
// 		if input.Example != nil {
// 			propertySchema.Example = input.Example
// 		}

// 		schema.Properties[input.Name] = &openapi3.SchemaRef{Value: propertySchema}

// 		// 如果是必填参数，添加到required列表
// 		if input.Required {
// 			schema.Required = append(schema.Required, input.Name)
// 		}
// 	}

// 	// 创建内容
// 	content := openapi3.NewContent()
// 	content["application/json"] = &openapi3.MediaType{
// 		Schema: &openapi3.SchemaRef{Value: schema},
// 	}

// 	return &interfaces.RequestBody{
// 		Description: "函数输入参数",
// 		Content:     content,
// 		Required:    true,
// 	}
// }
