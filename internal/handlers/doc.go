package handlers

import (
	"context"
	"fmt"
	"net/http"

	"lark-integration-skill/internal/models"
	"lark-integration-skill/pkg/larkclient"

	"github.com/gin-gonic/gin"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
)

type DocHandler struct {
	Client *larkclient.ClientWrapper
}

func NewDocHandler(client *larkclient.ClientWrapper) *DocHandler {
	return &DocHandler{Client: client}
}

// CreateDoc creates a new Docx file
func (h *DocHandler) CreateDoc(c *gin.Context) {
	var req models.CreateDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}

	// Use Docx V1 to create document
	input := larkdocx.NewCreateDocumentReqBuilder().
		Body(larkdocx.NewCreateDocumentReqBodyBuilder().
			Title(req.Title).
			FolderToken(req.FolderToken).
			Build()).
		Build()

	// Note: Client.Docx gives access to Docx V1 Service
	resp, err := h.Client.Client.Docx.Document.Create(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Status: "success",
		Data: models.DocResponse{
			DocToken: *resp.Data.Document.DocumentId,
			URL:      fmt.Sprintf("https://open.larksuite.com/docx/%s", *resp.Data.Document.DocumentId), // Construct URL manually as SDK might not return full URL
			Title:    *resp.Data.Document.Title,
		},
	})
}

// GetDocument retrieves basic information about a Docx file (using Drive Meta API)
func (h *DocHandler) GetDocument(c *gin.Context) {
	docToken := c.Param("doc_token")
	if docToken == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: "Doc Token is required"})
		return
	}

	// Use Drive Meta API to batch query (single item) for rich metadata
	input := larkdrive.NewBatchQueryMetaReqBuilder().
		MetaRequest(larkdrive.NewMetaRequestBuilder().
			RequestDocs([]*larkdrive.RequestDoc{
				larkdrive.NewRequestDocBuilder().
					DocToken(docToken).
					DocType("docx").
					Build(),
			}).
			Build()).
		Build()

	// Note: client.Drive.Meta -> V1 Meta service
	resp, err := h.Client.Client.Drive.Meta.BatchQuery(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		// Even if API call "succeeds", logic might fail
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	if len(resp.Data.Metas) == 0 {
		c.JSON(http.StatusNotFound, models.APIResponse{Status: "error", Message: "Document not found"})
		return
	}

	meta := resp.Data.Metas[0]

	data := models.DocInfoResponse{
		DocToken:   *meta.DocToken,
		Title:      *meta.Title,
		CreateTime: "",
		UpdateTime: "",
	}

	// Safely dereference optional fields
	if meta.CreateTime != nil {
		data.CreateTime = *meta.CreateTime // It's a string timestamp (ms)
	}
	if meta.LatestModifyTime != nil {
		data.UpdateTime = *meta.LatestModifyTime
	}
	if meta.OwnerId != nil {
		data.OwnerUserID = *meta.OwnerId
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Status: "success",
		Data:   data,
	})
}

// GetDocumentRawContent retrieves the raw text content of a Docx file
func (h *DocHandler) GetDocumentRawContent(c *gin.Context) {
	docToken := c.Param("doc_token")
	if docToken == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: "Doc Token is required"})
		return
	}

	input := larkdocx.NewRawContentDocumentReqBuilder().
		DocumentId(docToken).
		Build()

	resp, err := h.Client.Client.Docx.Document.RawContent(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Status: "success",
		Data: models.DocRawContentResponse{
			Content: *resp.Data.Content,
		},
	})
}

// GetDocumentBlocks retrieves all blocks (or paginated) of a Docx file
func (h *DocHandler) GetDocumentBlocks(c *gin.Context) {
	docToken := c.Param("doc_token")
	pageToken := c.Query("page_token")
	pageSizeStr := c.DefaultQuery("page_size", "500") // Default to 500 blocks

	if docToken == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: "Doc Token is required"})
		return
	}

	// Convert pageSize
	var pageSize int
	fmt.Sscanf(pageSizeStr, "%d", &pageSize)

	input := larkdocx.NewListDocumentBlockReqBuilder().
		DocumentId(docToken).
		PageSize(pageSize).
		PageToken(pageToken).
		Build()

	resp, err := h.Client.Client.Docx.DocumentBlock.List(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	// Determine HasMore safely
	hasMore := false
	if resp.Data.HasMore != nil {
		hasMore = *resp.Data.HasMore
	}
	nextPageToken := ""
	if resp.Data.PageToken != nil {
		nextPageToken = *resp.Data.PageToken
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Status: "success",
		Data: models.DocBlocksResponse{
			Blocks:    resp.Data.Items,
			HasMore:   hasMore,
			PageToken: nextPageToken,
		},
	})
}

