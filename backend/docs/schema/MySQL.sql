CREATE TABLE `organizations` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID юрлица. Передаётся в JWT сотрудника. Все запросы персонала фильтруются через него',
  `name` varchar(255) NOT NULL COMMENT 'Название для договора. Официальное название для договоров и инвойсов. Не путать с venues.name — публичным названием',
  `inn` varchar(12) COMMENT 'ИНН для проверки контрагента (РФ). ИНН: 10 цифр для ООО, 12 для ИП. Валидировать контрольную сумму. NULL разрешён в MVP',
  `is_active` boolean DEFAULT true COMMENT 'Глобальное выключение аккаунта. false = ручная блокировка саппортом. Все venues этой организации перестают работать',
  `created_at` timestamptz DEFAULT (now()) COMMENT 'Дата регистрации. Только для чтения. Используется в аналитике роста клиентской базы',
  `plan_id` uuid NOT NULL COMMENT 'FK → subscription_plans.id'
);

CREATE TABLE `venues` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID точки. Главный FK всей системы — почти каждая таблица ссылается на venue_id',
  `organization_id` uuid NOT NULL COMMENT 'Владелец точки. При авторизации: staff_roles.venue_id → venues → organizations → проверяем billing_status',
  `name` varchar(255) NOT NULL COMMENT 'Публичное название точки для клиентов и шапки PWA. Напр. "Кофе-пойнт ТЦ Мега"',
  `slug` varchar(255) UNIQUE NOT NULL COMMENT 'URL-часть: brand-kazan.foodapp.ru. Только [a-z0-9-], валидировать регуляркой. Менять осторожно — старые QR-коды перестанут работать',
  `address` text COMMENT 'Физический адрес. Используется для геокодинга. Для отображения в PWA может быть перекрыт app_configs.address_manual',
  `latitude` numeric(10,7) COMMENT 'Координаты для карты. GPS-широта. Заполнять через геокодер при сохранении адреса. Нужна для карты и функции "рядом со мной"',
  `longitude` numeric(10,7) COMMENT 'GPS-долгота. Заполнять через геокодер при сохранении адреса',
  `is_active` boolean DEFAULT true COMMENT 'false = точка закрыта навсегда или в ремонте. PWA возвращает 404. Не путать с is_emergency_stop',
  `is_emergency_stop` boolean DEFAULT false COMMENT 'Кнопка "Запара": мгновенно отключает приём заказов. 1 клик в админке. При true: кнопка оплаты в PWA неактивна, API заказов возвращает 503',
  `min_order_amount` numeric(12,2) DEFAULT 0 COMMENT 'Минимальная сумма заказа из ТЗ блок 4. 0 = нет ограничений. Проверяется на бэкенде при нажатии "Оплатить"',
  `avg_cooking_minutes` integer DEFAULT 20 COMMENT 'Среднее время готовности (мин) — показывается клиенту. "Время ожидания: ~20 мин". Не влияет на логику',
  `work_hours_json` jsonb COMMENT 'Расписание-оверрайд. Формат: {"mon":{"open":"09:00","close":"22:00"}, ...}. Null = берём из iiko.',
  `billing_status` ENUM ('ACTIVE', 'PAST_DUE', 'BLOCKED') DEFAULT 'ACTIVE' COMMENT 'Синхронизировать с organization_billing при каждом изменении. Fron-cron проверяет paid_until каждый час',
  `paid_until` timestamptz COMMENT 'До какой даты оплачен сервис. При успешной оплате: paid_until += 1 month. NULL = триальный период',
  `created_at` timestamptz DEFAULT (now()) COMMENT 'Дата создания точки. Только для чтения',
  `updated_at` timestamptz COMMENT 'Проставлять вручную: SET updated_at = now() при каждом UPDATE. Используется для инвалидации кэша PWA',
  `deleted_at` timestamptz COMMENT 'Soft delete. NULL = активна. Никогда не делать DELETE FROM venues. Все SELECT: WHERE deleted_at IS NULL'
);

CREATE TABLE `organization_billing` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID записи биллинга. Текущий план = MAX(created_at) для данной организации',
  `organization_id` uuid NOT NULL COMMENT 'При смене плана: НЕ обновлять старую запись, создавать новую — для хранения истории',
  `plan_id` uuid NOT NULL COMMENT 'JOIN на subscription_plans для получения цены и лимитов при генерации инвойса',
  `billing_status` ENUM ('ACTIVE', 'PAST_DUE', 'BLOCKED') DEFAULT 'ACTIVE' COMMENT 'Дублирует venues.billing_status. Синхронизировать оба поля при любом изменении статуса',
  `paid_until` timestamptz COMMENT 'При оплате: paid_until = MAX(paid_until, now()) + 1 month. Не затирать будущую дату при досрочной оплате',
  `payment_method_token` text COMMENT 'Токен привязанной карты для рекуррентных платежей',
  `last_invoice_at` timestamptz COMMENT 'Когда последний раз генерировался инвойс. Обновлять при каждой генерации. Для отладки авто-биллинга',
  `created_at` timestamptz DEFAULT (now()) COMMENT 'Дата начала действия этой записи биллинга'
);

