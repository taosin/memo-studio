package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"memo-studio/backend/database"
	"memo-studio/backend/handlers"
	"memo-studio/backend/middleware"
	"memo-studio/backend/models"
	"memo-studio/backend/utils"

	"github.com/gin-gonic/gin"
)

func setup(t *testing.T) (router *gin.Engine, adminID int, storageDir string) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "notes.db")
	storageDir = filepath.Join(tmp, "storage")

	t.Setenv("MEMO_DB_PATH", dbPath)
	t.Setenv("MEMO_STORAGE_DIR", storageDir)
	t.Setenv("MEMO_ADMIN_PASSWORD", "AdminPass123!")

	if err := database.Init(); err != nil {
		t.Fatalf("database.Init: %v", err)
	}

	if err := database.DB.QueryRow(`SELECT id FROM users WHERE username='admin'`).Scan(&adminID); err != nil {
		t.Fatalf("query admin id: %v", err)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	// public auth routes
	public := r.Group("/api")
	{
		public.POST("/auth/login", handlers.Login)
		public.POST("/auth/register", handlers.Register)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/memos", handlers.ListMemos)
		api.POST("/memos", handlers.CreateMemo)
		api.PUT("/memos/:id", handlers.UpdateMemo)
		api.DELETE("/memos/:id", handlers.DeleteMemo)

		api.GET("/tags", handlers.GetTags)
		api.POST("/tags", handlers.CreateTag)

		api.GET("/review/random", handlers.RandomReview)

		api.POST("/resources", handlers.UploadResource)

		api.GET("/users/me", handlers.GetMe)
		api.PUT("/users/me", handlers.UpdateMe)
		api.PUT("/users/me/password", handlers.ChangeMyPassword)

		admin := api.Group("/users")
		admin.Use(middleware.AdminOnly())
		{
			admin.GET("", handlers.AdminListUsers)
			admin.POST("", handlers.AdminCreateUser)
			admin.PUT("/:id", handlers.AdminUpdateUser)
			admin.DELETE("/:id", handlers.AdminDeleteUser)
		}
	}

	return r, adminID, storageDir
}

func authHeader(t *testing.T, userID int, username string, isAdmin bool) string {
	t.Helper()
	token, err := utils.GenerateToken(userID, username, isAdmin)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	return "Bearer " + token
}

