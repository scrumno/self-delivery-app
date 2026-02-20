CREATE TYPE "user_role" AS ENUM (
  'admin',
  'manager',
  'cashier',
  'customer'
);

CREATE TYPE "order_status" AS ENUM (
  'NEW',
  'WAITING_PAYMENT',
  'COOKING',
  'READY',
  'COMPLETED',
  'CANCELLED'
);

CREATE TYPE "pos_sync_status" AS ENUM (
  'NOT_SENT',
  'SENT',
  'ERROR'
);

CREATE TYPE "billing_status" AS ENUM (
  'ACTIVE',
  'PAST_DUE',
  'BLOCKED'
);

CREATE TYPE "payment_status" AS ENUM (
  'PENDING',
  'SUCCESS',
  'FAILED',
  'REFUNDED',
  'CANCELLED'
);

CREATE TYPE "theme_preset" AS ENUM (
  'light',
  'dark',
  'coffee',
  'fastfood'
);

CREATE TYPE "discount_type" AS ENUM (
  'PERCENT',
  'FIXED'
);

CREATE TYPE "order_type" AS ENUM (
  'takeaway',
  'dine_in',
  'delivery'
);

CREATE TABLE "organizations" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "inn" varchar(12),
  "is_active" boolean DEFAULT true,
  "created_at" timestamptz DEFAULT (now()),
  "plan_id" uuid NOT NULL
);

CREATE TABLE "venues" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "organization_id" uuid NOT NULL,
  "name" varchar NOT NULL,
  "slug" varchar UNIQUE NOT NULL,
  "address" text,
  "latitude" numeric(10,7),
  "longitude" numeric(10,7),
  "is_active" boolean DEFAULT true,
  "is_emergency_stop" boolean DEFAULT false,
  "min_order_amount" numeric(12,2) DEFAULT 0,
  "avg_cooking_minutes" integer DEFAULT 20,
  "work_hours_json" jsonb,
  "billing_status" billing_status DEFAULT 'ACTIVE',
  "paid_until" timestamptz,
  "created_at" timestamptz DEFAULT (now()),
  "updated_at" timestamptz,
  "deleted_at" timestamptz
);

CREATE TABLE "organization_billing" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "organization_id" uuid NOT NULL,
  "plan_id" uuid NOT NULL,
  "billing_status" billing_status DEFAULT 'ACTIVE',
  "paid_until" timestamptz,
  "payment_method_token" text,
  "last_invoice_at" timestamptz,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "subscription_plans" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "price_per_month" numeric(12,2) NOT NULL,
  "features_json" jsonb,
  "is_active" boolean DEFAULT true,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "languages" (
  "code" varchar(5) PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "app_configs" (
  "venue_id" uuid PRIMARY KEY,
  "theme_preset" theme_preset DEFAULT 'light',
  "accent_color" varchar(7) DEFAULT '#000000',
  "logo_url" text,
  "banner_url" text,
  "address_manual" text,
  "updated_at" timestamptz
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "phone" varchar UNIQUE NOT NULL,
  "full_name" varchar,
  "birth_date" date,
  "iiko_guest_id" uuid,
  "is_active" boolean DEFAULT true,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "staff_roles" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "venue_id" uuid NOT NULL,
  "role" user_role NOT NULL,
  "is_active" boolean DEFAULT true,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "user_sessions" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "token" text UNIQUE NOT NULL,
  "device_info" jsonb,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "categories" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "venue_id" uuid NOT NULL,
  "external_id" uuid,
  "parent_id" uuid,
  "sort_order" integer DEFAULT 0,
  "is_visible" boolean DEFAULT true,
  "deleted_at" timestamptz
);

CREATE TABLE "category_translations" (
  "category_id" uuid,
  "language_code" varchar,
  "title" varchar NOT NULL,
  PRIMARY KEY ("category_id", "language_code")
);

CREATE TABLE "products" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "venue_id" uuid NOT NULL,
  "category_id" uuid,
  "external_id" uuid NOT NULL,
  "price" numeric(12,2) NOT NULL,
  "image_url" text,
  "is_available" boolean DEFAULT true,
  "sort_order" integer DEFAULT 0,
  "weight" integer,
  "calories" numeric(8,2),
  "protein" numeric(8,2),
  "fat" numeric(8,2),
  "carbs" numeric(8,2),
  "deleted_at" timestamptz
);

CREATE TABLE "product_translations" (
  "product_id" uuid,
  "language_code" varchar,
  "name" varchar NOT NULL,
  "description" text,
  PRIMARY KEY ("product_id", "language_code")
);

CREATE TABLE "product_modifiers" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "product_id" uuid NOT NULL,
  "external_id" uuid,
  "extra_price" numeric(12,2) DEFAULT 0,
  "is_required" boolean DEFAULT false,
  "max_quantity" integer DEFAULT 1,
  "group_name" varchar
);

CREATE TABLE "modifier_translations" (
  "modifier_id" uuid,
  "language_code" varchar,
  "name" varchar NOT NULL,
  PRIMARY KEY ("modifier_id", "language_code")
);

CREATE TABLE "product_recommendations" (
  "product_id" uuid,
  "recommended_id" uuid,
  "sort_order" integer DEFAULT 0,
  PRIMARY KEY ("product_id", "recommended_id")
);

CREATE TABLE "orders" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "venue_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "status" order_status DEFAULT 'NEW',
  "order_type" order_type DEFAULT 'takeaway',
  "total_amount" numeric(12,2) NOT NULL,
  "scheduled_at" timestamptz,
  "comment" text,
  "cutlery_needed" boolean DEFAULT false,
  "cancelled_reason" text,
  "refund_id" text,
  "refunded_at" timestamptz,
  "payment_id" uuid,
  "pos_order_id" uuid,
  "pos_sync_status" pos_sync_status DEFAULT 'NOT_SENT',
  "order_source" varchar DEFAULT 'PWA',
  "promocode_id" uuid,
  "discount_amount" numeric(12,2) DEFAULT 0,
  "created_at" timestamptz DEFAULT (now()),
  "deleted_at" timestamptz
);

