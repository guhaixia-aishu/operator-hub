I will perform the open-source migration for `operator-app` following the reference scheme.

### 1. Copy `ormhelper` Library
Copy the `ormhelper` library from `operator-integration` to `operator-app` to replace the dependency on the internal `operator-go-lib`.
- Source: `operator-integration/server/infra/common/ormhelper`
- Destination: `operator-app/server/infra/common/ormhelper`

### 2. Modify `go.mod`
Update the module definition and dependencies in `operator-app/go.mod`:
- **Module Path**: Change `devops.aishu.cn/AISHUDevOps/DIP/_git/agent-operator-app` to `github.com/kweaver-ai/operator-hub/operator-app`.
- **Dependencies**:
    - Remove `devops.aishu.cn/AISHUDevOps/DIP/_git/operator-go-lib`.
    - Replace `TelemetrySDK-Go` with `github.com/kweaver-ai/TelemetrySDK-Go/span/v2` (v2.10.0).
    - Replace `proton-rds-sdk-go` with `github.com/kweaver-ai/proton-rds-sdk-go` (v1.4.2).

### 3. Code Reference Correction
Perform global text replacement in the `operator-app` directory to update import paths:
- Replace internal module references with `github.com/kweaver-ai/operator-hub/operator-app`.
- Replace `operator-go-lib/ormhelper` imports with the local path `github.com/kweaver-ai/operator-hub/operator-app/server/infra/common/ormhelper`.
- Replace `TelemetrySDK-Go` and `proton-rds-sdk-go` imports with their GitHub counterparts.

### 4. Verification
- Run `go mod tidy` to clean up dependencies.
- (Optional) Attempt to build or verify `go.mod` integrity.