CREATE TABLE `subscription_plans` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID плана. Зашит в organizations.plan_id',
  `name` varchar(255) NOT NULL COMMENT 'Название для отображения: Starter, Pro, Enterprise',
  `price_per_month` numeric(12,2) NOT NULL COMMENT 'Цена в рублях. При изменении цены — создавать новый план, не обновлять. Старые организации остаются на старой цене',
  `features_json` jsonb COMMENT 'Лимиты: макс. блюд, наличие iiko, кол-во точек и т.д. Лимиты и флаги. Формат: {"max_products":100,"iiko_integration":true,"max_venues":1,"push_marketing":false}',
  `is_active` boolean DEFAULT true COMMENT 'Архивирование устаревших планов. false = план архивирован, не показывать новым клиентам. Существующие продолжают им пользоваться',
  `created_at` timestamptz DEFAULT (now()) COMMENT 'Дата создания плана'
);

CREATE TABLE `languages` (
  `code` varchar(5) PRIMARY KEY COMMENT 'ISO 639-1 код языка: ru, en, ge. Заполняется при деплое seed-скриптом. Не изменять в runtime',
  `name` varchar(255) NOT NULL COMMENT 'Читаемое название для выпадающего списка в админке: Русский, English, ქართული'
);

CREATE TABLE `app_configs` (
  `venue_id` uuid PRIMARY KEY COMMENT 'Дизайн настраивается для каждой точки. 1:1 с venues. Создаётся автоматически при создании точки с дефолтными значениями',
  `theme_preset` ENUM ('light', 'dark', 'coffee', 'fastfood') DEFAULT 'light' COMMENT 'Фронт читает при загрузке PWA и применяет CSS-переменные темы. Менять только через конструктор в админке',
  `accent_color` varchar(7) DEFAULT '#000000' COMMENT 'HEX цвет кнопок и акцентов. Валидировать: /^#[0-9A-Fa-f]{6}$/. Напр. #E74C3C',
  `logo_url` text COMMENT 'URL логотипа в S3/CDN (PNG или SVG). NULL = заглушка с первой буквой названия точки. Не хранить base64 в БД',
  `banner_url` text COMMENT 'URL главного баннера на главной PWA. NULL = без баннера. Рекомендуемый размер: 1200×400px',
  `address_manual` text COMMENT 'Ручная правка адреса (переопределяет venues.address). Перекрывает venues.address только в UI PWA. venues.address по-прежнему используется для геокодинга',
  `updated_at` timestamptz COMMENT 'Проставлять при каждом изменении. Используется для инвалидации SW-кэша браузера'
);

CREATE TABLE `users` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID пользователя. Хранится в JWT. Из сессии: sessions.token → sessions.user_id',
  `phone` varchar(255) UNIQUE NOT NULL COMMENT 'E.164 формат: +79161234567. Нормализовывать при сохранении (убирать пробелы, скобки). Индекс idx_users_phone обязателен',
  `full_name` varchar(255) COMMENT 'Имя клиента. NULL при регистрации — запрашиваем после первого входа. Используется в истории заказов и базе клиентов',
  `birth_date` date COMMENT 'Дата рождения. NULL — заполняется опционально в профиле. Для именинных акций. Хранить без времени (тип DATE)',
  `iiko_guest_id` uuid COMMENT 'ID гостя в iikoCard. NULL до первого заказа. При первом заказе: создать гостя в iiko → сохранить ID → создать заказ',
  `is_active` boolean DEFAULT true COMMENT 'false = мягкая блокировка или удаление по запросу (152-ФЗ). При false: 401 на все запросы. Данные не удалять',
  `created_at` timestamptz DEFAULT (now()) COMMENT 'Дата регистрации. Только для чтения'
);

CREATE TABLE `staff_roles` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID записи роли',
  `user_id` uuid NOT NULL COMMENT 'FK на пользователя. Один человек может иметь несколько записей для разных точек',
  `venue_id` uuid NOT NULL COMMENT 'FK на точку. При авторизации: читаем все staff_roles пользователя, строим список доступных точек',
  `role` ENUM ('admin', 'manager', 'cashier', 'customer') NOT NULL COMMENT 'admin: полный доступ к точке включая биллинг. manager: стоп-листы и заказы. cashier: только список заказов',
  `is_active` boolean DEFAULT true COMMENT 'false = сотрудник уволен. НЕ удалять запись. При авторизации: WHERE is_active = true',
  `created_at` timestamptz DEFAULT (now()) COMMENT 'Дата назначения на роль'
);