CREATE TABLE "order_items" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "order_id" uuid NOT NULL,
  "product_id" uuid NOT NULL,
  "quantity" integer NOT NULL,
  "unit_price" numeric(12,2) NOT NULL,
  "product_name_snapshot" varchar
);

CREATE TABLE "order_item_modifiers" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "order_item_id" uuid NOT NULL,
  "modifier_id" uuid,
  "name_snapshot" varchar NOT NULL,
  "extra_price" numeric(12,2) NOT NULL DEFAULT 0
);

CREATE TABLE "payments" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "order_id" uuid UNIQUE NOT NULL,
  "external_id" varchar NOT NULL,
  "status" payment_status DEFAULT 'PENDING',
  "amount" numeric(12,2) NOT NULL,
  "provider" varchar,
  "paid_at" timestamptz,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "pos_sync_logs" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "venue_id" uuid NOT NULL,
  "order_id" uuid NOT NULL,
  "request_payload" jsonb,
  "response_payload" jsonb,
  "http_status_code" integer,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "pos_settings" (
  "venue_id" uuid PRIMARY KEY,
  "api_key_encrypted" text NOT NULL,
  "terminal_id" uuid NOT NULL,
  "external_org_id" uuid NOT NULL,
  "sync_interval_minutes" integer DEFAULT 1,
  "stoplist_poll_seconds" integer DEFAULT 60,
  "last_menu_sync_at" timestamptz,
  "last_stoplist_sync_at" timestamptz,
  "updated_at" timestamptz
);

CREATE TABLE "promocodes" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "venue_id" uuid NOT NULL,
  "code" varchar NOT NULL,
  "discount_type" discount_type NOT NULL,
  "discount_value" numeric(12,2) NOT NULL,
  "min_order_amount" numeric(12,2) DEFAULT 0,
  "max_uses" integer,
  "used_count" integer NOT NULL DEFAULT 0,
  "expires_at" timestamptz,
  "is_active" boolean DEFAULT true,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "order_reviews" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "order_id" uuid UNIQUE NOT NULL,
  "user_id" uuid NOT NULL,
  "venue_id" uuid NOT NULL,
  "rating" smallint NOT NULL,
  "comment" text,
  "owner_reply" text,
  "replied_at" timestamptz,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "push_subscriptions" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "venue_id" uuid NOT NULL,
  "endpoint" text NOT NULL,
  "p256dh" text NOT NULL,
  "auth" text NOT NULL,
  "is_active" boolean DEFAULT true,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "telegram_bindings" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "telegram_chat_id" bigint NOT NULL,
  "is_active" boolean DEFAULT true,
  "created_at" timestamptz DEFAULT (now())
);

CREATE INDEX "idx_users_phone" ON "users" ("phone");

CREATE UNIQUE INDEX "idx_staff_roles_user_venue" ON "staff_roles" ("user_id", "venue_id");

CREATE INDEX "idx_sessions_token" ON "user_sessions" ("token");

CREATE INDEX "idx_sessions_user" ON "user_sessions" ("user_id");

CREATE INDEX "idx_sessions_expires" ON "user_sessions" ("expires_at");

CREATE INDEX "idx_categories_venue_visible" ON "categories" ("venue_id", "is_visible");

CREATE UNIQUE INDEX "idx_categories_external_unique" ON "categories" ("venue_id", "external_id");

CREATE INDEX "idx_products_pwa_main" ON "products" ("venue_id", "category_id", "is_available", "sort_order");

CREATE UNIQUE INDEX "idx_products_external_unique" ON "products" ("venue_id", "external_id");

CREATE INDEX "idx_recs_sort" ON "product_recommendations" ("product_id", "sort_order");

CREATE INDEX "idx_orders_venue_status_date" ON "orders" ("venue_id", "status", "created_at");

CREATE INDEX "idx_orders_user_date" ON "orders" ("user_id", "created_at");

CREATE INDEX "idx_orders_payment" ON "orders" ("payment_id");

CREATE INDEX "idx_orders_sync" ON "orders" ("pos_sync_status");

CREATE INDEX "idx_order_items_order" ON "order_items" ("order_id");

CREATE INDEX "idx_oim_order_item" ON "order_item_modifiers" ("order_item_id");

CREATE INDEX "idx_payments_order" ON "payments" ("order_id");

CREATE UNIQUE INDEX "idx_payments_external_unique" ON "payments" ("external_id");

CREATE INDEX "idx_payments_status" ON "payments" ("status");

CREATE INDEX "idx_pos_logs_order_date" ON "pos_sync_logs" ("order_id", "created_at");

CREATE INDEX "idx_pos_logs_venue_date" ON "pos_sync_logs" ("venue_id", "created_at");

CREATE UNIQUE INDEX "idx_promocodes_venue_code" ON "promocodes" ("venue_id", "code");

CREATE INDEX "idx_promocodes_active" ON "promocodes" ("venue_id", "is_active", "expires_at");

CREATE INDEX "idx_reviews_venue_date" ON "order_reviews" ("venue_id", "created_at");

CREATE INDEX "idx_reviews_user" ON "order_reviews" ("user_id");

CREATE INDEX "idx_push_user_venue" ON "push_subscriptions" ("user_id", "venue_id");

CREATE UNIQUE INDEX "idx_tg_user" ON "telegram_bindings" ("user_id");

CREATE INDEX "idx_tg_chat" ON "telegram_bindings" ("telegram_chat_id");

COMMENT ON TABLE "organizations" IS 'Юридическое лицо. Одна организация может иметь несколько точек (venues). Не хранит операционных данных.';

COMMENT ON COLUMN "organizations"."id" IS 'UUID юрлица. Передаётся в JWT сотрудника. Все запросы персонала фильтруются через него';

