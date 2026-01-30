# Development Notes & SDK Patterns

During the development of this skill using the `larksuite/oapi-sdk-go/v3` SDK, several patterns were identified and documented to ensure compilation and correct functionality.

## SDK Request Patterns
The SDK uses a `ReqBuilder` pattern. The common format is:
`larkservice.New[Action][Resource]ReqBuilder().[Field](value).Build()`

### 1. Tasks (V1)
- **Create**: `larktask.NewCreateTaskReqBuilder().Task(taskBody).UserIdType("open_id").Build()`
- **Note**: Body is passed via `.Task()` method, not `.Body()`.
- **Fields**: Uses `Id` and `Summary`.

### 2. Documents (Docx V1)
- **Create**: `larkdocx.NewCreateDocumentReqBuilder().Body(body).Build()`
- **Get Info**: `larkdocx.NewGetDocumentReqBuilder().DocumentId(id).Build()`
- **Note**: `docx.v1.Get` returns limited metadata. Use `drive.v1.Meta` for rich metadata.

### 3. Drive Meta (V1)
- **Batch Query**: `larkdrive.NewBatchQueryMetaReqBuilder().MetaRequest(metaReq).Build()`
- **Note**: Request body is passed via `.MetaRequest()`, not `.Body()`.
- **Field Mappings**: 
  - `CreateTime` is the creation timestamp string.
  - `LatestModifyTime` is used for UpdateTime.

### 4. Wiki (V2)
- **Create Node**: `larkwiki.NewCreateSpaceNodeReqBuilder().SpaceId(id).Node(nodeBody).Build()`
- **Note**: Body is passed via `.Node()`.

## Docker Build Notes
- **Go Version**: `go.mod` must match the version in `Dockerfile` (currently using 1.24).
- **Build Command**: Use `--load` if using `docker-container` driver to ensure the image is available to the local daemon.
  - `docker build -t lark-integration-skill --load .`