CREATE TABLE `user_sessions` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID сессии',
  `user_id` uuid NOT NULL COMMENT 'FK на пользователя. При выходе: DELETE WHERE user_id = ? AND id = ?',
  `token` text UNIQUE NOT NULL COMMENT 'Случайный непрозрачный токен (не JWT). Можно мгновенно отозвать. Индекс idx_sessions_token критически важен',
  `device_info` jsonb COMMENT 'PWA, UA и платформа устройства. Формат: {"ua":"Mozilla...","platform":"Android","pwa":true}. Для аналитики и списка устройств',
  `expires_at` timestamptz NOT NULL COMMENT 'Клиенты: now()+30d, персонал: now()+8h. Fоновый job: DELETE WHERE expires_at < now(). Продлевать при каждом запросе',
  `created_at` timestamptz DEFAULT (now()) COMMENT 'Дата создания сессии (момент входа)'
);

CREATE TABLE `categories` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID категории',
  `venue_id` uuid NOT NULL COMMENT 'Привязка к точке. При запросе меню: WHERE venue_id = ? AND is_visible = true AND deleted_at IS NULL',
  `external_id` uuid COMMENT 'UUID категории в iiko Cloud. Индекс: idx_categories_external. Матчинг при синхронизации. NULL = категория создана вручную в нашей админке',
  `parent_id` uuid COMMENT 'Для вложенных категорий. FK на родительскую категорию. NULL = верхний уровень. Для дерева: рекурсивный WITH RECURSIVE запрос',
  `sort_order` integer DEFAULT 0 COMMENT 'Порядок среди соседних категорий. ORDER BY sort_order ASC. При drag&drop: batch UPDATE',
  `is_visible` boolean DEFAULT true COMMENT 'Чекбокс в админке: Скрыть/Показать. false = скрыта от клиентов (напр. "Хозтовары"). Только для чтения — видна в админке',
  `deleted_at` timestamptz COMMENT 'Soft delete. При удалении из iiko: SET deleted_at = now(). Все SELECT: WHERE deleted_at IS NULL'
);

CREATE TABLE `category_translations` (
  `category_id` uuid COMMENT 'FK на категорию',
  `language_code` varchar(255) COMMENT 'Код языка. При запросе меню: JOIN WHERE language_code = язык из Accept-Language',
  `title` varchar(255) NOT NULL COMMENT 'Название категории на данном языке. Напр. "Горячие напитки", "Hot Beverages". ru заполняется из iiko, остальные вручную',
  PRIMARY KEY (`category_id`, `language_code`)
);

CREATE TABLE `products` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID товара',
  `venue_id` uuid NOT NULL COMMENT 'FK на точку. Индекс idx_products_venue_cat_avail используется при каждом запросе меню',
  `category_id` uuid COMMENT 'FK на категорию. ON DELETE RESTRICT — нельзя удалить категорию пока в ней есть товары',
  `external_id` uuid NOT NULL COMMENT 'UUID из iiko. МАТЧИНГ СТРОГО ПО UUID при синхронизации. Нет в нашей БД → INSERT, есть → UPDATE, нет в iiko → soft delete. Индекс: idx_products_external',
  `price` numeric(12,2) NOT NULL COMMENT 'Текущая цена в рублях из iiko. Обновляется при каждой синхронизации. На цену в заказе не влияет (снапшот в order_items)',
  `image_url` text COMMENT 'URL фото из iiko Cloud или нашего S3. NULL = заглушка с логотипом точки. Не хранить base64',
  `is_available` boolean DEFAULT true COMMENT 'Стоп-лист из iiko или ручной. false = стоп-лист. Обновляется каждые 60 сек из iiko. Проверяется дважды: при показе меню и при нажатии "Оплатить"',
  `sort_order` integer DEFAULT 0 COMMENT 'Порядок отображения внутри категории. Берётся из iiko при синхронизации. ORDER BY sort_order ASC',
  `weight` integer COMMENT 'Вес порции в граммах из iiko. NULL если не заполнено в iiko. Отображается в карточке: "250 г"',
  `calories` numeric(8,2) COMMENT 'ккал на 100г из iiko. NULL если не заполнено. Блок КБЖУ показываем только если calories IS NOT NULL',
  `protein` numeric(8,2) COMMENT 'Белки на 100г из iiko. Обновляется при ежедневной полной синхронизации',
  `fat` numeric(8,2) COMMENT 'Жиры на 100г из iiko',
  `carbs` numeric(8,2) COMMENT 'Углеводы на 100г из iiko',
  `deleted_at` timestamptz COMMENT 'Soft delete. При исчезновении из iiko: SET deleted_at = now(), is_available = false. Все SELECT: WHERE deleted_at IS NULL'
);

CREATE TABLE `product_translations` (
  `product_id` uuid COMMENT 'FK на товар',
  `language_code` varchar(255) COMMENT 'Код языка перевода',
  `name` varchar(255) NOT NULL COMMENT 'Название товара. Используется в поиске. GIN-индекс на to_tsvector создаётся отдельной миграцией',
  `description` text COMMENT 'Описание товара. NULL разрешён. Отображается в карточке товара при нажатии на него',
  PRIMARY KEY (`product_id`, `language_code`)
);

