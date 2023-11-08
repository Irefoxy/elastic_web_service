package service

import (
	"elastic_web_service/internal/auth"
	"elastic_web_service/internal/model"
	"elastic_web_service/internal/repo"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

var _ Service = (*Handler)(nil)

type Handler struct {
	Repo    repo.Repo
	Auth    auth.Auth
	Address string
}

func NewHandler(repo repo.Repo, auth auth.Auth, address string) Handler {
	return Handler{repo, auth, address}
}

func (h Handler) Run() error {
	r := h.init()
	return r.Run(h.Address)
}

func (h Handler) init() *gin.Engine {
	r := gin.Default()
	r.Use(h.errorMiddleware)
	r.GET("/api/places", h.handleAPIPlaces)
	r.GET("/api/recommend", h.jwtMiddleware, h.handleAPIRecommend)
	r.GET("/api/get_token", h.getTokenHandler)

	return r
}

func (h Handler) handleAPIPlaces(c *gin.Context) {
	pageStr, _ := c.GetQuery("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("wrong page")).SetType(gin.ErrorTypePublic)
		return
	}
	limit := 10
	offset := (page - 1) * limit

	places, total, err := h.Repo.GetPlaces(limit, offset)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	totalPages := total/limit + 1
	prevPage := page - 1
	nextPage := page + 1
	showPrevious := page > 1
	showNext := page < totalPages

	if err != nil || page < 1 || page > totalPages {
		c.AbortWithError(http.StatusBadRequest, errors.New("wrong page value")).SetType(gin.ErrorTypePublic)
		return
	}

	response := struct {
		Name     string        `json:"name"`
		Total    int           `json:"total"`
		Places   []model.Place `json:"places"`
		PrevPage int           `json:"prev_page,omitempty"`
		NextPage int           `json:"next_page,omitempty"`
		LastPage int           `json:"last_page"`
	}{
		Name:     "Places",
		Total:    total,
		Places:   places,
		PrevPage: prevPage,
		NextPage: nextPage,
		LastPage: totalPages,
	}

	if !showPrevious {
		response.PrevPage = 0
	}

	if !showNext {
		response.NextPage = 0
	}
	c.Header("Content-Type", "application/json")
	err = json.NewEncoder(c.Writer).Encode(response)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (h Handler) getTokenHandler(c *gin.Context) {
	token, err := h.Auth.GenerateToken()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response := map[string]string{
		"token": token,
	}
	c.Header("Content-Type", "application/json")
	err = json.NewEncoder(c.Writer).Encode(response)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (h Handler) handleAPIRecommend(c *gin.Context) {
	latStr, _ := c.GetQuery("lat")
	lonStr, _ := c.GetQuery("lon")

	if latStr == "" || lonStr == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid coordinates")).SetType(gin.ErrorTypePublic)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid latitude")).SetType(gin.ErrorTypePublic)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid longitude")).SetType(gin.ErrorTypePublic)
		return
	}

	places, err := h.Repo.GetRecommendPlaces(lon, lat)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"name":   "Recommendation",
		"places": places,
	}

	c.Header("Content-Type", "application/json")
	if err := json.NewEncoder(c.Writer).Encode(response); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (h Handler) jwtMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithError(http.StatusUnauthorized, errors.New("no authorization token")).SetType(gin.ErrorTypePublic)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		c.AbortWithError(http.StatusUnauthorized, errors.New("wrong authorization token format")).SetType(gin.ErrorTypePublic)
		return
	}
	ok, err := h.Auth.VerifyToken(tokenString)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err).SetType(gin.ErrorTypePublic)
		return
	}
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, errors.New("invalid authorization token")).SetType(gin.ErrorTypePublic)
		return
	}
	c.Next()
}

func (h Handler) errorMiddleware(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		switch c.Errors[0].Type {
		case gin.ErrorTypePublic:
			c.JSON(-1, gin.H{"error": c.Errors[0].Error()})
		default:
			c.JSON(-1, gin.H{"error": "Something went wrong"})
		}
	}
}
