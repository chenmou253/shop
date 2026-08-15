package store

import "time"

type SiteSetting struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	GroupName string    `json:"group_name"`
	Label     string    `json:"label"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	ValueType string    `json:"value_type"`
	Sort      int       `json:"sort"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SiteSetting) TableName() string { return "site_settings" }

type NavItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ParentID  uint      `json:"parent_id"`
	Label     string    `json:"label"`
	Path      string    `json:"path"`
	Location  string    `json:"location"`
	OpenNew   bool      `json:"open_new"`
	Sort      int       `json:"sort"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Children  []NavItem `json:"children,omitempty" gorm:"-"`
}

func (NavItem) TableName() string { return "site_nav_items" }

type Banner struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle"`
	ImageURL  string    `json:"image_url"`
	LinkURL   string    `json:"link_url"`
	Position  string    `json:"position"`
	Sort      int       `json:"sort"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Banner) TableName() string { return "site_banners" }

type Benefit struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Icon      string    `json:"icon"`
	Sort      int       `json:"sort"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Benefit) TableName() string { return "site_benefits" }

type ProductCategory struct {
	ID           uint              `json:"id" gorm:"primaryKey"`
	ParentID     uint              `json:"parent_id"`
	Name         string            `json:"name"`
	Slug         string            `json:"slug"`
	ImageURL     string            `json:"image_url"`
	Summary      string            `json:"summary"`
	HomeFeatured bool              `json:"home_featured"`
	Sort         int               `json:"sort"`
	Status       bool              `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Children     []ProductCategory `json:"children,omitempty" gorm:"-"`
}

func (ProductCategory) TableName() string { return "product_categories" }

type Product struct {
	ID               uint            `json:"id" gorm:"primaryKey"`
	CategoryID       uint            `json:"category_id"`
	Name             string          `json:"name"`
	Slug             string          `json:"slug"`
	SKU              string          `json:"sku"`
	Summary          string          `json:"summary"`
	Description      string          `json:"description"`
	MainImage        string          `json:"main_image"`
	VideoURL         string          `json:"video_url"`
	Material         string          `json:"material"`
	Size             string          `json:"size"`
	ThreadStandard   string          `json:"thread_standard"`
	PressureRating   string          `json:"pressure_rating"`
	TemperatureRange string          `json:"temperature_range"`
	MOQ              string          `json:"moq"`
	Standard         string          `json:"standard"`
	Application      string          `json:"application"`
	HotTags          string          `json:"hot_tags"`
	Featured         bool            `json:"featured"`
	Latest           bool            `json:"latest"`
	Sort             int             `json:"sort"`
	Status           bool            `json:"status"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	Category         ProductCategory `json:"category" gorm:"foreignKey:CategoryID"`
	Images           []ProductImage  `json:"images" gorm:"foreignKey:ProductID"`
}

func (Product) TableName() string { return "products" }

type ProductImage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id"`
	ImageURL  string    `json:"image_url"`
	AltText   string    `json:"alt_text"`
	Sort      int       `json:"sort"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProductImage) TableName() string { return "product_images" }

type Page struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	Subtitle   string    `json:"subtitle"`
	Content    string    `json:"content"`
	CoverImage string    `json:"cover_image"`
	Template   string    `json:"template"`
	Sort       int       `json:"sort"`
	Status     bool      `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Page) TableName() string { return "content_pages" }

type Article struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ContentType string    `json:"content_type"`
	Category    string    `json:"category"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Summary     string    `json:"summary"`
	Content     string    `json:"content"`
	CoverImage  string    `json:"cover_image"`
	PublishedAt time.Time `json:"published_at"`
	Featured    bool      `json:"featured"`
	Sort        int       `json:"sort"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Article) TableName() string { return "content_articles" }

type TeamMember struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
	ImageURL  string    `json:"image_url"`
	Bio       string    `json:"bio"`
	Sort      int       `json:"sort"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TeamMember) TableName() string { return "team_members" }

type Partner struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name"`
	LogoURL    string    `json:"logo_url"`
	WebsiteURL string    `json:"website_url"`
	Sort       int       `json:"sort"`
	Status     bool      `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Partner) TableName() string { return "partners" }

type Industry struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	LinkURL     string    `json:"link_url"`
	Sort        int       `json:"sort"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Industry) TableName() string { return "industries" }

type Video struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CoverURL    string    `json:"cover_url"`
	VideoURL    string    `json:"video_url"`
	Sort        int       `json:"sort"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Video) TableName() string { return "videos" }

type Inquiry struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Message       string    `json:"message"`
	AttachmentURL string    `json:"attachment_url"`
	ProductIDs    string    `json:"product_ids"`
	Source        string    `json:"source"`
	LeadStatus    string    `json:"lead_status"`
	AdminNote     string    `json:"admin_note"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Inquiry) TableName() string { return "inquiries" }

type Subscriber struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Subscriber) TableName() string { return "email_subscribers" }
