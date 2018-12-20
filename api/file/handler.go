package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type HttpHandler struct {
	usecase Usecase
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	repo := user.NewRepository(session)
	u := NewUsecase(repo)
	h := &HttpHandler{u}
	r = r.Group("/files")
	r.POST("/transcript", h.UploadTranscript)
	r.POST("/image", h.UploadImageProfile)
	r.GET("/transcript/:id", h.DownloadTranscript)
}

// UploadTranscript godoc
// @Summary Upload transcript file
// @Description Upload transcript file
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file body int true "file"
// @Success 200 {object} models.CommonResponse
// @Failure 500 {object} utils.HTTPError
// @Router /files/transcript [post]
func (h *HttpHandler) UploadTranscript(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}

	fn := file.Filename
	t := fn[len(fn)-3:]
	if strings.ToUpper(t) != "PDF" {
		return utils.NewError(c, http.StatusInternalServerError, utils.ErrNotPDFFile)
	}

	src, err := file.Open()
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	defer src.Close()

	u := getUserFromToken(c)
	filename := getTranscriptFilename(u)

	// Destination
	dst, err := os.Create(filename)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	defer dst.Close()

	// Copy
	_, err = io.Copy(dst, src)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}

	err = h.usecase.UpdateUser(u.ID.Hex(), filename)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, models.CommonResponse{Message: "Upload transcript success"})
}

func getUserFromToken(c echo.Context) *models.User {
	t := c.Get("user").(*jwt.Token)
	claims := t.Claims.(*models.JwtCustomClaims)
	return claims.User
}

func getTranscriptFilename(u *models.User) (filename string) {
	path := "files/transcripts"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	r := utils.RandStringBytes(12)
	filename = fmt.Sprintf("%s/%s_%s_%s.pdf", path, strings.ToLower(u.FirstName), strings.ToLower(u.LastName), r)
	return
}

// DownloadTranscript godoc
// @Summary Download transcript file
// @Description Download transcript file
// @Tags files
// @Produce json
// @Param id path string true "user id"
// @Success 200 {array} string
// @Failure 400 {object} utils.HTTPError
// @Failure 401 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /files/transcript/{id} [get]
func (h *HttpHandler) DownloadTranscript(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, utils.ErrInvalidPath)
	}
	user := getUserFromToken(c)
	if user.ID.Hex() != id && !user.IsAdmin() {
		return utils.NewError(c, http.StatusUnauthorized, utils.ErrPermissionDenied)
	}

	filename, err := h.usecase.GetPathTranscript(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.Attachment(filename, filename)
}

// UploadImageProfile godoc
// @Summary Upload image profile
// @Description Upload image profile
// @Tags files
// @Accept image/png,image/gif,image/jpeg
// @Produce json
// @Param image-profile body int true "file"
// @Success 200 {object} models.CommonResponse
// @Failure 500 {object} utils.HTTPError
// @Router /files/image [post]
func (h *HttpHandler) UploadImageProfile(c echo.Context) error {
	file, err := c.FormFile("image-profile")
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	src, err := file.Open()
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	defer src.Close()

	u := getUserFromToken(c)
	filename := getImageFilename(u)

	dst, err := os.Create(filename)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}

	err = h.usecase.UpdateImageProfileUser(u.ID.Hex(), filename)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, models.CommonResponse{Message: "Upload image profile success"})
}

func getImageFilename(u *models.User) (filename string) {
	path := "files/images"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
	filename = fmt.Sprintf("%s/%s_%s.png", path, strings.ToLower(u.FirstName), strings.ToLower(u.LastName))
	return
}