CREATE TABLE `product_modifiers` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID модификатора',
  `product_id` uuid NOT NULL COMMENT 'FK на товар. При загрузке карточки: SELECT * FROM product_modifiers WHERE product_id = ?',
  `external_id` uuid COMMENT 'UUID модификатора в iiko. NULL если создан вручную. Матчинг при синхронизации аналогично products.external_id',
  `extra_price` numeric(12,2) DEFAULT 0 COMMENT '0 = бесплатный модификатор (напр. "Без сахара"). Прибавляется к цене. Снапшот хранится в order_item_modifiers',
  `is_required` boolean DEFAULT false COMMENT 'true = клиент не может добавить товар в корзину без выбора. Напр. группа "Размер": S/M/L обязательна',
  `max_quantity` integer DEFAULT 1 COMMENT '1 = чекбокс. >1 = можно выбрать несколько (напр. количество сиропов до 3). Валидировать на бэкенде',
  `group_name` varchar(255) COMMENT 'Название группы для UI: "Размер", "Добавки", "Основа". NULL = одиночный модификатор. Одинаковый group_name = одна визуальная группа'
);

CREATE TABLE `modifier_translations` (
  `modifier_id` uuid COMMENT 'FK на модификатор',
  `language_code` varchar(255) COMMENT 'Код языка перевода',
  `name` varchar(255) NOT NULL COMMENT 'Название добавки на данном языке. Напр. "Большой", "Large", "Без сахара"',
  PRIMARY KEY (`modifier_id`, `language_code`)
);

CREATE TABLE `product_recommendations` (
  `product_id` uuid COMMENT 'Товар, на странице которого показывается блок. SELECT recommended_id WHERE product_id = ? ORDER BY sort_order',
  `recommended_id` uuid COMMENT 'Рекомендуемый товар. При запросе: проверять is_available = true — не показывать стоп-лист в upsell',
  `sort_order` integer DEFAULT 0 COMMENT 'Порядок в блоке "С этим также берут". Обновляется drag&drop в админке. ORDER BY sort_order ASC',
  PRIMARY KEY (`product_id`, `recommended_id`)
);

CREATE TABLE `orders` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID заказа. Показывается клиенту как номер: последние 4 символа UUID',
  `venue_id` uuid NOT NULL COMMENT 'FK на точку. Все заказы точки: WHERE venue_id = ? AND deleted_at IS NULL',
  `user_id` uuid NOT NULL COMMENT 'FK на клиента. История клиента: WHERE user_id = ? ORDER BY created_at DESC LIMIT 20',
  `status` ENUM ('NEW', 'WAITING_PAYMENT', 'COOKING', 'READY', 'COMPLETED', 'CANCELLED') DEFAULT 'NEW' COMMENT 'Менять только через сервисные методы, не UPDATE напрямую — иначе пуши и логи не сработают',
  `order_type` ENUM ('takeaway', 'dine_in', 'delivery') DEFAULT 'takeaway' COMMENT 'В MVP всегда takeaway. dine_in и delivery — для будущих версий',
  `total_amount` numeric(12,2) NOT NULL COMMENT 'Итог после скидки. Формула: SUM(unit_price*qty) + SUM(modifier.extra_price) - discount_amount. Фиксируется при создании',
  `scheduled_at` timestamptz COMMENT 'Предзаказ на время.NULL = "как можно скорее". Если задано: шаг 15 мин, только текущий или следующий день. Передаётся в iiko как время доставки',
  `comment` text COMMENT 'Комментарий клиента к заказу, "Без лука, пожалуйста". Передаётся в iiko в поле комментария. Показывается сотруднику в карточке заказа',
  `cutlery_needed` boolean DEFAULT false COMMENT 'Чекбокс "Нужны приборы" в корзине. Передаётся в iiko как тег или комментарий',
  `cancelled_reason` text COMMENT 'NULL если заказ не отменён. Заполняется автоматически ("Таймаут оплаты") или вручную сотрудником. Для аналитики потерь',
  `refund_id` text COMMENT 'ID возврата в платёжной системе. NULL пока возврата не было. При наличии: показываем статус "Возврат выполнен"',
  `refunded_at` timestamptz COMMENT 'Момент выполнения возврата. NULL если возврата не было',
  `payment_id` uuid COMMENT 'FK → payments.id (nullable, проставляется после создания платежа). NULL до инициации платежа. Порядок: создать заказ → создать payment → обновить orders.payment_id',
  `pos_order_id` uuid COMMENT 'UUID заказа от iiko. NULL до успешной отправки. Используется для отмены в iiko и запроса статуса готовности',
  `pos_sync_status` ENUM ('NOT_SENT', 'SENT', 'ERROR') DEFAULT 'NOT_SENT' COMMENT 'Воркер каждые 5 сек: WHERE pos_sync_status = NOT_SENT AND status = COOKING. Retry до 3 раз при ERROR',
  `order_source` varchar(255) DEFAULT 'PWA' COMMENT 'Передаётся в iiko как orderSource для маркетинговой аналитики. По ТЗ = название нашего приложения',
  `promocode_id` uuid COMMENT 'NULL если промокод не применялся. При применении: проверить валидность → INCREMENT used_count → сохранить ID',
  `discount_amount` numeric(12,2) DEFAULT 0 COMMENT 'Сумма скидки по промокоду. total_amount уже с учётом скидки. Поле хранится для прозрачности аналитики промокодов',
  `created_at` timestamptz DEFAULT (now()) COMMENT 'Время создания заказа клиентом',
  `deleted_at` timestamptz COMMENT 'Soft delete. Заказы НИКОГДА не удалять физически — финансовая история. deleted_at только для технически ошибочных записей'
);

