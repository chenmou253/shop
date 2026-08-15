package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"hbfittings-front/internal/store"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (a *App) settings() (map[string]string, error) {
	var rows []store.SiteSetting
	if err := a.store.DB.Where("status = ?", true).Order("sort ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rows))
	for _, row := range rows {
		out[row.Key] = row.Value
	}
	return out, nil
}

func (a *App) bootstrap(c *gin.Context) {
	settings, err := a.settings()
	if err != nil {
		serverError(c, err)
		return
	}
	var categories []store.ProductCategory
	a.store.DB.Where("status = ?", true).Order("sort ASC, id ASC").Find(&categories)
	var navItems []store.NavItem
	a.store.DB.Where("status = ?", true).Order("sort ASC, id ASC").Find(&navItems)
	c.JSON(http.StatusOK, gin.H{"settings": settings, "categories": categoryTree(categories), "navigation": navTree(navItems)})
}

func (a *App) home(c *gin.Context) {
	settings, err := a.settings()
	if err != nil {
		serverError(c, err)
		return
	}
	var banners []store.Banner
	var benefits []store.Benefit
	var categories []store.ProductCategory
	var products []store.Product
	var partners []store.Partner
	var team []store.TeamMember
	var industries []store.Industry
	var articles []store.Article
	db := a.store.DB
	db.Where("status = ? AND position = ?", true, "home").Order("sort ASC, id ASC").Find(&banners)
	db.Where("status = ?", true).Order("sort ASC, id ASC").Find(&benefits)
	db.Where("status = ? AND home_featured = ?", true, true).Order("sort ASC, id ASC").Limit(3).Find(&categories)
	db.Preload("Category").Where("products.status = ? AND featured = ?", true, true).Order("products.sort ASC, products.id DESC").Limit(12).Find(&products)
	db.Where("status = ?", true).Order("sort ASC, id ASC").Find(&partners)
	db.Where("status = ?", true).Order("sort ASC, id ASC").Find(&team)
	db.Where("status = ?", true).Order("sort ASC, id ASC").Find(&industries)
	db.Where("status = ? AND content_type = ?", true, "news").Order("featured DESC, published_at DESC, sort ASC").Limit(4).Find(&articles)
	c.JSON(http.StatusOK, gin.H{
		"settings": settings, "banners": banners, "benefits": benefits,
		"featured_categories": categories, "hot_products": products, "partners": partners,
		"team": team, "industries": industries, "news": articles,
	})
}

func (a *App) categories(c *gin.Context) {
	var rows []store.ProductCategory
	if err := a.store.DB.Where("status = ?", true).Order("sort ASC, id ASC").Find(&rows).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categoryTree(rows)})
}

func categoryTree(rows []store.ProductCategory) []store.ProductCategory {
	byParent := make(map[uint][]store.ProductCategory)
	for _, row := range rows {
		byParent[row.ParentID] = append(byParent[row.ParentID], row)
	}
	var walk func(uint) []store.ProductCategory
	walk = func(parent uint) []store.ProductCategory {
		children := byParent[parent]
		for i := range children {
			children[i].Children = walk(children[i].ID)
		}
		return children
	}
	return walk(0)
}

func navTree(rows []store.NavItem) []store.NavItem {
	byParent := make(map[uint][]store.NavItem)
	for _, row := range rows {
		byParent[row.ParentID] = append(byParent[row.ParentID], row)
	}
	var walk func(uint) []store.NavItem
	walk = func(parent uint) []store.NavItem {
		children := byParent[parent]
		for i := range children {
			children[i].Children = walk(children[i].ID)
		}
		return children
	}
	return walk(0)
}

func (a *App) products(c *gin.Context) {
	page, pageSize := pagination(c)
	q := a.store.DB.Model(&store.Product{}).Where("products.status = ?", true)
	if category := strings.TrimSpace(c.Query("category")); category != "" {
		var cat store.ProductCategory
		if err := a.store.DB.Where("slug = ? AND status = ?", category, true).First(&cat).Error; err == nil {
			ids := a.categoryDescendantIDs(cat.ID)
			q = q.Where("category_id IN ?", ids)
		}
	}
	if keyword := strings.TrimSpace(c.Query("q")); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name LIKE ? OR summary LIKE ? OR material LIKE ? OR hot_tags LIKE ?", like, like, like, like)
	}
	filters := map[string]string{"material": "material", "size": "size", "standard": "standard", "application": "application"}
	for param, column := range filters {
		if value := strings.TrimSpace(c.Query(param)); value != "" {
			q = q.Where(column+" LIKE ?", "%"+value+"%")
		}
	}
	var total int64
	q.Count(&total)
	var rows []store.Product
	if err := q.Preload("Category").Order("products.sort ASC, products.id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": rows, "total": total, "page": page, "page_size": pageSize})
}

func (a *App) categoryDescendantIDs(root uint) []uint {
	ids := []uint{root}
	frontier := []uint{root}
	for len(frontier) > 0 {
		var children []store.ProductCategory
		if err := a.store.DB.Select("id").Where("parent_id IN ? AND status = ?", frontier, true).Find(&children).Error; err != nil {
			break
		}
		frontier = frontier[:0]
		for _, child := range children {
			ids = append(ids, child.ID)
			frontier = append(frontier, child.ID)
		}
	}
	return ids
}

