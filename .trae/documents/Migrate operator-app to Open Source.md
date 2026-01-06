我将按以下步骤对 `operator-app` 进行开源迁移：

1.  **修改 `go.mod`**:
    *   将 module path 从 `devops.aishu.cn/AISHUDevOps/DIP/_git/agent-operator-app` 更新为 `github.com/kweaver-ai/operator-hub/operator-app`。
    *   将内网依赖替换为 GitHub 对应的版本：
        *   `devops.aishu.cn/AISHUDevOps/DIP/_git/operator-go-lib` -> `github.com/kweaver-ai/kweaver-go-lib v1.0.2`
        *   `devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/v2` -> `github.com/kweaver-ai/TelemetrySDK-Go/span/v2 v2.10.3`
        *   `devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/proton-rds-sdk-go` -> `github.com/kweaver-ai/proton-rds-sdk-go v1.4.2`
    *   移除任何内部的 `replace` 指令（如果有）。

2.  **更新导入路径**:
    *   在所有 `.go` 文件中，全局将 `devops.aishu.cn/AISHUDevOps/DIP/_git/agent-operator-app` 替换为 `github.com/kweaver-ai/operator-hub/operator-app`。
    *   在所有 `.go` 文件中，全局将 `devops.aishu.cn/AISHUDevOps/DIP/_git/operator-go-lib` 替换为 `github.com/kweaver-ai/kweaver-go-lib`。
    *   在所有 `.go` 文件中，全局将 `devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git` 替换为 `github.com/kweaver-ai/TelemetrySDK-Go`。
    *   在所有 `.go` 文件中，全局将 `devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/proton-rds-sdk-go` 替换为 `github.com/kweaver-ai/proton-rds-sdk-go`。

3.  **验证**:
    *   运行 `go mod tidy` 清理依赖。
    *   运行 `go build ./...` 确保项目编译成功。
