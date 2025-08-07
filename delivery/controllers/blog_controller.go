package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/InkForge/Blog_Website/delivery/controllers/dto"
	"github.com/InkForge/Blog_Website/domain"
)

type BlogController struct {
	BlogUsecase domain.IBlogUseCase
}

func NewBlogController(usecase domain.IBlogUseCase) *BlogController {
	return &BlogController{
		BlogUsecase: usecase,
	}
}

func (bc *BlogController) CreateBlog(c *gin.Context) {
	ogCtx := c.Request.Context()

	//Binding request to Blog struct in json format
	var input dto.BlogJson
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Payload", "details": err.Error()})
		return
	}

	// obtaining  userID from auth , set by jwt

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: user ID is of wrong type"})
		return
	}
	input.UserID = userIDStr

	// creating context to pass down the layers
	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()

	//delegate to usecase
	blogID, err := bc.BlogUsecase.CreateBlog(ctx, input.ToDomainBlog())
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out during blog creation"})
		case errors.Is(err, domain.ErrBlogRequired):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Blog is required"})
		case errors.Is(err, domain.ErrEmptyTitle):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		case errors.Is(err, domain.ErrEmptyContent):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		case errors.Is(err, domain.ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		case errors.Is(err, domain.ErrInsertingDocuments):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert blog document", "details": err.Error()})
		case errors.Is(err, domain.ErrInvalidBlogIdFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID format"})
		case errors.Is(err, domain.ErrQueryFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query", "details": err.Error()})
		case errors.Is(err, domain.ErrDocumentDecoding):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode document", "details": err.Error()})
		default:
			log.Printf("Error creating blog: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Blog created successfully", "blog_id": blogID})
}

func (bc *BlogController) GetAllBlogs(c *gin.Context) {
	ogCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()

	// default values for pagination
	page := 1
	limit := 10

	if p, err := strconv.Atoi(c.DefaultQuery("p", "1")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(c.DefaultQuery("l", "10")); err == nil && l > 0 {
		limit = l
	}

	paginatedBlogs, err := bc.BlogUsecase.GetAllBlogs(ctx, page, limit)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out during blog retrieval"})
		case errors.Is(err, domain.ErrRetrievingDocuments):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blogs", "details": err.Error()})
		case errors.Is(err, domain.ErrDecodingDocument):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode blog document", "details": err.Error()})
		case errors.Is(err, domain.ErrCursorIteration):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor iteration error", "details": err.Error()})
		case errors.Is(err, domain.ErrCursorFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor failed", "details": err.Error()})
		case errors.Is(err, domain.ErrBlogNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "No blogs found"})
		default:
			log.Printf("Error fetching all blogs: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blogs", "details": err.Error()})
		}
		return
	}
	jsonResponse := dto.FromDomainPaginatedBlogs(*paginatedBlogs)
	c.JSON(http.StatusOK, gin.H{"blogs": jsonResponse})
}

