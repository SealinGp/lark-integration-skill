package handlers

import (
	"context"
	"fmt"
	"net/http"

	"lark-integration-skill/internal/models"
	"lark-integration-skill/pkg/larkclient"

	"github.com/gin-gonic/gin"
	larksearch "github.com/larksuite/oapi-sdk-go/v3/service/search/v2"
	larkwiki "github.com/larksuite/oapi-sdk-go/v3/service/wiki/v2"
)

type WikiHandler struct {
	Client *larkclient.ClientWrapper
}

func NewWikiHandler(client *larkclient.ClientWrapper) *WikiHandler {
	return &WikiHandler{Client: client}
}

// SearchWikiNode searches for Wiki nodes
func (h *WikiHandler) SearchWikiNode(c *gin.Context) {
	var req models.WikiSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}

	input := larksearch.NewSearchDocWikiReqBuilder().
		Body(larksearch.NewSearchDocWikiReqBodyBuilder().
			Query(req.Query).
			PageSize(req.PageSize).
			PageToken(req.PageToken).
			Build()).
		Build()

	resp, err := h.Client.Client.Search.DocWiki.Search(context.Background(), input)
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
		Data: models.WikiSearchResponse{
			Items:     resp.Data.ResUnits,
			HasMore:   *resp.Data.HasMore,
			PageToken: *resp.Data.PageToken,
		},
	})
}

// GetWikiNodeInfo retrieves information about a specific Wiki node
func (h *WikiHandler) GetWikiNodeInfo(c *gin.Context) {
	nodeToken := c.Param("node_token")
	if nodeToken == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: "Node Token is required"})
		return
	}

	objType := c.Query("obj_type")

	builder := larkwiki.NewGetNodeSpaceReqBuilder().
		Token(nodeToken)

	if objType != "" {
		builder.ObjType(objType)
	}

	input := builder.Build()

	resp, err := h.Client.Client.Wiki.Space.GetNode(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	node := resp.Data.Node
	c.JSON(http.StatusOK, models.APIResponse{
		Status: "success",
		Data: models.WikiNodeInfoResponse{
			NodeToken:       *node.NodeToken,
			ObjToken:        *node.ObjToken,
			ObjType:         *node.ObjType,
			ParentNodeToken: *node.ParentNodeToken,
			NodeType:        *node.NodeType,
			Title:           *node.Title,
			HasChild:        *node.HasChild,
		},
	})
}

// GetWikiNodeList retrieves a list of sub-nodes for a given space and parent node
func (h *WikiHandler) GetWikiNodeList(c *gin.Context) {
	spaceID := c.Param("space_id")
	if spaceID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: "Space ID is required"})
		return
	}

	parentNodeToken := c.Query("parent_node_token")
	pageToken := c.Query("page_token")
	pageSizeStr := c.Query("page_size")
	var pageSize int
	if pageSizeStr != "" {
		fmt.Sscanf(pageSizeStr, "%d", &pageSize)
	}

	builder := larkwiki.NewListSpaceNodeReqBuilder().
		SpaceId(spaceID).
		PageToken(pageToken)

	if pageSize > 0 {
		builder.PageSize(pageSize)
	}

	if parentNodeToken != "" {
		builder.ParentNodeToken(parentNodeToken)
	}

	input := builder.Build()

	resp, err := h.Client.Client.Wiki.SpaceNode.List(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	// Safely dereference potentially nil fields
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
		Data: models.WikiNodeListResponse{
			Items:     resp.Data.Items,
			HasMore:   hasMore,
			PageToken: nextPageToken,
		},
	})
}

// CreateWikiNode creates a new node in a Knowledge Base (Space)
func (h *WikiHandler) CreateWikiNode(c *gin.Context) {
	var req models.CreateWikiNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}

	objType := req.ObjType
	if objType == "" {
		objType = "docx"
	}

	// Use Wiki V2
	input := larkwiki.NewCreateSpaceNodeReqBuilder().
		SpaceId(req.SpaceID).
		Node(larkwiki.NewNodeBuilder().
			ObjType(objType).
			ParentNodeToken(req.ParentNode).
			Title(req.Title).
			Build()).
		Build()

	resp, err := h.Client.Client.Wiki.SpaceNode.Create(context.Background(), input)
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
		Data: models.WikiNodeResponse{
			NodeToken: *resp.Data.Node.NodeToken,
			ObjToken:  *resp.Data.Node.ObjToken,
			Title:     *resp.Data.Node.Title,
		},
	})
}