COMMENT ON COLUMN "organizations"."name" IS 'Название для договора. Официальное название для договоров и инвойсов. Не путать с venues.name — публичным названием';

COMMENT ON COLUMN "organizations"."inn" IS 'ИНН для проверки контрагента (РФ). ИНН: 10 цифр для ООО, 12 для ИП. Валидировать контрольную сумму. NULL разрешён в MVP';

COMMENT ON COLUMN "organizations"."is_active" IS 'Глобальное выключение аккаунта. false = ручная блокировка саппортом. Все venues этой организации перестают работать';

COMMENT ON COLUMN "organizations"."created_at" IS 'Дата регистрации. Только для чтения. Используется в аналитике роста клиентской базы';

COMMENT ON COLUMN "organizations"."plan_id" IS 'FK → subscription_plans.id';

COMMENT ON TABLE "venues" IS 'Центральная сущность системы. Торговая точка. Один organization может иметь много venues. Вся iiko/PWA логика привязана к venue_id.';

COMMENT ON COLUMN "venues"."id" IS 'UUID точки. Главный FK всей системы — почти каждая таблица ссылается на venue_id';

COMMENT ON COLUMN "venues"."organization_id" IS 'Владелец точки. При авторизации: staff_roles.venue_id → venues → organizations → проверяем billing_status';

COMMENT ON COLUMN "venues"."name" IS 'Публичное название точки для клиентов и шапки PWA. Напр. "Кофе-пойнт ТЦ Мега"';

COMMENT ON COLUMN "venues"."slug" IS 'URL-часть: brand-kazan.foodapp.ru. Только [a-z0-9-], валидировать регуляркой. Менять осторожно — старые QR-коды перестанут работать';

COMMENT ON COLUMN "venues"."address" IS 'Физический адрес. Используется для геокодинга. Для отображения в PWA может быть перекрыт app_configs.address_manual';

COMMENT ON COLUMN "venues"."latitude" IS 'Координаты для карты. GPS-широта. Заполнять через геокодер при сохранении адреса. Нужна для карты и функции "рядом со мной"';

COMMENT ON COLUMN "venues"."longitude" IS 'GPS-долгота. Заполнять через геокодер при сохранении адреса';

COMMENT ON COLUMN "venues"."is_active" IS 'false = точка закрыта навсегда или в ремонте. PWA возвращает 404. Не путать с is_emergency_stop';

COMMENT ON COLUMN "venues"."is_emergency_stop" IS 'Кнопка "Запара": мгновенно отключает приём заказов. 1 клик в админке. При true: кнопка оплаты в PWA неактивна, API заказов возвращает 503';

COMMENT ON COLUMN "venues"."min_order_amount" IS 'Минимальная сумма заказа из ТЗ блок 4. 0 = нет ограничений. Проверяется на бэкенде при нажатии "Оплатить"';

COMMENT ON COLUMN "venues"."avg_cooking_minutes" IS 'Среднее время готовности (мин) — показывается клиенту. "Время ожидания: ~20 мин". Не влияет на логику';

COMMENT ON COLUMN "venues"."work_hours_json" IS 'Расписание-оверрайд. Формат: {"mon":{"open":"09:00","close":"22:00"}, ...}. Null = берём из iiko.';

COMMENT ON COLUMN "venues"."billing_status" IS 'Синхронизировать с organization_billing при каждом изменении. Fron-cron проверяет paid_until каждый час';

COMMENT ON COLUMN "venues"."paid_until" IS 'До какой даты оплачен сервис. При успешной оплате: paid_until += 1 month. NULL = триальный период';

COMMENT ON COLUMN "venues"."created_at" IS 'Дата создания точки. Только для чтения';

COMMENT ON COLUMN "venues"."updated_at" IS 'Проставлять вручную: SET updated_at = now() при каждом UPDATE. Используется для инвалидации кэша PWA';

COMMENT ON COLUMN "venues"."deleted_at" IS 'Soft delete. NULL = активна. Никогда не делать DELETE FROM venues. Все SELECT: WHERE deleted_at IS NULL';

COMMENT ON TABLE "organization_billing" IS 'История биллинга организации. Смена плана = новая запись. Старые записи не удалять.';

COMMENT ON COLUMN "organization_billing"."id" IS 'UUID записи биллинга. Текущий план = MAX(created_at) для данной организации';

COMMENT ON COLUMN "organization_billing"."organization_id" IS 'При смене плана: НЕ обновлять старую запись, создавать новую — для хранения истории';

COMMENT ON COLUMN "organization_billing"."plan_id" IS 'JOIN на subscription_plans для получения цены и лимитов при генерации инвойса';

COMMENT ON COLUMN "organization_billing"."billing_status" IS 'Дублирует venues.billing_status. Синхронизировать оба поля при любом изменении статуса';

COMMENT ON COLUMN "organization_billing"."paid_until" IS 'При оплате: paid_until = MAX(paid_until, now()) + 1 month. Не затирать будущую дату при досрочной оплате';

COMMENT ON COLUMN "organization_billing"."payment_method_token" IS 'Токен привязанной карты для рекуррентных платежей';

COMMENT ON COLUMN "organization_billing"."last_invoice_at" IS 'Когда последний раз генерировался инвойс. Обновлять при каждой генерации. Для отладки авто-биллинга';

COMMENT ON COLUMN "organization_billing"."created_at" IS 'Дата начала действия этой записи биллинга';

COMMENT ON COLUMN "subscription_plans"."id" IS 'UUID плана. Зашит в organizations.plan_id';

COMMENT ON COLUMN "subscription_plans"."name" IS 'Название для отображения: Starter, Pro, Enterprise';

COMMENT ON COLUMN "subscription_plans"."price_per_month" IS 'Цена в рублях. При изменении цены — создавать новый план, не обновлять. Старые организации остаются на старой цене';