CREATE TABLE `order_items` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  `order_id` uuid NOT NULL COMMENT 'FK на заказ. Все позиции заказа: WHERE order_id = ?',
  `product_id` uuid NOT NULL COMMENT 'FK на товар. Использовать LEFT JOIN — товар может быть удалён из меню',
  `quantity` integer NOT NULL COMMENT 'Количество единиц. Минимум 1. Нельзя изменить после создания заказа',
  `unit_price` numeric(12,2) NOT NULL COMMENT 'Цена 1 единицы товара на момент оплаты (снапшот из products.price). Не зависит от будущих изменений цены',
  `product_name_snapshot` varchar(255) COMMENT 'Название товара на языке клиента на момент заказа. ОБЯЗАТЕЛЬНО заполнять. Без снапшота история пустеет при удалении товара'
);

CREATE TABLE `order_item_modifiers` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID записи модификатора заказа',
  `order_item_id` uuid NOT NULL COMMENT 'FK на позицию заказа. Все добавки позиции: WHERE order_item_id = ?',
  `modifier_id` uuid COMMENT 'FK на модификатор. NULLABLE — модификатор может быть удалён из меню. Всегда LEFT JOIN, не INNER',
  `name_snapshot` varchar(255) NOT NULL COMMENT 'Название добавки на момент заказа. ОБЯЗАТЕЛЬНО заполнять. Аналогично product_name_snapshot — защита от удаления',
  `extra_price` numeric(12,2) NOT NULL DEFAULT 0 COMMENT 'Цена добавки на момент заказа (снапшот). Итог позиции: (unit_price + SUM(extra_price)) * quantity'
);

CREATE TABLE `payments` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID платежа в нашей системе',
  `order_id` uuid UNIQUE NOT NULL COMMENT '1:1 с заказом. UNIQUE гарантирует один платёж на заказ',
  `external_id` varchar(255) NOT NULL COMMENT 'ID транзакции в Тинькофф/ЮKassa. Используется для идентификации входящих webhook: WHERE external_id = ?',
  `status` ENUM ('PENDING', 'SUCCESS', 'FAILED', 'REFUNDED', 'CANCELLED') DEFAULT 'PENDING' COMMENT 'Обновлять ТОЛЬКО через webhook от банка или polling. Никогда не менять вручную без подтверждения банка',
  `amount` numeric(12,2) NOT NULL COMMENT 'Сумма платежа в рублях. Должна совпадать с orders.total_amount',
  `provider` varchar(255) COMMENT 'Название провайдера: tinkoff, yookassa, sbp. Нужен при возврате — используем правильный API',
  `paid_at` timestamptz COMMENT 'Момент фактического списания от банка (из данных webhook). NULL до SUCCESS. Не ставить своё время — только из данных банка',
  `created_at` timestamptz DEFAULT (now()) COMMENT 'Время создания платёжной записи (инициации платежа)'
);

CREATE TABLE `pos_sync_logs` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID записи лога',
  `venue_id` uuid NOT NULL COMMENT 'FK на точку. Для дебага: WHERE venue_id = ? ORDER BY created_at DESC',
  `order_id` uuid NOT NULL COMMENT 'FK на заказ. Все попытки отправки заказа: WHERE order_id = ? ORDER BY created_at DESC',
  `request_payload` jsonb COMMENT 'Полное тело запроса отправленного в iiko. Логировать ДО отправки. НЕ включать api_key в payload',
  `response_payload` jsonb COMMENT 'Полный ответ iiko (JSON). NULL если нет ответа (таймаут). При ERROR статусе — главный источник для диагностики',
  `http_status_code` integer COMMENT 'HTTP статус ответа iiko. 200 = успех → SENT. 500/503 = retry до 3 раз через 30 сек. 401 = невалидный api_key → оповестить владельца. NULL = таймаут',
  `created_at` timestamptz DEFAULT (now()) COMMENT 'Время попытки синхронизации. TTL 14 дней: pg_cron DELETE WHERE created_at < NOW() - INTERVAL 14 days'
);

