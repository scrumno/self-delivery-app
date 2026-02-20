CREATE TABLE [organizations] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [name] nvarchar(255) NOT NULL,
  [inn] varchar(12),
  [is_active] boolean DEFAULT (true),
  [created_at] timestamptz DEFAULT (now()),
  [plan_id] uuid NOT NULL
)
GO

CREATE TABLE [venues] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [organization_id] uuid NOT NULL,
  [name] nvarchar(255) NOT NULL,
  [slug] nvarchar(255) UNIQUE NOT NULL,
  [address] text,
  [latitude] numeric(10,7),
  [longitude] numeric(10,7),
  [is_active] boolean DEFAULT (true),
  [is_emergency_stop] boolean DEFAULT (false),
  [min_order_amount] numeric(12,2) DEFAULT (0),
  [avg_cooking_minutes] integer DEFAULT (20),
  [work_hours_json] jsonb,
  [billing_status] nvarchar(255) NOT NULL CHECK ([billing_status] IN ('ACTIVE', 'PAST_DUE', 'BLOCKED')) DEFAULT 'ACTIVE',
  [paid_until] timestamptz,
  [created_at] timestamptz DEFAULT (now()),
  [updated_at] timestamptz,
  [deleted_at] timestamptz
)
GO

CREATE TABLE [organization_billing] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [organization_id] uuid NOT NULL,
  [plan_id] uuid NOT NULL,
  [billing_status] nvarchar(255) NOT NULL CHECK ([billing_status] IN ('ACTIVE', 'PAST_DUE', 'BLOCKED')) DEFAULT 'ACTIVE',
  [paid_until] timestamptz,
  [payment_method_token] text,
  [last_invoice_at] timestamptz,
  [created_at] timestamptz DEFAULT (now())
)
GO

CREATE TABLE [subscription_plans] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [name] nvarchar(255) NOT NULL,
  [price_per_month] numeric(12,2) NOT NULL,
  [features_json] jsonb,
  [is_active] boolean DEFAULT (true),
  [created_at] timestamptz DEFAULT (now())
)
GO

CREATE TABLE [languages] (
  [code] varchar(5) PRIMARY KEY,
  [name] nvarchar(255) NOT NULL
)
GO

CREATE TABLE [app_configs] (
  [venue_id] uuid PRIMARY KEY,
  [theme_preset] nvarchar(255) NOT NULL CHECK ([theme_preset] IN ('light', 'dark', 'coffee', 'fastfood')) DEFAULT 'light',
  [accent_color] varchar(7) DEFAULT '#000000',
  [logo_url] text,
  [banner_url] text,
  [address_manual] text,
  [updated_at] timestamptz
)
GO

CREATE TABLE [users] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [phone] nvarchar(255) UNIQUE NOT NULL,
  [full_name] nvarchar(255),
  [birth_date] date,
  [iiko_guest_id] uuid,
  [is_active] boolean DEFAULT (true),
  [created_at] timestamptz DEFAULT (now())
)
GO

CREATE TABLE [staff_roles] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [user_id] uuid NOT NULL,
  [venue_id] uuid NOT NULL,
  [role] nvarchar(255) NOT NULL CHECK ([role] IN ('admin', 'manager', 'cashier', 'customer')) NOT NULL,
  [is_active] boolean DEFAULT (true),
  [created_at] timestamptz DEFAULT (now())
)
GO

CREATE TABLE [user_sessions] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [user_id] uuid NOT NULL,
  [token] text UNIQUE NOT NULL,
  [device_info] jsonb,
  [expires_at] timestamptz NOT NULL,
  [created_at] timestamptz DEFAULT (now())
)
GO

CREATE TABLE [categories] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [venue_id] uuid NOT NULL,
  [external_id] uuid,
  [parent_id] uuid,
  [sort_order] integer DEFAULT (0),
  [is_visible] boolean DEFAULT (true),
  [deleted_at] timestamptz
)
GO

CREATE TABLE [category_translations] (
  [category_id] uuid,
  [language_code] nvarchar(255),
  [title] nvarchar(255) NOT NULL,
  PRIMARY KEY ([category_id], [language_code])
)
GO

CREATE TABLE [products] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [venue_id] uuid NOT NULL,
  [category_id] uuid,
  [external_id] uuid NOT NULL,
  [price] numeric(12,2) NOT NULL,
  [image_url] text,
  [is_available] boolean DEFAULT (true),
  [sort_order] integer DEFAULT (0),
  [weight] integer,
  [calories] numeric(8,2),
  [protein] numeric(8,2),
  [fat] numeric(8,2),
  [carbs] numeric(8,2),
  [deleted_at] timestamptz
)
GO

CREATE TABLE [product_translations] (
  [product_id] uuid,
  [language_code] nvarchar(255),
  [name] nvarchar(255) NOT NULL,
  [description] text,
  PRIMARY KEY ([product_id], [language_code])
)
GO

CREATE TABLE [product_modifiers] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [product_id] uuid NOT NULL,
  [external_id] uuid,
  [extra_price] numeric(12,2) DEFAULT (0),
  [is_required] boolean DEFAULT (false),
  [max_quantity] integer DEFAULT (1),
  [group_name] nvarchar(255)
)
GO

CREATE TABLE [modifier_translations] (
  [modifier_id] uuid,
  [language_code] nvarchar(255),
  [name] nvarchar(255) NOT NULL,
  PRIMARY KEY ([modifier_id], [language_code])
)
GO

CREATE TABLE [product_recommendations] (
  [product_id] uuid,
  [recommended_id] uuid,
  [sort_order] integer DEFAULT (0),
  PRIMARY KEY ([product_id], [recommended_id])
)
GO