// CreateDocBlock creates children blocks in a document
func (h *DocHandler) CreateDocBlock(c *gin.Context) {
	documentID := c.Param("document_id")
	blockID := c.Param("block_id")

	if documentID == "" || blockID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: "Document ID and Block ID are required"})
		return
	}

	var req models.CreateDocBlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}

	bodyBuilder := larkdocx.NewCreateDocumentBlockChildrenReqBodyBuilder().
		Children(req.Children)

	if req.Index != nil {
		bodyBuilder.Index(*req.Index)
	}

	input := larkdocx.NewCreateDocumentBlockChildrenReqBuilder().
		DocumentId(documentID).
		BlockId(blockID).
		Body(bodyBuilder.Build()).
		Build()

	resp, err := h.Client.Client.Docx.DocumentBlockChildren.Create(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Status: "success",
		Data: models.CreateDocBlockResponse{
			Blocks: resp.Data.Children,
		},
	})
}

// UpdateDocBlock updates a block in a document
func (h *DocHandler) UpdateDocBlock(c *gin.Context) {
	documentID := c.Param("document_id")
	blockID := c.Param("block_id")

	if documentID == "" || blockID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: "Document ID and Block ID are required"})
		return
	}

	var req models.UpdateDocBlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}

	// Ensure BlockId is set correctly from URL param if not in body, though typically body has details
	req.UpdateBlockRequest.BlockId = &blockID

	input := larkdocx.NewPatchDocumentBlockReqBuilder().
		DocumentId(documentID).
		BlockId(blockID).
		UpdateBlockRequest(req.UpdateBlockRequest).
		Build()

	resp, err := h.Client.Client.Docx.DocumentBlock.Patch(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Status: "success",
		Data: models.UpdateDocBlockResponse{
			Block: resp.Data.Block,
		},
	})
}

// GetDocBlock retrieves a specific block in a document
func (h *DocHandler) GetDocBlock(c *gin.Context) {
	documentID := c.Param("document_id")
	blockID := c.Param("block_id")

	if documentID == "" || blockID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: "Document ID and Block ID are required"})
		return
	}

	input := larkdocx.NewGetDocumentBlockReqBuilder().
		DocumentId(documentID).
		BlockId(blockID).
		Build()

	resp, err := h.Client.Client.Docx.DocumentBlock.Get(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Status: "success",
		Data: models.GetDocBlockResponse{
			Block: resp.Data.Block,
		},
	})
}

// GetDocBlockChildren retrieves child blocks of a specific block
func (h *DocHandler) GetDocBlockChildren(c *gin.Context) {
	documentID := c.Param("document_id")
	blockID := c.Param("block_id")
	pageToken := c.Query("page_token")
	pageSizeStr := c.DefaultQuery("page_size", "500")

	if documentID == "" || blockID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: "Document ID and Block ID are required"})
		return
	}

	var pageSize int
	fmt.Sscanf(pageSizeStr, "%d", &pageSize)

	inputBuilder := larkdocx.NewGetDocumentBlockChildrenReqBuilder().
		DocumentId(documentID).
		BlockId(blockID).
		PageSize(pageSize)

	if pageToken != "" {
		inputBuilder.PageToken(pageToken)
	}

	resp, err := h.Client.Client.Docx.DocumentBlockChildren.Get(context.Background(), inputBuilder.Build())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	hasMore := false
	if resp.Data.HasMore != nil {
		hasMore = *resp.Data.HasMore
	}
	nextPageToken := ""
	if resp.Data.PageToken != nil {
		nextPageToken = *resp.Data.PageToken
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Status: "success",
		Data: models.GetDocBlockChildrenResponse{
			Items:     resp.Data.Items,
			HasMore:   hasMore,
			PageToken: nextPageToken,
		},
	})
}

// DeleteDocBlockChildren batch deletes child blocks
func (h *DocHandler) DeleteDocBlockChildren(c *gin.Context) {
	documentID := c.Param("document_id")
	blockID := c.Param("block_id")

	if documentID == "" || blockID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: "Document ID and Block ID are required"})
		return
	}

	var req models.DeleteDocBlockChildrenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}

	bodyBuilder := larkdocx.NewBatchDeleteDocumentBlockChildrenReqBodyBuilder()
	if req.StartIndex != nil {
		bodyBuilder.StartIndex(*req.StartIndex)
	}
	if req.EndIndex != nil {
		bodyBuilder.EndIndex(*req.EndIndex)
	}

	input := larkdocx.NewBatchDeleteDocumentBlockChildrenReqBuilder().
		DocumentId(documentID).
		BlockId(blockID).
		Body(bodyBuilder.Build()).
		Build()

	resp, err := h.Client.Client.Docx.DocumentBlockChildren.BatchDelete(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Status: "success",
		Data: models.DeleteDocBlockChildrenResponse{
			DocumentRevisionId: resp.Data.DocumentRevisionId,
		},
	})
}

// ConvertContentToBlocks converts Markdown/HTML to Blocks
func (h *DocHandler) ConvertContentToBlocks(c *gin.Context) {
	var req models.ConvertContentToBlocksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}

	bodyBuilder := larkdocx.NewConvertDocumentReqBodyBuilder().
		Content(req.Content).
		ContentType(req.ContentType)

	if req.ContentType == "" {
		bodyBuilder.ContentType("markdown")
	}

	input := larkdocx.NewConvertDocumentReqBuilder().
		Body(bodyBuilder.Build()).
		Build()

	resp, err := h.Client.Client.Docx.Document.Convert(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Status: "success",
		Data: models.ConvertContentToBlocksResponse{
			Blocks: resp.Data.Blocks,
		},
	})
}
