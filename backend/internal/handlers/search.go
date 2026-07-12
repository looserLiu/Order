package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ledger/backend/internal/service"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"
)

type SearchHandler struct {
	searchService *service.SearchService
}

func NewSearchHandler(searchService *service.SearchService) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

func (h *SearchHandler) Search(c *gin.Context) {
	userID := middleware.GetUserID(c)
	keyword := c.Query("q")
	searchType := c.DefaultQuery("type", "all")

	if keyword == "" {
		response.Error(c, http.StatusBadRequest, "Keyword is required")
		return
	}

	result, err := h.searchService.Search(c.Request.Context(), userID, keyword, searchType, 20)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}
