package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ─────────────────────────────────────────────────────────────────────
// ENUMS
// ─────────────────────────────────────────────────────────────────────

type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleManager  UserRole = "manager"
	UserRoleCashier  UserRole = "cashier"
	UserRoleCustomer UserRole = "customer"
)

type OrderStatus string

const (
	OrderStatusNew            OrderStatus = "NEW"
	OrderStatusWaitingPayment OrderStatus = "WAITING_PAYMENT"
	OrderStatusCooking        OrderStatus = "COOKING"
	OrderStatusReady          OrderStatus = "READY"
	OrderStatusCompleted      OrderStatus = "COMPLETED"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type PosSyncStatus string

const (
	PosSyncStatusNotSent PosSyncStatus = "NOT_SENT"
	PosSyncStatusSent    PosSyncStatus = "SENT"
	PosSyncStatusError   PosSyncStatus = "ERROR"
)

type BillingStatus string

const (
	BillingStatusActive  BillingStatus = "ACTIVE"
	BillingStatusPastDue BillingStatus = "PAST_DUE"
	BillingStatusBlocked BillingStatus = "BLOCKED"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusSuccess   PaymentStatus = "SUCCESS"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
)

type ThemePreset string

const (
	ThemePresetLight    ThemePreset = "light"
	ThemePresetDark     ThemePreset = "dark"
	ThemePresetCoffee   ThemePreset = "coffee"
	ThemePresetFastfood ThemePreset = "fastfood"
)

type DiscountType string

const (
	DiscountTypePercent DiscountType = "PERCENT"
	DiscountTypeFixed   DiscountType = "FIXED"
)

type OrderType string

const (
	OrderTypeTakeaway OrderType = "takeaway"
	OrderTypeDineIn   OrderType = "dine_in"
	OrderTypeDelivery OrderType = "delivery"
)

// PaymentProvider — поддерживаемые платёжные провайдеры.
// При добавлении нового: добавить константу + реализовать интерфейс
// PaymentGateway в сервисном слое.
type PaymentProvider string

const (
	PaymentProviderTinkoff  PaymentProvider = "tinkoff"
	PaymentProviderYookassa PaymentProvider = "yookassa"
	PaymentProviderSBP      PaymentProvider = "sbp"
	PaymentProviderStripe   PaymentProvider = "stripe" // для будущих международных точек
)

// POSProvider — поддерживаемые кассовые / ресторанные системы.
type POSProvider string

const (
	POSProviderIiko    POSProvider = "iiko"
	POSProviderSyrve   POSProvider = "syrve" // бывший iiko cloud, отдельный API
	POSProviderRkeeper POSProvider = "rkeeper"
	POSProviderManual  POSProvider = "manual" // без интеграции, ручное управление меню
)

// ─────────────────────────────────────────────────────────────────────
// БЛОК 1: ОРГАНИЗАЦИЯ И БИЛЛИНГ
// ─────────────────────────────────────────────────────────────────────

// Organization — юридическое лицо / владелец аккаунта.
// Не хранит операционных данных — только регистрационная информация.
type Organization struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name     string    `gorm:"not null"                                        json:"name"`
	INN      *string   `gorm:"type:varchar(12)"                                json:"inn,omitempty"`
	IsActive bool      `gorm:"default:true"                                    json:"is_active"`
	// PlanID — текущий тарифный план. При смене плана обновляется здесь
	// и создаётся новая запись в organization_billing (для истории).
	PlanID    uuid.UUID `gorm:"type:uuid;not null" json:"plan_id"`
	CreatedAt time.Time `gorm:"autoCreateTime"     json:"created_at"`

	Plan    SubscriptionPlan      `gorm:"foreignKey:PlanID"         json:"plan,omitempty"`
	Venues  []Venue               `gorm:"foreignKey:OrganizationID" json:"venues,omitempty"`
	Billing []OrganizationBilling `gorm:"foreignKey:OrganizationID" json:"billing,omitempty"`
}

