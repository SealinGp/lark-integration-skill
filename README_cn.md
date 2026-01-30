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

本服务旨在作为 **Skill (技能)** 供运行在 **Clawdbot** 上的 AI Agent 使用。

### 配置方法

要使用此技能，请将本项目所在的路径添加到您的 `clawdbot.json` (或 `openclaw.json`) 配置文件中。

**`clawdbot.json` 配置示例:**

```json
{
  "skills": {
    "load": {
      "extraDirs": [
        "/absolute/path/to/lark-integration-skill"
      ]
    }
  },
  "tools": {
    "allow": ["*"],
    "deny": []
  }
}
```

*   请将 `/absolute/path/to/lark-integration-skill` 替换为您克隆此仓库的实际绝对路径。
*   重启 Clawdbot/OpenClaw 以加载该技能。


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