func doJSON(t *testing.T, r http.Handler, method, path string, auth string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode json: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func TestUploadResourceAndMemoCRUD(t *testing.T) {
	r, adminID, storageDir := setup(t)
	auth := authHeader(t, adminID, "admin", true)

	// 1) upload resource
	var mp bytes.Buffer
	w := multipart.NewWriter(&mp)
	fw, err := w.CreateFormFile("file", "hello.png")
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := io.Copy(fw, bytes.NewReader([]byte("fakepngdata"))); err != nil {
		t.Fatalf("write file: %v", err)
	}
	_ = w.Close()

	req := httptest.NewRequest("POST", "/api/resources", &mp)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", auth)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("upload status=%d body=%s", rr.Code, rr.Body.String())
	}
	var res models.Resource
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("unmarshal resource: %v", err)
	}
	if res.ID <= 0 || res.StoragePath == "" || res.URL == "" {
		t.Fatalf("bad resource: %+v", res)
	}
	if got := filepath.Join(storageDir, filepath.FromSlash(res.StoragePath)); func() bool {
		_, err := os.Stat(got)
		return err == nil
	}() == false {
		t.Fatalf("file not found on disk: %s", got)
	}

	// 2) create memo with tag + resource
	createBody := map[string]any{
		"title":        "t1",
		"content":      "c1 #tagA",
		"tags":         []string{"tagA"},
		"pinned":       true,
		"content_type": "markdown",
		"resource_ids": []int{res.ID},
	}
	rr = doJSON(t, r, "POST", "/api/memos", auth, createBody)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create memo status=%d body=%s", rr.Code, rr.Body.String())
	}
	var note models.Note
	if err := json.Unmarshal(rr.Body.Bytes(), &note); err != nil {
		t.Fatalf("unmarshal note: %v", err)
	}
	if note.ID <= 0 || note.ContentType != "markdown" || note.Pinned != true {
		t.Fatalf("bad note: %+v", note)
	}
	if len(note.Tags) != 1 || note.Tags[0].Name != "tagA" {
		t.Fatalf("bad tags: %+v", note.Tags)
	}
	if len(note.Resources) != 1 || note.Resources[0].ID != res.ID {
		t.Fatalf("bad resources: %+v", note.Resources)
	}

	// 3) list memos
	rr = doJSON(t, r, "GET", "/api/memos?limit=50&offset=0", auth, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("list memos status=%d body=%s", rr.Code, rr.Body.String())
	}
	var list []models.Note
	if err := json.Unmarshal(rr.Body.Bytes(), &list); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("expected non-empty list")
	}

	// 4) update memo
	updateBody := map[string]any{
		"title":        "t2",
		"content":      "c2",
		"tags":         []string{"tagB"},
		"pinned":       false,
		"content_type": "markdown",
		"resource_ids": []int{},
	}
	rr = doJSON(t, r, "PUT", "/api/memos/"+itoa(note.ID), auth, updateBody)
	if rr.Code != http.StatusOK {
		t.Fatalf("update memo status=%d body=%s", rr.Code, rr.Body.String())
	}
	var updated models.Note
	_ = json.Unmarshal(rr.Body.Bytes(), &updated)
	if updated.Title != "t2" || updated.Pinned != false {
		t.Fatalf("unexpected updated: %+v", updated)
	}

	// 5) delete memo
	rr = doJSON(t, r, "DELETE", "/api/memos/"+itoa(note.ID), auth, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("delete memo status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestUserIsolationForTags(t *testing.T) {
	r, adminID, _ := setup(t)
	adminAuth := authHeader(t, adminID, "admin", true)

	// create another user
	u2, err := models.CreateUser("u2user", "password1", "u2@example.com")
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	u2Auth := authHeader(t, u2.ID, u2.Username, false)

	// admin creates memo with tagX
	rr := doJSON(t, r, "POST", "/api/memos", adminAuth, map[string]any{
		"title":        "",
		"content":      "admin memo",
		"tags":         []string{"tagX"},
		"pinned":       false,
		"content_type": "markdown",
		"resource_ids": []int{},
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("admin create memo status=%d body=%s", rr.Code, rr.Body.String())
	}

	// u2 creates memo with same tagX, should create separate tag under u2
	rr = doJSON(t, r, "POST", "/api/memos", u2Auth, map[string]any{
		"title":        "",
		"content":      "u2 memo",
		"tags":         []string{"tagX"},
		"pinned":       false,
		"content_type": "markdown",
		"resource_ids": []int{},
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("u2 create memo status=%d body=%s", rr.Code, rr.Body.String())
	}

	// admin list tags should include tagX with user_id=admin
	rr = doJSON(t, r, "GET", "/api/tags", adminAuth, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("admin list tags status=%d body=%s", rr.Code, rr.Body.String())
	}
	var adminTags []models.Tag
	if err := json.Unmarshal(rr.Body.Bytes(), &adminTags); err != nil {
		t.Fatalf("unmarshal admin tags: %v", err)
	}
	if !containsTag(adminTags, "tagX") {
		t.Fatalf("admin tags missing tagX: %+v", adminTags)
	}

	// u2 list tags should also include tagX (its own)
	rr = doJSON(t, r, "GET", "/api/tags", u2Auth, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("u2 list tags status=%d body=%s", rr.Code, rr.Body.String())
	}
	var u2Tags []models.Tag
	if err := json.Unmarshal(rr.Body.Bytes(), &u2Tags); err != nil {
		t.Fatalf("unmarshal u2 tags: %v", err)
	}
	if !containsTag(u2Tags, "tagX") {
		t.Fatalf("u2 tags missing tagX: %+v", u2Tags)
	}

	// sanity: tags are distinct rows (different id) per user
	adminTagID := findTagID(adminTags, "tagX")
	u2TagID := findTagID(u2Tags, "tagX")
	if adminTagID == 0 || u2TagID == 0 || adminTagID == u2TagID {
		t.Fatalf("expected distinct tag IDs. admin=%d u2=%d", adminTagID, u2TagID)
	}
}

func TestAuthRegisterLoginAndChangePassword(t *testing.T) {
	r, _, _ := setup(t)

	// register
	rr := doJSON(t, r, "POST", "/api/auth/register", "", map[string]any{
		"username": "alice",
		"password": "alice123",
		"email":    "a@example.com",
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("register status=%d body=%s", rr.Code, rr.Body.String())
	}
	var reg struct {
		Token string      `json:"token"`
		User  models.User `json:"user"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &reg); err != nil {
		t.Fatalf("unmarshal register: %v", err)
	}
	if reg.Token == "" || reg.User.ID == 0 || reg.User.Username != "alice" {
		t.Fatalf("bad register resp: %+v", reg)
	}

	// login
	rr = doJSON(t, r, "POST", "/api/auth/login", "", map[string]any{
		"username": "alice",
		"password": "alice123",
	})
	if rr.Code != http.StatusOK {
		t.Fatalf("login status=%d body=%s", rr.Code, rr.Body.String())
	}
	var loginResp struct {
		Token string      `json:"token"`
		User  models.User `json:"user"`
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &loginResp)
	if loginResp.Token == "" || loginResp.User.ID == 0 {
		t.Fatalf("bad login resp: %s", rr.Body.String())
	}
	auth := "Bearer " + loginResp.Token

	// me
	rr = doJSON(t, r, "GET", "/api/users/me", auth, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("me status=%d body=%s", rr.Code, rr.Body.String())
	}

	// change password (wrong old)
	rr = doJSON(t, r, "PUT", "/api/users/me/password", auth, map[string]any{
		"old_password": "wrong",
		"new_password": "alice456",
	})
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for wrong old password, got %d body=%s", rr.Code, rr.Body.String())
	}

	// change password (ok)
	rr = doJSON(t, r, "PUT", "/api/users/me/password", auth, map[string]any{
		"old_password": "alice123",
		"new_password": "alice456",
	})
	if rr.Code != http.StatusOK {
		t.Fatalf("change password status=%d body=%s", rr.Code, rr.Body.String())
	}

	// login with new password should work
	rr = doJSON(t, r, "POST", "/api/auth/login", "", map[string]any{
		"username": "alice",
		"password": "alice456",
	})
	if rr.Code != http.StatusOK {
		t.Fatalf("login new pwd status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestAdminUserCRUD(t *testing.T) {
	r, adminID, _ := setup(t)
	adminAuth := authHeader(t, adminID, "admin", true)

	// create user
	rr := doJSON(t, r, "POST", "/api/users", adminAuth, map[string]any{
		"username": "bob",
		"password": "bob12345",
		"email":    "b@example.com",
		"is_admin": false,
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("admin create user status=%d body=%s", rr.Code, rr.Body.String())
	}
	var bob models.User
	_ = json.Unmarshal(rr.Body.Bytes(), &bob)
	if bob.ID == 0 || bob.Username != "bob" {
		t.Fatalf("bad created user: %s", rr.Body.String())
	}

	// list users
	rr = doJSON(t, r, "GET", "/api/users", adminAuth, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("admin list users status=%d body=%s", rr.Code, rr.Body.String())
	}

	// update user
	rr = doJSON(t, r, "PUT", "/api/users/"+itoa(bob.ID), adminAuth, map[string]any{
		"username": "bob2",
		"email":    "b2@example.com",
		"is_admin": true,
	})
	if rr.Code != http.StatusOK {
		t.Fatalf("admin update user status=%d body=%s", rr.Code, rr.Body.String())
	}

	// delete user
	rr = doJSON(t, r, "DELETE", "/api/users/"+itoa(bob.ID), adminAuth, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("admin delete user status=%d body=%s", rr.Code, rr.Body.String())
	}

	// non-admin cannot access
	u, err := models.CreateUser("u3user", "password1", "")
	if err != nil {
		t.Fatalf("CreateUser u3: %v", err)
	}
	uAuth := authHeader(t, u.ID, u.Username, false)
	rr = doJSON(t, r, "GET", "/api/users", uAuth, nil)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-admin list users, got %d body=%s", rr.Code, rr.Body.String())
	}
}

func containsTag(tags []models.Tag, name string) bool {
	for _, t := range tags {
		if t.Name == name {
			return true
		}
	}
	return false
}

func findTagID(tags []models.Tag, name string) int {
	for _, t := range tags {
		if t.Name == name {
			return t.ID
		}
	}
	return 0
}

func itoa(i int) string {
	// small helper to avoid strconv import spam in tests
	return fmt.Sprintf("%d", i)
}