// Venue — торговая точка. Центральная сущность всей системы.
// Почти каждая таблица ссылается на venue_id.
type Venue struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index"                        json:"organization_id"`
	Name           string    `gorm:"not null"                                        json:"name"`
	// Slug используется в URL: brand-kazan.foodapp.ru
	// Только [a-z0-9-]. Менять осторожно — старые QR-коды перестанут работать.
	Slug      string   `gorm:"uniqueIndex;not null" json:"slug"`
	Address   *string  `gorm:"type:text"            json:"address,omitempty"`
	Latitude  *float64 `gorm:"type:numeric(10,7)"   json:"latitude,omitempty"`
	Longitude *float64 `gorm:"type:numeric(10,7)"   json:"longitude,omitempty"`

	// is_active = false → точка закрыта навсегда (PWA → 404).
	// is_emergency_stop = true → временная пауза (PWA → 503, 1 клик в админке).
	IsActive          bool    `gorm:"default:true"                 json:"is_active"`
	IsEmergencyStop   bool    `gorm:"default:false"                json:"is_emergency_stop"`
	MinOrderAmount    float64 `gorm:"type:numeric(12,2);default:0" json:"min_order_amount"`
	AvgCookingMinutes int     `gorm:"default:20"                   json:"avg_cooking_minutes"`

	// Расписание работы (оверрайд над POS).
	// Формат: {"mon":{"open":"09:00","close":"22:00"}, ...}. Null = берём из POS.
	WorkHoursJSON datatypes.JSON `gorm:"type:jsonb" json:"work_hours_json,omitempty"`

	// billing_status — кэш из organization_billing для быстрой проверки на каждый запрос.
	// Источник правды — organization_billing. Синхронизируется cron-задачей каждый час
	// и при каждом изменении статуса в organization_billing.
	BillingStatus BillingStatus `gorm:"type:varchar(20);default:'ACTIVE'" json:"billing_status"`
	PaidUntil     *time.Time    `                                         json:"paid_until,omitempty"`

	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `                       json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index"          json:"deleted_at,omitempty"`

	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
}

// OrganizationBilling — история биллинга.
// Смена плана = новая запись (старые не удалять).
// Текущий план = запись с MAX(created_at) для данной организации.
type OrganizationBilling struct {
	ID             uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrganizationID uuid.UUID     `gorm:"type:uuid;not null;index"                        json:"organization_id"`
	PlanID         uuid.UUID     `gorm:"type:uuid;not null"                              json:"plan_id"`
	BillingStatus  BillingStatus `gorm:"type:varchar(20);default:'ACTIVE'"               json:"billing_status"`
	// При оплате: paid_until = MAX(paid_until, now()) + 1 month.
	// Не затирать будущую дату при досрочной оплате.
	PaidUntil          *time.Time `json:"paid_until,omitempty"`
	PaymentMethodToken *string    `gorm:"type:text" json:"payment_method_token,omitempty"`
	LastInvoiceAt      *time.Time `json:"last_invoice_at,omitempty"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Organization Organization     `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	Plan         SubscriptionPlan `gorm:"foreignKey:PlanID"         json:"plan,omitempty"`
}

// SubscriptionPlan — тарифный план.
// При изменении цены создавать новый план, не обновлять существующий.
type SubscriptionPlan struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name          string    `gorm:"not null"                                        json:"name"`
	PricePerMonth float64   `gorm:"type:numeric(12,2);not null"                     json:"price_per_month"`
	// Лимиты и флаги плана.
	// Формат: {"max_products":100,"max_venues":1,"push_marketing":false}
	FeaturesJSON datatypes.JSON `gorm:"type:jsonb"     json:"features_json,omitempty"`
	IsActive     bool           `gorm:"default:true"   json:"is_active"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
}

// Language — справочник языков (заполняется seed-скриптом при деплое).
type Language struct {
	Code string `gorm:"type:varchar(5);primaryKey" json:"code"`
	Name string `gorm:"not null"                   json:"name"`
}

// ─────────────────────────────────────────────────────────────────────
// БЛОК 2: КОНФИГУРАЦИЯ PWA
// ─────────────────────────────────────────────────────────────────────

// AppConfig — настройки дизайна PWA для конкретной точки (1:1 с Venue).
// Создаётся автоматически при создании точки с дефолтными значениями.
type AppConfig struct {
	VenueID     uuid.UUID   `gorm:"type:uuid;primaryKey"              json:"venue_id"`
	ThemePreset ThemePreset `gorm:"type:varchar(20);default:'light'"  json:"theme_preset"`
	// HEX цвет кнопок и акцентов. Валидировать: /^#[0-9A-Fa-f]{6}$/
	AccentColor string  `gorm:"type:varchar(7);default:'#000000'" json:"accent_color"`
	LogoURL     *string `gorm:"type:text"                         json:"logo_url,omitempty"`
	BannerURL   *string `gorm:"type:text"                         json:"banner_url,omitempty"`
	// Перекрывает venues.address только в UI PWA.
	// venues.address по-прежнему используется для геокодинга.
	AddressManual *string    `gorm:"type:text" json:"address_manual,omitempty"`
	UpdatedAt     *time.Time `                  json:"updated_at,omitempty"`

	Venue Venue `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
}

