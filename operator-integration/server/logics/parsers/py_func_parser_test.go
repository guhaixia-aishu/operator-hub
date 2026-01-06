package parsers

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/kweaver-ai/operator-hub/operator-integration/server/interfaces"
	jsoniter "github.com/json-iterator/go"
	. "github.com/smartystreets/goconvey/convey"
)

// 测试Function转换成OpenAPI3.0规范的Schema

func TestFunctionToOpenAPISchema(t *testing.T) {
	Convey("TestFunctionToOpenAPISchema: 传参为空时", t, func() {
		input := &interfaces.FunctionInput{
			Name:        "传参为空时",
			Description: "test function",
			Inputs:      []interfaces.ParameterDef{},
			Outputs:     []interfaces.ParameterDef{},
		}
		schema := convertToPathItemContent(input)
		So(schema, ShouldNotBeNil)
		// 输出到文件中
		data, _ := jsoniter.Marshal(schema)
		filename := input.Name + ".json"
		err := os.WriteFile(filename, data, 0644)
		So(err, ShouldBeNil)
		t.Logf("Successfully wrote API metadata to %s", filename)
		_ = os.Remove(filename)
	})
	Convey("TestFunctionToOpenAPISchema: 传参为nil时", t, func() {
		input := &interfaces.FunctionInput{
			Name:        "传参为nil时",
			Description: "test function",
			Inputs: []interfaces.ParameterDef{
				{
					Name:        "id",
					Description: "任务ID",
					Type:        "string",
					Required:    true,
				},
				{
					Name:        "params",
					Description: "任务参数",
					Type:        "object",
					Required:    true,
				},
			},
			Outputs: []interfaces.ParameterDef{
				{
					Name:        "status",
					Description: "任务执行状态",
					Type:        "string",
					Required:    true,
				},
				{
					Name:        "output",
					Description: "任务执行结果",
					Type:        "object",
					Required:    true,
				},
			},
		}
		schema := convertToPathItemContent(input)
		So(schema, ShouldNotBeNil)
		// 输出到文件中
		data, _ := jsoniter.Marshal(schema)
		filename := input.Name + ".json"
		err := os.WriteFile(filename, data, 0644)
		So(err, ShouldBeNil)
		t.Logf("Successfully wrote API metadata to %s", filename)
		_ = os.Remove(filename)
	})
}

func TestCheckHandler(t *testing.T) {
	Convey("TestCheckHandler: 检查是否包含入口函数handler", t, func() {
		code := `def handler(event):
    return {"status": "success"}`
		err := checkRegexpHandler(context.Background(), code)
		So(err, ShouldBeNil)
		code2 := `def handler(event, context):
    return {"status": "success"}`
		err = checkRegexpHandler(context.Background(), code2)
		So(err, ShouldBeNil)
		code3 := `def main(event, context):
    return {"status": "success"}`
		err = checkRegexpHandler(context.Background(), code3)
		So(err, ShouldNotBeNil)
		fmt.Println(err)
		code4 := `def handler(text,event):
    return {"status": "success"}`
		err = checkRegexpHandler(context.Background(), code4)
		So(err, ShouldNotBeNil)
	})
}