COMMENT ON COLUMN "subscription_plans"."features_json" IS 'Лимиты: макс. блюд, наличие iiko, кол-во точек и т.д. Лимиты и флаги. Формат: {"max_products":100,"iiko_integration":true,"max_venues":1,"push_marketing":false}';

COMMENT ON COLUMN "subscription_plans"."is_active" IS 'Архивирование устаревших планов. false = план архивирован, не показывать новым клиентам. Существующие продолжают им пользоваться';

COMMENT ON COLUMN "subscription_plans"."created_at" IS 'Дата создания плана';

COMMENT ON COLUMN "languages"."code" IS 'ISO 639-1 код языка: ru, en, ge. Заполняется при деплое seed-скриптом. Не изменять в runtime';

COMMENT ON COLUMN "languages"."name" IS 'Читаемое название для выпадающего списка в админке: Русский, English, ქართული';

COMMENT ON COLUMN "app_configs"."venue_id" IS 'Дизайн настраивается для каждой точки. 1:1 с venues. Создаётся автоматически при создании точки с дефолтными значениями';

COMMENT ON COLUMN "app_configs"."theme_preset" IS 'Фронт читает при загрузке PWA и применяет CSS-переменные темы. Менять только через конструктор в админке';

COMMENT ON COLUMN "app_configs"."accent_color" IS 'HEX цвет кнопок и акцентов. Валидировать: /^#[0-9A-Fa-f]{6}$/. Напр. #E74C3C';

COMMENT ON COLUMN "app_configs"."logo_url" IS 'URL логотипа в S3/CDN (PNG или SVG). NULL = заглушка с первой буквой названия точки. Не хранить base64 в БД';

COMMENT ON COLUMN "app_configs"."banner_url" IS 'URL главного баннера на главной PWA. NULL = без баннера. Рекомендуемый размер: 1200×400px';

COMMENT ON COLUMN "app_configs"."address_manual" IS 'Ручная правка адреса (переопределяет venues.address). Перекрывает venues.address только в UI PWA. venues.address по-прежнему используется для геокодинга';

COMMENT ON COLUMN "app_configs"."updated_at" IS 'Проставлять при каждом изменении. Используется для инвалидации SW-кэша браузера';

COMMENT ON TABLE "users" IS 'Только идентификация. Роли и привязка к точкам — в таблице staff_roles. Клиент = users без записи в staff_roles.';

COMMENT ON COLUMN "users"."id" IS 'UUID пользователя. Хранится в JWT. Из сессии: sessions.token → sessions.user_id';

COMMENT ON COLUMN "users"."phone" IS 'E.164 формат: +79161234567. Нормализовывать при сохранении (убирать пробелы, скобки). Индекс idx_users_phone обязателен';

COMMENT ON COLUMN "users"."full_name" IS 'Имя клиента. NULL при регистрации — запрашиваем после первого входа. Используется в истории заказов и базе клиентов';

COMMENT ON COLUMN "users"."birth_date" IS 'Дата рождения. NULL — заполняется опционально в профиле. Для именинных акций. Хранить без времени (тип DATE)';

COMMENT ON COLUMN "users"."iiko_guest_id" IS 'ID гостя в iikoCard. NULL до первого заказа. При первом заказе: создать гостя в iiko → сохранить ID → создать заказ';

COMMENT ON COLUMN "users"."is_active" IS 'false = мягкая блокировка или удаление по запросу (152-ФЗ). При false: 401 на все запросы. Данные не удалять';

COMMENT ON COLUMN "users"."created_at" IS 'Дата регистрации. Только для чтения';

COMMENT ON TABLE "staff_roles" IS 'Персонал: один user может быть admin в одном заведении и manager в другом.';

COMMENT ON COLUMN "staff_roles"."id" IS 'UUID записи роли';

COMMENT ON COLUMN "staff_roles"."user_id" IS 'FK на пользователя. Один человек может иметь несколько записей для разных точек';

COMMENT ON COLUMN "staff_roles"."venue_id" IS 'FK на точку. При авторизации: читаем все staff_roles пользователя, строим список доступных точек';

COMMENT ON COLUMN "staff_roles"."role" IS 'admin: полный доступ к точке включая биллинг. manager: стоп-листы и заказы. cashier: только список заказов';

COMMENT ON COLUMN "staff_roles"."is_active" IS 'false = сотрудник уволен. НЕ удалять запись. При авторизации: WHERE is_active = true';

COMMENT ON COLUMN "staff_roles"."created_at" IS 'Дата назначения на роль';

COMMENT ON COLUMN "user_sessions"."id" IS 'UUID сессии';

COMMENT ON COLUMN "user_sessions"."user_id" IS 'FK на пользователя. При выходе: DELETE WHERE user_id = ? AND id = ?';

COMMENT ON COLUMN "user_sessions"."token" IS 'Случайный непрозрачный токен (не JWT). Можно мгновенно отозвать. Индекс idx_sessions_token критически важен';

COMMENT ON COLUMN "user_sessions"."device_info" IS 'PWA, UA и платформа устройства. Формат: {"ua":"Mozilla...","platform":"Android","pwa":true}. Для аналитики и списка устройств';

COMMENT ON COLUMN "user_sessions"."expires_at" IS 'Клиенты: now()+30d, персонал: now()+8h. Fоновый job: DELETE WHERE expires_at < now(). Продлевать при каждом запросе';

COMMENT ON COLUMN "user_sessions"."created_at" IS 'Дата создания сессии (момент входа)';

COMMENT ON COLUMN "categories"."id" IS 'UUID категории';

COMMENT ON COLUMN "categories"."venue_id" IS 'Привязка к точке. При запросе меню: WHERE venue_id = ? AND is_visible = true AND deleted_at IS NULL';

COMMENT ON COLUMN "categories"."external_id" IS 'UUID категории в iiko Cloud. Индекс: idx_categories_external. Матчинг при синхронизации. NULL = категория создана вручную в нашей админке';