// ─────────────────────────────────────────────────────────────────────
// БЛОК 3: ПОЛЬЗОВАТЕЛИ И АВТОРИЗАЦИЯ
// ─────────────────────────────────────────────────────────────────────

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Phone     string     `gorm:"uniqueIndex:idx_users_phone;not null"            json:"phone"`
	FullName  *string    `                                                       json:"full_name,omitempty"`
	BirthDate *time.Time `gorm:"type:date"                                       json:"birth_date,omitempty"`
	IsActive  bool       `gorm:"default:true"                                    json:"is_active"`
	CreatedAt time.Time  `gorm:"autoCreateTime"                                  json:"created_at"`
}

type UserPOSProfile struct {
	ID              uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"             json:"id"`
	UserID          uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex:idx_user_pos,priority:1"      json:"user_id"`
	VenueID         uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex:idx_user_pos,priority:2"      json:"venue_id"`
	Provider        POSProvider `gorm:"type:varchar(20);not null;uniqueIndex:idx_user_pos,priority:3" json:"provider"`
	ExternalGuestID string      `gorm:"not null"                                                    json:"external_guest_id"`
	CreatedAt       time.Time   `gorm:"autoCreateTime"                                              json:"created_at"`

	User  User  `gorm:"foreignKey:UserID"  json:"user,omitempty"`
	Venue Venue `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
}

// StaffRole — роль сотрудника в конкретной точке.
// Один пользователь может быть admin в одном заведении и manager в другом.
type StaffRole struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"           json:"id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_staff_roles_user_venue" json:"user_id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_staff_roles_user_venue" json:"venue_id"`
	Role    UserRole  `gorm:"type:varchar(20);not null"                                 json:"role"`
	// false = сотрудник уволен. НЕ удалять запись.
	IsActive  bool      `gorm:"default:true"   json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	User  User  `gorm:"foreignKey:UserID"  json:"user,omitempty"`
	Venue Venue `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
}

// UserSession — сессия пользователя.
// Клиенты: expires_at = now()+30d. Персонал: now()+8h.
// Продлевать при каждом запросе. Фоновый job чистит протухшие.
type UserSession struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index:idx_sessions_user"      json:"user_id"`
	// Случайный непрозрачный токен (не JWT). Можно мгновенно отозвать.
	Token string `gorm:"uniqueIndex:idx_sessions_token;not null" json:"token"`
	// Формат: {"ua":"Mozilla...","platform":"Android","pwa":true}
	DeviceInfo datatypes.JSON `gorm:"type:jsonb"                          json:"device_info,omitempty"`
	ExpiresAt  time.Time      `gorm:"index:idx_sessions_expires;not null" json:"expires_at"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"                      json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (UserSession) TableName() string { return "user_sessions" }

// ─────────────────────────────────────────────────────────────────────
// БЛОК 4: МЕНЮ И КАТАЛОГ
// ─────────────────────────────────────────────────────────────────────

// Category — категория меню точки.
type Category struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"                  json:"id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;index:idx_categories_venue_visible,priority:1" json:"venue_id"`
	// UUID категории в POS-системе. NULL = создана вручную в нашей админке.
	ExternalID *uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_categories_external_unique,priority:2"  json:"external_id,omitempty"`
	// Для вложенных категорий. NULL = верхний уровень.
	ParentID  *uuid.UUID `gorm:"type:uuid"                                                        json:"parent_id,omitempty"`
	SortOrder int        `gorm:"default:0"                                                        json:"sort_order"`
	IsVisible bool       `gorm:"default:true;index:idx_categories_venue_visible,priority:2"       json:"is_visible"`
	DeletedAt *time.Time `gorm:"index"                                                            json:"deleted_at,omitempty"`

	Venue        Venue                 `gorm:"foreignKey:VenueID"    json:"venue,omitempty"`
	Parent       *Category             `gorm:"foreignKey:ParentID"   json:"parent,omitempty"`
	Translations []CategoryTranslation `gorm:"foreignKey:CategoryID" json:"translations,omitempty"`
}

// CategoryTranslation — перевод названия категории.
type CategoryTranslation struct {
	CategoryID   uuid.UUID `gorm:"type:uuid;primaryKey"       json:"category_id"`
	LanguageCode string    `gorm:"type:varchar(5);primaryKey" json:"language_code"`
	Title        string    `gorm:"not null"                   json:"title"`

	Category Category `gorm:"foreignKey:CategoryID"                   json:"category,omitempty"`
	Language Language `gorm:"foreignKey:LanguageCode;references:Code" json:"language,omitempty"`
}