func (bc *BlogController) GetBlogByID(c *gin.Context) {
	ogCtx := c.Request.Context()
	blogID := c.Param("id")

	// obtaining  userID from auth , set by jwt
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: user ID is of wrong type"})
		return
	}

	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()

	blog, err := bc.BlogUsecase.GetBlogByID(ctx, blogID, userIDStr)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out during blog retrieval"})
		case errors.Is(err, domain.ErrBlogIDRequired):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		case errors.Is(err, domain.ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		case errors.Is(err, domain.ErrBlogNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		case errors.Is(err, domain.ErrViewRecordAlreadyExists):
			log.Printf("View record already exists for blog %s by user %s. Returning blog anyway.", blogID, userID)
			c.JSON(http.StatusOK, gin.H{"blog": blog})
		case errors.Is(err, domain.ErrCreateViewRecordFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create view record", "details": err.Error()})
		case errors.Is(err, domain.ErrIncrementViewFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment blog view count", "details": err.Error()})
		case errors.Is(err, domain.ErrRetrievingDocuments):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blog document", "details": err.Error()})
		case errors.Is(err, domain.ErrDocumentDecoding):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode blog document", "details": err.Error()})
		case errors.Is(err, domain.ErrQueryFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query", "details": err.Error()})
		default:
			log.Printf("Error fetching blog by ID %s: %v", blogID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blog", "details": err.Error()})
		}
		return
	}
	jsonBlog := dto.FromDomainBlog(blog)

	c.JSON(http.StatusOK, gin.H{"blog": jsonBlog})
}

func (bc *BlogController) UpdateBlog(c *gin.Context) {

	ogCtx := c.Request.Context()
	blogID := c.Param("id")

	var jsonBlog dto.BlogJson
	if err := c.ShouldBindJSON(&jsonBlog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Payload", "details": err.Error()})
		return
	}

	// obtaining  userID from auth , set by jwt
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: user ID is of wrong type"})
		return
	}
	jsonBlog.BlogID = blogID

	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()


	err := bc.BlogUsecase.UpdateBlog(ctx, jsonBlog.ToDomainBlog(), userIDStr)

	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out during blog update"})
		case errors.Is(err, domain.ErrBlogRequired):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Blog is required"})
		case errors.Is(err, domain.ErrBlogIDRequired):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
		case errors.Is(err, domain.ErrBlogNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		case errors.Is(err, domain.ErrNoBlogChangesMade):
			c.JSON(http.StatusNotModified, gin.H{"error": "No changes were made to the blog"})
		case errors.Is(err, domain.ErrInvalidBlogID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		case errors.Is(err, domain.ErrInvalidBlogIdFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID format"})
		case errors.Is(err, domain.ErrRetrievingDocuments):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blog document", "details": err.Error()})
		case errors.Is(err, domain.ErrUpdatingDocument):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog document", "details": err.Error()})
		case errors.Is(err, domain.ErrDocumentDecoding):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode blog document", "details": err.Error()})
		case errors.Is(err, domain.ErrQueryFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query", "details": err.Error()})
		case errors.Is(err, domain.ErrInsertingDocuments):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert document(s)", "details": err.Error()})
		case errors.Is(err, domain.ErrCursorIteration):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor iteration error", "details": err.Error()})
		case errors.Is(err, domain.ErrCursorFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor failed", "details": err.Error()})
		default:
			log.Printf("Error updating blog %s: %v", blogID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog", "details": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog updated successfully"})
}

func (bc *BlogController) DeleteBlog(c *gin.Context) {
	ogCtx := c.Request.Context()
	blogID := c.Param("id")

	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()

	err := bc.BlogUsecase.DeleteBlog(ctx, blogID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out during blog deletion"})
			return
		}
		switch {
		case errors.Is(err, domain.ErrBlogIDRequired):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrBlogNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			log.Printf("Error deleting blog %s: %v", blogID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog", "details": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func (bc *BlogController) Search(c *gin.Context) {
	ogCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()

	title := c.DefaultQuery("title", "")
	author := c.DefaultQuery("author", "")
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("l", "10"))

	paginatedBlogs, err := bc.BlogUsecase.SearchBlogs(ctx, title, author, page, limit)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out during blog search"})
		case errors.Is(err, domain.ErrRetrievingDocuments):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blogs", "details": err.Error()})
		case errors.Is(err, domain.ErrQueryFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query", "details": err.Error()})
		case errors.Is(err, domain.ErrDocumentDecoding):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode blog document", "details": err.Error()})
		case errors.Is(err, domain.ErrBlogNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "No blogs found"})
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "No users found for author search"})
		default:
			log.Printf("Error searching blogs: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search blogs", "details": err.Error()})
		}
		return
	}
	jsonResponse := dto.FromDomainPaginatedBlogs(*paginatedBlogs)
	c.JSON(http.StatusOK, gin.H{"blogs": jsonResponse})
}

func (bc *BlogController) FilterBlogs(c *gin.Context) {
	ogCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ogCtx, 5*time.Second)
	defer cancel()

	tags := c.QueryArray("tag")
	popularity := c.DefaultQuery("popularity", "")
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("l", "10"))

	params := domain.FilterParams{
		TagIDs:     tags,
		Popularity: popularity,
		Page:       page,
		Limit:      limit,
	}

	paginatedBlogs, err := bc.BlogUsecase.FilterBlogs(ctx, params)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out during blog filtering"})
		case errors.Is(err, domain.ErrRetrievingDocuments):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blogs", "details": err.Error()})
		case errors.Is(err, domain.ErrQueryFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query", "details": err.Error()})
		case errors.Is(err, domain.ErrDocumentDecoding):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode blog document", "details": err.Error()})
		case errors.Is(err, domain.ErrBlogNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "No blogs found"})
		default:
			log.Printf("Error filtering blogs: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to filter blogs", "details": err.Error()})
		}
		return
	}
	jsonResponse := dto.FromDomainPaginatedBlogs(*paginatedBlogs)
	c.JSON(http.StatusOK, gin.H{"blogs": jsonResponse})
}