CREATE TABLE [orders] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [venue_id] uuid NOT NULL,
  [user_id] uuid NOT NULL,
  [status] nvarchar(255) NOT NULL CHECK ([status] IN ('NEW', 'WAITING_PAYMENT', 'COOKING', 'READY', 'COMPLETED', 'CANCELLED')) DEFAULT 'NEW',
  [order_type] nvarchar(255) NOT NULL CHECK ([order_type] IN ('takeaway', 'dine_in', 'delivery')) DEFAULT 'takeaway',
  [total_amount] numeric(12,2) NOT NULL,
  [scheduled_at] timestamptz,
  [comment] text,
  [cutlery_needed] boolean DEFAULT (false),
  [cancelled_reason] text,
  [refund_id] text,
  [refunded_at] timestamptz,
  [payment_id] uuid,
  [pos_order_id] uuid,
  [pos_sync_status] nvarchar(255) NOT NULL CHECK ([pos_sync_status] IN ('NOT_SENT', 'SENT', 'ERROR')) DEFAULT 'NOT_SENT',
  [order_source] nvarchar(255) DEFAULT 'PWA',
  [promocode_id] uuid,
  [discount_amount] numeric(12,2) DEFAULT (0),
  [created_at] timestamptz DEFAULT (now()),
  [deleted_at] timestamptz
)
GO

CREATE TABLE [order_items] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [order_id] uuid NOT NULL,
  [product_id] uuid NOT NULL,
  [quantity] integer NOT NULL,
  [unit_price] numeric(12,2) NOT NULL,
  [product_name_snapshot] nvarchar(255)
)
GO

CREATE TABLE [order_item_modifiers] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [order_item_id] uuid NOT NULL,
  [modifier_id] uuid,
  [name_snapshot] nvarchar(255) NOT NULL,
  [extra_price] numeric(12,2) NOT NULL DEFAULT (0)
)
GO

CREATE TABLE [payments] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [order_id] uuid UNIQUE NOT NULL,
  [external_id] nvarchar(255) NOT NULL,
  [status] nvarchar(255) NOT NULL CHECK ([status] IN ('PENDING', 'SUCCESS', 'FAILED', 'REFUNDED', 'CANCELLED')) DEFAULT 'PENDING',
  [amount] numeric(12,2) NOT NULL,
  [provider] nvarchar(255),
  [paid_at] timestamptz,
  [created_at] timestamptz DEFAULT (now())
)
GO

CREATE TABLE [pos_sync_logs] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [venue_id] uuid NOT NULL,
  [order_id] uuid NOT NULL,
  [request_payload] jsonb,
  [response_payload] jsonb,
  [http_status_code] integer,
  [created_at] timestamptz DEFAULT (now())
)
GO

CREATE TABLE [pos_settings] (
  [venue_id] uuid PRIMARY KEY,
  [api_key_encrypted] text NOT NULL,
  [terminal_id] uuid NOT NULL,
  [external_org_id] uuid NOT NULL,
  [sync_interval_minutes] integer DEFAULT (1),
  [stoplist_poll_seconds] integer DEFAULT (60),
  [last_menu_sync_at] timestamptz,
  [last_stoplist_sync_at] timestamptz,
  [updated_at] timestamptz
)
GO

CREATE TABLE [promocodes] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [venue_id] uuid NOT NULL,
  [code] nvarchar(255) NOT NULL,
  [discount_type] nvarchar(255) NOT NULL CHECK ([discount_type] IN ('PERCENT', 'FIXED')) NOT NULL,
  [discount_value] numeric(12,2) NOT NULL,
  [min_order_amount] numeric(12,2) DEFAULT (0),
  [max_uses] integer,
  [used_count] integer NOT NULL DEFAULT (0),
  [expires_at] timestamptz,
  [is_active] boolean DEFAULT (true),
  [created_at] timestamptz DEFAULT (now())
)
GO

CREATE TABLE [order_reviews] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [order_id] uuid UNIQUE NOT NULL,
  [user_id] uuid NOT NULL,
  [venue_id] uuid NOT NULL,
  [rating] smallint NOT NULL,
  [comment] text,
  [owner_reply] text,
  [replied_at] timestamptz,
  [created_at] timestamptz DEFAULT (now())
)
GO

CREATE TABLE [push_subscriptions] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [user_id] uuid NOT NULL,
  [venue_id] uuid NOT NULL,
  [endpoint] text NOT NULL,
  [p256dh] text NOT NULL,
  [auth] text NOT NULL,
  [is_active] boolean DEFAULT (true),
  [created_at] timestamptz DEFAULT (now())
)
GO

CREATE TABLE [telegram_bindings] (
  [id] uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  [user_id] uuid NOT NULL,
  [telegram_chat_id] bigint NOT NULL,
  [is_active] boolean DEFAULT (true),
  [created_at] timestamptz DEFAULT (now())
)
GO

CREATE INDEX [idx_users_phone] ON [users] ("phone")
GO

CREATE UNIQUE INDEX [idx_staff_roles_user_venue] ON [staff_roles] ("user_id", "venue_id")
GO

CREATE INDEX [idx_sessions_token] ON [user_sessions] ("token")
GO

CREATE INDEX [idx_sessions_user] ON [user_sessions] ("user_id")
GO

CREATE INDEX [idx_sessions_expires] ON [user_sessions] ("expires_at")
GO

CREATE INDEX [idx_categories_venue_visible] ON [categories] ("venue_id", "is_visible")
GO

CREATE UNIQUE INDEX [idx_categories_external_unique] ON [categories] ("venue_id", "external_id")
GO

CREATE INDEX [idx_products_pwa_main] ON [products] ("venue_id", "category_id", "is_available", "sort_order")
GO

CREATE UNIQUE INDEX [idx_products_external_unique] ON [products] ("venue_id", "external_id")
GO

CREATE INDEX [idx_recs_sort] ON [product_recommendations] ("product_id", "sort_order")
GO

CREATE INDEX [idx_orders_venue_status_date] ON [orders] ("venue_id", "status", "created_at")
GO

CREATE INDEX [idx_orders_user_date] ON [orders] ("user_id", "created_at")
GO

CREATE INDEX [idx_orders_payment] ON [orders] ("payment_id")
GO

CREATE INDEX [idx_orders_sync] ON [orders] ("pos_sync_status")
GO

