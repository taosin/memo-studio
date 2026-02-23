package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// WhisperConfig 语音转文本配置
type WhisperConfig struct {
	APIKey    string
	BaseURL   string
	Model     string
	Timeout   time.Duration
}

// SpeechToTextRequest 语音转文本请求
type SpeechToTextRequest struct {
	Language   string `json:"language"`   // 可选：语言代码，如 "zh"
	Prompt     string `json:"prompt"`     // 可选：提示文本
	Temperature float64 `json:"temperature"` // 可选：温度（0-2）
}

// SpeechToTextResponse 语音转文本响应
type SpeechToTextResponse struct {
	Text     string  `json:"text"`
	Duration float64 `json:"duration"`
	Language string  `json:"language"`
}

var defaultConfig = WhisperConfig{
	BaseURL: "https://api.openai.com/v1",
	Model:   "whisper-1",
	Timeout: 60 * time.Second,
}

// SpeechToText 将语音转换为文本
func SpeechToText(audioData []byte, filename string, req SpeechToTextRequest) (*SpeechToTextResponse, error) {
	config := WhisperConfig{
		APIKey:  os.Getenv("OPENAI_API_KEY"),
		BaseURL: os.Getenv("OPENAI_BASE_URL"),
		Model:   os.Getenv("WHISPER_MODEL"),
		Timeout: defaultConfig.Timeout,
	}

	if config.APIKey == "" {
		return nil, nil // 未配置，不进行转录
	}
	if config.BaseURL == "" {
		config.BaseURL = defaultConfig.BaseURL
	}
	if config.Model == "" {
		config.Model = defaultConfig.Model
	}

	// 根据文件名检测 MIME 类型
	mimeType := detectMimeType(filename)
	form := bytes.NewBuffer(nil)
	writer := NewSectionsWriter(form)

	// 添加文件
	partName := "file"
	if strings.HasSuffix(filename, ".mp3") {
		partName = "audio.mp3"
	} else if strings.HasSuffix(filename, ".wav") {
		partName = "audio.wav"
	} else if strings.HasSuffix(filename, ".m4a") {
		partName = "audio.m4a"
	} else if strings.HasSuffix(filename, ".ogg") {
		partName = "audio.ogg"
	} else if strings.HasSuffix(filename, ".webm") {
		partName = "audio.webm"
	}

	writer.CreateFormFile(partName, filename, mimeType, audioData)

	// 添加模型
	writer.WriteField("model", config.Model)

	// 添加可选参数
	if req.Language != "" {
		writer.WriteField("language", req.Language)
	}
	if req.Prompt != "" {
		writer.WriteField("prompt", req.Prompt)
	}
	if req.Temperature > 0 {
		writer.WriteField("temperature", fmt.Sprint(req.Temperature))
	}

	// 添加响应格式
	writer.WriteField("response_format", "verbose_json")

	writer.Close()

	// 创建请求
	url := config.BaseURL + "/audio/transcriptions"
	httpReq, err := http.NewRequest("POST", url, form)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+config.APIKey)
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{Timeout: config.Timeout}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &APIError{Code: resp.StatusCode, Message: string(body)}
	}

	// 解析响应
	var result WhisperResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &SpeechToTextResponse{
		Text:     result.Text,
		Duration: result.Duration,
		Language: result.Language,
	}, nil
}

// WhisperResponse OpenAI Whisper API 响应
type WhisperResponse struct {
	Text     string  `json:"text"`
	Duration float64 `json:"duration"`
	Language string  `json:"language"`
}

// APIError API 错误
type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

// SectionsWriter io.Writer 的简单实现
type SectionsWriter struct {
	writer io.Writer
	boundary string
}

func NewSectionsWriter(w io.Writer) *SectionsWriter {
	return &SectionsWriter{
		writer:   w,
		boundary: "----WebKitFormBoundary" + randomHex(16),
	}
}

func (w *SectionsWriter) WriteField(name, value string) {
	w.writer.Write([]byte("--" + w.boundary + "\r\n"))
	w.writer.Write([]byte("Content-Disposition: form-data; name=\"" + name + "\"\r\n\r\n"))
	w.writer.Write([]byte(value + "\r\n"))
}

func (w *SectionsWriter) CreateFormFile(fieldName, filename, mimeType string, data []byte) {
	w.writer.Write([]byte("--" + w.boundary + "\r\n"))
	w.writer.Write([]byte("Content-Disposition: form-data; name=\"" + fieldName + "\"; filename=\"" + filename + "\"\r\n"))
	w.writer.Write([]byte("Content-Type: " + mimeType + "\r\n\r\n"))
	w.writer.Write(data)
	w.writer.Write([]byte("\r\n"))
}

func (w *SectionsWriter) Close() {
	w.writer.Write([]byte("--" + w.boundary + "--\r\n"))
}

func (w *SectionsWriter) FormDataContentType() string {
	return "multipart/form-data; boundary=" + w.boundary
}

func randomHex(n int) string {
	const letters = "0123456789abcdef"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%16]
	}
	return string(b)
}

func detectMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".m4a":
		return "audio/mp4"
	case ".ogg":
		return "audio/ogg"
	case ".webm":
		return "audio/webm"
	case ".mp4":
		return "audio/mp4"
	case ".flac":
		return "audio/flac"
	default:
		return "audio/*"
	}
}

// TranscribeWithFile 从文件转录
func TranscribeWithFile(filePath string, req SpeechToTextRequest) (*SpeechToTextResponse, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return SpeechToText(data, filepath.Base(filePath), req)
}

// TranscribeWithGin 从 Gin Context 转录
func TranscribeWithGin(c *gin.Context, filename string) (*SpeechToTextResponse, error) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var req SpeechToTextRequest
	if lang := c.PostForm("language"); lang != "" {
		req.Language = lang
	}
	if prompt := c.PostForm("prompt"); prompt != "" {
		req.Prompt = prompt
	}

	return SpeechToText(data, filename, req)
}
