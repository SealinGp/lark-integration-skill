package models

import larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"

// Task Models
type CreateTaskRequest struct {
	Summary     string `json:"summary" binding:"required"`
	Description string `json:"description"`
	DueTime     int64  `json:"due_time"` // Unix timestamp
	UserID      string `json:"user_id"`  // Optional: Assign to user (OpenID)
}

type TaskResponse struct {
	TaskID  string `json:"task_id"`
	Summary string `json:"summary"`
	URL     string `json:"url"`
}

type QueryTaskRequest struct {
	UserID    string `json:"user_id"`    // OpenID
	PageToken string `json:"page_token"` // Pagination
	PageSize  int    `json:"page_size"`
}

// Doc Models
type CreateDocRequest struct {
	Title       string `json:"title" binding:"required"`
	FolderToken string `json:"folder_token"` // Optional: Create in specific folder
	Content     string `json:"content"`      // Initial content (simple text for now)
}

type DocResponse struct {
	DocToken string `json:"doc_token"`
	URL      string `json:"url"`
	Title    string `json:"title"`
}

type DocInfoResponse struct {
	DocToken    string `json:"doc_token"`
	Title       string `json:"title"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	OwnerUserID string `json:"owner_user_id"`
}

type DocRawContentResponse struct {
	Content string `json:"content"`
}

type DocBlocksResponse struct {
	Blocks    interface{} `json:"blocks"` // Using interface{} to pass through SDK block structure or simplified map
	HasMore   bool        `json:"has_more"`
	PageToken string      `json:"page_token"`
}

type CreateDocBlockRequest struct {
	Children []*larkdocx.Block `json:"children" binding:"required"`
	Index    *int              `json:"index"` // Optional
}

type CreateDocBlockResponse struct {
	Blocks []*larkdocx.Block `json:"blocks"`
}

type UpdateDocBlockRequest struct {
	*larkdocx.UpdateBlockRequest
}

type UpdateDocBlockResponse struct {
	Block *larkdocx.Block `json:"block"`
}

type GetDocBlockRequest struct {
	// No body params required, derived from URL
}

type GetDocBlockResponse struct {
	Block *larkdocx.Block `json:"block"`
}

type GetDocBlockChildrenRequest struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type GetDocBlockChildrenResponse struct {
	Items     []*larkdocx.Block `json:"items"`
	HasMore   bool              `json:"has_more"`
	PageToken string            `json:"page_token"`
}

type DeleteDocBlockChildrenRequest struct {
	StartIndex *int `json:"start_index"` // Optional
	EndIndex   *int `json:"end_index"`   // Optional
}

type DeleteDocBlockChildrenResponse struct {
	DocumentRevisionId *int `json:"document_revision_id"`
}

type ConvertContentToBlocksRequest struct {
	Content     string `json:"content" binding:"required"`
	ContentType string `json:"content_type"` // "markdown" or "html", default "markdown"
}

type ConvertContentToBlocksResponse struct {
	Blocks []*larkdocx.Block `json:"blocks"`
}

// Wiki Models
type CreateWikiNodeRequest struct {
	SpaceID    string `json:"space_id" binding:"required"`
	ParentNode string `json:"parent_node_token"` // Optional
	Title      string `json:"title" binding:"required"`
	ObjType    string `json:"obj_type"` // "doc", "docx", "sheet", etc. Default "docx"
}

type WikiNodeInfoResponse struct {
	NodeToken       string `json:"node_token"`
	ObjToken        string `json:"obj_token"`
	ObjType         string `json:"obj_type"`
	ParentNodeToken string `json:"parent_node_token"`
	NodeType        string `json:"node_type"`
	Title           string `json:"title"`
	HasChild        bool   `json:"has_child"`
}

type WikiNodeListResponse struct {
	Items     interface{} `json:"items"` // Pass through SDK items
	HasMore   bool        `json:"has_more"`
	PageToken string      `json:"page_token"`
}

type MoveWikiNodeRequest struct {
	TargetParentToken string `json:"target_parent_token"`
	TargetSpaceID     string `json:"target_space_id"`
}

type MoveWikiNodeResponse struct {
	NodeToken string `json:"node_token"`
	ObjToken  string `json:"obj_token"`
}

type UpdateWikiNodeTitleRequest struct {
	Title string `json:"title" binding:"required"`
}

type UpdateWikiNodeTitleResponse struct {
	// Currently empty, just need success status
}

type MoveDocsToWikiRequest struct {
	ParentWikiToken string `json:"parent_wiki_token"` // Optional
	ObjType         string `json:"obj_type" binding:"required"`
	ObjToken        string `json:"obj_token" binding:"required"`
	Apply           bool   `json:"apply"` // Optional
}

type MoveDocsToWikiResponse struct {
	WikiToken string `json:"wiki_token"`
	TaskID    string `json:"task_id"`
	Applied   bool   `json:"applied"`
}

type WikiSearchRequest struct {
	Query     string `json:"query" binding:"required"`
	PageSize  int    `json:"page_size"`
	PageToken string `json:"page_token"`
}

type WikiSearchResponse struct {
	Items     interface{} `json:"items"`
	HasMore   bool        `json:"has_more"`
	PageToken string      `json:"page_token"`
}

type WikiNodeResponse struct {
	NodeToken string `json:"node_token"`
	ObjToken  string `json:"obj_token"`
	Title     string `json:"title"`
	URL       string `json:"url"`
}

// Common Response
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
