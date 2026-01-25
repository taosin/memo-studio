package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"memo-studio/backend/models"

	"github.com/gin-gonic/gin"
)

const maxUploadSize = 20 << 20 // 20MB

func storageBaseDir() string {
	if v := strings.TrimSpace(os.Getenv("MEMO_STORAGE_DIR")); v != "" {
		return v
	}
	return "./storage"
}

func randomHex(n int) string {
	if n <= 0 {
		n = 16
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// 最坏情况下退化为时间戳（仍可用）
		return hex.EncodeToString([]byte(time.Now().Format("20060102150405.000000000")))
	}
	return hex.EncodeToString(b)
}

func ensureDir(p string) error {
	return os.MkdirAll(p, 0o755)
}

func saveMultipartFile(file multipart.File, dst string) (size int64, sha string, err error) {
	h := sha256.New()
	if err := ensureDir(filepath.Dir(dst)); err != nil {
		return 0, "", err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return 0, "", err
	}
	defer out.Close()

	w := io.MultiWriter(out, h)
	n, err := io.Copy(w, file)
	if err != nil {
		return 0, "", err
	}
	return n, hex.EncodeToString(h.Sum(nil)), nil
}

func getOptionalUserID(c *gin.Context) *int {
	v, ok := c.Get("userID")
	if !ok {
		return nil
	}
	if id, ok := v.(int); ok {
		return &id
	}
	return nil
}

// UploadResource 上传附件
// POST /api/resources (multipart/form-data)
// form field: file
func UploadResource(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请使用 multipart/form-data 并提供 file 字段"})
		return
	}

	if fh.Size <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件为空"})
		return
	}

	file, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取上传文件"})
		return
	}
	defer file.Close()

	userID := getOptionalUserID(c)
	userSeg := "public"
	if userID != nil && *userID > 0 {
		userSeg = "u" + strconv.Itoa(*userID)
	}
	now := time.Now()
	dateSeg := filepath.Join(
		now.Format("2006"),
		now.Format("01"),
		now.Format("02"),
	)

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	name := strings.TrimSuffix(fh.Filename, filepath.Ext(fh.Filename))
	name = strings.TrimSpace(name)
	if name == "" {
		name = "upload"
	}
	safeName := sanitizeFilename(name)
	filename := safeName + "_" + randomHex(8) + ext

	relPath := filepath.ToSlash(filepath.Join(userSeg, dateSeg, filename))
	dst := filepath.Join(storageBaseDir(), filepath.FromSlash(relPath))

	size, sha, err := saveMultipartFile(file, dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}

	mimeType := strings.TrimSpace(fh.Header.Get("Content-Type"))
	res, err := models.CreateResource(userID, fh.Filename, relPath, mimeType, size, sha)
	if err != nil {
		// DB 写入失败：尽量清理文件
		_ = os.Remove(dst)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "写入资源记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func sanitizeFilename(s string) string {
	s = strings.TrimSpace(s)
	// 仅保留常见安全字符，其他转成下划线
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_' || r == '.':
			b.WriteRune(r)
		case r >= 0x4e00 && r <= 0x9fff: // 简单支持中文
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	out := strings.Trim(b.String(), "._-")
	if out == "" {
		return "upload"
	}
	return out
}