CREATE INDEX [idx_order_items_order] ON [order_items] ("order_id")
GO

CREATE INDEX [idx_oim_order_item] ON [order_item_modifiers] ("order_item_id")
GO

CREATE INDEX [idx_payments_order] ON [payments] ("order_id")
GO

CREATE UNIQUE INDEX [idx_payments_external_unique] ON [payments] ("external_id")
GO

CREATE INDEX [idx_payments_status] ON [payments] ("status")
GO

CREATE INDEX [idx_pos_logs_order_date] ON [pos_sync_logs] ("order_id", "created_at")
GO

CREATE INDEX [idx_pos_logs_venue_date] ON [pos_sync_logs] ("venue_id", "created_at")
GO

CREATE UNIQUE INDEX [idx_promocodes_venue_code] ON [promocodes] ("venue_id", "code")
GO

CREATE INDEX [idx_promocodes_active] ON [promocodes] ("venue_id", "is_active", "expires_at")
GO

CREATE INDEX [idx_reviews_venue_date] ON [order_reviews] ("venue_id", "created_at")
GO

CREATE INDEX [idx_reviews_user] ON [order_reviews] ("user_id")
GO

CREATE INDEX [idx_push_user_venue] ON [push_subscriptions] ("user_id", "venue_id")
GO

CREATE UNIQUE INDEX [idx_tg_user] ON [telegram_bindings] ("user_id")
GO

