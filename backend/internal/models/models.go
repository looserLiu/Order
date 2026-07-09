package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Email        string         `gorm:"uniqueIndex;size:255" json:"email"`
	Phone        string         `gorm:"uniqueIndex;size:20" json:"phone"`
	PasswordHash string         `gorm:"size:255" json:"-"`
	Nickname     string         `gorm:"size:100" json:"nickname"`
	AvatarURL    string         `gorm:"size:500" json:"avatar_url"`
	Currency     string         `gorm:"size:10;default:CNY" json:"currency"`
	Timezone     string         `gorm:"size:50;default:Asia/Shanghai" json:"timezone"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Accounts     []Account      `gorm:"foreignKey:UserID" json:"accounts,omitempty"`
	Categories   []Category     `gorm:"foreignKey:UserID" json:"categories,omitempty"`
	Transactions []Transaction  `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
	Budgets      []Budget       `gorm:"foreignKey:UserID" json:"budgets,omitempty"`
	Tags         []Tag          `gorm:"foreignKey:UserID" json:"tags,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type Account struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Type        string         `gorm:"size:20;not null" json:"type"` // bank, credit, cash, wechat, alipay, investment
	Balance     float64        `gorm:"default:0" json:"balance"`
	Currency    string         `gorm:"size:10;default:CNY" json:"currency"`
	Icon        string         `gorm:"size:100" json:"icon"`
	Color       string         `gorm:"size:20" json:"color"`
	IsDefault   bool           `gorm:"default:false" json:"is_default"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User         User          `gorm:"foreignKey:UserID" json:"-"`
	Transactions []Transaction `gorm:"foreignKey:AccountID" json:"-"`
}

func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

type Category struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;index" json:"user_id"`
	ParentID   *uuid.UUID     `gorm:"type:uuid" json:"parent_id"`
	Name       string         `gorm:"size:100;not null" json:"name"`
	Icon       string         `gorm:"size:100" json:"icon"`
	Color      string         `gorm:"size:20" json:"color"`
	Type       string         `gorm:"size:20;not null" json:"type"` // income, expense, transfer
	SortOrder  int            `gorm:"default:0" json:"sort_order"`
	IsSystem   bool           `gorm:"default:false" json:"is_system"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	User        User        `gorm:"foreignKey:UserID" json:"-"`
	Parent      *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:CategoryID" json:"-"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type Transaction struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	AccountID       uuid.UUID      `gorm:"type:uuid;index;not null" json:"account_id"`
	TargetAccountID *uuid.UUID     `gorm:"type:uuid" json:"target_account_id"`
	CategoryID      *uuid.UUID     `gorm:"type:uuid;index" json:"category_id"`
	Type            string         `gorm:"size:20;not null" json:"type"` // income, expense, transfer
	Amount          float64        `gorm:"not null" json:"amount"`
	Currency        string         `gorm:"size:10;default:CNY" json:"currency"`
	ExchangeRate    float64        `gorm:"default:1" json:"exchange_rate"`
	Tags            []uuid.UUID    `gorm:"-" json:"tags"`
	Merchant        string         `gorm:"size:255" json:"merchant"`
	Note            string         `gorm:"type:text" json:"note"`
	AttachmentURLs []string       `gorm:"-" json:"attachment_urls"`
	BillDate        time.Time      `gorm:"index;not null" json:"bill_date"`
	IsRecurring     bool           `gorm:"default:false" json:"is_recurring"`
	RecurringRule   string         `gorm:"type:jsonb" json:"recurring_rule"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	User          User          `gorm:"foreignKey:UserID" json:"-"`
	Account       Account       `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	TargetAccount Account       `gorm:"foreignKey:TargetAccountID" json:"target_account,omitempty"`
	Category      Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type Budget struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	CategoryID     *uuid.UUID     `gorm:"type:uuid" json:"category_id"`
	Amount          float64        `gorm:"not null" json:"amount"`
	Period          string         `gorm:"size:20;not null" json:"period"` // monthly, yearly
	StartDate       time.Time      `gorm:"not null" json:"start_date"`
	EndDate         *time.Time     `json:"end_date"`
	AlertThreshold  float64        `gorm:"default:0.8" json:"alert_threshold"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	User     User      `gorm:"foreignKey:UserID" json:"-"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (b *Budget) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type Tag struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Color     string         `gorm:"size:20" json:"color"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type Device struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	DeviceID    string    `gorm:"size:255;not null" json:"device_id"`
	DeviceName  string    `gorm:"size:100" json:"device_name"`
	DeviceType  string    `gorm:"size:20" json:"device_type"` // ios, android, web, desktop
	PushToken   string    `gorm:"size:500" json:"push_token"`
	LastLoginAt time.Time `json:"last_login_at"`
	CreatedAt   time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (d *Device) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

type Family struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	OwnerID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"owner_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Owner   User           `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members []FamilyMember `gorm:"foreignKey:FamilyID" json:"members,omitempty"`
}

func (f *Family) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

type FamilyMember struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	FamilyID  uuid.UUID `gorm:"type:uuid;index;not null" json:"family_id"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	Role      string    `gorm:"size:20;default:member" json:"role"` // owner, admin, member
	JoinedAt  time.Time `json:"joined_at"`

	Family Family `gorm:"foreignKey:FamilyID" json:"-"`
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (fm *FamilyMember) BeforeCreate(tx *gorm.DB) error {
	if fm.ID == uuid.Nil {
		fm.ID = uuid.New()
	}
	return nil
}

type FamilyTransaction struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	FamilyID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"family_id"`
	UserID     uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	AccountID  uuid.UUID      `gorm:"type:uuid;index;not null" json:"account_id"`
	CategoryID *uuid.UUID     `gorm:"type:uuid" json:"category_id"`
	Type       string         `gorm:"size:20;not null" json:"type"`
	Amount     float64        `gorm:"not null" json:"amount"`
	Note       string         `gorm:"type:text" json:"note"`
	BillDate   time.Time      `gorm:"index;not null" json:"bill_date"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Family   Family    `gorm:"foreignKey:FamilyID" json:"-"`
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Account  Account   `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Category Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (ft *FamilyTransaction) BeforeCreate(tx *gorm.DB) error {
	if ft.ID == uuid.Nil {
		ft.ID = uuid.New()
	}
	return nil
}

