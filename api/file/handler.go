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
	r.GET("/transcript/:id", h.DownloadTranscript)
	r.DELETE("/transcript/:id", h.RemoveTranscript)
	r.POST("/image", h.UploadImageProfile)
	r.GET("/image/:id", h.DownloadImageProfile)
	r.DELETE("/image/:id", h.RemoveImageFile)
	r.POST("/degreecertificate", h.UploadDegreeCertificate)
	r.GET("/degreecertificate/:id", h.DownloadDegreeCertificate)
	r.DELETE("/degreecertificate/:id", h.RemoveDegreeCertificate)
}

// UploadDegreeCertificate godoc
// @Summary Upload Degree Certificate file
// @Description Upload Degree Certificate  file
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file body int true "file"
// @Success 200 {object} models.Response
// @Failure 500 {object} utils.HTTPError
// @Router /files/degreecertificate [post]
func (h *HttpHandler) UploadDegreeCertificate(c echo.Context) error {
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
	user, err := h.usecase.GetUserByID(u.ID.Hex())
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	filename := getDegreeCertificateFilename(user)

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

	err = h.usecase.UpdateDegreeCertificate(u.ID.Hex(), filename)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, models.Response{Message: "Upload degree certificate success."})
}

// UploadTranscript godoc
// @Summary Upload transcript file
// @Description Upload transcript file
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file body int true "file"
// @Success 200 {object} models.Response
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
	user, err := h.usecase.GetUserByID(u.ID.Hex())
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	filename := getTranscriptFilename(user)

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

	return c.JSON(http.StatusOK, models.Response{Message: "Upload transcript success."})
}

// DownloadTranscript godoc
// @Summary Download transcript file
// @Description Download transcript file
// @Tags files
// @Produce json
// @Param id path string true "user id"
// @Success 200 {array} string
// @Failure 403 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /files/transcript/{id} [get]
func (h *HttpHandler) DownloadTranscript(c echo.Context) error {
	id := c.Param("id")
	user := getUserFromToken(c)
	if user.ID.Hex() != id && !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	filename, err := h.usecase.GetPathTranscript(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.Attachment(filename, filename)
}

// DownloadTranscript godoc
// @Summary Download degree certificate file
// @Description Download degree certificate file
// @Tags files
// @Produce json
// @Param id path string true "user id"
// @Success 200 {array} string
// @Failure 403 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /files/degreecertificate/{id} [get]
func (h *HttpHandler) DownloadDegreeCertificate(c echo.Context) error {
	id := c.Param("id")
	user := getUserFromToken(c)
	if user.ID.Hex() != id && !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	filename, err := h.usecase.GetPathDegreeCertificate(id)
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
// @Success 200 {object} models.Response
// @Failure 500 {object} utils.HTTPError
// @Router /files/image [post]
func (h *HttpHandler) UploadImageProfile(c echo.Context) error {
	file, err := c.FormFile("image-profile")
	if err != nil {
		return utils.NewError(c, http.StatusBadRequest, err)
	}
	src, err := file.Open()
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	defer src.Close()

	u := getUserFromToken(c)
	user, err := h.usecase.GetUserByID(u.ID.Hex())
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	filename := getImageFilename(user, file.Filename)

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

	return c.JSON(http.StatusOK, models.Response{Message: "Upload image profile success."})
}

// DownloadImageProfile godoc
// @Summary Download image file
// @Description Download image file
// @Tags files
// @Produce json
// @Param id path string true "user id"
// @Success 200 {array} string
// @Failure 403 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /files/image/{id} [get]
func (h *HttpHandler) DownloadImageProfile(c echo.Context) error {
	id := c.Param("id")
	user := getUserFromToken(c)
	if user.ID.Hex() != id && !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	filename, err := h.usecase.GetPathImageProfile(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.Attachment(filename, filename)
}

// RemoveTranscript godoc
// @Summary Remove transcript file
// @Description remove transcript file
// @Tags files
// @Produce json
// @Param id path string true "user id"
// @Success 200 {array} models.Response
// @Failure 403 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /files/transcript/{id} [delete]
func (h *HttpHandler) RemoveTranscript(c echo.Context) error {
	id := c.Param("id")
	user := getUserFromToken(c)
	if user.ID.Hex() != id && !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	filename, err := h.usecase.GetPathTranscript(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, utils.ErrNoTranscriptFile)
	}

	err = h.usecase.RemoveTranscript(filename)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}

	err = h.usecase.UpdateUser(id, "")
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, models.Response{Message: "Remove transcript success"})
}

// RemoveImageFile godoc
// @Summary Remove image file
// @Description remove image file
// @Tags files
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} models.Response
// @Failure 403 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /files/image/{id} [delete]
func (h *HttpHandler) RemoveImageFile(c echo.Context) error {
	id := c.Param("id")
	user := getUserFromToken(c)
	if user.ID.Hex() != id && !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	filename, err := h.usecase.GetPathImageProfile(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, utils.ErrNoTranscriptFile)
	}

	err = h.usecase.RemoveImage(filename)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}

	err = h.usecase.UpdateImageProfileUser(id, "")
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, models.Response{Message: "Remove image success"})
}

// RemoveDegreeCertificate godoc
// @Summary Remove degree certificate file
// @Description remove degree certificate file
// @Tags files
// @Produce json
// @Param id path string true "user id"
// @Success 200 {array} models.Response
// @Failure 403 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /files/degreecertificate/{id} [delete]
func (h *HttpHandler) RemoveDegreeCertificate(c echo.Context) error {
	id := c.Param("id")
	user := getUserFromToken(c)
	if user.ID.Hex() != id && !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	filename, err := h.usecase.GetPathDegreeCertificate(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, utils.ErrNoTranscriptFile)
	}

	err = h.usecase.RemoveDegreeCertificate(filename)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}

	err = h.usecase.UpdateDegreeCertificate(id, "")
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, models.Response{Message: "Remove transcript success"})
}

func getImageFilename(u *models.User, oldName string) (filename string) {
	path := "files/images"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	t := ".png"
	arr := strings.Split(oldName, ".")
	if len(arr) == 2 {
		t = arr[1]
	}
	r := utils.RandStringBytes(12)
	filename = fmt.Sprintf("%s/%s_%s_%s.%s", path, strings.ToLower(u.FirstName), strings.ToLower(u.LastName), r, t)
	return
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

func getDegreeCertificateFilename(u *models.User) (filename string) {
	path := "files/degreecertificate"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	r := utils.RandStringBytes(12)
	filename = fmt.Sprintf("%s/%s_%s_%s.pdf", path, strings.ToLower(u.FirstName), strings.ToLower(u.LastName), r)
	return
}