CREATE INDEX [idx_tg_chat] ON [telegram_bindings] ("telegram_chat_id")
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Юридическое лицо. Одна организация может иметь несколько точек (venues). Не хранит операционных данных.',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organizations';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID юрлица. Передаётся в JWT сотрудника. Все запросы персонала фильтруются через него',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organizations',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Название для договора. Официальное название для договоров и инвойсов. Не путать с venues.name — публичным названием',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organizations',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ИНН для проверки контрагента (РФ). ИНН: 10 цифр для ООО, 12 для ИП. Валидировать контрольную сумму. NULL разрешён в MVP',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organizations',
@level2type = N'Column', @level2name = 'inn';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Глобальное выключение аккаунта. false = ручная блокировка саппортом. Все venues этой организации перестают работать',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organizations',
@level2type = N'Column', @level2name = 'is_active';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Дата регистрации. Только для чтения. Используется в аналитике роста клиентской базы',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organizations',
@level2type = N'Column', @level2name = 'created_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK → subscription_plans.id',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organizations',
@level2type = N'Column', @level2name = 'plan_id';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Центральная сущность системы. Торговая точка. Один organization может иметь много venues. Вся iiko/PWA логика привязана к venue_id.',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID точки. Главный FK всей системы — почти каждая таблица ссылается на venue_id',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Владелец точки. При авторизации: staff_roles.venue_id → venues → organizations → проверяем billing_status',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'organization_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Публичное название точки для клиентов и шапки PWA. Напр. "Кофе-пойнт ТЦ Мега"',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'URL-часть: brand-kazan.foodapp.ru. Только [a-z0-9-], валидировать регуляркой. Менять осторожно — старые QR-коды перестанут работать',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'slug';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Физический адрес. Используется для геокодинга. Для отображения в PWA может быть перекрыт app_configs.address_manual',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'address';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Координаты для карты. GPS-широта. Заполнять через геокодер при сохранении адреса. Нужна для карты и функции "рядом со мной"',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'latitude';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'GPS-долгота. Заполнять через геокодер при сохранении адреса',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'longitude';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'false = точка закрыта навсегда или в ремонте. PWA возвращает 404. Не путать с is_emergency_stop',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'is_active';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Кнопка "Запара": мгновенно отключает приём заказов. 1 клик в админке. При true: кнопка оплаты в PWA неактивна, API заказов возвращает 503',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'is_emergency_stop';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Минимальная сумма заказа из ТЗ блок 4. 0 = нет ограничений. Проверяется на бэкенде при нажатии "Оплатить"',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'min_order_amount';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Среднее время готовности (мин) — показывается клиенту. "Время ожидания: ~20 мин". Не влияет на логику',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'avg_cooking_minutes';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Расписание-оверрайд. Формат: {"mon":{"open":"09:00","close":"22:00"}, ...}. Null = берём из iiko.',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'work_hours_json';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Синхронизировать с organization_billing при каждом изменении. Fron-cron проверяет paid_until каждый час',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'billing_status';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'До какой даты оплачен сервис. При успешной оплате: paid_until += 1 month. NULL = триальный период',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'paid_until';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Дата создания точки. Только для чтения',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'created_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Проставлять вручную: SET updated_at = now() при каждом UPDATE. Используется для инвалидации кэша PWA',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'updated_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Soft delete. NULL = активна. Никогда не делать DELETE FROM venues. Все SELECT: WHERE deleted_at IS NULL',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'venues',
@level2type = N'Column', @level2name = 'deleted_at';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'История биллинга организации. Смена плана = новая запись. Старые записи не удалять.',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organization_billing';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID записи биллинга. Текущий план = MAX(created_at) для данной организации',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organization_billing',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'При смене плана: НЕ обновлять старую запись, создавать новую — для хранения истории',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organization_billing',
@level2type = N'Column', @level2name = 'organization_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'JOIN на subscription_plans для получения цены и лимитов при генерации инвойса',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organization_billing',
@level2type = N'Column', @level2name = 'plan_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Дублирует venues.billing_status. Синхронизировать оба поля при любом изменении статуса',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organization_billing',
@level2type = N'Column', @level2name = 'billing_status';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'При оплате: paid_until = MAX(paid_until, now()) + 1 month. Не затирать будущую дату при досрочной оплате',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organization_billing',
@level2type = N'Column', @level2name = 'paid_until';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Токен привязанной карты для рекуррентных платежей',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organization_billing',
@level2type = N'Column', @level2name = 'payment_method_token';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Когда последний раз генерировался инвойс. Обновлять при каждой генерации. Для отладки авто-биллинга',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organization_billing',
@level2type = N'Column', @level2name = 'last_invoice_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Дата начала действия этой записи биллинга',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'organization_billing',
@level2type = N'Column', @level2name = 'created_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID плана. Зашит в organizations.plan_id',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'subscription_plans',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Название для отображения: Starter, Pro, Enterprise',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'subscription_plans',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Цена в рублях. При изменении цены — создавать новый план, не обновлять. Старые организации остаются на старой цене',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'subscription_plans',
@level2type = N'Column', @level2name = 'price_per_month';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Лимиты: макс. блюд, наличие iiko, кол-во точек и т.д. Лимиты и флаги. Формат: {"max_products":100,"iiko_integration":true,"max_venues":1,"push_marketing":false}',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'subscription_plans',
@level2type = N'Column', @level2name = 'features_json';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Архивирование устаревших планов. false = план архивирован, не показывать новым клиентам. Существующие продолжают им пользоваться',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'subscription_plans',
@level2type = N'Column', @level2name = 'is_active';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Дата создания плана',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'subscription_plans',
@level2type = N'Column', @level2name = 'created_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ISO 639-1 код языка: ru, en, ge. Заполняется при деплое seed-скриптом. Не изменять в runtime',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'languages',
@level2type = N'Column', @level2name = 'code';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Читаемое название для выпадающего списка в админке: Русский, English, ქართული',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'languages',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Дизайн настраивается для каждой точки. 1:1 с venues. Создаётся автоматически при создании точки с дефолтными значениями',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'app_configs',
@level2type = N'Column', @level2name = 'venue_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Фронт читает при загрузке PWA и применяет CSS-переменные темы. Менять только через конструктор в админке',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'app_configs',
@level2type = N'Column', @level2name = 'theme_preset';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'HEX цвет кнопок и акцентов. Валидировать: /^#[0-9A-Fa-f]{6}$/. Напр. #E74C3C',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'app_configs',
@level2type = N'Column', @level2name = 'accent_color';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'URL логотипа в S3/CDN (PNG или SVG). NULL = заглушка с первой буквой названия точки. Не хранить base64 в БД',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'app_configs',
@level2type = N'Column', @level2name = 'logo_url';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'URL главного баннера на главной PWA. NULL = без баннера. Рекомендуемый размер: 1200×400px',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'app_configs',
@level2type = N'Column', @level2name = 'banner_url';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Ручная правка адреса (переопределяет venues.address). Перекрывает venues.address только в UI PWA. venues.address по-прежнему используется для геокодинга',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'app_configs',
@level2type = N'Column', @level2name = 'address_manual';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Проставлять при каждом изменении. Используется для инвалидации SW-кэша браузера',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'app_configs',
@level2type = N'Column', @level2name = 'updated_at';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Только идентификация. Роли и привязка к точкам — в таблице staff_roles. Клиент = users без записи в staff_roles.',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID пользователя. Хранится в JWT. Из сессии: sessions.token → sessions.user_id',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'E.164 формат: +79161234567. Нормализовывать при сохранении (убирать пробелы, скобки). Индекс idx_users_phone обязателен',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users',
@level2type = N'Column', @level2name = 'phone';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Имя клиента. NULL при регистрации — запрашиваем после первого входа. Используется в истории заказов и базе клиентов',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users',
@level2type = N'Column', @level2name = 'full_name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Дата рождения. NULL — заполняется опционально в профиле. Для именинных акций. Хранить без времени (тип DATE)',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users',
@level2type = N'Column', @level2name = 'birth_date';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ID гостя в iikoCard. NULL до первого заказа. При первом заказе: создать гостя в iiko → сохранить ID → создать заказ',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users',
@level2type = N'Column', @level2name = 'iiko_guest_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'false = мягкая блокировка или удаление по запросу (152-ФЗ). При false: 401 на все запросы. Данные не удалять',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users',
@level2type = N'Column', @level2name = 'is_active';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Дата регистрации. Только для чтения',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users',
@level2type = N'Column', @level2name = 'created_at';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Персонал: один user может быть admin в одном заведении и manager в другом.',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'staff_roles';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID записи роли',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'staff_roles',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на пользователя. Один человек может иметь несколько записей для разных точек',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'staff_roles',
@level2type = N'Column', @level2name = 'user_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на точку. При авторизации: читаем все staff_roles пользователя, строим список доступных точек',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'staff_roles',
@level2type = N'Column', @level2name = 'venue_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'admin: полный доступ к точке включая биллинг. manager: стоп-листы и заказы. cashier: только список заказов',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'staff_roles',
@level2type = N'Column', @level2name = 'role';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'false = сотрудник уволен. НЕ удалять запись. При авторизации: WHERE is_active = true',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'staff_roles',
@level2type = N'Column', @level2name = 'is_active';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Дата назначения на роль',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'staff_roles',
@level2type = N'Column', @level2name = 'created_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID сессии',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'user_sessions',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на пользователя. При выходе: DELETE WHERE user_id = ? AND id = ?',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'user_sessions',
@level2type = N'Column', @level2name = 'user_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Случайный непрозрачный токен (не JWT). Можно мгновенно отозвать. Индекс idx_sessions_token критически важен',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'user_sessions',
@level2type = N'Column', @level2name = 'token';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'PWA, UA и платформа устройства. Формат: {"ua":"Mozilla...","platform":"Android","pwa":true}. Для аналитики и списка устройств',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'user_sessions',
@level2type = N'Column', @level2name = 'device_info';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Клиенты: now()+30d, персонал: now()+8h. Fоновый job: DELETE WHERE expires_at < now(). Продлевать при каждом запросе',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'user_sessions',
@level2type = N'Column', @level2name = 'expires_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Дата создания сессии (момент входа)',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'user_sessions',
@level2type = N'Column', @level2name = 'created_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID категории',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'categories',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Привязка к точке. При запросе меню: WHERE venue_id = ? AND is_visible = true AND deleted_at IS NULL',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'categories',
@level2type = N'Column', @level2name = 'venue_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID категории в iiko Cloud. Индекс: idx_categories_external. Матчинг при синхронизации. NULL = категория создана вручную в нашей админке',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'categories',
@level2type = N'Column', @level2name = 'external_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Для вложенных категорий. FK на родительскую категорию. NULL = верхний уровень. Для дерева: рекурсивный WITH RECURSIVE запрос',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'categories',
@level2type = N'Column', @level2name = 'parent_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Порядок среди соседних категорий. ORDER BY sort_order ASC. При drag&drop: batch UPDATE',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'categories',
@level2type = N'Column', @level2name = 'sort_order';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Чекбокс в админке: Скрыть/Показать. false = скрыта от клиентов (напр. "Хозтовары"). Только для чтения — видна в админке',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'categories',
@level2type = N'Column', @level2name = 'is_visible';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Soft delete. При удалении из iiko: SET deleted_at = now(). Все SELECT: WHERE deleted_at IS NULL',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'categories',
@level2type = N'Column', @level2name = 'deleted_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на категорию',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'category_translations',
@level2type = N'Column', @level2name = 'category_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Код языка. При запросе меню: JOIN WHERE language_code = язык из Accept-Language',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'category_translations',
@level2type = N'Column', @level2name = 'language_code';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Название категории на данном языке. Напр. "Горячие напитки", "Hot Beverages". ru заполняется из iiko, остальные вручную',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'category_translations',
@level2type = N'Column', @level2name = 'title';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID товара',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на точку. Индекс idx_products_venue_cat_avail используется при каждом запросе меню',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'venue_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на категорию. ON DELETE RESTRICT — нельзя удалить категорию пока в ней есть товары',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'category_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID из iiko. МАТЧИНГ СТРОГО ПО UUID при синхронизации. Нет в нашей БД → INSERT, есть → UPDATE, нет в iiko → soft delete. Индекс: idx_products_external',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'external_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Текущая цена в рублях из iiko. Обновляется при каждой синхронизации. На цену в заказе не влияет (снапшот в order_items)',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'price';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'URL фото из iiko Cloud или нашего S3. NULL = заглушка с логотипом точки. Не хранить base64',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'image_url';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Стоп-лист из iiko или ручной. false = стоп-лист. Обновляется каждые 60 сек из iiko. Проверяется дважды: при показе меню и при нажатии "Оплатить"',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'is_available';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Порядок отображения внутри категории. Берётся из iiko при синхронизации. ORDER BY sort_order ASC',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'sort_order';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Вес порции в граммах из iiko. NULL если не заполнено в iiko. Отображается в карточке: "250 г"',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'weight';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ккал на 100г из iiko. NULL если не заполнено. Блок КБЖУ показываем только если calories IS NOT NULL',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'calories';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Белки на 100г из iiko. Обновляется при ежедневной полной синхронизации',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'protein';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Жиры на 100г из iiko',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'fat';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Углеводы на 100г из iiko',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'carbs';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Soft delete. При исчезновении из iiko: SET deleted_at = now(), is_available = false. Все SELECT: WHERE deleted_at IS NULL',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'deleted_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на товар',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_translations',
@level2type = N'Column', @level2name = 'product_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Код языка перевода',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_translations',
@level2type = N'Column', @level2name = 'language_code';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Название товара. Используется в поиске. GIN-индекс на to_tsvector создаётся отдельной миграцией',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_translations',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Описание товара. NULL разрешён. Отображается в карточке товара при нажатии на него',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_translations',
@level2type = N'Column', @level2name = 'description';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID модификатора',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_modifiers',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на товар. При загрузке карточки: SELECT * FROM product_modifiers WHERE product_id = ?',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_modifiers',
@level2type = N'Column', @level2name = 'product_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID модификатора в iiko. NULL если создан вручную. Матчинг при синхронизации аналогично products.external_id',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_modifiers',
@level2type = N'Column', @level2name = 'external_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '0 = бесплатный модификатор (напр. "Без сахара"). Прибавляется к цене. Снапшот хранится в order_item_modifiers',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_modifiers',
@level2type = N'Column', @level2name = 'extra_price';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'true = клиент не может добавить товар в корзину без выбора. Напр. группа "Размер": S/M/L обязательна',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_modifiers',
@level2type = N'Column', @level2name = 'is_required';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '1 = чекбокс. >1 = можно выбрать несколько (напр. количество сиропов до 3). Валидировать на бэкенде',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_modifiers',
@level2type = N'Column', @level2name = 'max_quantity';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Название группы для UI: "Размер", "Добавки", "Основа". NULL = одиночный модификатор. Одинаковый group_name = одна визуальная группа',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_modifiers',
@level2type = N'Column', @level2name = 'group_name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на модификатор',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'modifier_translations',
@level2type = N'Column', @level2name = 'modifier_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Код языка перевода',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'modifier_translations',
@level2type = N'Column', @level2name = 'language_code';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Название добавки на данном языке. Напр. "Большой", "Large", "Без сахара"',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'modifier_translations',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Настраивается вручную в админке. Блок "С этим также берут" на странице товара.',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_recommendations';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Товар, на странице которого показывается блок. SELECT recommended_id WHERE product_id = ? ORDER BY sort_order',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_recommendations',
@level2type = N'Column', @level2name = 'product_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Рекомендуемый товар. При запросе: проверять is_available = true — не показывать стоп-лист в upsell',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_recommendations',
@level2type = N'Column', @level2name = 'recommended_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Порядок в блоке "С этим также берут". Обновляется drag&drop в админке. ORDER BY sort_order ASC',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product_recommendations',
@level2type = N'Column', @level2name = 'sort_order';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID заказа. Показывается клиенту как номер: последние 4 символа UUID',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на точку. Все заказы точки: WHERE venue_id = ? AND deleted_at IS NULL',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'venue_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на клиента. История клиента: WHERE user_id = ? ORDER BY created_at DESC LIMIT 20',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'user_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Менять только через сервисные методы, не UPDATE напрямую — иначе пуши и логи не сработают',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'status';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'В MVP всегда takeaway. dine_in и delivery — для будущих версий',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'order_type';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Итог после скидки. Формула: SUM(unit_price*qty) + SUM(modifier.extra_price) - discount_amount. Фиксируется при создании',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'total_amount';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Предзаказ на время.NULL = "как можно скорее". Если задано: шаг 15 мин, только текущий или следующий день. Передаётся в iiko как время доставки',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'scheduled_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Комментарий клиента к заказу, "Без лука, пожалуйста". Передаётся в iiko в поле комментария. Показывается сотруднику в карточке заказа',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'comment';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Чекбокс "Нужны приборы" в корзине. Передаётся в iiko как тег или комментарий',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'cutlery_needed';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'NULL если заказ не отменён. Заполняется автоматически ("Таймаут оплаты") или вручную сотрудником. Для аналитики потерь',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'cancelled_reason';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ID возврата в платёжной системе. NULL пока возврата не было. При наличии: показываем статус "Возврат выполнен"',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'refund_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Момент выполнения возврата. NULL если возврата не было',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'refunded_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK → payments.id (nullable, проставляется после создания платежа). NULL до инициации платежа. Порядок: создать заказ → создать payment → обновить orders.payment_id',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'payment_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID заказа от iiko. NULL до успешной отправки. Используется для отмены в iiko и запроса статуса готовности',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'pos_order_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Воркер каждые 5 сек: WHERE pos_sync_status = NOT_SENT AND status = COOKING. Retry до 3 раз при ERROR',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'pos_sync_status';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Передаётся в iiko как orderSource для маркетинговой аналитики. По ТЗ = название нашего приложения',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'order_source';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'NULL если промокод не применялся. При применении: проверить валидность → INCREMENT used_count → сохранить ID',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'promocode_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Сумма скидки по промокоду. total_amount уже с учётом скидки. Поле хранится для прозрачности аналитики промокодов',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'discount_amount';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Время создания заказа клиентом',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'created_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Soft delete. Заказы НИКОГДА не удалять физически — финансовая история. deleted_at только для технически ошибочных записей',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'deleted_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на заказ. Все позиции заказа: WHERE order_id = ?',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_items',
@level2type = N'Column', @level2name = 'order_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на товар. Использовать LEFT JOIN — товар может быть удалён из меню',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_items',
@level2type = N'Column', @level2name = 'product_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Количество единиц. Минимум 1. Нельзя изменить после создания заказа',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_items',
@level2type = N'Column', @level2name = 'quantity';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Цена 1 единицы товара на момент оплаты (снапшот из products.price). Не зависит от будущих изменений цены',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_items',
@level2type = N'Column', @level2name = 'unit_price';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Название товара на языке клиента на момент заказа. ОБЯЗАТЕЛЬНО заполнять. Без снапшота история пустеет при удалении товара',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_items',
@level2type = N'Column', @level2name = 'product_name_snapshot';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Нормализованная замена jsonb selected_modifiers. Позволяет строить аналитику по добавкам.',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_item_modifiers';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID записи модификатора заказа',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_item_modifiers',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на позицию заказа. Все добавки позиции: WHERE order_item_id = ?',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_item_modifiers',
@level2type = N'Column', @level2name = 'order_item_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на модификатор. NULLABLE — модификатор может быть удалён из меню. Всегда LEFT JOIN, не INNER',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_item_modifiers',
@level2type = N'Column', @level2name = 'modifier_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Название добавки на момент заказа. ОБЯЗАТЕЛЬНО заполнять. Аналогично product_name_snapshot — защита от удаления',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_item_modifiers',
@level2type = N'Column', @level2name = 'name_snapshot';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Цена добавки на момент заказа (снапшот). Итог позиции: (unit_price + SUM(extra_price)) * quantity',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_item_modifiers',
@level2type = N'Column', @level2name = 'extra_price';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID платежа в нашей системе',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'payments',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '1:1 с заказом. UNIQUE гарантирует один платёж на заказ',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'payments',
@level2type = N'Column', @level2name = 'order_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ID транзакции в Тинькофф/ЮKassa. Используется для идентификации входящих webhook: WHERE external_id = ?',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'payments',
@level2type = N'Column', @level2name = 'external_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Обновлять ТОЛЬКО через webhook от банка или polling. Никогда не менять вручную без подтверждения банка',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'payments',
@level2type = N'Column', @level2name = 'status';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Сумма платежа в рублях. Должна совпадать с orders.total_amount',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'payments',
@level2type = N'Column', @level2name = 'amount';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Название провайдера: tinkoff, yookassa, sbp. Нужен при возврате — используем правильный API',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'payments',
@level2type = N'Column', @level2name = 'provider';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Момент фактического списания от банка (из данных webhook). NULL до SUCCESS. Не ставить своё время — только из данных банка',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'payments',
@level2type = N'Column', @level2name = 'paid_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Время создания платёжной записи (инициации платежа)',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'payments',
@level2type = N'Column', @level2name = 'created_at';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'TTL: DELETE WHERE created_at < NOW() - INTERVAL 14 days (pg_cron job).',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_sync_logs';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID записи лога',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_sync_logs',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на точку. Для дебага: WHERE venue_id = ? ORDER BY created_at DESC',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_sync_logs',
@level2type = N'Column', @level2name = 'venue_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на заказ. Все попытки отправки заказа: WHERE order_id = ? ORDER BY created_at DESC',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_sync_logs',
@level2type = N'Column', @level2name = 'order_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Полное тело запроса отправленного в iiko. Логировать ДО отправки. НЕ включать api_key в payload',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_sync_logs',
@level2type = N'Column', @level2name = 'request_payload';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Полный ответ iiko (JSON). NULL если нет ответа (таймаут). При ERROR статусе — главный источник для диагностики',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_sync_logs',
@level2type = N'Column', @level2name = 'response_payload';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'HTTP статус ответа iiko. 200 = успех → SENT. 500/503 = retry до 3 раз через 30 сек. 401 = невалидный api_key → оповестить владельца. NULL = таймаут',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_sync_logs',
@level2type = N'Column', @level2name = 'http_status_code';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Время попытки синхронизации. TTL 14 дней: pg_cron DELETE WHERE created_at < NOW() - INTERVAL 14 days',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_sync_logs',
@level2type = N'Column', @level2name = 'created_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '1:1 с точкой. Создаётся при подключении iiko в онбординге',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_settings',
@level2type = N'Column', @level2name = 'venue_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'API-ключ iikoTransport. ЗАШИФРОВАТЬ через pgcrypto/Vault. В БД только blob. Ключ шифрования в env. Никогда не логировать',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_settings',
@level2type = N'Column', @level2name = 'api_key_encrypted';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID терминала iiko, куда "падают" заказы. Если несколько терминалов — владелец выбирает нужный при настройке',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_settings',
@level2type = N'Column', @level2name = 'terminal_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID организации в облаке iiko (≠ наш organizations.id). Передаётся в каждый запрос к iikoCloud как organizationId',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_settings',
@level2type = N'Column', @level2name = 'external_org_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Интервал полной синхронизации меню в минутах. Планировщик запускает каждые N минут. Полная (картинки) — раз в сутки',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_settings',
@level2type = N'Column', @level2name = 'sync_interval_minutes';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Интервал опроса стоп-листа в секундах. Отдельный частый воркер обновляет products.is_available всех товаров точки',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_settings',
@level2type = N'Column', @level2name = 'stoplist_poll_seconds';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Timestamp последней успешной синхронизации меню. Обновлять при каждом успехе. Отображать в админке: "Обновлено 2 мин назад"',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_settings',
@level2type = N'Column', @level2name = 'last_menu_sync_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Timestamp последнего опроса стоп-листа. Если > 5 мин назад — алерт в Sentry. NULL = синхронизации ещё не было',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_settings',
@level2type = N'Column', @level2name = 'last_stoplist_sync_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Дата последнего изменения настроек. Проставлять при каждом UPDATE',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'pos_settings',
@level2type = N'Column', @level2name = 'updated_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID промокода',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'promocodes',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на точку. Промокод действует только в этой точке. UNIQUE(venue_id, code)',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'promocodes',
@level2type = N'Column', @level2name = 'venue_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Код промокода: WELCOME, SUMMER20. Нормализовывать: UPPER(TRIM(code)) при записи и при поиске',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'promocodes',
@level2type = N'Column', @level2name = 'code';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'PERCENT или FIXED. Определяет логику вычисления discount_amount в заказе',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'promocodes',
@level2type = N'Column', @level2name = 'discount_type';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'PERCENT: 10 = 10%. FIXED: 150 = 150 руб. Для FIXED: discount = MIN(value, total) — защита от отрицательной суммы',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'promocodes',
@level2type = N'Column', @level2name = 'discount_value';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Минимальная сумма заказа для применения промокода. 0 = без ограничений. Проверять до применения скидки',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'promocodes',
@level2type = N'Column', @level2name = 'min_order_amount';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Лимит использований. NULL = безлимит. При применении: SELECT FOR UPDATE → проверить → INCREMENT. Защита от race condition',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'promocodes',
@level2type = N'Column', @level2name = 'max_uses';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Счётчик использований. INCREMENT в транзакции с SELECT FOR UPDATE. Никогда не уменьшать вручную',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'promocodes',
@level2type = N'Column', @level2name = 'used_count';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Срок действия. NULL = бессрочно. Проверка: WHERE (expires_at IS NULL OR expires_at > NOW()) AND is_active = true',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'promocodes',
@level2type = N'Column', @level2name = 'expires_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'false = промокод отключён вручную. При деактивации существующие заказы с промокодом не затрагиваются',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'promocodes',
@level2type = N'Column', @level2name = 'is_active';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Дата создания промокода',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'promocodes',
@level2type = N'Column', @level2name = 'created_at';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Внутренние отзывы — видны только владельцу в админке (ТЗ блок 6).',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_reviews';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID отзыва',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_reviews',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '1:1 с заказом. UNIQUE гарантирует один отзыв на заказ. Предлагать оценку через пуш спустя 1 час',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_reviews',
@level2type = N'Column', @level2name = 'order_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на клиента',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_reviews',
@level2type = N'Column', @level2name = 'user_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на точку. Средний рейтинг: SELECT AVG(rating) WHERE venue_id = ?',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_reviews',
@level2type = N'Column', @level2name = 'venue_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Оценка 1-5. Валидировать: CHECK (rating BETWEEN 1 AND 5). 1 звезда = очень плохо, 5 = отлично',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_reviews',
@level2type = N'Column', @level2name = 'rating';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Текст отзыва. NULL если клиент поставил только оценку без текста. Показывать ТОЛЬКО в админке владельца — не публиковать в PWA',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_reviews',
@level2type = N'Column', @level2name = 'comment';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Ответ владельца. NULL если не ответил. Заполняется в разделе "Отзывы" в админке. В MVP клиенту не доставляется',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_reviews',
@level2type = N'Column', @level2name = 'owner_reply';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Время ответа владельца. Проставлять: SET replied_at = now() вместе с owner_reply',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_reviews',
@level2type = N'Column', @level2name = 'replied_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Время создания отзыва',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'order_reviews',
@level2type = N'Column', @level2name = 'created_at';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID подписки',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'push_subscriptions',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на пользователя',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'push_subscriptions',
@level2type = N'Column', @level2name = 'user_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на точку. Маркетинговые пуши точки: WHERE venue_id = ? AND is_active = true',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'push_subscriptions',
@level2type = N'Column', @level2name = 'venue_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'URL push-сервиса браузера (Google FCM, Apple и др.). Получаем из ServiceWorker.pushManager.subscribe(). Уникален для устройства',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'push_subscriptions',
@level2type = N'Column', @level2name = 'endpoint';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Публичный ключ шифрования устройства. Не секрет. Передаётся в web-push библиотеку для шифрования payload',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'push_subscriptions',
@level2type = N'Column', @level2name = 'p256dh';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Auth-секрет устройства. Не секрет. Передаётся в web-push: sendNotification({endpoint, keys:{p256dh,auth}}, payload)',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'push_subscriptions',
@level2type = N'Column', @level2name = 'auth';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'false = пользователь отозвал разрешение. Ставить false если браузер вернул 410 Gone при отправке пуша',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'push_subscriptions',
@level2type = N'Column', @level2name = 'is_active';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Время создания подписки',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'push_subscriptions',
@level2type = N'Column', @level2name = 'created_at';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Привязка аккаунта к Telegram-боту для уведомлений персоналу и клиентам.',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'telegram_bindings';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'UUID привязки',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'telegram_bindings',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'FK на пользователя. UNIQUE гарантирует один Telegram на аккаунт',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'telegram_bindings',
@level2type = N'Column', @level2name = 'user_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Числовой chat_id из Telegram Bot API (bigint — может быть > 2^31). Получаем из message.chat.id при команде /start',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'telegram_bindings',
@level2type = N'Column', @level2name = 'telegram_chat_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'false = пользователь написал /stop боту. При отправке пушей: WHERE is_active = true',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'telegram_bindings',
@level2type = N'Column', @level2name = 'is_active';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Дата привязки аккаунта к боту',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'telegram_bindings',
@level2type = N'Column', @level2name = 'created_at';
GO