// Product — блюдо / товар в меню точки.
type Product struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"                        json:"id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;index:idx_products_pwa_main,priority:1"             json:"venue_id"`
	// ON DELETE RESTRICT — нельзя удалить категорию пока в ней есть товары.
	CategoryID uuid.UUID `gorm:"type:uuid;index:idx_products_pwa_main,priority:2"                    json:"category_id"`
	// UUID из POS-системы. Матчинг строго по UUID при синхронизации.
	// Нет в нашей БД → INSERT, есть → UPDATE, нет в POS → soft delete.
	ExternalID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_products_external_unique,priority:2" json:"external_id"`

	// Текущая цена из POS. Обновляется при синхронизации.
	// На цену в заказе не влияет — там снапшот в order_items.unit_price.
	Price       float64 `gorm:"type:numeric(12,2);not null"                              json:"price"`
	ImageURL    *string `gorm:"type:text"                                                json:"image_url,omitempty"`
	IsAvailable bool    `gorm:"default:true;index:idx_products_pwa_main,priority:3"      json:"is_available"`
	SortOrder   int     `gorm:"default:0;index:idx_products_pwa_main,priority:4"         json:"sort_order"`

	// Физические параметры и КБЖУ из POS.
	// Блок КБЖУ показываем только если calories IS NOT NULL.
	Weight   *int     `gorm:"type:integer"      json:"weight,omitempty"`
	Calories *float64 `gorm:"type:numeric(8,2)" json:"calories,omitempty"`
	Protein  *float64 `gorm:"type:numeric(8,2)" json:"protein,omitempty"`
	Fat      *float64 `gorm:"type:numeric(8,2)" json:"fat,omitempty"`
	Carbs    *float64 `gorm:"type:numeric(8,2)" json:"carbs,omitempty"`

	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	Venue        Venue                `gorm:"foreignKey:VenueID"    json:"venue,omitempty"`
	Category     Category             `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Translations []ProductTranslation `gorm:"foreignKey:ProductID"  json:"translations,omitempty"`
	Modifiers    []ProductModifier    `gorm:"foreignKey:ProductID"  json:"modifiers,omitempty"`
}

// ProductTranslation — перевод названия и описания товара.
// GIN-индекс для full-text search создаётся отдельной миграцией:
// CREATE INDEX idx_product_fts ON product_translations USING GIN (to_tsvector('russian', name));
type ProductTranslation struct {
	ProductID    uuid.UUID `gorm:"type:uuid;primaryKey"       json:"product_id"`
	LanguageCode string    `gorm:"type:varchar(5);primaryKey" json:"language_code"`
	Name         string    `gorm:"not null"                   json:"name"`
	Description  *string   `gorm:"type:text"                  json:"description,omitempty"`

	Product  Product  `gorm:"foreignKey:ProductID"                    json:"product,omitempty"`
	Language Language `gorm:"foreignKey:LanguageCode;references:Code" json:"language,omitempty"`
}

// ProductModifier — добавка / модификатор к товару.
type ProductModifier struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProductID  uuid.UUID  `gorm:"type:uuid;not null;index"                        json:"product_id"`
	ExternalID *uuid.UUID `gorm:"type:uuid"                                       json:"external_id,omitempty"`
	ExtraPrice float64    `gorm:"type:numeric(12,2);default:0"                    json:"extra_price"`
	// true = клиент не может добавить товар в корзину без выбора этого модификатора.
	IsRequired bool `gorm:"default:false" json:"is_required"`
	// 1 = чекбокс. >1 = можно выбрать несколько (напр. до 3 сиропов).
	MaxQuantity int `gorm:"default:1" json:"max_quantity"`
	// Одинаковый group_name = одна визуальная группа в UI. NULL = одиночный модификатор.
	GroupName *string `json:"group_name,omitempty"`

	Product      Product               `gorm:"foreignKey:ProductID"  json:"product,omitempty"`
	Translations []ModifierTranslation `gorm:"foreignKey:ModifierID" json:"translations,omitempty"`
}

// ModifierTranslation — перевод названия модификатора.
type ModifierTranslation struct {
	ModifierID   uuid.UUID `gorm:"type:uuid;primaryKey"       json:"modifier_id"`
	LanguageCode string    `gorm:"type:varchar(5);primaryKey" json:"language_code"`
	Name         string    `gorm:"not null"                   json:"name"`

	Modifier ProductModifier `gorm:"foreignKey:ModifierID"                   json:"modifier,omitempty"`
	Language Language        `gorm:"foreignKey:LanguageCode;references:Code" json:"language,omitempty"`
}

// ProductRecommendation — связь "С этим также берут".
// Настраивается вручную в админке.
// При запросе: проверять is_available = true — не показывать стоп-лист в upsell.
type ProductRecommendation struct {
	ProductID     uuid.UUID `gorm:"type:uuid;primaryKey;index:idx_recs_sort,priority:1" json:"product_id"`
	RecommendedID uuid.UUID `gorm:"type:uuid;primaryKey"                                json:"recommended_id"`
	SortOrder     int       `gorm:"default:0;index:idx_recs_sort,priority:2"           json:"sort_order"`

	Product     Product `gorm:"foreignKey:ProductID"     json:"product,omitempty"`
	Recommended Product `gorm:"foreignKey:RecommendedID" json:"recommended,omitempty"`
}

// ─────────────────────────────────────────────────────────────────────
// БЛОК 5: ЗАКАЗЫ
// ─────────────────────────────────────────────────────────────────────

// Order — заказ клиента.
// Статус менять только через сервисные методы — иначе пуши и логи не сработают.
// Заказы НИКОГДА не удалять физически — финансовая история.
type Order struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"                   json:"id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;index:idx_orders_venue_status_date,priority:1" json:"venue_id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;index:idx_orders_user_date,priority:1"         json:"user_id"`

	Status    OrderStatus `gorm:"type:varchar(20);default:'NEW';index:idx_orders_venue_status_date,priority:2" json:"status"`
	OrderType OrderType   `gorm:"type:varchar(20);default:'takeaway'"                                           json:"order_type"`
	// Итог после скидки. Фиксируется при создании заказа.
	TotalAmount float64 `gorm:"type:numeric(12,2);not null" json:"total_amount"`
	// NULL = "как можно скорее". Шаг 15 мин, текущий или следующий день.
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`

	Comment       *string `gorm:"type:text"     json:"comment,omitempty"`
	CutleryNeeded bool    `gorm:"default:false" json:"cutlery_needed"`

	CancelledReason *string    `gorm:"type:text" json:"cancelled_reason,omitempty"`
	RefundID        *string    `gorm:"type:text" json:"refund_id,omitempty"`
	RefundedAt      *time.Time `                 json:"refunded_at,omitempty"`

	// payment_id — денормализованный кэш для быстрого JOIN.
	// Источник правды — payments.order_id.
	// Порядок: создать order → создать payment → обновить payment_id транзакционно.
	PaymentID *uuid.UUID `gorm:"type:uuid;index:idx_orders_payment" json:"payment_id,omitempty"`

	PosOrderID    *uuid.UUID    `gorm:"type:uuid"                                                json:"pos_order_id,omitempty"`
	PosSyncStatus PosSyncStatus `gorm:"type:varchar(20);default:'NOT_SENT';index:idx_orders_sync" json:"pos_sync_status"`

	// Передаётся в POS как orderSource для маркетинговой аналитики.
	OrderSource string `gorm:"default:'PWA'" json:"order_source"`

	PromocodeID *uuid.UUID `gorm:"type:uuid" json:"promocode_id,omitempty"`
	// Хранится для аналитики. total_amount уже включает скидку.
	DiscountAmount float64 `gorm:"type:numeric(12,2);default:0" json:"discount_amount"`

	CreatedAt time.Time  `gorm:"autoCreateTime;index:idx_orders_venue_status_date,priority:3;index:idx_orders_user_date,priority:2" json:"created_at"`
	DeletedAt *time.Time `gorm:"index"                                                                                              json:"deleted_at,omitempty"`

	Venue Venue `gorm:"foreignKey:VenueID"               json:"venue,omitempty"`
	User  User  `gorm:"foreignKey:UserID"                json:"user,omitempty"`
	// Payment ищем через payments.order_id (не через orders.payment_id).
	// orders.payment_id — денормализованный кэш, не FK в смысле GORM association.
	Payment    *Payment         `gorm:"foreignKey:OrderID;references:ID" json:"payment,omitempty"`
	Promocode  *Promocode       `gorm:"foreignKey:PromocodeID"           json:"promocode,omitempty"`
	Items      []OrderItem      `gorm:"foreignKey:OrderID"               json:"items,omitempty"`
	StatusLogs []OrderStatusLog `gorm:"foreignKey:OrderID"               json:"status_logs,omitempty"`
}

// OrderStatusLog — аудит-лог смены статусов заказа.
// Позволяет восстановить хронологию и отлаживать зависания.
// ChangedBy = NULL — автоматическое изменение (воркер, webhook, таймаут).
type OrderStatusLog struct {
	ID        uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID   uuid.UUID    `gorm:"type:uuid;not null;index:idx_order_status_log"   json:"order_id"`
	OldStatus *OrderStatus `gorm:"type:varchar(20)"                                json:"old_status,omitempty"`
	NewStatus OrderStatus  `gorm:"type:varchar(20);not null"                       json:"new_status"`
	// NULL = автоматическое изменение (воркер, webhook). UUID сотрудника если ручное.
	ChangedBy *uuid.UUID `gorm:"type:uuid" json:"changed_by,omitempty"`
	Reason    *string    `gorm:"type:text" json:"reason,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Order Order `gorm:"foreignKey:OrderID"  json:"order,omitempty"`
	Staff *User `gorm:"foreignKey:ChangedBy" json:"staff,omitempty"`
}

