package service

import (
	"elastic_web_service/internal/auth"
	"elastic_web_service/internal/model"
	"elastic_web_service/internal/repo"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
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
	r.GET("/places", h.handleHTMLPlaces)
	r.GET("/api/places", h.handleAPIPlaces)
	r.GET("/api/recommend", h.jwtMiddleware, h.handleAPIRecommend)
	r.GET("/api/get_token", h.getTokenHandler)

	return r
}

func (h Handler) handleHTMLPlaces(c *gin.Context) {
	// Parse page parameter from the URL query
	pageStr, _ := c.GetQuery("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("wrong page"))
		return
	}
	limit := 10                  // Number of places per page
	offset := (page - 1) * limit // Calculate offset

	// Get places from the store
	places, total, err := h.Repo.GetPlaces(limit, offset)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	totalPages := total/limit + 1
	showPrevious := page > 1
	showNext := page < totalPages

	if err != nil || page < 1 || page > totalPages {
		c.AbortWithError(http.StatusBadRequest, errors.New("wrong page value"))
		return
	}

	tmpl := template.Must(template.New("index.html").ParseFiles("./internal/templates/index.html"))
	data := struct {
		Total        int
		Places       []model.Place
		PrevURL      string
		NextURL      string
		LastURL      string
		ShowNext     bool
		ShowPrevious bool
	}{
		Total:        total,
		Places:       places,
		PrevURL:      fmt.Sprintf("/places?page=%d", page-1),
		NextURL:      fmt.Sprintf("/places?page=%d", page+1),
		LastURL:      fmt.Sprintf("/places?page=%d", total/limit+1),
		ShowNext:     showNext,
		ShowPrevious: showPrevious,
	}

	tmpl.Execute(c.Writer, data)
}

func (h Handler) handleAPIPlaces(c *gin.Context) {
	pageStr, _ := c.GetQuery("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("wrong page"))
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
		c.AbortWithError(http.StatusBadRequest, errors.New("wrong page value"))
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
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid coordinates"))
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid latitude"))
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid longitude"))
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
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ok, err := h.Auth.VerifyToken(tokenString)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}

func (h Handler) errorMiddleware(c *gin.Context) {
	if len(c.Errors) > 0 {
	}
}