ALTER TABLE [organizations] ADD FOREIGN KEY ([plan_id]) REFERENCES [subscription_plans] ([id])
GO

ALTER TABLE [venues] ADD FOREIGN KEY ([organization_id]) REFERENCES [organizations] ([id])
GO

ALTER TABLE [organization_billing] ADD FOREIGN KEY ([organization_id]) REFERENCES [organizations] ([id])
GO

ALTER TABLE [organization_billing] ADD FOREIGN KEY ([plan_id]) REFERENCES [subscription_plans] ([id])
GO

ALTER TABLE [app_configs] ADD FOREIGN KEY ([venue_id]) REFERENCES [venues] ([id])
GO

ALTER TABLE [staff_roles] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [staff_roles] ADD FOREIGN KEY ([venue_id]) REFERENCES [venues] ([id])
GO

ALTER TABLE [user_sessions] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [categories] ADD FOREIGN KEY ([venue_id]) REFERENCES [venues] ([id])
GO

ALTER TABLE [categories] ADD FOREIGN KEY ([parent_id]) REFERENCES [categories] ([id])
GO

ALTER TABLE [category_translations] ADD FOREIGN KEY ([category_id]) REFERENCES [categories] ([id])
GO

ALTER TABLE [category_translations] ADD FOREIGN KEY ([language_code]) REFERENCES [languages] ([code])
GO

