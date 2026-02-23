package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// WhisperConfig 配置
type WhisperConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

// SpeechToTextForm 语音转文本表单
type SpeechToTextForm struct {
	Language   string  `form:"language"`
	Prompt     string  `form:"prompt"`
	Temperature float64 `form:"temperature"`
}

// UploadResourceAndTranscribe 上传并转录
// POST /api/resources/transcribe
func UploadResourceAndTranscribe(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 20<<20)

	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请使用 multipart/form-data 提供 file 字段"})
		return
	}

	if !isAudioFile(fh.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持音频文件（mp3, wav, m4a, ogg, webm, flac）"})
		return
	}

	file, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取文件"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取文件失败"})
		return
	}

	// 保存文件
	storagePath, err := saveAudioFile(data, fh.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}

	// 获取转录参数
	var req struct {
		Language   string
		Prompt     string
		Temperature float64
	}
	req.Language = c.PostForm("language")
	req.Prompt = c.PostForm("prompt")
	if tempStr := c.PostForm("temperature"); tempStr != "" {
		req.Temperature, _ = strconv.ParseFloat(tempStr, 64)
	}

	// 尝试转录
	text, duration, language := "", 0.0, ""
	config := getWhisperConfig()
	if config.APIKey != "" {
		t, d, lang, err := callWhisperAPI(data, fh.Filename, config, req.Language, req.Prompt, req.Temperature)
		if err == nil {
			text, duration, language = t, d, lang
		}
	}

	// 获取文件 URL
	url := "/uploads/" + storagePath

	c.JSON(http.StatusCreated, gin.H{
		"id":          0, // 简化：实际应保存到数据库
		"filename":    fh.Filename,
		"storage_path": storagePath,
		"url":         url,
		"mime_type":   fh.Header.Get("Content-Type"),
		"size":        fh.Size,
		"transcript":  text,
		"duration":    duration,
		"language":    language,
		"created_at":  time.Now().Format(time.RFC3339),
	})
}

// SpeechToTextOnly 仅转录（不保存）
// POST /api/speech-to-text
func SpeechToTextOnly(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供 file 字段"})
		return
	}

	if !isAudioFile(fh.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持音频文件"})
		return
	}

	file, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取文件"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取文件失败"})
		return
	}

	config := getWhisperConfig()
	if config.APIKey == "" {
		c.JSON(http.StatusOK, gin.H{
			"text":       "",
			"configured": false,
			"message":    "请设置 OPENAI_API_KEY 环境变量以启用语音转文本",
		})
		return
	}

	lang := c.PostForm("language")
	prompt := c.PostForm("prompt")
	var temp float64
	if tStr := c.PostForm("temperature"); tStr != "" {
		temp, _ = strconv.ParseFloat(tStr, 64)
	}

	text, duration, language, err := callWhisperAPI(data, fh.Filename, config, lang, prompt, temp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"text":       text,
		"duration":   duration,
		"language":  language,
		"configured": true,
	})
}

// ========== 内部函数 ==========

func isAudioFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	audioExts := []string{".mp3", ".wav", ".m4a", ".ogg", ".webm", ".flac", ".mp4"}
	for _, e := range audioExts {
		if ext == e {
			return true
		}
	}
	return false
}

func saveAudioFile(data []byte, filename string) (string, error) {
	userID := getOptionalUserIDForAudio()
	userSeg := "public"
	if userID != nil && *userID > 0 {
		userSeg = "u" + strconv.Itoa(*userID)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	name := strings.TrimSuffix(filename, ext)
	name = strings.TrimSpace(name)
	if name == "" {
		name = "audio"
	}
	safeName := sanitizeFilename(name)
	newFilename := safeName + "_" + randomHex(8) + ext

	now := time.Now()
	dateSeg := now.Format("2006/01/02")
	relPath := filepath.ToSlash(filepath.Join(userSeg, dateSeg, newFilename))

	dst := storageBaseDir() + "/" + relPath

	if err := ensureDirForAudio(filepath.Dir(dst)); err != nil {
		return "", err
	}

	if err := os.WriteFile(dst, data, 0644); err != nil {
		return "", err
	}

	return relPath, nil
}

func getWhisperConfig() WhisperConfig {
	return WhisperConfig{
		APIKey:  os.Getenv("OPENAI_API_KEY"),
		BaseURL: os.Getenv("OPENAI_BASE_URL"),
		Model:   os.Getenv("WHISPER_MODEL"),
	}
}

func callWhisperAPI(audioData []byte, filename string, config WhisperConfig, language, prompt string, temperature float64) (text string, duration float64, respLang string, err error) {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.openai.com/v1"
	}
	if config.Model == "" {
		config.Model = "whisper-1"
	}

	// 构建 multipart 请求
	body, contentType := buildMultipart(audioData, filename, config.Model, language, prompt, temperature)

	// 发送请求
	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest("POST", config.BaseURL+"/audio/transcriptions", body)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		err = readErr
		return
	}
	if resp.StatusCode != 200 {
		err = &APIError{Code: resp.StatusCode, Message: string(respBody)}
		return
	}

	// 解析响应
	var result WhisperResponse
	if err = json.Unmarshal(respBody, &result); err != nil {
		return
	}

	return result.Text, result.Duration, result.Language, nil
}

func buildMultipart(audioData []byte, filename, model, language, prompt string, temperature float64) (*bytes.Reader, string) {
	boundary := randomHex(16)
	var buf bytes.Buffer

	// 文件部分
	buf.WriteString("--" + boundary + "\r\n")
	buf.WriteString("Content-Disposition: form-data; name=\"file\"; filename=\"" + filename + "\"\r\n")
	buf.WriteString("Content-Type: audio/*\r\n\r\n")
	buf.Write(audioData)
	buf.WriteString("\r\n")

	// 模型
	buf.WriteString("--" + boundary + "\r\n")
	buf.WriteString("Content-Disposition: form-data; name=\"model\"\r\n\r\n")
	buf.WriteString(model + "\r\n")

	// 语言
	if language != "" {
		buf.WriteString("--" + boundary + "\r\n")
		buf.WriteString("Content-Disposition: form-data; name=\"language\"\r\n\r\n")
		buf.WriteString(language + "\r\n")
	}

	// 提示
	if prompt != "" {
		buf.WriteString("--" + boundary + "\r\n")
		buf.WriteString("Content-Disposition: form-data; name=\"prompt\"\r\n\r\n")
		buf.WriteString(prompt + "\r\n")
	}

	// 温度
	if temperature > 0 {
		buf.WriteString("--" + boundary + "\r\n")
		buf.WriteString("Content-Disposition: form-data; name=\"temperature\"\r\n\r\n")
		buf.WriteString(strconv.FormatFloat(temperature, 'f', -1, 64) + "\r\n")
	}

	buf.WriteString("--" + boundary + "--\r\n")

	return bytes.NewReader(buf.Bytes()), "multipart/form-data; boundary=" + boundary
}

func ensureDirForAudio(path string) error {
	return os.MkdirAll(path, 0755)
}

func getOptionalUserIDForAudio() *int {
	// 从上下文获取（简化实现）
	return nil
}

// WhisperResponse Whisper API 响应
type WhisperResponse struct {
	Text     string  `json:"text"`
	Duration float64 `json:"duration"`
	Language string  `json:"language"`
}

// APIError 错误
type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}