COMMENT ON COLUMN "categories"."parent_id" IS 'Для вложенных категорий. FK на родительскую категорию. NULL = верхний уровень. Для дерева: рекурсивный WITH RECURSIVE запрос';

COMMENT ON COLUMN "categories"."sort_order" IS 'Порядок среди соседних категорий. ORDER BY sort_order ASC. При drag&drop: batch UPDATE';

COMMENT ON COLUMN "categories"."is_visible" IS 'Чекбокс в админке: Скрыть/Показать. false = скрыта от клиентов (напр. "Хозтовары"). Только для чтения — видна в админке';

COMMENT ON COLUMN "categories"."deleted_at" IS 'Soft delete. При удалении из iiko: SET deleted_at = now(). Все SELECT: WHERE deleted_at IS NULL';

COMMENT ON COLUMN "category_translations"."category_id" IS 'FK на категорию';

COMMENT ON COLUMN "category_translations"."language_code" IS 'Код языка. При запросе меню: JOIN WHERE language_code = язык из Accept-Language';

COMMENT ON COLUMN "category_translations"."title" IS 'Название категории на данном языке. Напр. "Горячие напитки", "Hot Beverages". ru заполняется из iiko, остальные вручную';

COMMENT ON COLUMN "products"."id" IS 'UUID товара';

COMMENT ON COLUMN "products"."venue_id" IS 'FK на точку. Индекс idx_products_venue_cat_avail используется при каждом запросе меню';

COMMENT ON COLUMN "products"."category_id" IS 'FK на категорию. ON DELETE RESTRICT — нельзя удалить категорию пока в ней есть товары';

COMMENT ON COLUMN "products"."external_id" IS 'UUID из iiko. МАТЧИНГ СТРОГО ПО UUID при синхронизации. Нет в нашей БД → INSERT, есть → UPDATE, нет в iiko → soft delete. Индекс: idx_products_external';

COMMENT ON COLUMN "products"."price" IS 'Текущая цена в рублях из iiko. Обновляется при каждой синхронизации. На цену в заказе не влияет (снапшот в order_items)';

COMMENT ON COLUMN "products"."image_url" IS 'URL фото из iiko Cloud или нашего S3. NULL = заглушка с логотипом точки. Не хранить base64';

COMMENT ON COLUMN "products"."is_available" IS 'Стоп-лист из iiko или ручной. false = стоп-лист. Обновляется каждые 60 сек из iiko. Проверяется дважды: при показе меню и при нажатии "Оплатить"';

COMMENT ON COLUMN "products"."sort_order" IS 'Порядок отображения внутри категории. Берётся из iiko при синхронизации. ORDER BY sort_order ASC';

COMMENT ON COLUMN "products"."weight" IS 'Вес порции в граммах из iiko. NULL если не заполнено в iiko. Отображается в карточке: "250 г"';

COMMENT ON COLUMN "products"."calories" IS 'ккал на 100г из iiko. NULL если не заполнено. Блок КБЖУ показываем только если calories IS NOT NULL';

COMMENT ON COLUMN "products"."protein" IS 'Белки на 100г из iiko. Обновляется при ежедневной полной синхронизации';

COMMENT ON COLUMN "products"."fat" IS 'Жиры на 100г из iiko';

COMMENT ON COLUMN "products"."carbs" IS 'Углеводы на 100г из iiko';

COMMENT ON COLUMN "products"."deleted_at" IS 'Soft delete. При исчезновении из iiko: SET deleted_at = now(), is_available = false. Все SELECT: WHERE deleted_at IS NULL';

COMMENT ON COLUMN "product_translations"."product_id" IS 'FK на товар';

COMMENT ON COLUMN "product_translations"."language_code" IS 'Код языка перевода';

COMMENT ON COLUMN "product_translations"."name" IS 'Название товара. Используется в поиске. GIN-индекс на to_tsvector создаётся отдельной миграцией';

COMMENT ON COLUMN "product_translations"."description" IS 'Описание товара. NULL разрешён. Отображается в карточке товара при нажатии на него';

COMMENT ON COLUMN "product_modifiers"."id" IS 'UUID модификатора';

COMMENT ON COLUMN "product_modifiers"."product_id" IS 'FK на товар. При загрузке карточки: SELECT * FROM product_modifiers WHERE product_id = ?';

COMMENT ON COLUMN "product_modifiers"."external_id" IS 'UUID модификатора в iiko. NULL если создан вручную. Матчинг при синхронизации аналогично products.external_id';

COMMENT ON COLUMN "product_modifiers"."extra_price" IS '0 = бесплатный модификатор (напр. "Без сахара"). Прибавляется к цене. Снапшот хранится в order_item_modifiers';

COMMENT ON COLUMN "product_modifiers"."is_required" IS 'true = клиент не может добавить товар в корзину без выбора. Напр. группа "Размер": S/M/L обязательна';

COMMENT ON COLUMN "product_modifiers"."max_quantity" IS '1 = чекбокс. >1 = можно выбрать несколько (напр. количество сиропов до 3). Валидировать на бэкенде';

COMMENT ON COLUMN "product_modifiers"."group_name" IS 'Название группы для UI: "Размер", "Добавки", "Основа". NULL = одиночный модификатор. Одинаковый group_name = одна визуальная группа';

COMMENT ON COLUMN "modifier_translations"."modifier_id" IS 'FK на модификатор';

COMMENT ON COLUMN "modifier_translations"."language_code" IS 'Код языка перевода';

COMMENT ON COLUMN "modifier_translations"."name" IS 'Название добавки на данном языке. Напр. "Большой", "Large", "Без сахара"';

COMMENT ON TABLE "product_recommendations" IS 'Настраивается вручную в админке. Блок "С этим также берут" на странице товара.';

COMMENT ON COLUMN "product_recommendations"."product_id" IS 'Товар, на странице которого показывается блок. SELECT recommended_id WHERE product_id = ? ORDER BY sort_order';

