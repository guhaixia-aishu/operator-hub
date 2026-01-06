package interfaces

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/kweaver-ai/operator-hub/operator-integration/server/interfaces/model"
)

// Metadata 元数据通用接口
type Metadata interface {
	GetType() string
	GetSummary() string
	SetSummary(summary string)
	GetDescription() string
	SetDescription(description string)
	GetVersion() string
	SetVersion(version string)
	GetScriptType() string
	SetScriptType(scriptType string)
	GetServerURL() string
	SetServerURL(serverURL string)
	GetAPISpec() string
	SetAPISpec(apiSpec string)
	GetMethod() string
	SetMethod(method string)
	GetPath() string
	SetPath(path string)
	Validate(ctx context.Context) error
	GetUpdateUser() (user string)
	SetUpdateInfo(user string)
	GetCreateUser() (user string)
	SetCreateInfo(user string)
	// UpdataMetadata(metadata interface{}) error
	// 获取ErrMessage信息
	GetErrMessage() string
	GetFunctionContent() (code, scriptType, dependencies string)
	SetFunctionContent(code, scriptType, dependencies string)
}

type ScriptType string

const (
	ScriptTypePython ScriptType = "python" // Python 脚本类型
)

// FunctionContent 函数内容定义
type FunctionContent struct {
	ScriptType   ScriptType `json:"script_type" form:"script_type" default:"python" validate:"required,oneof=python"` // 脚本类型
	Code         string     `json:"code" form:"code" validate:"required"`                                             // Python 代码（必填）
	Dependencies []string   `json:"dependencies,omitempty" form:"dependencies"`                                       // 依赖库列表
}

// ParameterDef 参数定义
type ParameterDef struct {
	Name        string `json:"name"`                                                                               // 参数名称
	Description string `json:"description,omitempty"`                                                              // 参数描述
	Type        string `json:"type" default:"string" validate:"required,oneof=string number boolean array object"` // 参数类型
	Required    bool   `json:"required"`                                                                           // 是否必填
	Default     any    `json:"default,omitempty"`                                                                  // 默认值
	Enum        []any  `json:"enum,omitempty"`                                                                     // 枚举值（可选）
	Example     any    `json:"example,omitempty"`                                                                  // 示例值
}

// FunctionInput  函数输入定义
type FunctionInput struct {
	// 基础信息
	Name        string `json:"name" form:"name"`                         // 函数名称
	Description string `json:"description,omitempty" form:"description"` // 函数描述，用于说明函数的功能和行为
	// 参数定义
	Inputs  []ParameterDef `json:"inputs,omitempty" form:"inputs"`   // 输入参数列表
	Outputs []ParameterDef `json:"outputs,omitempty" form:"outputs"` // 输出参数列表
	// 代码相关
	ScriptType   ScriptType `json:"script_type" form:"script_type" default:"python" validate:"required,oneof=python"` // 脚本类型
	Code         string     `json:"code" form:"code"`                                                                 // Python 代码（必填）
	Dependencies []string   `json:"dependencies,omitempty" form:"dependencies"`                                       // 依赖库列表
}

// FunctionInputEdit 函数输入编辑定义
type FunctionInputEdit struct {
	// 参数定义
	Inputs  []ParameterDef `json:"inputs,omitempty" form:"inputs"`   // 输入参数列表
	Outputs []ParameterDef `json:"outputs,omitempty" form:"outputs"` // 输出参数列表
	// 代码相关
	ScriptType   ScriptType `json:"script_type" form:"script_type" default:"python" validate:"required,oneof=python"` // 脚本类型
	Code         string     `json:"code" form:"code"`                                                                 // Python 代码（必填）
	Dependencies []string   `json:"dependencies,omitempty" form:"dependencies"`                                       // 依赖库列表
}

// OpenAPIInput OpenAPI 输入定义
type OpenAPIInput struct {
	// 基础信息
	Data json.RawMessage `json:"data" form:"data"` // 原始内容（OpenAPI JSON/YAML）
}

// IMetadataService 统一元数据管理接口
type IMetadataService interface {
	// 注册元数据
	RegisterMetadata(ctx context.Context, tx *sql.Tx, metadata Metadata) (version string, err error)
	// 批量注册元数据
	BatchRegisterMetadata(ctx context.Context, tx *sql.Tx, metadatas []Metadata) (versions []string, err error)
	// 根据版本查询元数据
	GetMetadataByVersion(ctx context.Context, metadataType MetadataType, version string) (Metadata, error)
	// 批量查询元数据
	BatchGetMetadata(ctx context.Context, apiVersions, funcVersions []string) ([]Metadata, error)
	// 更新元数据
	UpdateMetadata(ctx context.Context, tx *sql.Tx, metadata Metadata) error
	// 删除元数据
	DeleteMetadata(ctx context.Context, tx *sql.Tx, metadataType MetadataType, version string) error
	// 批量删除元数据
	BatchDeleteMetadata(ctx context.Context, tx *sql.Tx, metadataType MetadataType, versions []string) error
	// 验证元数据格式
	ValidateMetadata(ctx context.Context, metadata Metadata) error
	// 元数据解析
	ParseMetadata(ctx context.Context, metadataType MetadataType, input any) ([]Metadata, error)
	// 获取解析后的原始内容
	ParseRawContent(ctx context.Context, metadataType MetadataType, input any) (content any, err error)
	// 根据SourceID、SourceType查询元数据
	GetMetadataBySource(ctx context.Context, sourceID string, sourceType model.SourceType) (bool, Metadata, error)
	// 批量根据SourceID、SourceType查询元数据
	BatchGetMetadataBySourceIDs(ctx context.Context, sourceMap map[model.SourceType][]string) (sourceIDToMetadataMap map[string]Metadata, err error)
	// 检查并返回元数据是否存在
	CheckMetadataExists(ctx context.Context, metadataType MetadataType, version string) (bool, Metadata, error)
}

// // MetadataParser 元数据解析器接口
// type MetadataParser interface {
// 	// 支持的元数据类型
// 	SupportedType() MetadataType
// 	// 解析原始数据为目标结构
// 	Parse(ctx context.Context, rawData []byte) ([]Metadata, error)
// 	// 验证原始数据格式
// 	ValidateRawData(ctx context.Context, rawData []byte) error
// 	// 转换为统一元数据格式
// 	ToUnifiedFormat(metadata Metadata) (map[string]interface{}, error)
// }
