# 飞书集成 Skill (Lark Integration Skill)

[English Documentation](README.md)

这是一个专为 **clawdbot** 设计的轻量级后端服务 (Skill)，旨在实现与飞书 (Lark) 的集成。

该服务提供了一套 RESTful API，用于管理飞书中的任务 (Tasks)、文档 (Docs) 和知识库 (Wiki) 节点。

## 功能特性

- **任务 (Tasks)**: 创建、查询、删除任务。
- **文档 (Docs)**: 创建新的多维文档 (Docx)，获取文档信息、原始内容及文档块 (Blocks)。
- **知识库 (Wiki)**: 创建节点、搜索节点、移动节点、移动文档到知识库、更新节点标题。
- **多维文档 (Docx)**: 详细的块管理 (获取、创建、更新、删除子块、内容转换)。

## 前置要求

1.  **飞书应用**: 在 [飞书开放平台](https://open.feishu.cn/) 创建一个企业自建应用。
2.  **权限配置**: 确保应用拥有以下权限：
    -   `task:task` (任务管理)
    -   `drive:drive` (查看、评论、编辑和管理云空间所有文件)
    -   `wiki:wiki` (知识库管理)
    -   `contact:user.id:readonly` (通过手机号或邮箱获取用户 ID)
    -   `docx:document:read_write` (编辑多维文档)
    -   `docx:document:read` (阅读多维文档)

## 与 Clawdbot (OpenClaw) 集成

本服务旨在作为 "Skill" 或 "Tool" 供运行在 **Clawdbot** 上的 AI Agent 使用。

### 方法 1: OpenAPI 导入 (推荐)

Clawdbot 支持通过 OpenAPI 规范导入工具。

1.  部署本服务 (例如部署到 `http://your-server:8000`)。
2.  在 Clawdbot 中，导航至 **Tools > Import**。
3.  选择 **OpenAPI / Swagger**。
4.  输入 `openapi.yaml` 文件的 URL (例如 `https://raw.githubusercontent.com/SealinGp/lark-integration-skill/main/docs/openapi.yaml` 或您自行托管的地址)。
5.  Clawdbot 将自动识别并导入所有 15+ 个工具。

### 方法 2: 手动注册

如果需要手动注册工具，可以将它们定义为指向 [API 使用](#api-使用) 中列出的端点的 `HTTP` 工具。

**示例: 创建任务工具**
- **Method**: `POST`
- **URL**: `http://your-server:8000/api/v1/tasks`
- **Body**:
  ```json
  {
    "summary": "{{summary}}",
    "description": "{{description}}",
    "due_time": {{due_time}},
    "user_id": "{{user_id}}"
  }
  ```

## 安装与部署

### 1. 配置

复制 `.env.example` 为 `.env` 并填写您的凭证：

```bash
cp .env.example .env
```

编辑 `.env` 文件:
```env
LARK_APP_ID="your_app_id"
LARK_APP_SECRET="your_app_secret"
```

### 2. 使用 Docker 运行 (推荐)

```bash
# 构建镜像
docker build -t lark-integration-skill .

# 运行容器
docker run -d -p 8000:8000 --env-file .env --name lark-skill lark-integration-skill
```

### 3. 本地运行

```bash
go run cmd/server/main.go
```

## API 使用

### 创建任务 (Create Task)
**POST** `/api/v1/tasks`
```json
{
  "summary": "修复紧急 Bug",
  "description": "检查日志并修复空指针异常",
  "due_time": 1678888888,
  "user_id": "ou_xxxxxx" 
}
```

### 创建文档 (Create Document)
**POST** `/api/v1/docs`
```json
{
  "title": "项目会议纪要",
  "folder_token": "" 
}
```

### 创建知识库节点 (Create Wiki Node)
**POST** `/api/v1/wiki`
```json
{
  "space_id": "698888888888",
  "title": "新知识页面",
  "obj_type": "docx"
}
```

### 多维文档块操作 (Docx Block Operations)
**GET** `/api/v1/docx/v1/documents/:document_id/blocks/:block_id`
**GET** `/api/v1/docx/v1/documents/:document_id/blocks/:block_id/children`
**POST** `/api/v1/docx/v1/documents/:document_id/blocks/:block_id/children`
**PATCH** `/api/v1/docx/v1/documents/:document_id/blocks/:block_id`
**DELETE** `/api/v1/docx/v1/documents/:document_id/blocks/:block_id/children/batch_delete`
**POST** `/api/v1/docx/v1/documents/blocks/convert`

查看 [docs/API_SPEC.md](docs/API_SPEC.md) 获取完整的 API 详情。