CREATE TABLE `pos_settings` (
  `venue_id` uuid PRIMARY KEY COMMENT '1:1 с точкой. Создаётся при подключении iiko в онбординге',
  `api_key_encrypted` text NOT NULL COMMENT 'API-ключ iikoTransport. ЗАШИФРОВАТЬ через pgcrypto/Vault. В БД только blob. Ключ шифрования в env. Никогда не логировать',
  `terminal_id` uuid NOT NULL COMMENT 'UUID терминала iiko, куда "падают" заказы. Если несколько терминалов — владелец выбирает нужный при настройке',
  `external_org_id` uuid NOT NULL COMMENT 'UUID организации в облаке iiko (≠ наш organizations.id). Передаётся в каждый запрос к iikoCloud как organizationId',
  `sync_interval_minutes` integer DEFAULT 1 COMMENT 'Интервал полной синхронизации меню в минутах. Планировщик запускает каждые N минут. Полная (картинки) — раз в сутки',
  `stoplist_poll_seconds` integer DEFAULT 60 COMMENT 'Интервал опроса стоп-листа в секундах. Отдельный частый воркер обновляет products.is_available всех товаров точки',
  `last_menu_sync_at` timestamptz COMMENT 'Timestamp последней успешной синхронизации меню. Обновлять при каждом успехе. Отображать в админке: "Обновлено 2 мин назад"',
  `last_stoplist_sync_at` timestamptz COMMENT 'Timestamp последнего опроса стоп-листа. Если > 5 мин назад — алерт в Sentry. NULL = синхронизации ещё не было',
  `updated_at` timestamptz COMMENT 'Дата последнего изменения настроек. Проставлять при каждом UPDATE'
);

CREATE TABLE `promocodes` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID промокода',
  `venue_id` uuid NOT NULL COMMENT 'FK на точку. Промокод действует только в этой точке. UNIQUE(venue_id, code)',
  `code` varchar(255) NOT NULL COMMENT 'Код промокода: WELCOME, SUMMER20. Нормализовывать: UPPER(TRIM(code)) при записи и при поиске',
  `discount_type` ENUM ('PERCENT', 'FIXED') NOT NULL COMMENT 'PERCENT или FIXED. Определяет логику вычисления discount_amount в заказе',
  `discount_value` numeric(12,2) NOT NULL COMMENT 'PERCENT: 10 = 10%. FIXED: 150 = 150 руб. Для FIXED: discount = MIN(value, total) — защита от отрицательной суммы',
  `min_order_amount` numeric(12,2) DEFAULT 0 COMMENT 'Минимальная сумма заказа для применения промокода. 0 = без ограничений. Проверять до применения скидки',
  `max_uses` integer COMMENT 'Лимит использований. NULL = безлимит. При применении: SELECT FOR UPDATE → проверить → INCREMENT. Защита от race condition',
  `used_count` integer NOT NULL DEFAULT 0 COMMENT 'Счётчик использований. INCREMENT в транзакции с SELECT FOR UPDATE. Никогда не уменьшать вручную',
  `expires_at` timestamptz COMMENT 'Срок действия. NULL = бессрочно. Проверка: WHERE (expires_at IS NULL OR expires_at > NOW()) AND is_active = true',
  `is_active` boolean DEFAULT true COMMENT 'false = промокод отключён вручную. При деактивации существующие заказы с промокодом не затрагиваются',
  `created_at` timestamptz DEFAULT (now()) COMMENT 'Дата создания промокода'
);

CREATE TABLE `order_reviews` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID отзыва',
  `order_id` uuid UNIQUE NOT NULL COMMENT '1:1 с заказом. UNIQUE гарантирует один отзыв на заказ. Предлагать оценку через пуш спустя 1 час',
  `user_id` uuid NOT NULL COMMENT 'FK на клиента',
  `venue_id` uuid NOT NULL COMMENT 'FK на точку. Средний рейтинг: SELECT AVG(rating) WHERE venue_id = ?',
  `rating` smallint NOT NULL COMMENT 'Оценка 1-5. Валидировать: CHECK (rating BETWEEN 1 AND 5). 1 звезда = очень плохо, 5 = отлично',
  `comment` text COMMENT 'Текст отзыва. NULL если клиент поставил только оценку без текста. Показывать ТОЛЬКО в админке владельца — не публиковать в PWA',
  `owner_reply` text COMMENT 'Ответ владельца. NULL если не ответил. Заполняется в разделе "Отзывы" в админке. В MVP клиенту не доставляется',
  `replied_at` timestamptz COMMENT 'Время ответа владельца. Проставлять: SET replied_at = now() вместе с owner_reply',
  `created_at` timestamptz DEFAULT (now()) COMMENT 'Время создания отзыва'
);

