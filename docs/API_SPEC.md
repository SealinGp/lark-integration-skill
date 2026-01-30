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
- `POST /wiki/search`
  - Search for wiki nodes.
  - Body: `WikiSearchRequest` (Query)
- `GET /wiki/nodes/:node_token`
  - Get wiki node information.
- `GET /wiki/spaces/:space_id/nodes`
  - List nodes in a wiki space.
- `POST /wiki/spaces/:space_id/nodes/:node_token/move`
  - Move a wiki node.
- `POST /wiki/spaces/:space_id/nodes/:node_token/update_title`
  - Update a wiki node's title.
- `POST /wiki/spaces/:space_id/nodes/move_docs_to_wiki`
  - Move an existing Doc/Docx to Wiki.

## Docx (V1)
- `GET /docx/v1/documents/:document_id/blocks/:block_id`
  - Get a specific block.
- `GET /docx/v1/documents/:document_id/blocks/:block_id/children`
  - Get children blocks of a specific block.
- `POST /docx/v1/documents/:document_id/blocks/:block_id/children`
  - Create children blocks.
  - Body: `CreateDocBlockRequest` (Children)
- `PATCH /docx/v1/documents/:document_id/blocks/:block_id`
  - Update a specific block.
  - Body: `UpdateDocBlockRequest`
- `DELETE /docx/v1/documents/:document_id/blocks/:block_id/children/batch_delete`
  - Batch delete children blocks.
  - Body: `DeleteDocBlockChildrenRequest` (StartIndex, EndIndex)
- `POST /docx/v1/documents/blocks/convert`
  - Convert Markdown/HTML content to blocks.
  - Body: `ConvertContentToBlocksRequest` (Content, ContentType)
