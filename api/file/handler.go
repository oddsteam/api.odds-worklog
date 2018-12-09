package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type HttpHandler struct {
}

func NewHttpHandler(r *echo.Group) {
	h := &HttpHandler{}
	r.POST("/files/transcript", h.UploadTranscript)
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
