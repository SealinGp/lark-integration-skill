---
name: lark-integration
description: Access and manage Lark/Feishu documents and blocks via REST API.
---
The Lark Integration service is running locally at http://localhost:8000.

This skill provides access to the following capabilities:

1.  **Task Management**:
    -   Create Task: `POST /api/v1/tasks`
    -   Get Task: `GET /api/v1/tasks/:task_id`
    -   Delete Task: `DELETE /api/v1/tasks/:task_id`

2.  **Document Management (Docx)**:
    -   Create Document: `POST /api/v1/docs`
    -   Get Document Info: `GET /api/v1/docs/:doc_token`
    -   Get Raw Content: `GET /api/v1/docs/:doc_token/raw`
    -   Get Blocks: `GET /api/v1/docs/:doc_token/blocks`

3.  **Wiki Management**:
    -   Create Node: `POST /api/v1/wiki`
    -   Search Nodes: `POST /api/v1/wiki/search`
    -   Get Node Info: `GET /api/v1/wiki/nodes/:node_token`
    -   List Nodes: `GET /api/v1/wiki/spaces/:space_id/nodes`
    -   Move Node: `POST /api/v1/wiki/spaces/:space_id/nodes/:node_token/move`
    -   Update Title: `POST /api/v1/wiki/spaces/:space_id/nodes/:node_token/update_title`
    -   Move Docs to Wiki: `POST /api/v1/wiki/spaces/:space_id/nodes/move_docs_to_wiki`

4.  **Docx Block Operations**:
    -   Get Block: `GET /api/v1/docx/v1/documents/:document_id/blocks/:block_id`
    -   Get Children: `GET /api/v1/docx/v1/documents/:document_id/blocks/:block_id/children`
    -   Create Children: `POST /api/v1/docx/v1/documents/:document_id/blocks/:block_id/children`
    -   Update Block: `PATCH /api/v1/docx/v1/documents/:document_id/blocks/:block_id`
    -   Delete Children: `DELETE /api/v1/docx/v1/documents/:document_id/blocks/:block_id/children/batch_delete`
    -   Convert Content: `POST /api/v1/docx/v1/documents/blocks/convert`

Refer to `docs/openapi.yaml` or `README.md` for payload details.