COMMENT ON COLUMN "product_recommendations"."recommended_id" IS 'Рекомендуемый товар. При запросе: проверять is_available = true — не показывать стоп-лист в upsell';

COMMENT ON COLUMN "product_recommendations"."sort_order" IS 'Порядок в блоке "С этим также берут". Обновляется drag&drop в админке. ORDER BY sort_order ASC';

COMMENT ON COLUMN "orders"."id" IS 'UUID заказа. Показывается клиенту как номер: последние 4 символа UUID';

COMMENT ON COLUMN "orders"."venue_id" IS 'FK на точку. Все заказы точки: WHERE venue_id = ? AND deleted_at IS NULL';

COMMENT ON COLUMN "orders"."user_id" IS 'FK на клиента. История клиента: WHERE user_id = ? ORDER BY created_at DESC LIMIT 20';

COMMENT ON COLUMN "orders"."status" IS 'Менять только через сервисные методы, не UPDATE напрямую — иначе пуши и логи не сработают';

COMMENT ON COLUMN "orders"."order_type" IS 'В MVP всегда takeaway. dine_in и delivery — для будущих версий';

COMMENT ON COLUMN "orders"."total_amount" IS 'Итог после скидки. Формула: SUM(unit_price*qty) + SUM(modifier.extra_price) - discount_amount. Фиксируется при создании';

COMMENT ON COLUMN "orders"."scheduled_at" IS 'Предзаказ на время.NULL = "как можно скорее". Если задано: шаг 15 мин, только текущий или следующий день. Передаётся в iiko как время доставки';

COMMENT ON COLUMN "orders"."comment" IS 'Комментарий клиента к заказу, "Без лука, пожалуйста". Передаётся в iiko в поле комментария. Показывается сотруднику в карточке заказа';

COMMENT ON COLUMN "orders"."cutlery_needed" IS 'Чекбокс "Нужны приборы" в корзине. Передаётся в iiko как тег или комментарий';

COMMENT ON COLUMN "orders"."cancelled_reason" IS 'NULL если заказ не отменён. Заполняется автоматически ("Таймаут оплаты") или вручную сотрудником. Для аналитики потерь';

COMMENT ON COLUMN "orders"."refund_id" IS 'ID возврата в платёжной системе. NULL пока возврата не было. При наличии: показываем статус "Возврат выполнен"';

COMMENT ON COLUMN "orders"."refunded_at" IS 'Момент выполнения возврата. NULL если возврата не было';

COMMENT ON COLUMN "orders"."payment_id" IS 'FK → payments.id (nullable, проставляется после создания платежа). NULL до инициации платежа. Порядок: создать заказ → создать payment → обновить orders.payment_id';

COMMENT ON COLUMN "orders"."pos_order_id" IS 'UUID заказа от iiko. NULL до успешной отправки. Используется для отмены в iiko и запроса статуса готовности';

COMMENT ON COLUMN "orders"."pos_sync_status" IS 'Воркер каждые 5 сек: WHERE pos_sync_status = NOT_SENT AND status = COOKING. Retry до 3 раз при ERROR';

COMMENT ON COLUMN "orders"."order_source" IS 'Передаётся в iiko как orderSource для маркетинговой аналитики. По ТЗ = название нашего приложения';

COMMENT ON COLUMN "orders"."promocode_id" IS 'NULL если промокод не применялся. При применении: проверить валидность → INCREMENT used_count → сохранить ID';

COMMENT ON COLUMN "orders"."discount_amount" IS 'Сумма скидки по промокоду. total_amount уже с учётом скидки. Поле хранится для прозрачности аналитики промокодов';

COMMENT ON COLUMN "orders"."created_at" IS 'Время создания заказа клиентом';

COMMENT ON COLUMN "orders"."deleted_at" IS 'Soft delete. Заказы НИКОГДА не удалять физически — финансовая история. deleted_at только для технически ошибочных записей';

COMMENT ON COLUMN "order_items"."order_id" IS 'FK на заказ. Все позиции заказа: WHERE order_id = ?';

COMMENT ON COLUMN "order_items"."product_id" IS 'FK на товар. Использовать LEFT JOIN — товар может быть удалён из меню';

COMMENT ON COLUMN "order_items"."quantity" IS 'Количество единиц. Минимум 1. Нельзя изменить после создания заказа';

COMMENT ON COLUMN "order_items"."unit_price" IS 'Цена 1 единицы товара на момент оплаты (снапшот из products.price). Не зависит от будущих изменений цены';

COMMENT ON COLUMN "order_items"."product_name_snapshot" IS 'Название товара на языке клиента на момент заказа. ОБЯЗАТЕЛЬНО заполнять. Без снапшота история пустеет при удалении товара';

COMMENT ON TABLE "order_item_modifiers" IS 'Нормализованная замена jsonb selected_modifiers. Позволяет строить аналитику по добавкам.';

COMMENT ON COLUMN "order_item_modifiers"."id" IS 'UUID записи модификатора заказа';

COMMENT ON COLUMN "order_item_modifiers"."order_item_id" IS 'FK на позицию заказа. Все добавки позиции: WHERE order_item_id = ?';

COMMENT ON COLUMN "order_item_modifiers"."modifier_id" IS 'FK на модификатор. NULLABLE — модификатор может быть удалён из меню. Всегда LEFT JOIN, не INNER';

COMMENT ON COLUMN "order_item_modifiers"."name_snapshot" IS 'Название добавки на момент заказа. ОБЯЗАТЕЛЬНО заполнять. Аналогично product_name_snapshot — защита от удаления';

COMMENT ON COLUMN "order_item_modifiers"."extra_price" IS 'Цена добавки на момент заказа (снапшот). Итог позиции: (unit_price + SUM(extra_price)) * quantity';

COMMENT ON COLUMN "payments"."id" IS 'UUID платежа в нашей системе';

COMMENT ON COLUMN "payments"."order_id" IS '1:1 с заказом. UNIQUE гарантирует один платёж на заказ';