CREATE TABLE `push_subscriptions` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID подписки',
  `user_id` uuid NOT NULL COMMENT 'FK на пользователя',
  `venue_id` uuid NOT NULL COMMENT 'FK на точку. Маркетинговые пуши точки: WHERE venue_id = ? AND is_active = true',
  `endpoint` text NOT NULL COMMENT 'URL push-сервиса браузера (Google FCM, Apple и др.). Получаем из ServiceWorker.pushManager.subscribe(). Уникален для устройства',
  `p256dh` text NOT NULL COMMENT 'Публичный ключ шифрования устройства. Не секрет. Передаётся в web-push библиотеку для шифрования payload',
  `auth` text NOT NULL COMMENT 'Auth-секрет устройства. Не секрет. Передаётся в web-push: sendNotification({endpoint, keys:{p256dh,auth}}, payload)',
  `is_active` boolean DEFAULT true COMMENT 'false = пользователь отозвал разрешение. Ставить false если браузер вернул 410 Gone при отправке пуша',
  `created_at` timestamptz DEFAULT (now()) COMMENT 'Время создания подписки'
);

CREATE TABLE `telegram_bindings` (
  `id` uuid PRIMARY KEY DEFAULT (uuid_generate_v4()) COMMENT 'UUID привязки',
  `user_id` uuid NOT NULL COMMENT 'FK на пользователя. UNIQUE гарантирует один Telegram на аккаунт',
  `telegram_chat_id` bigint NOT NULL COMMENT 'Числовой chat_id из Telegram Bot API (bigint — может быть > 2^31). Получаем из message.chat.id при команде /start',
  `is_active` boolean DEFAULT true COMMENT 'false = пользователь написал /stop боту. При отправке пушей: WHERE is_active = true',
  `created_at` timestamptz DEFAULT (now()) COMMENT 'Дата привязки аккаунта к боту'
);

CREATE INDEX `idx_users_phone` ON `users` (`phone`);

CREATE UNIQUE INDEX `idx_staff_roles_user_venue` ON `staff_roles` (`user_id`, `venue_id`);

CREATE INDEX `idx_sessions_token` ON `user_sessions` (`token`);

CREATE INDEX `idx_sessions_user` ON `user_sessions` (`user_id`);

CREATE INDEX `idx_sessions_expires` ON `user_sessions` (`expires_at`);

CREATE INDEX `idx_categories_venue_visible` ON `categories` (`venue_id`, `is_visible`);

CREATE UNIQUE INDEX `idx_categories_external_unique` ON `categories` (`venue_id`, `external_id`);

CREATE INDEX `idx_products_pwa_main` ON `products` (`venue_id`, `category_id`, `is_available`, `sort_order`);

CREATE UNIQUE INDEX `idx_products_external_unique` ON `products` (`venue_id`, `external_id`);

CREATE INDEX `idx_recs_sort` ON `product_recommendations` (`product_id`, `sort_order`);

CREATE INDEX `idx_orders_venue_status_date` ON `orders` (`venue_id`, `status`, `created_at`);

CREATE INDEX `idx_orders_user_date` ON `orders` (`user_id`, `created_at`);

CREATE INDEX `idx_orders_payment` ON `orders` (`payment_id`);

CREATE INDEX `idx_orders_sync` ON `orders` (`pos_sync_status`);

CREATE INDEX `idx_order_items_order` ON `order_items` (`order_id`);

CREATE INDEX `idx_oim_order_item` ON `order_item_modifiers` (`order_item_id`);

CREATE INDEX `idx_payments_order` ON `payments` (`order_id`);

CREATE UNIQUE INDEX `idx_payments_external_unique` ON `payments` (`external_id`);

CREATE INDEX `idx_payments_status` ON `payments` (`status`);

CREATE INDEX `idx_pos_logs_order_date` ON `pos_sync_logs` (`order_id`, `created_at`);

CREATE INDEX `idx_pos_logs_venue_date` ON `pos_sync_logs` (`venue_id`, `created_at`);

CREATE UNIQUE INDEX `idx_promocodes_venue_code` ON `promocodes` (`venue_id`, `code`);

CREATE INDEX `idx_promocodes_active` ON `promocodes` (`venue_id`, `is_active`, `expires_at`);

CREATE INDEX `idx_reviews_venue_date` ON `order_reviews` (`venue_id`, `created_at`);

CREATE INDEX `idx_reviews_user` ON `order_reviews` (`user_id`);

CREATE INDEX `idx_push_user_venue` ON `push_subscriptions` (`user_id`, `venue_id`);

CREATE UNIQUE INDEX `idx_tg_user` ON `telegram_bindings` (`user_id`);

CREATE INDEX `idx_tg_chat` ON `telegram_bindings` (`telegram_chat_id`);

ALTER TABLE `organizations` COMMENT = 'Юридическое лицо. Одна организация может иметь несколько точек (venues). Не хранит операционных данных.';

