# Lark Integration Skill

A lightweight backend service (Skill) designed for **clawdbot**, enabling integration with Lark (Feishu).

This service provides a RESTful API to manage Tasks, Documents, and Knowledge Base nodes in Lark.

## Features

- **Tasks**: Create, Retrieve, Delete tasks.
- **Docs**: Create new Documents (Docx), Retrieve document info, raw content, and blocks.
- **Wiki**: Create nodes, Search nodes, Move nodes, Move Docs to Wiki, Update node titles.
- **Docx**: Detailed block management (Get, Create, Update, Delete Children, Convert).

## Prerequisites

1.  **Lark App**: Create a custom app in [Lark Developer Console](https://open.larksuite.com/).
2.  **Permissions**: Ensure the app has the following permissions:
    -   `task:task` (Manage tasks)
    -   `drive:drive` (View/Edit files)
    -   `wiki:wiki` (Manage Knowledge Base)
    -   `contact:user.id:readonly` (To resolve OpenID)

## Setup & Deployment

### 1. Configuration

Copy `.env.example` to `.env` and fill in your credentials:

```bash
cp .env.example .env
```

Edit `.env`:
```env
LARK_APP_ID="your_app_id"
LARK_APP_SECRET="your_app_secret"
```

### 2. Run with Docker (Recommended)

```bash
# Build Image
docker build -t lark-integration-skill .

# Run Container
docker run -d -p 8000:8000 --env-file .env --name lark-skill lark-integration-skill
```

### 3. Run Locally

```bash
go run cmd/server/main.go
```

## API Usage

### Create Task
**POST** `/api/v1/tasks`
```json
{
  "summary": "Fix critical bug",
  "description": "Check logs and fix NPE",
  "due_time": 1678888888,
  "user_id": "ou_xxxxxx" 
}
```

### Create Document
**POST** `/api/v1/docs`
```json
{
  "title": "Project Meeting Notes",
  "folder_token": "" 
}
```

### Create Wiki Node
**POST** `/api/v1/wiki`
```json
{
  "space_id": "698888888888",
  "title": "New Knowledge Page",
  "obj_type": "docx"
}
```

### Docx Block Operations
**GET** `/api/v1/docx/v1/documents/:document_id/blocks/:block_id`
**GET** `/api/v1/docx/v1/documents/:document_id/blocks/:block_id/children`
**POST** `/api/v1/docx/v1/documents/:document_id/blocks/:block_id/children`
**PATCH** `/api/v1/docx/v1/documents/:document_id/blocks/:block_id`
**DELETE** `/api/v1/docx/v1/documents/:document_id/blocks/:block_id/children/batch_delete`
**POST** `/api/v1/docx/v1/documents/blocks/convert`

See `docs/API_SPEC.md` for full API details.