COMMENT ON COLUMN "payments"."external_id" IS 'ID транзакции в Тинькофф/ЮKassa. Используется для идентификации входящих webhook: WHERE external_id = ?';

COMMENT ON COLUMN "payments"."status" IS 'Обновлять ТОЛЬКО через webhook от банка или polling. Никогда не менять вручную без подтверждения банка';

COMMENT ON COLUMN "payments"."amount" IS 'Сумма платежа в рублях. Должна совпадать с orders.total_amount';

COMMENT ON COLUMN "payments"."provider" IS 'Название провайдера: tinkoff, yookassa, sbp. Нужен при возврате — используем правильный API';

COMMENT ON COLUMN "payments"."paid_at" IS 'Момент фактического списания от банка (из данных webhook). NULL до SUCCESS. Не ставить своё время — только из данных банка';

COMMENT ON COLUMN "payments"."created_at" IS 'Время создания платёжной записи (инициации платежа)';

COMMENT ON TABLE "pos_sync_logs" IS 'TTL: DELETE WHERE created_at < NOW() - INTERVAL 14 days (pg_cron job).';

COMMENT ON COLUMN "pos_sync_logs"."id" IS 'UUID записи лога';

COMMENT ON COLUMN "pos_sync_logs"."venue_id" IS 'FK на точку. Для дебага: WHERE venue_id = ? ORDER BY created_at DESC';

COMMENT ON COLUMN "pos_sync_logs"."order_id" IS 'FK на заказ. Все попытки отправки заказа: WHERE order_id = ? ORDER BY created_at DESC';

COMMENT ON COLUMN "pos_sync_logs"."request_payload" IS 'Полное тело запроса отправленного в iiko. Логировать ДО отправки. НЕ включать api_key в payload';

COMMENT ON COLUMN "pos_sync_logs"."response_payload" IS 'Полный ответ iiko (JSON). NULL если нет ответа (таймаут). При ERROR статусе — главный источник для диагностики';

COMMENT ON COLUMN "pos_sync_logs"."http_status_code" IS 'HTTP статус ответа iiko. 200 = успех → SENT. 500/503 = retry до 3 раз через 30 сек. 401 = невалидный api_key → оповестить владельца. NULL = таймаут';

COMMENT ON COLUMN "pos_sync_logs"."created_at" IS 'Время попытки синхронизации. TTL 14 дней: pg_cron DELETE WHERE created_at < NOW() - INTERVAL 14 days';

COMMENT ON COLUMN "pos_settings"."venue_id" IS '1:1 с точкой. Создаётся при подключении iiko в онбординге';

COMMENT ON COLUMN "pos_settings"."api_key_encrypted" IS 'API-ключ iikoTransport. ЗАШИФРОВАТЬ через pgcrypto/Vault. В БД только blob. Ключ шифрования в env. Никогда не логировать';

COMMENT ON COLUMN "pos_settings"."terminal_id" IS 'UUID терминала iiko, куда "падают" заказы. Если несколько терминалов — владелец выбирает нужный при настройке';

COMMENT ON COLUMN "pos_settings"."external_org_id" IS 'UUID организации в облаке iiko (≠ наш organizations.id). Передаётся в каждый запрос к iikoCloud как organizationId';

COMMENT ON COLUMN "pos_settings"."sync_interval_minutes" IS 'Интервал полной синхронизации меню в минутах. Планировщик запускает каждые N минут. Полная (картинки) — раз в сутки';

COMMENT ON COLUMN "pos_settings"."stoplist_poll_seconds" IS 'Интервал опроса стоп-листа в секундах. Отдельный частый воркер обновляет products.is_available всех товаров точки';

COMMENT ON COLUMN "pos_settings"."last_menu_sync_at" IS 'Timestamp последней успешной синхронизации меню. Обновлять при каждом успехе. Отображать в админке: "Обновлено 2 мин назад"';

COMMENT ON COLUMN "pos_settings"."last_stoplist_sync_at" IS 'Timestamp последнего опроса стоп-листа. Если > 5 мин назад — алерт в Sentry. NULL = синхронизации ещё не было';

COMMENT ON COLUMN "pos_settings"."updated_at" IS 'Дата последнего изменения настроек. Проставлять при каждом UPDATE';

COMMENT ON COLUMN "promocodes"."id" IS 'UUID промокода';

COMMENT ON COLUMN "promocodes"."venue_id" IS 'FK на точку. Промокод действует только в этой точке. UNIQUE(venue_id, code)';

COMMENT ON COLUMN "promocodes"."code" IS 'Код промокода: WELCOME, SUMMER20. Нормализовывать: UPPER(TRIM(code)) при записи и при поиске';

COMMENT ON COLUMN "promocodes"."discount_type" IS 'PERCENT или FIXED. Определяет логику вычисления discount_amount в заказе';

COMMENT ON COLUMN "promocodes"."discount_value" IS 'PERCENT: 10 = 10%. FIXED: 150 = 150 руб. Для FIXED: discount = MIN(value, total) — защита от отрицательной суммы';

COMMENT ON COLUMN "promocodes"."min_order_amount" IS 'Минимальная сумма заказа для применения промокода. 0 = без ограничений. Проверять до применения скидки';

COMMENT ON COLUMN "promocodes"."max_uses" IS 'Лимит использований. NULL = безлимит. При применении: SELECT FOR UPDATE → проверить → INCREMENT. Защита от race condition';

COMMENT ON COLUMN "promocodes"."used_count" IS 'Счётчик использований. INCREMENT в транзакции с SELECT FOR UPDATE. Никогда не уменьшать вручную';

COMMENT ON COLUMN "promocodes"."expires_at" IS 'Срок действия. NULL = бессрочно. Проверка: WHERE (expires_at IS NULL OR expires_at > NOW()) AND is_active = true';