func (OrderStatusLog) TableName() string { return "order_status_logs" }

// OrderItem — позиция в заказе (снапшот цены на момент оплаты).
type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index:idx_order_items_order"  json:"order_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"                              json:"product_id"`
	Quantity  int       `gorm:"not null"                                        json:"quantity"`
	// Снапшот цены на момент оплаты. Не меняется при изменении products.price.
	UnitPrice float64 `gorm:"type:numeric(12,2);not null" json:"unit_price"`
	// Снапшот названия — история не пустеет при удалении товара из меню.
	ProductNameSnapshot string `gorm:"not null" json:"product_name_snapshot"`

	Order     Order               `gorm:"foreignKey:OrderID"     json:"order,omitempty"`
	Product   Product             `gorm:"foreignKey:ProductID"   json:"product,omitempty"`
	Modifiers []OrderItemModifier `gorm:"foreignKey:OrderItemID" json:"modifiers,omitempty"`
}

// OrderItemModifier — модификатор в позиции заказа (снапшот на момент оплаты).
// Итог позиции: (unit_price + SUM(extra_price)) * quantity.
type OrderItemModifier struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderItemID uuid.UUID `gorm:"type:uuid;not null;index:idx_oim_order_item"     json:"order_item_id"`
	// NULLABLE — модификатор может быть удалён из меню. Всегда LEFT JOIN.
	ModifierID   *uuid.UUID `gorm:"type:uuid"                             json:"modifier_id,omitempty"`
	NameSnapshot string     `gorm:"not null"                              json:"name_snapshot"`
	ExtraPrice   float64    `gorm:"type:numeric(12,2);not null;default:0" json:"extra_price"`

	OrderItem OrderItem        `gorm:"foreignKey:OrderItemID" json:"order_item,omitempty"`
	Modifier  *ProductModifier `gorm:"foreignKey:ModifierID"  json:"modifier,omitempty"`
}