ALTER TABLE `venues` COMMENT = 'Центральная сущность системы. Торговая точка. Один organization может иметь много venues. Вся iiko/PWA логика привязана к venue_id.';

ALTER TABLE `organization_billing` COMMENT = 'История биллинга организации. Смена плана = новая запись. Старые записи не удалять.';

ALTER TABLE `users` COMMENT = 'Только идентификация. Роли и привязка к точкам — в таблице staff_roles. Клиент = users без записи в staff_roles.';

ALTER TABLE `staff_roles` COMMENT = 'Персонал: один user может быть admin в одном заведении и manager в другом.';

ALTER TABLE `product_recommendations` COMMENT = 'Настраивается вручную в админке. Блок "С этим также берут" на странице товара.';

ALTER TABLE `order_item_modifiers` COMMENT = 'Нормализованная замена jsonb selected_modifiers. Позволяет строить аналитику по добавкам.';

ALTER TABLE `pos_sync_logs` COMMENT = 'TTL: DELETE WHERE created_at < NOW() - INTERVAL 14 days (pg_cron job).';

ALTER TABLE `order_reviews` COMMENT = 'Внутренние отзывы — видны только владельцу в админке (ТЗ блок 6).';

ALTER TABLE `telegram_bindings` COMMENT = 'Привязка аккаунта к Telegram-боту для уведомлений персоналу и клиентам.';

ALTER TABLE `organizations` ADD FOREIGN KEY (`plan_id`) REFERENCES `subscription_plans` (`id`);

ALTER TABLE `venues` ADD FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`);

ALTER TABLE `organization_billing` ADD FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`);

ALTER TABLE `organization_billing` ADD FOREIGN KEY (`plan_id`) REFERENCES `subscription_plans` (`id`);

ALTER TABLE `app_configs` ADD FOREIGN KEY (`venue_id`) REFERENCES `venues` (`id`);

ALTER TABLE `staff_roles` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `staff_roles` ADD FOREIGN KEY (`venue_id`) REFERENCES `venues` (`id`);

ALTER TABLE `user_sessions` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `categories` ADD FOREIGN KEY (`venue_id`) REFERENCES `venues` (`id`);

ALTER TABLE `categories` ADD FOREIGN KEY (`parent_id`) REFERENCES `categories` (`id`);

ALTER TABLE `category_translations` ADD FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`);

ALTER TABLE `category_translations` ADD FOREIGN KEY (`language_code`) REFERENCES `languages` (`code`);

ALTER TABLE `products` ADD FOREIGN KEY (`venue_id`) REFERENCES `venues` (`id`);

ALTER TABLE `products` ADD FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`);

ALTER TABLE `product_translations` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);

ALTER TABLE `product_translations` ADD FOREIGN KEY (`language_code`) REFERENCES `languages` (`code`);

ALTER TABLE `product_modifiers` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);

ALTER TABLE `modifier_translations` ADD FOREIGN KEY (`modifier_id`) REFERENCES `product_modifiers` (`id`);

ALTER TABLE `modifier_translations` ADD FOREIGN KEY (`language_code`) REFERENCES `languages` (`code`);

ALTER TABLE `product_recommendations` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);

ALTER TABLE `product_recommendations` ADD FOREIGN KEY (`recommended_id`) REFERENCES `products` (`id`);

ALTER TABLE `orders` ADD FOREIGN KEY (`venue_id`) REFERENCES `venues` (`id`);

ALTER TABLE `orders` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `orders` ADD FOREIGN KEY (`payment_id`) REFERENCES `payments` (`id`);

ALTER TABLE `orders` ADD FOREIGN KEY (`promocode_id`) REFERENCES `promocodes` (`id`);

ALTER TABLE `order_items` ADD FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`);

ALTER TABLE `order_items` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);

ALTER TABLE `order_item_modifiers` ADD FOREIGN KEY (`order_item_id`) REFERENCES `order_items` (`id`);

ALTER TABLE `order_item_modifiers` ADD FOREIGN KEY (`modifier_id`) REFERENCES `product_modifiers` (`id`);

ALTER TABLE `payments` ADD FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`);

ALTER TABLE `pos_sync_logs` ADD FOREIGN KEY (`venue_id`) REFERENCES `venues` (`id`);

ALTER TABLE `pos_sync_logs` ADD FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`);

ALTER TABLE `pos_settings` ADD FOREIGN KEY (`venue_id`) REFERENCES `venues` (`id`);

ALTER TABLE `promocodes` ADD FOREIGN KEY (`venue_id`) REFERENCES `venues` (`id`);

ALTER TABLE `order_reviews` ADD FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`);

ALTER TABLE `order_reviews` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `order_reviews` ADD FOREIGN KEY (`venue_id`) REFERENCES `venues` (`id`);

ALTER TABLE `push_subscriptions` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `push_subscriptions` ADD FOREIGN KEY (`venue_id`) REFERENCES `venues` (`id`);

ALTER TABLE `telegram_bindings` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