COMMENT ON COLUMN "promocodes"."is_active" IS 'false = промокод отключён вручную. При деактивации существующие заказы с промокодом не затрагиваются';

COMMENT ON COLUMN "promocodes"."created_at" IS 'Дата создания промокода';

COMMENT ON TABLE "order_reviews" IS 'Внутренние отзывы — видны только владельцу в админке (ТЗ блок 6).';

COMMENT ON COLUMN "order_reviews"."id" IS 'UUID отзыва';

COMMENT ON COLUMN "order_reviews"."order_id" IS '1:1 с заказом. UNIQUE гарантирует один отзыв на заказ. Предлагать оценку через пуш спустя 1 час';

COMMENT ON COLUMN "order_reviews"."user_id" IS 'FK на клиента';

COMMENT ON COLUMN "order_reviews"."venue_id" IS 'FK на точку. Средний рейтинг: SELECT AVG(rating) WHERE venue_id = ?';

COMMENT ON COLUMN "order_reviews"."rating" IS 'Оценка 1-5. Валидировать: CHECK (rating BETWEEN 1 AND 5). 1 звезда = очень плохо, 5 = отлично';

COMMENT ON COLUMN "order_reviews"."comment" IS 'Текст отзыва. NULL если клиент поставил только оценку без текста. Показывать ТОЛЬКО в админке владельца — не публиковать в PWA';

COMMENT ON COLUMN "order_reviews"."owner_reply" IS 'Ответ владельца. NULL если не ответил. Заполняется в разделе "Отзывы" в админке. В MVP клиенту не доставляется';

COMMENT ON COLUMN "order_reviews"."replied_at" IS 'Время ответа владельца. Проставлять: SET replied_at = now() вместе с owner_reply';

COMMENT ON COLUMN "order_reviews"."created_at" IS 'Время создания отзыва';

COMMENT ON COLUMN "push_subscriptions"."id" IS 'UUID подписки';

COMMENT ON COLUMN "push_subscriptions"."user_id" IS 'FK на пользователя';

COMMENT ON COLUMN "push_subscriptions"."venue_id" IS 'FK на точку. Маркетинговые пуши точки: WHERE venue_id = ? AND is_active = true';

COMMENT ON COLUMN "push_subscriptions"."endpoint" IS 'URL push-сервиса браузера (Google FCM, Apple и др.). Получаем из ServiceWorker.pushManager.subscribe(). Уникален для устройства';

COMMENT ON COLUMN "push_subscriptions"."p256dh" IS 'Публичный ключ шифрования устройства. Не секрет. Передаётся в web-push библиотеку для шифрования payload';

COMMENT ON COLUMN "push_subscriptions"."auth" IS 'Auth-секрет устройства. Не секрет. Передаётся в web-push: sendNotification({endpoint, keys:{p256dh,auth}}, payload)';

COMMENT ON COLUMN "push_subscriptions"."is_active" IS 'false = пользователь отозвал разрешение. Ставить false если браузер вернул 410 Gone при отправке пуша';

COMMENT ON COLUMN "push_subscriptions"."created_at" IS 'Время создания подписки';

COMMENT ON TABLE "telegram_bindings" IS 'Привязка аккаунта к Telegram-боту для уведомлений персоналу и клиентам.';

COMMENT ON COLUMN "telegram_bindings"."id" IS 'UUID привязки';

COMMENT ON COLUMN "telegram_bindings"."user_id" IS 'FK на пользователя. UNIQUE гарантирует один Telegram на аккаунт';

COMMENT ON COLUMN "telegram_bindings"."telegram_chat_id" IS 'Числовой chat_id из Telegram Bot API (bigint — может быть > 2^31). Получаем из message.chat.id при команде /start';

COMMENT ON COLUMN "telegram_bindings"."is_active" IS 'false = пользователь написал /stop боту. При отправке пушей: WHERE is_active = true';

COMMENT ON COLUMN "telegram_bindings"."created_at" IS 'Дата привязки аккаунта к боту';

ALTER TABLE "organizations" ADD FOREIGN KEY ("plan_id") REFERENCES "subscription_plans" ("id");

ALTER TABLE "venues" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");

ALTER TABLE "organization_billing" ADD FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id");

ALTER TABLE "organization_billing" ADD FOREIGN KEY ("plan_id") REFERENCES "subscription_plans" ("id");

ALTER TABLE "app_configs" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

ALTER TABLE "staff_roles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "staff_roles" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

ALTER TABLE "user_sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "categories" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

ALTER TABLE "categories" ADD FOREIGN KEY ("parent_id") REFERENCES "categories" ("id");

ALTER TABLE "category_translations" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "category_translations" ADD FOREIGN KEY ("language_code") REFERENCES "languages" ("code");

ALTER TABLE "products" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "product_translations" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "product_translations" ADD FOREIGN KEY ("language_code") REFERENCES "languages" ("code");

ALTER TABLE "product_modifiers" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "modifier_translations" ADD FOREIGN KEY ("modifier_id") REFERENCES "product_modifiers" ("id");

ALTER TABLE "modifier_translations" ADD FOREIGN KEY ("language_code") REFERENCES "languages" ("code");

ALTER TABLE "product_recommendations" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "product_recommendations" ADD FOREIGN KEY ("recommended_id") REFERENCES "products" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("payment_id") REFERENCES "payments" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("promocode_id") REFERENCES "promocodes" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "order_item_modifiers" ADD FOREIGN KEY ("order_item_id") REFERENCES "order_items" ("id");

ALTER TABLE "order_item_modifiers" ADD FOREIGN KEY ("modifier_id") REFERENCES "product_modifiers" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "pos_sync_logs" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

ALTER TABLE "pos_sync_logs" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "pos_settings" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

ALTER TABLE "promocodes" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

ALTER TABLE "order_reviews" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_reviews" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "order_reviews" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

ALTER TABLE "push_subscriptions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "push_subscriptions" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

ALTER TABLE "telegram_bindings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