func (OrderItemModifier) TableName() string { return "order_item_modifiers" }

// ─────────────────────────────────────────────────────────────────────
// БЛОК 6: ПЛАТЕЖИ
// ─────────────────────────────────────────────────────────────────────

// Payment — платёж по заказу (1:1 с Order).
// Статус обновлять ТОЛЬКО через webhook от провайдера или polling.
// Никогда не менять вручную без подтверждения от платёжной системы.
type Payment struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"          json:"id"`
	OrderID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null;index:idx_payments_order"  json:"order_id"`
	// ID транзакции у провайдера. Используется для идентификации входящих webhook.
	ExternalID string        `gorm:"not null;uniqueIndex:idx_payments_external_unique"  json:"external_id"`
	Status     PaymentStatus `gorm:"type:varchar(20);default:'PENDING';index:idx_payments_status" json:"status"`
	Amount     float64       `gorm:"type:numeric(12,2);not null"                        json:"amount"`
	// Провайдер нужен при возврате — выбираем правильный API.
	// Значения: tinkoff, yookassa, sbp, stripe.
	Provider PaymentProvider `gorm:"type:varchar(30);not null" json:"provider"`
	// Момент фактического списания из данных webhook (не наше время).
	PaidAt    *time.Time `json:"paid_at,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Order Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// VenuePaymentSettings — настройки платёжного провайдера для конкретной точки (1:1).
//
// credentials_encrypted содержит провайдер-специфичный JSON (шифровать через pgcrypto/Vault):
//
//	tinkoff:  {"terminal_key":"...", "password":"..."}
//	yookassa: {"shop_id":"...", "secret_key":"..."}
//	sbp:      {"merchant_id":"...", "api_key":"..."}
//	stripe:   {"publishable_key":"...", "secret_key":"..."}
//
// Ключ шифрования — в env. Никогда не логировать расшифрованный payload.
type VenuePaymentSettings struct {
	ID       uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	VenueID  uuid.UUID       `gorm:"type:uuid;not null;uniqueIndex"                  json:"venue_id"`
	Provider PaymentProvider `gorm:"type:varchar(30);not null"                      json:"provider"`
	// Зашифрованный JSON с credentials провайдера.
	CredentialsEncrypted string `gorm:"type:text;not null" json:"credentials_encrypted"`
	// false = настройки сохранены но провайдер не используется (напр. при смене).
	IsActive  bool       `gorm:"default:true"   json:"is_active"`
	UpdatedAt *time.Time `                       json:"updated_at,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Venue Venue `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
}