ALTER TABLE [products] ADD FOREIGN KEY ([venue_id]) REFERENCES [venues] ([id])
GO

ALTER TABLE [products] ADD FOREIGN KEY ([category_id]) REFERENCES [categories] ([id])
GO

ALTER TABLE [product_translations] ADD FOREIGN KEY ([product_id]) REFERENCES [products] ([id])
GO

ALTER TABLE [product_translations] ADD FOREIGN KEY ([language_code]) REFERENCES [languages] ([code])
GO

ALTER TABLE [product_modifiers] ADD FOREIGN KEY ([product_id]) REFERENCES [products] ([id])
GO

ALTER TABLE [modifier_translations] ADD FOREIGN KEY ([modifier_id]) REFERENCES [product_modifiers] ([id])
GO

ALTER TABLE [modifier_translations] ADD FOREIGN KEY ([language_code]) REFERENCES [languages] ([code])
GO

ALTER TABLE [product_recommendations] ADD FOREIGN KEY ([product_id]) REFERENCES [products] ([id])
GO

ALTER TABLE [product_recommendations] ADD FOREIGN KEY ([recommended_id]) REFERENCES [products] ([id])
GO

ALTER TABLE [orders] ADD FOREIGN KEY ([venue_id]) REFERENCES [venues] ([id])
GO

