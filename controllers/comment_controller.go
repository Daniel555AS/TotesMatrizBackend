package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"totesbackend/config"
	"totesbackend/controllers/utilities"
	"totesbackend/dtos"
	"totesbackend/models"
	"totesbackend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentController struct {
	Service *services.CommentService
	Auth    *utilities.AuthorizationUtil
	Log     *utilities.LogUtil
}

func NewCommentController(service *services.CommentService, auth *utilities.AuthorizationUtil, log *utilities.LogUtil) *CommentController {
	return &CommentController{Service: service, Auth: auth, Log: log}
}

// GetCommentByID godoc
// @Summary      Get comment by ID
// @Description  Retrieves a comment by its unique ID. Requires permission.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Comment ID"
// @Success      200  {object}  dtos.GetCommentDTO       "The retrieved comment"
// @Failure      400  {object}  models.ErrorResponse     "Invalid comment ID"
// @Failure      401  {object}  models.ErrorResponse     "Unauthorized or permission denied"
// @Failure      404  {object}  models.ErrorResponse     "Comment not found"
// @Failure      500  {object}  models.ErrorResponse     "Error retrieving comment or registering log"
// @Security     ApiKeyAuth
// @Router       /comments/{id} [get]
func (cc *CommentController) GetCommentByID(c *gin.Context) {
	idParam := c.Param("id")

	if cc.Log.RegisterLog(c, "Attempting to retrieve Comment with ID: "+idParam) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_COMMENT_BY_ID
	if !cc.Auth.CheckPermission(c, permissionId) {
		_ = cc.Log.RegisterLog(c, "Access denied for GetCommentByID")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		_ = cc.Log.RegisterLog(c, "Invalid comment ID format: "+idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	comment, err := cc.Service.GetCommentByID(id)
	if err != nil {
		_ = cc.Log.RegisterLog(c, "Error retrieving Comment with ID "+idParam+": "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving comment"})
		return
	}

	if comment == nil {
		_ = cc.Log.RegisterLog(c, "Comment with ID "+idParam+" not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	_ = cc.Log.RegisterLog(c, "Successfully retrieved Comment with ID: "+idParam)

	commentDTO := dtos.GetCommentDTO{
		ID:             comment.ID,
		Name:           comment.Name,
		LastName:       comment.LastName,
		Email:          comment.Email,
		Phone:          comment.Phone,
		ResidenceState: comment.ResidenceState,
		ResidenceCity:  comment.ResidenceCity,
		Comment:        comment.Comment,
	}

	c.JSON(http.StatusOK, commentDTO)
}

// GetAllComments godoc
// @Summary      Get all comments
// @Description  Retrieves a list of all submitted comments. Requires permission.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Success      200  {array}   dtos.GetCommentDTO       "List of all comments"
// @Failure      401  {object}  models.ErrorResponse     "Unauthorized or permission denied"
// @Failure      500  {object}  models.ErrorResponse     "Failed to fetch comments or register log"
// @Security     ApiKeyAuth
// @Router       /comments [get]
func (cc *CommentController) GetAllComments(c *gin.Context) {
	if err := cc.Log.RegisterLog(c, "Attempting to retrieve all comments"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_GET_ALL_COMMENTS
	if !cc.Auth.CheckPermission(c, permissionId) {
		_ = cc.Log.RegisterLog(c, "Access denied for GetAllComments")
		return
	}

	comments, err := cc.Service.GetAllComments()
	if err != nil {
		_ = cc.Log.RegisterLog(c, "Error retrieving all comments: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	var commentsDTO []dtos.GetCommentDTO
	for _, comment := range comments {
		commentsDTO = append(commentsDTO, dtos.GetCommentDTO{
			ID:             comment.ID,
			Name:           comment.Name,
			LastName:       comment.LastName,
			Email:          comment.Email,
			Phone:          comment.Phone,
			ResidenceState: comment.ResidenceState,
			ResidenceCity:  comment.ResidenceCity,
			Comment:        comment.Comment,
		})
	}

	_ = cc.Log.RegisterLog(c, "Successfully retrieved all comments")

	c.JSON(http.StatusOK, commentsDTO)
}

// SearchCommentsByEmail godoc
// @Summary      Search comments by email
// @Description  Retrieves all comments that match the provided email address. Requires permission.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        email  query     string                  true  "Email address to search comments by"
// @Success      200    {array}   dtos.GetCommentDTO      "List of matching comments"
// @Failure      400    {object}  models.ErrorResponse    "Email parameter is required"
// @Failure      401    {object}  models.ErrorResponse    "Unauthorized or permission denied"
// @Failure      500    {object}  models.ErrorResponse    "Failed to search comments or register log"
// @Security     ApiKeyAuth
// @Router       /comments/searchByEmail [get]
func (cc *CommentController) SearchCommentsByEmail(c *gin.Context) {
	if err := cc.Log.RegisterLog(c, "Attempting to search comments by email"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_SEARCH_COMMENTS_BY_EMAIL
	if !cc.Auth.CheckPermission(c, permissionId) {
		_ = cc.Log.RegisterLog(c, "Access denied for SearchCommentsByEmail")
		return
	}

	email := c.Query("email")
	if email == "" {
		_ = cc.Log.RegisterLog(c, "Missing 'email' query parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email parameter is required"})
		return
	}

	comments, err := cc.Service.SearchCommentsByEmail(email)
	if err != nil {
		_ = cc.Log.RegisterLog(c, "Error searching comments by email '"+email+"': "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search comments"})
		return
	}

	var commentsDTO []dtos.GetCommentDTO
	for _, comment := range comments {
		commentsDTO = append(commentsDTO, dtos.GetCommentDTO{
			ID:             comment.ID,
			Name:           comment.Name,
			LastName:       comment.LastName,
			Email:          comment.Email,
			Phone:          comment.Phone,
			ResidenceState: comment.ResidenceState,
			ResidenceCity:  comment.ResidenceCity,
			Comment:        comment.Comment,
		})
	}

	_ = cc.Log.RegisterLog(c, "Successfully searched comments by email: "+email)

	c.JSON(http.StatusOK, commentsDTO)
}

// CreateComment godoc
// @Summary      Create a comment
// @Description  Creates a new comment. Requires permission.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        comment  body      dtos.CreateCommentDTO  true  "Comment data to create"
// @Success      201      {object}  dtos.GetCommentDTO     "Created comment"
// @Failure      400      {object}  models.ErrorResponse   "Invalid request data"
// @Failure      401      {object}  models.ErrorResponse   "Unauthorized or permission denied"
// @Failure      500      {object}  models.ErrorResponse   "Failed to create comment or register log"
// @Security     ApiKeyAuth
// @Router       /comments [post]
func (cc *CommentController) CreateComment(c *gin.Context) {
	if err := cc.Log.RegisterLog(c, "Attempting to create a comment"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering log"})
		return
	}

	permissionId := config.PERMISSION_CREATE_COMMENT
	if !cc.Auth.CheckPermission(c, permissionId) {
		_ = cc.Log.RegisterLog(c, "Access denied for CreateComment")
		return
	}

	var dto dtos.CreateCommentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		_ = cc.Log.RegisterLog(c, "Invalid input for CreateComment: "+err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := models.Comment{
		Name:           dto.Name,
		LastName:       dto.LastName,
		Email:          dto.Email,
		Phone:          dto.Phone,
		ResidenceState: dto.ResidenceState,
		ResidenceCity:  dto.ResidenceCity,
		Comment:        dto.Comment,
	}

	createdComment, err := cc.Service.CreateComment(comment)
	if err != nil {
		_ = cc.Log.RegisterLog(c, "Error creating comment: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	commentDTO := dtos.GetCommentDTO{
		ID:             createdComment.ID,
		Name:           createdComment.Name,
		LastName:       createdComment.LastName,
		Email:          createdComment.Email,
		Phone:          createdComment.Phone,
		ResidenceState: createdComment.ResidenceState,
		ResidenceCity:  createdComment.ResidenceCity,
		Comment:        createdComment.Comment,
	}

	_ = cc.Log.RegisterLog(c, "Successfully created comment with ID: "+strconv.Itoa(createdComment.ID))
	c.JSON(http.StatusCreated, commentDTO)
}

// UpdateComment godoc
// @Summary      Update a comment
// @Description  Updates an existing comment by ID. Requires permission.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        id       path      int                  true  "Comment ID"
// @Param        comment  body      dtos.UpdateCommentDTO  true  "Updated comment data"
// @Success      200      {object}  dtos.GetCommentDTO     "Updated comment"
// @Failure      400      {object}  models.ErrorResponse   "Invalid ID or request data"
// @Failure      401      {object}  models.ErrorResponse   "Unauthorized or permission denied"
// @Failure      404      {object}  models.ErrorResponse   "Comment not found"
// @Failure      500      {object}  models.ErrorResponse   "Internal server error or failed update"
// @Security     ApiKeyAuth
// @Router       /comments/{id} [put]
func (cc *CommentController) UpdateComment(c *gin.Context) {
	permissionId := config.PERMISSION_UPDATE_COMMENT

	if !cc.Auth.CheckPermission(c, permissionId) {
		_ = cc.Log.RegisterLog(c, "Access denied for UpdateComment")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = cc.Log.RegisterLog(c, "Invalid comment ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var dto dtos.UpdateCommentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		_ = cc.Log.RegisterLog(c, "Failed to bind JSON in UpdateComment: "+err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := cc.Service.GetCommentByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = cc.Log.RegisterLog(c, "Comment with ID "+strconv.Itoa(id)+" not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		_ = cc.Log.RegisterLog(c, "Internal error retrieving comment with ID "+strconv.Itoa(id)+": "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	comment.Name = dto.Name
	comment.LastName = dto.LastName
	comment.Email = dto.Email
	comment.Phone = dto.Phone
	comment.ResidenceState = dto.ResidenceState
	comment.ResidenceCity = dto.ResidenceCity
	comment.Comment = dto.Comment

	err = cc.Service.UpdateComment(comment)
	if err != nil {
		_ = cc.Log.RegisterLog(c, "Failed to update comment with ID "+strconv.Itoa(id)+": "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	_ = cc.Log.RegisterLog(c, "Successfully updated comment with ID: "+strconv.Itoa(id))

	updatedCommentDTO := dtos.GetCommentDTO{
		ID:             comment.ID,
		Name:           comment.Name,
		LastName:       comment.LastName,
		Email:          comment.Email,
		Phone:          comment.Phone,
		ResidenceState: comment.ResidenceState,
		ResidenceCity:  comment.ResidenceCity,
		Comment:        comment.Comment,
	}

	c.JSON(http.StatusOK, updatedCommentDTO)
}

// SearchCommentsByID godoc
// @Summary      Search comments by ID
// @Description  Searches for comments using a given ID. Requires appropriate permission.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        id       query     string               true  "ID to search for comments"
// @Success      200      {array}   dtos.GetCommentDTO   "List of comments matching the ID"
// @Failure      400      {object}  models.ErrorResponse "Invalid request parameters"
// @Failure      401      {object}  models.ErrorResponse "Unauthorized or permission denied"
// @Failure      404      {object}  models.ErrorResponse "No comments found for the given ID"
// @Failure      500      {object}  models.ErrorResponse "Internal server error or failure in processing"
// @Security     ApiKeyAuth
// @Router       /comments/searchByID [get]
func (cc *CommentController) SearchCommentsByID(c *gin.Context) {
	query := c.Query("id")

	permissionId := config.PERMISSION_SEARCH_COMMENTS_BY_ID

	if !cc.Auth.CheckPermission(c, permissionId) {
		_ = cc.Log.RegisterLog(c, "Access denied for SearchCommentsByID")
		return
	}

	comments, err := cc.Service.SearchCommentsByID(query)
	if err != nil {
		_ = cc.Log.RegisterLog(c, "Error retrieving comments with ID "+query+": "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving comments"})
		return
	}

	if len(comments) == 0 {
		_ = cc.Log.RegisterLog(c, "No comments found for ID "+query)
		c.JSON(http.StatusNotFound, gin.H{"message": "No comments found"})
		return
	}

	var commentsDTO []dtos.GetCommentDTO
	for _, comment := range comments {
		commentsDTO = append(commentsDTO, dtos.GetCommentDTO{
			ID:             comment.ID,
			Name:           comment.Name,
			LastName:       comment.LastName,
			Email:          comment.Email,
			Phone:          comment.Phone,
			ResidenceState: comment.ResidenceState,
			ResidenceCity:  comment.ResidenceCity,
			Comment:        comment.Comment,
		})
	}

	_ = cc.Log.RegisterLog(c, "Successfully retrieved comments with ID: "+query)
	c.JSON(http.StatusOK, commentsDTO)
}

// SearchCommentsByName godoc
// @Summary      Search comments by name
// @Description  Searches for comments using a given name. Requires appropriate permission.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        name     query     string               true  "Name to search for comments"
// @Success      200      {array}   dtos.GetCommentDTO   "List of comments matching the name"
// @Failure      400      {object}  models.ErrorResponse "Invalid request parameters"
// @Failure      401      {object}  models.ErrorResponse "Unauthorized or permission denied"
// @Failure      404      {object}  models.ErrorResponse "No comments found for the given name"
// @Failure      500      {object}  models.ErrorResponse "Internal server error or failure in processing"
// @Security     ApiKeyAuth
// @Router       /comments/searchByName [get]
func (cc *CommentController) SearchCommentsByName(c *gin.Context) {
	query := c.Query("name")

	permissionId := config.PERMISSION_SEARCH_COMMENTS_BY_NAME

	if !cc.Auth.CheckPermission(c, permissionId) {
		_ = cc.Log.RegisterLog(c, "Access denied for SearchCommentsByName")
		return
	}

	comments, err := cc.Service.SearchCommentsByName(query)
	if err != nil {
		_ = cc.Log.RegisterLog(c, "Error retrieving comments with name "+query+": "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving comments"})
		return
	}

	if len(comments) == 0 {
		_ = cc.Log.RegisterLog(c, "No comments found for name "+query)
		c.JSON(http.StatusNotFound, gin.H{"message": "No comments found"})
		return
	}

	var commentsDTO []dtos.GetCommentDTO
	for _, comment := range comments {
		commentsDTO = append(commentsDTO, dtos.GetCommentDTO{
			ID:             comment.ID,
			Name:           comment.Name,
			LastName:       comment.LastName,
			Email:          comment.Email,
			Phone:          comment.Phone,
			ResidenceState: comment.ResidenceState,
			ResidenceCity:  comment.ResidenceCity,
			Comment:        comment.Comment,
		})
	}

	_ = cc.Log.RegisterLog(c, "Successfully retrieved comments with name: "+query)
	c.JSON(http.StatusOK, commentsDTO)
}