func (VenuePaymentSettings) TableName() string { return "venue_payment_settings" }

// ─────────────────────────────────────────────────────────────────────
// БЛОК 7: НАСТРОЙКИ POS И ЛОГИ СИНХРОНИЗАЦИИ
// ─────────────────────────────────────────────────────────────────────

// VenuePOSSettings — настройки подключения POS-системы для точки (1:1).
// Абстрагирован от конкретной системы через поле provider.
//
// credentials_encrypted содержит провайдер-специфичный JSON (шифровать через pgcrypto/Vault):
//
//	iiko:    {"api_key":"...", "terminal_id":"uuid", "org_id":"uuid"}
//	syrve:   {"api_login":"...", "org_id":"uuid"}
//	rkeeper: {"ws_url":"...", "token":"..."}
//	manual:  {} — пустой объект, без интеграции
//
// Никогда не логировать расшифрованный payload.
type VenuePOSSettings struct {
	VenueID  uuid.UUID   `gorm:"type:uuid;primaryKey"      json:"venue_id"`
	Provider POSProvider `gorm:"type:varchar(20);not null" json:"provider"`
	// Зашифрованный JSON с credentials POS-системы.
	CredentialsEncrypted string `gorm:"type:text;not null" json:"credentials_encrypted"`
	// Интервал полной синхронизации меню (минуты).
	SyncIntervalMinutes int `gorm:"default:1"  json:"sync_interval_minutes"`
	// Интервал опроса стоп-листа (секунды).
	StoplistPollSeconds int `gorm:"default:60" json:"stoplist_poll_seconds"`

	// Если LastMenuSyncAt > 5 мин назад — алерт в Sentry.
	// NULL = синхронизации ещё не было.
	LastMenuSyncAt     *time.Time `json:"last_menu_sync_at,omitempty"`
	LastStoplistSyncAt *time.Time `json:"last_stoplist_sync_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`

	Venue Venue `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
}

func (VenuePOSSettings) TableName() string { return "venue_pos_settings" }

// PosSyncLog — лог одной попытки синхронизации заказа с POS-системой.
// TTL = 14 дней: pg_cron DELETE WHERE created_at < NOW() - INTERVAL '14 days'.
type PosSyncLog struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"                          json:"id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;index:idx_pos_logs_venue_date,priority:1"              json:"venue_id"`
	OrderID uuid.UUID `gorm:"type:uuid;not null;index:idx_pos_logs_order_date,priority:1"              json:"order_id"`
	// Провайдер POS-системы в момент синхронизации.
	// Позволяет различать логи iiko и rkeeper в одной таблице.
	Provider POSProvider `gorm:"type:varchar(20);not null" json:"provider"`
	// Логировать ДО отправки. НЕ включать credentials в payload.
	RequestPayload datatypes.JSON `gorm:"type:jsonb" json:"request_payload,omitempty"`
	// NULL если нет ответа (таймаут). При ERROR — главный источник для диагностики.
	ResponsePayload datatypes.JSON `gorm:"type:jsonb" json:"response_payload,omitempty"`
	// 200 = успех → SENT. 500/503 = retry до 3 раз. 401 = невалидный ключ. NULL = таймаут.
	HTTPStatusCode *int      `json:"http_status_code,omitempty"`
	CreatedAt      time.Time `gorm:"autoCreateTime;index:idx_pos_logs_venue_date,priority:2;index:idx_pos_logs_order_date,priority:2" json:"created_at"`

	Venue Venue `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
	Order Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

func (PosSyncLog) TableName() string { return "pos_sync_logs" }

// ─────────────────────────────────────────────────────────────────────
// БЛОК 8: ПРОМОКОДЫ
// ─────────────────────────────────────────────────────────────────────

// Promocode — промокод для точки.
// Нормализовывать code: UPPER(TRIM(code)) при записи и при поиске.
// SELECT FOR UPDATE при применении — защита от race condition на used_count.
type Promocode struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"                                                          json:"id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_promocodes_venue_code,priority:1;index:idx_promocodes_active,priority:1" json:"venue_id"`
	Code    string    `gorm:"not null;uniqueIndex:idx_promocodes_venue_code,priority:2"                                                json:"code"`

	DiscountType  DiscountType `gorm:"type:varchar(10);not null"    json:"discount_type"`
	DiscountValue float64      `gorm:"type:numeric(12,2);not null"  json:"discount_value"`
	// 0 = без ограничений по сумме заказа.
	MinOrderAmount float64 `gorm:"type:numeric(12,2);default:0" json:"min_order_amount"`
	// NULL = безлимит. INCREMENT в транзакции с SELECT FOR UPDATE.
	MaxUses   *int `json:"max_uses,omitempty"`
	UsedCount int  `gorm:"default:0;not null" json:"used_count"`

	// NULL = бессрочно.
	ExpiresAt *time.Time `gorm:"index:idx_promocodes_active,priority:3"               json:"expires_at,omitempty"`
	IsActive  bool       `gorm:"default:true;index:idx_promocodes_active,priority:2" json:"is_active"`
	CreatedAt time.Time  `gorm:"autoCreateTime"                                      json:"created_at"`

	Venue  Venue            `gorm:"foreignKey:VenueID"     json:"venue,omitempty"`
	Usages []PromocodeUsage `gorm:"foreignKey:PromocodeID" json:"usages,omitempty"`
}

