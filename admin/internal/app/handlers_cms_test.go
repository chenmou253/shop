package app

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRetiredCMSResourcesAreNotRegistered(t *testing.T) {
	for _, resource := range []string{"media", "product-images"} {
		if _, exists := cmsResources[resource]; exists {
			t.Fatalf("retired resource %q is still registered", resource)
		}
	}
}

func TestImageFieldsUploadInOwningResource(t *testing.T) {
	expected := map[string]map[string]string{
		"settings":   {"value": "image_value"},
		"banners":    {"image_url": "image"},
		"categories": {"image_url": "image"},
		"products":   {"main_image": "image", "images": "images"},
		"pages":      {"cover_image": "image"},
		"articles":   {"cover_image": "image"},
		"team":       {"image_url": "image"},
		"partners":   {"logo_url": "image"},
		"industries": {"image_url": "image"},
		"videos":     {"cover_url": "image"},
	}

	for resource, fields := range expected {
		cfg := cmsResources[resource]
		actual := make(map[string]string)
		for _, field := range cfg.FormFields {
			actual[field.Name] = field.Type
		}
		for name, fieldType := range fields {
			if actual[name] != fieldType {
				t.Errorf("%s.%s type = %q, want %q", resource, name, actual[name], fieldType)
			}
		}
	}
}

func TestProductImagesFromPayload(t *testing.T) {
	images, present, err := productImagesFromPayload(cmsResources["products"], map[string]any{
		"images": []any{
			map[string]any{"image_url": " /uploads/products/one.webp ", "alt_text": " One "},
			"https://example.com/two.jpg",
			map[string]any{"image_url": ""},
		},
	})
	if err != nil {
		t.Fatalf("parse images: %v", err)
	}
	if !present {
		t.Fatal("images payload should be present")
	}
	if len(images) != 2 {
		t.Fatalf("image count = %d, want 2", len(images))
	}
	if images[0].URL != "/uploads/products/one.webp" || images[0].AltText != "One" {
		t.Fatalf("unexpected first image: %#v", images[0])
	}
	if images[1].URL != "https://example.com/two.jpg" {
		t.Fatalf("unexpected second image: %#v", images[1])
	}
}

func TestAPIImageUploadStoresImageAndReturnsThumbnailURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var imageData bytes.Buffer
	picture := image.NewRGBA(image.Rect(0, 0, 2, 2))
	picture.Set(0, 0, color.RGBA{R: 255, A: 255})
	if err := png.Encode(&imageData, picture); err != nil {
		t.Fatalf("encode png: %v", err)
	}

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	filePart, err := writer.CreateFormFile("file", "thumbnail.png")
	if err != nil {
		t.Fatalf("create file part: %v", err)
	}
	if _, err := filePart.Write(imageData.Bytes()); err != nil {
		t.Fatalf("write file part: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart body: %v", err)
	}

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/cms/banners/upload", &requestBody)
	context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	context.Set("cms_config", cmsResources["banners"])

	uploadDir := t.TempDir()
	app := &App{cfg: Config{UploadDir: uploadDir}}
	app.apiImageUpload(context)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	var response struct {
		Item struct {
			URL string `json:"url"`
		} `json:"item"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !strings.HasPrefix(response.Item.URL, "/uploads/banners/") {
		t.Fatalf("unexpected image URL: %q", response.Item.URL)
	}
	storedFile := filepath.Join(uploadDir, "banners", filepath.Base(response.Item.URL))
	if _, err := os.Stat(storedFile); err != nil {
		t.Fatalf("uploaded image was not stored: %v", err)
	}
}