type AssetChange struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	AssetType    string         `gorm:"size:20;not null" json:"asset_type"` // debt_owed, debt_owing, investment
	RelatedUser  string         `gorm:"size:255" json:"related_user"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Amount      float64        `gorm:"not null" json:"amount"`
	InterestRate float64       `gorm:"default:0" json:"interest_rate"`
	StartDate   *time.Time    `json:"start_date"`
	EndDate     *time.Time    `json:"end_date"`
	Status      string         `gorm:"size:20;default:active" json:"status"` // active, settled
	Note        string         `gorm:"type:text" json:"note"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (ac *AssetChange) BeforeCreate(tx *gorm.DB) error {
	if ac.ID == uuid.Nil {
		ac.ID = uuid.New()
	}
	return nil
}

type Reminder struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	Title       string         `gorm:"size:100;not null" json:"title"`
	Content     string         `gorm:"type:text" json:"content"`
	RemindTime  time.Time      `gorm:"index;not null" json:"remind_time"`
	RepeatType  string         `gorm:"size:20;default:none" json:"repeat_type"` // none, daily, weekly, monthly, yearly
	CategoryID *uuid.UUID     `gorm:"type:uuid" json:"category_id"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (r *Reminder) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type Notification struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	Type      string         `gorm:"size:30;not null" json:"type"` // budget_alert, reminder, system
	Title     string         `gorm:"size:100;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	RelatedID *uuid.UUID     `gorm:"type:uuid" json:"related_id"`
	CreatedAt time.Time      `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

type AAGroup struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	TotalAmount float64        `gorm:"default:0" json:"total_amount"`
	Status      string         `gorm:"size:20;default:active" json:"status"` // active, settled
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User     User        `gorm:"foreignKey:UserID" json:"-"`
	Members  []AAMember `gorm:"foreignKey:GroupID" json:"members,omitempty"`
	Settlements []AASettlement `gorm:"foreignKey:GroupID" json:"settlements,omitempty"`
}

func (g *AAGroup) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

type AAMember struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	GroupID   uuid.UUID `gorm:"type:uuid;index;not null" json:"group_id"`
	UserID   *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	Name     string    `gorm:"size:50;not null" json:"name"`
	Paid     float64   `gorm:"default:0" json:"paid"`
	Owe      float64   `gorm:"default:0" json:"owe"`
	IsPaid   bool      `gorm:"default:false" json:"is_paid"`
	JoinedAt time.Time `json:"joined_at"`

	Group AAGroup `gorm:"foreignKey:GroupID" json:"-"`
}

func (m *AAMember) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

type AASettlement struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	GroupID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"group_id"`
	FromID    uuid.UUID      `gorm:"type:uuid;not null" json:"from_id"`
	ToID      uuid.UUID      `gorm:"type:uuid;not null" json:"to_id"`
	Amount    float64        `gorm:"not null" json:"amount"`
	Status    string         `gorm:"size:20;default:pending" json:"status"` // pending, completed
	SettledAt *time.Time    `json:"settled_at"`
	CreatedAt time.Time      `json:"created_at"`

	Group AAGroup `gorm:"foreignKey:GroupID" json:"-"`
}

func (s *AASettlement) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type FinancialGoal struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	TargetAmount float64       `gorm:"not null" json:"target_amount"`
	CurrentAmount float64      `gorm:"default:0" json:"current_amount"`
	Deadline    *time.Time    `json:"deadline"`
	Category    string         `gorm:"size:50" json:"category"` // savings, investment, debt, purchase
	Status      string         `gorm:"size:20;default:in_progress" json:"status"` // in_progress, completed, failed
	Note        string         `gorm:"type:text" json:"note"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (g *FinancialGoal) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

type Insurance struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	InsuranceType string      `gorm:"size:50" json:"insurance_type"` // health, life, car, property, travel
	Company     string         `gorm:"size:100" json:"company"`
	Premium     float64        `gorm:"default:0" json:"premium"`
	PaymentType string         `gorm:"size:20" json:"payment_type"` // yearly, quarterly, monthly
	StartDate   time.Time      `json:"start_date"`
	EndDate     *time.Time    `json:"end_date"`
	Coverage    float64        `gorm:"default:0" json:"coverage"`
	Beneficiary string         `gorm:"size:100" json:"beneficiary"`
	Note        string         `gorm:"type:text" json:"note"`
	Status      string         `gorm:"size:20;default:active" json:"status"` // active, expired, cancelled
	NextPaymentDate *time.Time `json:"next_payment_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (i *Insurance) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

type Backup struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;index;not null" json:"user_id"`
	FileName    string         `gorm:"size:255;not null" json:"file_name"`
	FileSize    int64          `gorm:"default:0" json:"file_size"`
	FileURL     string         `gorm:"size:500" json:"file_url"`
	BackupType  string         `gorm:"size:20" json:"backup_type"` // manual, auto
	Note        string         `gorm:"type:text" json:"note"`
	CreatedAt   time.Time      `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (b *Backup) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