ALTER TABLE [orders] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [orders] ADD FOREIGN KEY ([payment_id]) REFERENCES [payments] ([id])
GO

ALTER TABLE [orders] ADD FOREIGN KEY ([promocode_id]) REFERENCES [promocodes] ([id])
GO

ALTER TABLE [order_items] ADD FOREIGN KEY ([order_id]) REFERENCES [orders] ([id])
GO

ALTER TABLE [order_items] ADD FOREIGN KEY ([product_id]) REFERENCES [products] ([id])
GO

ALTER TABLE [order_item_modifiers] ADD FOREIGN KEY ([order_item_id]) REFERENCES [order_items] ([id])
GO

ALTER TABLE [order_item_modifiers] ADD FOREIGN KEY ([modifier_id]) REFERENCES [product_modifiers] ([id])
GO

ALTER TABLE [payments] ADD FOREIGN KEY ([order_id]) REFERENCES [orders] ([id])
GO

ALTER TABLE [pos_sync_logs] ADD FOREIGN KEY ([venue_id]) REFERENCES [venues] ([id])
GO

ALTER TABLE [pos_sync_logs] ADD FOREIGN KEY ([order_id]) REFERENCES [orders] ([id])
GO

ALTER TABLE [pos_settings] ADD FOREIGN KEY ([venue_id]) REFERENCES [venues] ([id])
GO

ALTER TABLE [promocodes] ADD FOREIGN KEY ([venue_id]) REFERENCES [venues] ([id])
GO

ALTER TABLE [order_reviews] ADD FOREIGN KEY ([order_id]) REFERENCES [orders] ([id])
GO

ALTER TABLE [order_reviews] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [order_reviews] ADD FOREIGN KEY ([venue_id]) REFERENCES [venues] ([id])
GO

ALTER TABLE [push_subscriptions] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [push_subscriptions] ADD FOREIGN KEY ([venue_id]) REFERENCES [venues] ([id])
GO

ALTER TABLE [telegram_bindings] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO
