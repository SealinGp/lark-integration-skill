# API Specification

All endpoints are prefixed with `/api/v1`.

## Health Check
- `GET /health`
  - Returns `{"status": "ok"}`

## Tasks
- `POST /tasks`
  - Create a new task.
  - Body: `CreateTaskRequest` (Summary, Description, DueTime, UserID)
- `GET /tasks/:task_id`
  - Retrieve task summary.
- `DELETE /tasks/:task_id`
  - Delete a task.

## Documents (Docx)
- `POST /docs`
  - Create a new Docx file.
  - Body: `CreateDocRequest` (Title, FolderToken)
- `GET /docs/:doc_token`
  - Get document metadata (Title, CreateTime, UpdateTime, OwnerID).
  - Uses Drive Meta API.
- `GET /docs/:doc_token/raw`
  - Get raw text content of the document.
- `GET /docs/:doc_token/blocks`
  - List all blocks in the document.
  - Query Params: `page_token`, `page_size`.

## Wiki
- `POST /wiki`
  - Create a new node in a Wiki space.
  - Body: `CreateWikiNodeRequest` (SpaceID, ParentNodeToken, Title, ObjType)
