package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"rbac-admin/internal/store"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CMSOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type CMSField struct {
	Name         string      `json:"name"`
	Label        string      `json:"label"`
	Type         string      `json:"type"`
	Required     bool        `json:"required,omitempty"`
	ReadOnly     bool        `json:"read_only,omitempty"`
	Wide         bool        `json:"wide,omitempty"`
	Help         string      `json:"help,omitempty"`
	Options      []CMSOption `json:"options,omitempty"`
	Relation     string      `json:"relation,omitempty"`
	RelationID   string      `json:"relation_id,omitempty"`
	RelationText string      `json:"relation_text,omitempty"`
}

type CMSConfig struct {
	Key         string     `json:"key"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Permission  string     `json:"permission"`
	Table       string     `json:"-"`
	Order       string     `json:"-"`
	Search      []string   `json:"-"`
	ListFields  []CMSField `json:"list_fields"`
	FormFields  []CMSField `json:"form_fields"`
	AllowCreate bool       `json:"allow_create"`
	AllowDelete bool       `json:"allow_delete"`
	New         func() any `json:"-"`
	NewSlice    func() any `json:"-"`
}

func textField(name, label string, required ...bool) CMSField {
	return CMSField{Name: name, Label: label, Type: "text", Required: len(required) > 0 && required[0]}
}
func imageField(name, label string) CMSField {
	return CMSField{Name: name, Label: label, Type: "image", Wide: true, Help: "可上传图片，也可直接填写已有图片 URL。"}
}
func imageValueField(name, label string) CMSField {
	return CMSField{Name: name, Label: label, Type: "image_value", Wide: true, Help: "值类型为 image 时可直接上传图片。"}
}
func imagesField(name, label string) CMSField {
	return CMSField{Name: name, Label: label, Type: "images", Wide: true, Help: "支持批量上传、调整顺序和填写替代文本。"}
}
func boolField(name, label string) CMSField {
	return CMSField{Name: name, Label: label, Type: "boolean"}
}
func numField(name, label string) CMSField { return CMSField{Name: name, Label: label, Type: "number"} }
func areaField(name, label string) CMSField {
	return CMSField{Name: name, Label: label, Type: "textarea", Wide: true}
}
func richField(name, label string) CMSField {
	return CMSField{Name: name, Label: label, Type: "richtext", Wide: true, Help: "支持 HTML，用于前台富文本内容。"}
}
func selectField(name, label string, values ...string) CMSField {
	options := make([]CMSOption, 0, len(values))
	for _, value := range values {
		options = append(options, CMSOption{Label: value, Value: value})
	}
	return CMSField{Name: name, Label: label, Type: "select", Options: options}
}
func relationField(name, label, resource, text string) CMSField {
	return CMSField{Name: name, Label: label, Type: "relation", Relation: resource, RelationID: "id", RelationText: text}
}
func listField(name, label, kind string) CMSField {
	return CMSField{Name: name, Label: label, Type: kind}
}

var cmsResources = map[string]CMSConfig{
	"navigation": {
		Key: "navigation", Title: "前台导航", Description: "统一维护英文站页头和页脚导航，避免栏目命名不一致。", Permission: "navigation", Table: "site_nav_items", Order: "location ASC, sort ASC, id ASC", Search: []string{"label", "path"}, AllowCreate: true, AllowDelete: true,
		ListFields: []CMSField{listField("label", "名称", "text"), listField("path", "路径", "code"), listField("location", "位置", "badge"), listField("parent_id", "上级 ID", "number"), listField("sort", "排序", "number"), listField("status", "状态", "boolean")},
		FormFields: []CMSField{relationField("parent_id", "上级导航", "navigation", "label"), textField("label", "英文名称", true), textField("path", "前台路径", true), selectField("location", "展示位置", "header", "footer"), boolField("open_new", "新窗口打开"), numField("sort", "排序"), boolField("status", "启用")},
		New:        func() any { return &store.NavItem{} }, NewSlice: func() any { return &[]store.NavItem{} },
	},
	"settings": {
		Key: "settings", Title: "站点设置", Description: "维护 Logo、联系方式、页脚、二维码和前台全局配置。", Permission: "setting", Table: "site_settings", Order: "sort ASC, id ASC", Search: []string{"label", "key", "value"}, AllowCreate: true, AllowDelete: true,
		ListFields: []CMSField{listField("group_name", "分组", "text"), listField("label", "名称", "text"), listField("key", "键", "code"), listField("value", "值", "truncate"), listField("status", "状态", "boolean")},
		FormFields: []CMSField{textField("group_name", "分组", true), textField("label", "显示名称", true), textField("key", "配置键", true), imageValueField("value", "配置值"), selectField("value_type", "值类型", "text", "url", "image", "email", "phone", "html"), numField("sort", "排序"), boolField("status", "启用")},
		New:        func() any { return &store.SiteSetting{} }, NewSlice: func() any { return &[]store.SiteSetting{} },
	},
	"banners": {
		Key: "banners", Title: "Banner 管理", Description: "维护首页轮播图、标题、说明和跳转链接。", Permission: "banner", Table: "site_banners", Order: "sort ASC, id DESC", Search: []string{"title", "subtitle"}, AllowCreate: true, AllowDelete: true,
		ListFields: []CMSField{listField("image_url", "图片", "image"), listField("title", "标题", "text"), listField("position", "位置", "text"), listField("sort", "排序", "number"), listField("status", "状态", "boolean")},
		FormFields: []CMSField{textField("title", "标题", true), textField("subtitle", "副标题"), imageField("image_url", "Banner 图片"), textField("link_url", "跳转链接"), selectField("position", "展示位置", "home", "inner"), numField("sort", "排序"), boolField("status", "启用")},
		New:        func() any { return &store.Banner{} }, NewSlice: func() any { return &[]store.Banner{} },
	},
	"benefits": {
		Key: "benefits", Title: "Why Depo", Description: "维护首页 Why Depo 卖点列表。", Permission: "benefit", Table: "site_benefits", Order: "sort ASC, id ASC", Search: []string{"title"}, AllowCreate: true, AllowDelete: true,
		ListFields: []CMSField{listField("title", "卖点", "text"), listField("icon", "图标", "text"), listField("sort", "排序", "number"), listField("status", "状态", "boolean")},
		FormFields: []CMSField{textField("title", "卖点内容", true), textField("icon", "图标名"), numField("sort", "排序"), boolField("status", "启用")},
		New:        func() any { return &store.Benefit{} }, NewSlice: func() any { return &[]store.Benefit{} },
	},
	"categories": {
		Key: "categories", Title: "产品分类", Description: "维护前台产品多级分类、首页推荐图和排序。", Permission: "category", Table: "product_categories", Order: "sort ASC, id ASC", Search: []string{"name", "slug", "summary"}, AllowCreate: true, AllowDelete: true,
		ListFields: []CMSField{listField("image_url", "图片", "image"), listField("name", "分类名称", "text"), listField("slug", "Slug", "code"), listField("parent_id", "上级 ID", "number"), listField("home_featured", "首页推荐", "boolean"), listField("status", "状态", "boolean")},
		FormFields: []CMSField{relationField("parent_id", "上级分类", "categories", "name"), textField("name", "分类名称", true), textField("slug", "Slug", true), imageField("image_url", "分类图片"), areaField("summary", "分类简介"), boolField("home_featured", "首页核心分类"), numField("sort", "排序"), boolField("status", "启用")},
		New:        func() any { return &store.ProductCategory{} }, NewSlice: func() any { return &[]store.ProductCategory{} },
	},
	"products": {
		Key: "products", Title: "产品管理", Description: "维护前台产品、主图、产品图集、规格参数、描述、热销和最新状态。", Permission: "product", Table: "products", Order: "sort ASC, id DESC", Search: []string{"name", "sku", "slug", "summary", "hot_tags"}, AllowCreate: true, AllowDelete: true,
		ListFields: []CMSField{listField("main_image", "主图", "image"), listField("name", "产品名称", "text"), listField("sku", "SKU", "code"), listField("category_id", "分类 ID", "number"), listField("featured", "热销", "boolean"), listField("status", "状态", "boolean")},
		FormFields: []CMSField{relationField("category_id", "产品分类", "categories", "name"), textField("name", "产品名称", true), textField("slug", "Slug", true), textField("sku", "SKU"), areaField("summary", "产品摘要"), richField("description", "产品详情"), imageField("main_image", "产品主图"), imagesField("images", "产品图集"), textField("video_url", "视频 URL"), textField("material", "Material"), textField("size", "Size"), textField("thread_standard", "Thread"), textField("pressure_rating", "Pressure"), textField("temperature_range", "Temperature"), textField("moq", "MOQ"), textField("standard", "Standard"), textField("application", "Application"), textField("hot_tags", "Hot Tags（逗号分隔）"), boolField("featured", "首页热销"), boolField("latest", "最新产品"), numField("sort", "排序"), boolField("status", "启用")},
		New:        func() any { return &store.Product{} }, NewSlice: func() any { return &[]store.Product{} },
	},
	"pages": {
		Key: "pages", Title: "单页管理", Description: "维护 About Us、OEM、Materials、Service、FAQ 和隐私政策。", Permission: "page", Table: "content_pages", Order: "sort ASC, id ASC", Search: []string{"title", "slug", "subtitle", "content"}, AllowCreate: true, AllowDelete: true,
		ListFields: []CMSField{listField("cover_image", "封面", "image"), listField("title", "标题", "text"), listField("slug", "Slug", "code"), listField("template", "模板", "text"), listField("status", "状态", "boolean")},
		FormFields: []CMSField{textField("title", "页面标题", true), textField("slug", "Slug", true), textField("subtitle", "副标题"), richField("content", "页面内容"), imageField("cover_image", "封面图片"), selectField("template", "页面模板", "standard", "about", "policy"), numField("sort", "排序"), boolField("status", "启用")},
		New:        func() any { return &store.ContentPage{} }, NewSlice: func() any { return &[]store.ContentPage{} },
	},
	"articles": {
		Key: "articles", Title: "文章管理", Description: "统一维护 News、Knowledge 和 Blog 内容。", Permission: "article", Table: "content_articles", Order: "published_at DESC, sort ASC, id DESC", Search: []string{"title", "slug", "summary", "category"}, AllowCreate: true, AllowDelete: true,
		ListFields: []CMSField{listField("cover_image", "封面", "image"), listField("title", "标题", "text"), listField("content_type", "类型", "badge"), listField("category", "分类", "text"), listField("published_at", "发布时间", "datetime"), listField("status", "状态", "boolean")},
		FormFields: []CMSField{selectField("content_type", "内容类型", "news", "knowledge", "blog"), textField("category", "内容分类", true), textField("title", "标题", true), textField("slug", "Slug", true), areaField("summary", "摘要"), richField("content", "正文"), imageField("cover_image", "封面图片"), CMSField{Name: "published_at", Label: "发布时间", Type: "datetime", Required: true}, boolField("featured", "推荐"), numField("sort", "排序"), boolField("status", "启用")},
		New:        func() any { return &store.ContentArticle{} }, NewSlice: func() any { return &[]store.ContentArticle{} },
	},
	"team": {
		Key: "team", Title: "团队成员", Description: "维护首页和 About Us 的团队成员、职位、邮箱与照片。", Permission: "team", Table: "team_members", Order: "sort ASC, id ASC", Search: []string{"name", "role", "email"}, AllowCreate: true, AllowDelete: true,
		ListFields: []CMSField{listField("image_url", "照片", "image"), listField("name", "姓名", "text"), listField("role", "职位", "text"), listField("email", "邮箱", "text"), listField("status", "状态", "boolean")},
		FormFields: []CMSField{textField("name", "姓名", true), textField("role", "职位", true), textField("email", "邮箱", true), imageField("image_url", "成员照片"), areaField("bio", "简介"), numField("sort", "排序"), boolField("status", "启用")},
		New:        func() any { return &store.TeamMember{} }, NewSlice: func() any { return &[]store.TeamMember{} },
	},
	"partners": {
		Key: "partners", Title: "合作伙伴", Description: "维护首页合作伙伴 Logo 墙。", Permission: "partner", Table: "partners", Order: "sort ASC, id ASC", Search: []string{"name"}, AllowCreate: true, AllowDelete: true,
		ListFields: []CMSField{listField("logo_url", "Logo", "image"), listField("name", "伙伴名称", "text"), listField("website_url", "网站", "truncate"), listField("sort", "排序", "number"), listField("status", "状态", "boolean")},
		FormFields: []CMSField{textField("name", "伙伴名称", true), imageField("logo_url", "Logo"), textField("website_url", "网站 URL"), numField("sort", "排序"), boolField("status", "启用")},
		New:        func() any { return &store.Partner{} }, NewSlice: func() any { return &[]store.Partner{} },
	},
	"industries": {
		Key: "industries", Title: "行业应用", Description: "维护首页 Products Industry Application 标签内容。", Permission: "industry", Table: "industries", Order: "sort ASC, id ASC", Search: []string{"title", "subtitle", "description"}, AllowCreate: true, AllowDelete: true,
		ListFields: []CMSField{listField("image_url", "图片", "image"), listField("title", "行业", "text"), listField("subtitle", "副标题", "truncate"), listField("sort", "排序", "number"), listField("status", "状态", "boolean")},
		FormFields: []CMSField{textField("title", "行业名称", true), textField("subtitle", "副标题"), areaField("description", "行业说明"), imageField("image_url", "行业图片"), textField("link_url", "查看详情链接"), numField("sort", "排序"), boolField("status", "启用")},
		New:        func() any { return &store.Industry{} }, NewSlice: func() any { return &[]store.Industry{} },
	},
	"videos": {
		Key: "videos", Title: "视频管理", Description: "维护工厂、产品测试、包装发货和安装演示视频。", Permission: "video", Table: "videos", Order: "sort ASC, id DESC", Search: []string{"title", "description"}, AllowCreate: true, AllowDelete: true,
		ListFields: []CMSField{listField("cover_url", "封面", "image"), listField("title", "标题", "text"), listField("video_url", "视频 URL", "truncate"), listField("sort", "排序", "number"), listField("status", "状态", "boolean")},
		FormFields: []CMSField{textField("title", "视频标题", true), areaField("description", "视频说明"), imageField("cover_url", "封面图片"), textField("video_url", "视频 URL", true), numField("sort", "排序"), boolField("status", "启用")},
		New:        func() any { return &store.Video{} }, NewSlice: func() any { return &[]store.Video{} },
	},
	"inquiries": {
		Key: "inquiries", Title: "询盘管理", Description: "查看前台询盘产品、附件和联系信息，并维护跟进状态。", Permission: "inquiry", Table: "inquiries", Order: "created_at DESC, id DESC", Search: []string{"name", "email", "phone", "message", "product_ids"}, AllowCreate: false, AllowDelete: true,
		ListFields: []CMSField{listField("name", "客户", "text"), listField("email", "邮箱", "text"), listField("phone", "Phone / WhatsApp", "text"), listField("message", "需求", "truncate"), listField("lead_status", "跟进状态", "badge"), listField("created_at", "提交时间", "datetime")},
		FormFields: []CMSField{{Name: "name", Label: "客户姓名", Type: "text", ReadOnly: true}, {Name: "email", Label: "邮箱", Type: "text", ReadOnly: true}, {Name: "phone", Label: "Phone / WhatsApp", Type: "text", ReadOnly: true}, {Name: "product_ids", Label: "询盘产品 ID", Type: "text", ReadOnly: true}, {Name: "attachment_url", Label: "附件 URL", Type: "url", ReadOnly: true, Wide: true}, {Name: "message", Label: "客户需求", Type: "textarea", ReadOnly: true, Wide: true}, selectField("lead_status", "跟进状态", "new", "contacted", "quoted", "won", "lost"), areaField("admin_note", "跟进备注")},
		New:        func() any { return &store.Inquiry{} }, NewSlice: func() any { return &[]store.Inquiry{} },
	},
	"subscribers": {
		Key: "subscribers", Title: "邮件订阅", Description: "管理 Blog 邮件订阅地址和启用状态。", Permission: "subscriber", Table: "email_subscribers", Order: "created_at DESC, id DESC", Search: []string{"email"}, AllowCreate: true, AllowDelete: true,
		ListFields: []CMSField{listField("email", "邮箱", "text"), listField("status", "状态", "boolean"), listField("created_at", "订阅时间", "datetime")},
		FormFields: []CMSField{textField("email", "邮箱", true), boolField("status", "启用")},
		New:        func() any { return &store.EmailSubscriber{} }, NewSlice: func() any { return &[]store.EmailSubscriber{} },
	},
}

func (a *App) cmsPermission(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg, ok := cmsResources[c.Param("resource")]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "未知内容资源"})
			c.Abort()
			return
		}
		if action == "create" && !cfg.AllowCreate {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "该资源不允许后台新增"})
			c.Abort()
			return
		}
		if action == "delete" && !cfg.AllowDelete {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "该资源不允许删除"})
			c.Abort()
			return
		}
		if !a.store.UserHasPermission(currentUserID(c), cfg.Permission+":"+action) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
			c.Abort()
			return
		}
		c.Set("cms_config", cfg)
		c.Next()
	}
}

func (a *App) cmsUploadPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg, ok := cmsResources[c.Param("resource")]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "未知内容资源"})
			c.Abort()
			return
		}
		hasImageField := false
		for _, field := range cfg.FormFields {
			if field.Type == "image" || field.Type == "images" || field.Type == "image_value" {
				hasImageField = true
				break
			}
		}
		userID := currentUserID(c)
		canUpload := a.store.UserHasPermission(userID, cfg.Permission+":create") ||
			a.store.UserHasPermission(userID, cfg.Permission+":update")
		if !hasImageField || !canUpload {
			c.JSON(http.StatusForbidden, gin.H{"error": "无图片上传权限"})
			c.Abort()
			return
		}
		c.Set("cms_config", cfg)
		c.Next()
	}
}

func (a *App) cmsPage(c *gin.Context) {
	cfg, ok := cmsResources[c.Param("resource")]
	if !ok {
		c.HTML(http.StatusNotFound, "error/forbidden", a.viewData(c, "页面不存在", ""))
		return
	}
	if !a.store.UserHasPermission(currentUserID(c), cfg.Permission+":view") {
		c.HTML(http.StatusForbidden, "error/forbidden", a.viewData(c, "无权限", ""))
		return
	}
	data := a.viewData(c, cfg.Title, "cms/"+cfg.Key)
	data["Resource"] = cfg.Key
	data["Description"] = cfg.Description
	c.HTML(http.StatusOK, "cms/index", data)
}

func cmsConfig(c *gin.Context) CMSConfig {
	value, _ := c.Get("cms_config")
	cfg, _ := value.(CMSConfig)
	return cfg
}

func (a *App) apiCMSMeta(c *gin.Context) {
	cfg := cmsConfig(c)
	userID := currentUserID(c)
	c.JSON(http.StatusOK, gin.H{"resource": cfg, "abilities": gin.H{
		"create": cfg.AllowCreate && a.store.UserHasPermission(userID, cfg.Permission+":create"),
		"update": a.store.UserHasPermission(userID, cfg.Permission+":update"),
		"delete": cfg.AllowDelete && a.store.UserHasPermission(userID, cfg.Permission+":delete"),
	}})
}

func (a *App) apiCMSIndex(c *gin.Context) {
	cfg := cmsConfig(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 500 {
		pageSize = 500
	}
	q := a.store.DB.Model(cfg.New())
	if keyword := strings.TrimSpace(c.Query("q")); keyword != "" && len(cfg.Search) > 0 {
		parts := make([]string, len(cfg.Search))
		args := make([]any, len(cfg.Search))
		for i, column := range cfg.Search {
			parts[i] = column + " LIKE ?"
			args[i] = "%" + keyword + "%"
		}
		q = q.Where("("+strings.Join(parts, " OR ")+")", args...)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, err)
		return
	}
	rows := cfg.NewSlice()
	if err := q.Order(cfg.Order).Offset((page - 1) * pageSize).Limit(pageSize).Find(rows).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": rows, "total": total, "page": page, "page_size": pageSize})
}

func (a *App) apiCMSShow(c *gin.Context) {
	cfg := cmsConfig(c)
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	row := cfg.New()
	query := a.store.DB
	if cfg.Key == "products" {
		query = query.Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort ASC, id ASC")
		})
	}
	if err := query.First(row, id).Error; err != nil {
		jsonError(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": row})
}

func (a *App) apiCMSCreate(c *gin.Context) {
	cfg := cmsConfig(c)
	values, payload, ok := bindCMSValues(c, cfg)
	if !ok {
		return
	}
	images, hasImages, err := productImagesFromPayload(cfg, payload)
	if err != nil {
		jsonMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	values["created_at"] = time.Now()
	values["updated_at"] = time.Now()
	err = a.store.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(cfg.Table).Create(&values).Error; err != nil {
			return err
		}
		if cfg.Key != "products" || !hasImages {
			return nil
		}
		var created struct {
			ID uint `gorm:"column:id"`
		}
		if err := tx.Table(cfg.Table).Select("id").Where("slug = ?", values["slug"]).Take(&created).Error; err != nil {
			return err
		}
		values["id"] = created.ID
		return replaceProductImages(tx, created.ID, images)
	})
	if err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": cfg.Title + "已创建", "item": values})
}

func (a *App) apiCMSUpdate(c *gin.Context) {
	cfg := cmsConfig(c)
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	values, payload, ok := bindCMSValues(c, cfg)
	if !ok {
		return
	}
	images, hasImages, err := productImagesFromPayload(cfg, payload)
	if err != nil {
		jsonMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	values["updated_at"] = time.Now()
	found := false
	err = a.store.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Table(cfg.Table).Where("id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return nil
		}
		found = true
		if err := tx.Table(cfg.Table).Where("id = ?", id).Updates(values).Error; err != nil {
			return err
		}
		if cfg.Key == "products" && hasImages {
			return replaceProductImages(tx, id, images)
		}
		return nil
	})
	if err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	if !found {
		jsonMessage(c, http.StatusNotFound, "记录不存在")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": cfg.Title + "已更新"})
}

func (a *App) apiCMSDelete(c *gin.Context) {
	cfg := cmsConfig(c)
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	err := a.store.DB.Transaction(func(tx *gorm.DB) error {
		if cfg.Key == "products" {
			if err := tx.Where("product_id = ?", id).Delete(&store.ProductImage{}).Error; err != nil {
				return err
			}
		}
		return tx.Table(cfg.Table).Where("id = ?", id).Delete(cfg.New()).Error
	})
	if err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "记录已删除"})
}

func bindCMSValues(c *gin.Context, cfg CMSConfig) (map[string]any, map[string]any, bool) {
	var payload map[string]any
	if err := c.ShouldBindJSON(&payload); err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return nil, nil, false
	}
	out := make(map[string]any)
	for _, field := range cfg.FormFields {
		if field.ReadOnly || field.Type == "images" {
			continue
		}
		value, exists := payload[field.Name]
		if !exists {
			continue
		}
		if field.Required && strings.TrimSpace(fmt.Sprint(value)) == "" {
			jsonMessage(c, http.StatusBadRequest, field.Label+"不能为空")
			return nil, nil, false
		}
		switch field.Type {
		case "number", "relation":
			switch n := value.(type) {
			case float64:
				out[field.Name] = int64(n)
			case string:
				parsed, _ := strconv.ParseInt(n, 10, 64)
				out[field.Name] = parsed
			default:
				out[field.Name] = n
			}
		case "boolean":
			out[field.Name] = value == true || value == float64(1) || fmt.Sprint(value) == "1"
		case "datetime":
			text := strings.TrimSpace(fmt.Sprint(value))
			if len(text) == 16 {
				text += ":00"
			}
			out[field.Name] = strings.Replace(text, "T", " ", 1)
		default:
			out[field.Name] = strings.TrimSpace(fmt.Sprint(value))
		}
	}
	return out, payload, true
}

type productImageInput struct {
	URL     string
	AltText string
}

func productImagesFromPayload(cfg CMSConfig, payload map[string]any) ([]productImageInput, bool, error) {
	if cfg.Key != "products" {
		return nil, false, nil
	}
	raw, exists := payload["images"]
	if !exists {
		return nil, false, nil
	}
	items, ok := raw.([]any)
	if !ok {
		return nil, true, fmt.Errorf("产品图集格式不正确")
	}
	if len(items) > 50 {
		return nil, true, fmt.Errorf("产品图集最多上传 50 张图片")
	}
	images := make([]productImageInput, 0, len(items))
	for _, item := range items {
		var image productImageInput
		switch value := item.(type) {
		case string:
			image.URL = strings.TrimSpace(value)
		case map[string]any:
			if imageURL, ok := value["image_url"]; ok && imageURL != nil {
				image.URL = strings.TrimSpace(fmt.Sprint(imageURL))
			}
			if altText, ok := value["alt_text"]; ok && altText != nil {
				image.AltText = strings.TrimSpace(fmt.Sprint(altText))
			}
		default:
			return nil, true, fmt.Errorf("产品图集格式不正确")
		}
		if image.URL != "" {
			images = append(images, image)
		}
	}
	return images, true, nil
}

func replaceProductImages(tx *gorm.DB, productID uint, images []productImageInput) error {
	if err := tx.Where("product_id = ?", productID).Delete(&store.ProductImage{}).Error; err != nil {
		return err
	}
	if len(images) == 0 {
		return nil
	}
	now := time.Now()
	rows := make([]store.ProductImage, 0, len(images))
	for index, image := range images {
		rows = append(rows, store.ProductImage{
			ProductID: productID,
			ImageURL:  image.URL,
			AltText:   image.AltText,
			Sort:      (index + 1) * 10,
			Status:    true,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	return tx.Create(&rows).Error
}

var safeFolder = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

func (a *App) apiImageUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		jsonMessage(c, http.StatusBadRequest, "请选择图片文件")
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedTypes := map[string]string{
		".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png",
		".gif": "image/gif", ".webp": "image/webp",
	}
	if file.Size <= 0 || file.Size > 15*1024*1024 || allowedTypes[ext] == "" {
		jsonMessage(c, http.StatusBadRequest, "仅支持 JPG、PNG、GIF、WebP，最大 15 MB")
		return
	}
	source, err := file.Open()
	if err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	header := make([]byte, 512)
	n, _ := source.Read(header)
	_ = source.Close()
	if n == 0 || http.DetectContentType(header[:n]) != allowedTypes[ext] {
		jsonMessage(c, http.StatusBadRequest, "文件内容与图片格式不匹配")
		return
	}
	cfg := cmsConfig(c)
	folder := cfg.Key
	dir := filepath.Join(a.cfg.UploadDir, folder)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		jsonError(c, http.StatusInternalServerError, err)
		return
	}
	base := strings.TrimSuffix(filepath.Base(file.Filename), ext)
	base = safeFolder.ReplaceAllString(base, "-")
	base = strings.Trim(base, "-")
	if base == "" {
		base = "image"
	}
	filename := fmt.Sprintf("%s-%d%s", base, time.Now().UnixNano(), ext)
	if err := c.SaveUploadedFile(file, filepath.Join(dir, filename)); err != nil {
		jsonError(c, http.StatusInternalServerError, err)
		return
	}
	url := "/uploads/" + folder + "/" + filename
	c.JSON(http.StatusCreated, gin.H{"message": "图片上传成功", "item": gin.H{"url": url}})
}