func (a *App) product(c *gin.Context) {
	var product store.Product
	if err := a.store.DB.Preload("Category").Preload("Images", "status = ?", true).Where("slug = ? AND status = ?", c.Param("slug"), true).First(&product).Error; err != nil {
		notFound(c)
		return
	}
	var related []store.Product
	a.store.DB.Preload("Category").Where("status = ? AND category_id = ? AND id <> ?", true, product.CategoryID, product.ID).Order("featured DESC, sort ASC").Limit(4).Find(&related)
	var previous, next store.Product
	a.store.DB.Select("name", "slug").Where("status = ? AND id < ?", true, product.ID).Order("id DESC").First(&previous)
	a.store.DB.Select("name", "slug").Where("status = ? AND id > ?", true, product.ID).Order("id ASC").First(&next)
	c.JSON(http.StatusOK, gin.H{"product": product, "related": related, "previous": previous, "next": next})
}

func (a *App) articles(c *gin.Context) {
	page, pageSize := pagination(c)
	q := a.store.DB.Model(&store.Article{}).Where("status = ?", true)
	if contentType := strings.TrimSpace(c.Query("type")); contentType != "" {
		q = q.Where("content_type = ?", contentType)
	}
	if category := strings.TrimSpace(c.Query("category")); category != "" {
		q = q.Where("category = ?", category)
	}
	if keyword := strings.TrimSpace(c.Query("q")); keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("title LIKE ? OR summary LIKE ? OR content LIKE ?", like, like, like)
	}
	var total int64
	q.Count(&total)
	var rows []store.Article
	if err := q.Order("published_at DESC, sort ASC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": rows, "total": total, "page": page, "page_size": pageSize})
}

func (a *App) article(c *gin.Context) {
	var row store.Article
	if err := a.store.DB.Where("slug = ? AND status = ?", c.Param("slug"), true).First(&row).Error; err != nil {
		notFound(c)
		return
	}
	var latest []store.Article
	a.store.DB.Where("status = ? AND content_type = ? AND id <> ?", true, row.ContentType, row.ID).Order("published_at DESC").Limit(5).Find(&latest)
	c.JSON(http.StatusOK, gin.H{"article": row, "latest": latest})
}

func (a *App) page(c *gin.Context) {
	var row store.Page
	if err := a.store.DB.Where("slug = ? AND status = ?", c.Param("slug"), true).First(&row).Error; err != nil {
		notFound(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"page": row})
}

func (a *App) videos(c *gin.Context) {
	var rows []store.Video
	if err := a.store.DB.Where("status = ?", true).Order("sort ASC, id DESC").Find(&rows).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"videos": rows})
}

func (a *App) search(c *gin.Context) {
	keyword := strings.TrimSpace(c.Query("q"))
	if keyword == "" {
		c.JSON(http.StatusOK, gin.H{"products": []store.Product{}, "articles": []store.Article{}})
		return
	}
	like := "%" + keyword + "%"
	var products []store.Product
	var articles []store.Article
	a.store.DB.Preload("Category").Where("status = ? AND (name LIKE ? OR summary LIKE ? OR hot_tags LIKE ?)", true, like, like, like).Limit(20).Find(&products)
	a.store.DB.Where("status = ? AND (title LIKE ? OR summary LIKE ? OR content LIKE ?)", true, like, like, like).Limit(20).Find(&articles)
	c.JSON(http.StatusOK, gin.H{"products": products, "articles": articles})
}

func (a *App) createInquiry(c *gin.Context) {
	row := store.Inquiry{
		Name: strings.TrimSpace(c.PostForm("name")), Email: strings.TrimSpace(c.PostForm("email")),
		Phone: strings.TrimSpace(c.PostForm("phone")), Message: strings.TrimSpace(c.PostForm("message")),
		ProductIDs: strings.TrimSpace(c.PostForm("product_ids")), Source: strings.TrimSpace(c.PostForm("source")), LeadStatus: "new",
	}
	if row.Name == "" || row.Email == "" || row.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, email and message are required."})
		return
	}
	if file, err := c.FormFile("attachment"); err == nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowed := map[string]bool{".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true, ".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
		if !allowed[ext] || file.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Attachment must be an image or office/PDF file up to 10 MB."})
			return
		}
		dir := filepath.Join(a.uploadDir, "inquiries")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			serverError(c, err)
			return
		}
		name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		if err := c.SaveUploadedFile(file, filepath.Join(dir, name)); err != nil {
			serverError(c, err)
			return
		}
		row.AttachmentURL = "/uploads/inquiries/" + name
	}
	if err := a.store.DB.Create(&row).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Thank you. Our team will reply within 24 hours.", "id": row.ID})
}

func (a *App) subscribe(c *gin.Context) {
	var payload struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || strings.TrimSpace(payload.Email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid email address is required."})
		return
	}
	row := store.Subscriber{Email: strings.TrimSpace(payload.Email), Status: true}
	if err := a.store.DB.Where("email = ?", row.Email).FirstOrCreate(&row).Error; err != nil && err != gorm.ErrRecordNotFound {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Subscription confirmed."})
}

func pagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 12
	}
	if size > 100 {
		size = 100
	}
	return page, size
}

func serverError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
func notFound(c *gin.Context) { c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}) }