// MoveWikiNode moves a node to a new space or parent
func (h *WikiHandler) MoveWikiNode(c *gin.Context) {
	spaceID := c.Param("space_id")
	nodeToken := c.Param("node_token")

	if spaceID == "" || nodeToken == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: "Space ID and Node Token are required"})
		return
	}

	var req models.MoveWikiNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}

	input := larkwiki.NewMoveSpaceNodeReqBuilder().
		SpaceId(spaceID).
		NodeToken(nodeToken).
		Body(larkwiki.NewMoveSpaceNodeReqBodyBuilder().
			TargetParentToken(req.TargetParentToken).
			TargetSpaceId(req.TargetSpaceID).
			Build()).
		Build()

	resp, err := h.Client.Client.Wiki.SpaceNode.Move(context.Background(), input)
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
		Data: models.MoveWikiNodeResponse{
			NodeToken: *resp.Data.Node.NodeToken,
			ObjToken:  *resp.Data.Node.ObjToken,
		},
	})
}

// UpdateWikiNodeTitle updates the title of a Wiki node
func (h *WikiHandler) UpdateWikiNodeTitle(c *gin.Context) {
	spaceID := c.Param("space_id")
	nodeToken := c.Param("node_token")

	if spaceID == "" || nodeToken == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: "Space ID and Node Token are required"})
		return
	}

	var req models.UpdateWikiNodeTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}

	input := larkwiki.NewUpdateTitleSpaceNodeReqBuilder().
		SpaceId(spaceID).
		NodeToken(nodeToken).
		Body(larkwiki.NewUpdateTitleSpaceNodeReqBodyBuilder().
			Title(req.Title).
			Build()).
		Build()

	resp, err := h.Client.Client.Wiki.SpaceNode.UpdateTitle(context.Background(), input)
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
		Data:   models.UpdateWikiNodeTitleResponse{},
	})
}

// MoveDocsToWiki moves a cloud document to a Wiki space
func (h *WikiHandler) MoveDocsToWiki(c *gin.Context) {
	spaceID := c.Param("space_id")
	if spaceID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: "Space ID is required"})
		return
	}

	var req models.MoveDocsToWikiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}

	builder := larkwiki.NewMoveDocsToWikiSpaceNodeReqBuilder().
		SpaceId(spaceID).
		Body(larkwiki.NewMoveDocsToWikiSpaceNodeReqBodyBuilder().
			ObjType(req.ObjType).
			ObjToken(req.ObjToken).
			Apply(req.Apply).
			Build())

	if req.ParentWikiToken != "" {
		builder.Body(larkwiki.NewMoveDocsToWikiSpaceNodeReqBodyBuilder().
			ParentWikiToken(req.ParentWikiToken).
			ObjType(req.ObjType).
			ObjToken(req.ObjToken).
			Apply(req.Apply).
			Build())
	}

	// Re-build body correctly to avoid overwriting if ParentWikiToken is present
	bodyBuilder := larkwiki.NewMoveDocsToWikiSpaceNodeReqBodyBuilder().
		ObjType(req.ObjType).
		ObjToken(req.ObjToken).
		Apply(req.Apply)

	if req.ParentWikiToken != "" {
		bodyBuilder.ParentWikiToken(req.ParentWikiToken)
	}

	input := builder.Body(bodyBuilder.Build()).Build()

	resp, err := h.Client.Client.Wiki.SpaceNode.MoveDocsToWiki(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	// Safely handle nil pointers
	wikiToken := ""
	if resp.Data.WikiToken != nil {
		wikiToken = *resp.Data.WikiToken
	}
	taskID := ""
	if resp.Data.TaskId != nil {
		taskID = *resp.Data.TaskId
	}
	applied := false
	if resp.Data.Applied != nil {
		applied = *resp.Data.Applied
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Status: "success",
		Data: models.MoveDocsToWikiResponse{
			WikiToken: wikiToken,
			TaskID:    taskID,
			Applied:   applied,
		},
	})
}