// PromocodeUsage — лог использования промокода.
// Уникальный индекс (promocode_id, user_id) ограничивает промокод одним
// использованием на пользователя. Если промокод многоразовый — убрать этот индекс.
// Никогда не удалять записи — финансовая история.
type PromocodeUsage struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"                  json:"id"`
	PromocodeID uuid.UUID `gorm:"type:uuid;not null;index:idx_promo_usage_code"                    json:"promocode_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_promo_usage_user_code,priority:1" json:"user_id"`
	// Один заказ может использовать только один промокод (и наоборот).
	OrderID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;uniqueIndex:idx_promo_usage_user_code,priority:2" json:"order_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Promocode Promocode `gorm:"foreignKey:PromocodeID" json:"promocode,omitempty"`
	User      User      `gorm:"foreignKey:UserID"      json:"user,omitempty"`
	Order     Order     `gorm:"foreignKey:OrderID"     json:"order,omitempty"`
}

func (PromocodeUsage) TableName() string { return "promocode_usages" }

// ─────────────────────────────────────────────────────────────────────
// БЛОК 9: ОТЗЫВЫ И РЕЙТИНГИ
// ─────────────────────────────────────────────────────────────────────

// OrderReview — отзыв клиента на заказ (1:1 с Order).
// Видны только владельцу в админке — не публиковать в PWA.
// Предлагать оценку через пуш спустя 1 час после COMPLETED.
type OrderReview struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"             json:"id"`
	OrderID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"                              json:"order_id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;index:idx_reviews_user"                  json:"user_id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;index:idx_reviews_venue_date,priority:1" json:"venue_id"`
	// CHECK (rating BETWEEN 1 AND 5) — добавить в миграции.
	Rating     int16   `gorm:"not null"  json:"rating"`
	Comment    *string `gorm:"type:text" json:"comment,omitempty"`
	OwnerReply *string `gorm:"type:text" json:"owner_reply,omitempty"`
	// SET replied_at = now() вместе с owner_reply.
	RepliedAt *time.Time `json:"replied_at,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime;index:idx_reviews_venue_date,priority:2" json:"created_at"`

	Order Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	User  User  `gorm:"foreignKey:UserID"  json:"user,omitempty"`
	Venue Venue `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
}

// ─────────────────────────────────────────────────────────────────────
// БЛОК 10: PUSH-УВЕДОМЛЕНИЯ
// ─────────────────────────────────────────────────────────────────────

// PushSubscription — браузерная Web Push подписка клиента.
// Ставить is_active = false если браузер вернул 410 Gone при отправке.
type PushSubscription struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"          json:"id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;index:idx_push_user_venue,priority:1"  json:"user_id"`
	VenueID uuid.UUID `gorm:"type:uuid;not null;index:idx_push_user_venue,priority:2"  json:"venue_id"`
	// URL push-сервиса браузера (Google FCM, Apple и др.).
	Endpoint string `gorm:"type:text;not null" json:"endpoint"`
	// Публичный ключ и auth-секрет для web-push шифрования payload.
	P256DH    string    `gorm:"type:text;not null" json:"p256dh"`
	Auth      string    `gorm:"type:text;not null" json:"auth"`
	IsActive  bool      `gorm:"default:true"       json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime"     json:"created_at"`

	User  User  `gorm:"foreignKey:UserID"  json:"user,omitempty"`
	Venue Venue `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
}

// TelegramBinding — привязка аккаунта к Telegram-боту для уведомлений.
// Один пользователь — один Telegram (UNIQUE на user_id).
// Ставить is_active = false если пользователь написал /stop.
type TelegramBinding struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_tg_user"    json:"user_id"`
	// bigint — Telegram chat_id может быть > 2^31.
	TelegramChatID int64     `gorm:"not null;index:idx_tg_chat" json:"telegram_chat_id"`
	IsActive       bool      `gorm:"default:true"               json:"is_active"`
	CreatedAt      time.Time `gorm:"autoCreateTime"             json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
