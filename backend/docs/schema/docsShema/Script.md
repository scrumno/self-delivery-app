// ╔══════════════════════════════════════════════════════════════════╗
// ║   SAAS ДЛЯ РЕСТОРАНОВ — СХЕМА v2.0 (ПОСЛЕ АУДИТА)             ║
// ║   Исправлено: 6 критических + 9 важных + 10 рекомендаций       ║
// ║   Docs: https://dbdiagram.io/docs                               ║
// ╚══════════════════════════════════════════════════════════════════╝
// КРИТИЧЕСКОЕ ИСПРАВЛЕНИЕ: Инструкция для SQL-миграции
// Note: Перед запуском SQL-скрипта выполнить: CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
// Note: Весь скрипт должен быть обернут в BEGIN; ... COMMIT;
// ─────────────────────────────────────────────────────────────────────
// ENUMS
// ─────────────────────────────────────────────────────────────────────

Enum user_role {
  admin    [note: 'Владелец, полный доступ: биллинг, настройки, аналитика, управление персоналом']
  manager  [note: 'Менеджер точки: управляет стоп-листами, заказами и видимостью меню. Нет доступа к биллингу']
  cashier  [note: 'Кассир: Только список входящих заказов и кнопка "Выдан". Минимальные права']
  customer [note: 'Клиент: заказ еды — хранится в users, но роль определяется ОТСУТСТВИЕМ записи в staff_roles.']
}

Enum order_status {
  NEW             [note: 'Создан в PWA, Заказ создан, деньги не списаны. Можно отменить без возврата']
  WAITING_PAYMENT [note: 'Ожидает оплаты эквайрингом. Корзина зафиксирована, платёж инициирован. Ждём webhook от банка']
  COOKING         [note: 'Успешно ушёл в iiko на кухню. Платёж подтверждён, заказ передан в iiko. Отмена только через возврат']
  READY           [note: 'Готов к выдаче. iiko прислал статус готовности. PWA показывает клиенту крупный номер заказа']
  COMPLETED       [note: 'Завершён/ Клиент забрал заказ или прошло 2 часа с момента READY']
  CANCELLED       [note: 'Отменён (с возвратом или без). Причина в orders.cancelled_reason, возврат в orders.refund_id']
}

Enum pos_sync_status {
  NOT_SENT [note: 'Ожидает отправки в iiko. Воркер берёт эти записи каждые 5 сек и пытается отправить в iiko']
  SENT     [note: 'Успешно принято iiko. iiko вернул 200 OK и pos_order_id. Заказ принят кухней']
  ERROR    [note: 'Ошибка API (требует внимания админа). iiko вернул ошибку. Детали в pos_sync_logs.response_payload. Требует внимания']
}

Enum billing_status {
  ACTIVE   [note: 'Сервис оплачен. paid_until в будущем. Заказы принимаются в штатном режиме']
  PAST_DUE [note: 'Задолженность < 3 дней (предупреждение). paid_until просрочен менее 3 дней. Баннер-предупреждение в админке']
  BLOCKED  [note: 'Сервис отключён за неуплату. paid_until просрочен более 3 дней. PWA не принимает заказы, только страница оплаты']
}

// ИСПРАВЛЕНО #5: Enum вместо varchar для статусов платежей
Enum payment_status {
  PENDING   [note: 'Платёж инициирован. Платёж создан в БД и у банка. Ждём webhook. Деньги не списаны']
  SUCCESS   [note: 'Деньги списаны. Банк подтвердил списание. Заполняем paid_at, меняем order.status → COOKING']
  FAILED    [note: 'Ошибка оплаты. Банк отказал. Заказ возвращается в статус NEW для повторной оплаты']
  REFUNDED  [note: 'Возврат выполнен. Возврат выполнен через API банка. Заполняем orders.refunded_at']
  CANCELLED [note: 'Платёж отменён до списания (клиент закрыл страницу или истёк таймаут 15 мин)']
}

// ИСПРАВЛЕНО #21: Enum для тем оформления
Enum theme_preset {
  light    [note: 'Светлая (по умолчанию). Белый фон, тёмный текст. Универсально. Дефолт для всех новых точек']
  dark     [note: 'Тёмная. Тёмный фон. Подходит для баров и кофеен с атмосферой']
  coffee   [note: 'Кофейная. Тёплые коричневые тона. Для кофеен и пекарен']
  fastfood [note: 'Фастфуд. Яркие акценты, крупные кнопки. Для точек быстрого обслуживания']
}

Enum discount_type {
  PERCENT [note: 'Скидка в процентах. Значение = %, итог: total * (1 - value/100). Валидировать: 1-100']
  FIXED   [note: 'Скидка фиксированной суммой. Значение = рубли, итог: MAX(total - value, 0). Защита от отрицательной суммы']
}

Enum order_type {
  takeaway  [note: 'Самовывоз (по умолчанию в MVP). Единственный тип в MVP. Клиент называет номер заказа']
  dine_in   [note: 'В зале (будущее). Зарезервировано. Потребует поле table_number в orders']
  delivery  [note: 'Доставка (будущее). Зарезервировано. Потребует адрес и интеграцию с курьерами']
}


// ─────────────────────────────────────────────────────────────────────
// БЛОК 1: ОРГАНИЗАЦИЯ И БИЛЛИНГ
// ─────────────────────────────────────────────────────────────────────

// Юридическое лицо / владелец аккаунта
Table organizations {
  id         uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID юрлица. Передаётся в JWT сотрудника. Все запросы персонала фильтруются через него']
  name       varchar [not null, note: 'Название для договора. Официальное название для договоров и инвойсов. Не путать с venues.name — публичным названием']
  inn        varchar(12) [note: 'ИНН для проверки контрагента (РФ). ИНН: 10 цифр для ООО, 12 для ИП. Валидировать контрольную сумму. NULL разрешён в MVP']
  is_active  boolean [default: true, note: 'Глобальное выключение аккаунта. false = ручная блокировка саппортом. Все venues этой организации перестают работать']
  created_at timestamptz [default: `now()`, note: 'Дата регистрации. Только для чтения. Используется в аналитике роста клиентской базы']

  // ИСПРАВЛЕНО #2 и #15: Теперь plan_id явно привязан
  plan_id    uuid [not null, ref: > subscription_plans.id, note: 'FK → subscription_plans.id']

  Note: 'Юридическое лицо. Одна организация может иметь несколько точек (venues). Не хранит операционных данных.'
}

// ИСПРАВЛЕНО #1: Новая таблица — торговые точки (venues)
// Ключевое архитектурное решение: вся point-specific логика перенесена сюда
Table venues {
  id               uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID точки. Главный FK всей системы — почти каждая таблица ссылается на venue_id']
  organization_id  uuid    [not null, ref: > organizations.id, note: 'Владелец точки. При авторизации: staff_roles.venue_id → venues → organizations → проверяем billing_status']
  name             varchar [not null, note: 'Публичное название точки для клиентов и шапки PWA. Напр. "Кофе-пойнт ТЦ Мега"']

  // ИСПРАВЛЕНО #16: CHECK-ограничение на slug задокументировано в note
  slug             varchar [unique, not null, note: 'URL-часть: brand-kazan.foodapp.ru. Только [a-z0-9-], валидировать регуляркой. Менять осторожно — старые QR-коды перестанут работать']

  address          text [note: 'Физический адрес. Используется для геокодинга. Для отображения в PWA может быть перекрыт app_configs.address_manual']
  latitude         numeric(10,7) [note: 'Координаты для карты. GPS-широта. Заполнять через геокодер при сохранении адреса. Нужна для карты и функции "рядом со мной"']
  longitude        numeric(10,7) [note: 'GPS-долгота. Заполнять через геокодер при сохранении адреса']

  // Операционные настройки точки
  is_active           boolean [default: true, note: 'false = точка закрыта навсегда или в ремонте. PWA возвращает 404. Не путать с is_emergency_stop']
  is_emergency_stop   boolean [default: false, note: 'Кнопка "Запара": мгновенно отключает приём заказов. 1 клик в админке. При true: кнопка оплаты в PWA неактивна, API заказов возвращает 503']
  min_order_amount    numeric(12,2) [default: 0, note: 'Минимальная сумма заказа из ТЗ блок 4. 0 = нет ограничений. Проверяется на бэкенде при нажатии "Оплатить"']
  avg_cooking_minutes integer [default: 20, note: 'Среднее время готовности (мин) — показывается клиенту. "Время ожидания: ~20 мин". Не влияет на логику']

  // Рабочие часы (переопределение iiko)
  work_hours_json jsonb [note: 'Расписание-оверрайд. Формат: {"mon":{"open":"09:00","close":"22:00"}, ...}. Null = берём из iiko.']

  // Биллинг точки
  billing_status  billing_status [default: 'ACTIVE', note: 'Синхронизировать с organization_billing при каждом изменении. Fron-cron проверяет paid_until каждый час']
  paid_until      timestamptz [note: 'До какой даты оплачен сервис. При успешной оплате: paid_until += 1 month. NULL = триальный период']

  created_at      timestamptz [default: `now()`, note: 'Дата создания точки. Только для чтения']
  updated_at      timestamptz [note: 'Проставлять вручную: SET updated_at = now() при каждом UPDATE. Используется для инвалидации кэша PWA']
  deleted_at      timestamptz [note: 'Soft delete. NULL = активна. Никогда не делать DELETE FROM venues. Все SELECT: WHERE deleted_at IS NULL']

  Note: 'Центральная сущность системы. Торговая точка. Один organization может иметь много venues. Вся iiko/PWA логика привязана к venue_id.'
}

// ИСПРАВЛЕНО #2: Биллинг вынесен и расширен для хранения истории
Table organization_billing {
  id               uuid [primary key, default: `uuid_generate_v4()`, note: 'UUID записи биллинга. Текущий план = MAX(created_at) для данной организации']
  organization_id  uuid [not null, ref: > organizations.id , note: 'При смене плана: НЕ обновлять старую запись, создавать новую — для хранения истории']
  plan_id          uuid [not null, ref: > subscription_plans.id, note: 'JOIN на subscription_plans для получения цены и лимитов при генерации инвойса']
  billing_status   billing_status [default: 'ACTIVE', note: 'Дублирует venues.billing_status. Синхронизировать оба поля при любом изменении статуса']
  paid_until       timestamptz [note: 'При оплате: paid_until = MAX(paid_until, now()) + 1 month. Не затирать будущую дату при досрочной оплате']
  payment_method_token text [note: 'Токен привязанной карты для рекуррентных платежей']
  last_invoice_at  timestamptz [note: 'Когда последний раз генерировался инвойс. Обновлять при каждой генерации. Для отладки авто-биллинга']
  created_at       timestamptz [default: `now()`, note: 'Дата начала действия этой записи биллинга']

  Note: 'История биллинга организации. Смена плана = новая запись. Старые записи не удалять.'
}

Table subscription_plans {
  id              uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID плана. Зашит в organizations.plan_id']
  name            varchar [not null, note: 'Название для отображения: Starter, Pro, Enterprise']
  price_per_month numeric(12,2) [not null, note: 'Цена в рублях. При изменении цены — создавать новый план, не обновлять. Старые организации остаются на старой цене']
  features_json   jsonb   [note: 'Лимиты: макс. блюд, наличие iiko, кол-во точек и т.д. Лимиты и флаги. Формат: {"max_products":100,"iiko_integration":true,"max_venues":1,"push_marketing":false}']
  is_active       boolean [default: true, note: 'Архивирование устаревших планов. false = план архивирован, не показывать новым клиентам. Существующие продолжают им пользоваться']
  created_at      timestamptz [default: `now()`, note: 'Дата создания плана']
}

// Справочник языков
Table languages {
  code varchar(5) [primary key, note: 'ISO 639-1 код языка: ru, en, ge. Заполняется при деплое seed-скриптом. Не изменять в runtime']
  name varchar [not null, note: 'Читаемое название для выпадающего списка в админке: Русский, English, ქართული']
}


// ─────────────────────────────────────────────────────────────────────
// БЛОК 2: КОНФИГУРАЦИЯ PWA
// ─────────────────────────────────────────────────────────────────────

// ИСПРАВЛЕНО #1: Привязка к venue, а не к organization
Table app_configs {
  venue_id     uuid [primary key, ref: - venues.id, note: 'Дизайн настраивается для каждой точки. 1:1 с venues. Создаётся автоматически при создании точки с дефолтными значениями']

  // ИСПРАВЛЕНО #21: Enum вместо varchar
  theme_preset theme_preset [default: 'light', note: 'Фронт читает при загрузке PWA и применяет CSS-переменные темы. Менять только через конструктор в админке']
  accent_color varchar(7)   [default: '#000000', note: 'HEX цвет кнопок и акцентов. Валидировать: /^#[0-9A-Fa-f]{6}$/. Напр. #E74C3C']
  logo_url     text         [note: 'URL логотипа в S3/CDN (PNG или SVG). NULL = заглушка с первой буквой названия точки. Не хранить base64 в БД']
  banner_url   text         [note: 'URL главного баннера на главной PWA. NULL = без баннера. Рекомендуемый размер: 1200×400px']
  address_manual text       [note: 'Ручная правка адреса (переопределяет venues.address). Перекрывает venues.address только в UI PWA. venues.address по-прежнему используется для геокодинга']

  updated_at   timestamptz [note: 'Проставлять при каждом изменении. Используется для инвалидации SW-кэша браузера']
}


// ─────────────────────────────────────────────────────────────────────
// БЛОК 3: ПОЛЬЗОВАТЕЛИ И АВТОРИЗАЦИЯ
// ─────────────────────────────────────────────────────────────────────

// ИСПРАВЛЕНО #4: users — только идентификация клиента, БЕЗ organization_id
Table users {
  id            uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID пользователя. Хранится в JWT. Из сессии: sessions.token → sessions.user_id']
  phone         varchar [unique, not null, note: 'E.164 формат: +79161234567. Нормализовывать при сохранении (убирать пробелы, скобки). Индекс idx_users_phone обязателен']
  full_name     varchar [note: 'Имя клиента. NULL при регистрации — запрашиваем после первого входа. Используется в истории заказов и базе клиентов']
  birth_date    date    [note: 'Дата рождения. NULL — заполняется опционально в профиле. Для именинных акций. Хранить без времени (тип DATE)']
  iiko_guest_id uuid    [note: 'ID гостя в iikoCard. NULL до первого заказа. При первом заказе: создать гостя в iiko → сохранить ID → создать заказ']
  is_active     boolean [default: true, note: 'false = мягкая блокировка или удаление по запросу (152-ФЗ). При false: 401 на все запросы. Данные не удалять']
  created_at    timestamptz [default: `now()`, note: 'Дата регистрации. Только для чтения']
  indexes {
    phone [name: 'idx_users_phone', note: 'Критично для скорости авторизации']
  }
  Note: 'Только идентификация. Роли и привязка к точкам — в таблице staff_roles. Клиент = users без записи в staff_roles.'
}

// ИСПРАВЛЕНО #4: Новая таблица для персонала (admin/manager/cashier)
Table staff_roles {
  id         uuid [primary key, default: `uuid_generate_v4()`, note: 'UUID записи роли']
  user_id    uuid [not null, ref: > users.id, note: 'FK на пользователя. Один человек может иметь несколько записей для разных точек']
  venue_id   uuid [not null, ref: > venues.id, note: 'FK на точку. При авторизации: читаем все staff_roles пользователя, строим список доступных точек']
  role       user_role [not null, note: 'admin: полный доступ к точке включая биллинг. manager: стоп-листы и заказы. cashier: только список заказов']
  is_active  boolean [default: true, note: 'false = сотрудник уволен. НЕ удалять запись. При авторизации: WHERE is_active = true']
  created_at timestamptz [default: `now()`, note: 'Дата назначения на роль']

  indexes {
    (user_id, venue_id) [unique, name: 'idx_staff_roles_user_venue', note: 'Один пользователь — одна роль на точку']
  }

  Note: 'Персонал: один user может быть admin в одном заведении и manager в другом.'
}

Table user_sessions {
  id          uuid [primary key, default: `uuid_generate_v4()`, note: 'UUID сессии']
  user_id     uuid [not null, ref: > users.id, note: 'FK на пользователя. При выходе: DELETE WHERE user_id = ? AND id = ?']
  token       text [unique, not null, note: 'Случайный непрозрачный токен (не JWT). Можно мгновенно отозвать. Индекс idx_sessions_token критически важен']
  device_info jsonb [note: 'PWA, UA и платформа устройства. Формат: {"ua":"Mozilla...","platform":"Android","pwa":true}. Для аналитики и списка устройств']
  expires_at  timestamptz [not null, note: 'Клиенты: now()+30d, персонал: now()+8h. Fоновый job: DELETE WHERE expires_at < now(). Продлевать при каждом запросе']
  created_at  timestamptz [default: `now()`, note: 'Дата создания сессии (момент входа)']

  indexes {
    (token) [name: 'idx_sessions_token']
    (user_id) [name: 'idx_sessions_user']
    (expires_at) [name: 'idx_sessions_expires', note: 'Для фоновой очистки протухших сессий']
  }
}


// ─────────────────────────────────────────────────────────────────────
// БЛОК 4: МЕНЮ И КАТАЛОГ
// ─────────────────────────────────────────────────────────────────────

// ИСПРАВЛЕНО #1: Привязка к venue
Table categories {
  id              uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID категории']
  venue_id        uuid    [not null, ref: > venues.id, note: 'Привязка к точке. При запросе меню: WHERE venue_id = ? AND is_visible = true AND deleted_at IS NULL']
  external_id     uuid    [note: 'UUID категории в iiko Cloud. Индекс: idx_categories_external. Матчинг при синхронизации. NULL = категория создана вручную в нашей админке']
  parent_id       uuid    [ref: > categories.id, note: 'Для вложенных категорий. FK на родительскую категорию. NULL = верхний уровень. Для дерева: рекурсивный WITH RECURSIVE запрос']
  sort_order      integer [default: 0, note: 'Порядок среди соседних категорий. ORDER BY sort_order ASC. При drag&drop: batch UPDATE']
  is_visible      boolean [default: true, note: 'Чекбокс в админке: Скрыть/Показать. false = скрыта от клиентов (напр. "Хозтовары"). Только для чтения — видна в админке']
  deleted_at      timestamptz [note: 'Soft delete. При удалении из iiko: SET deleted_at = now(). Все SELECT: WHERE deleted_at IS NULL']

  indexes {
    (venue_id, is_visible) [name: 'idx_categories_venue_visible']
    (venue_id, external_id) [unique, name: 'idx_categories_external_unique', note: 'Защита от дублей iiko внутри одной точки']
  }
}

Table category_translations {
  category_id   uuid    [ref: > categories.id, note: 'FK на категорию']
  language_code varchar [ref: > languages.code, note: 'Код языка. При запросе меню: JOIN WHERE language_code = язык из Accept-Language']
  title         varchar [not null, note: 'Название категории на данном языке. Напр. "Горячие напитки", "Hot Beverages". ru заполняется из iiko, остальные вручную']

  indexes {
    (category_id, language_code) [pk]
  }
}

Table products {
  id              uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID товара']
  venue_id        uuid    [not null, ref: > venues.id, note: 'FK на точку. Индекс idx_products_venue_cat_avail используется при каждом запросе меню']

  // ИСПРАВЛЕНО #17: ON DELETE RESTRICT — нельзя удалить категорию с товарами
  category_id     uuid    [ref: > categories.id, note: 'FK на категорию. ON DELETE RESTRICT — нельзя удалить категорию пока в ней есть товары']
  external_id     uuid    [not null, note: 'UUID из iiko. МАТЧИНГ СТРОГО ПО UUID при синхронизации. Нет в нашей БД → INSERT, есть → UPDATE, нет в iiko → soft delete. Индекс: idx_products_external']

  price           numeric(12,2) [not null, note: 'Текущая цена в рублях из iiko. Обновляется при каждой синхронизации. На цену в заказе не влияет (снапшот в order_items)']
  image_url       text [note: 'URL фото из iiko Cloud или нашего S3. NULL = заглушка с логотипом точки. Не хранить base64']
  is_available    boolean [default: true, note: 'Стоп-лист из iiko или ручной. false = стоп-лист. Обновляется каждые 60 сек из iiko. Проверяется дважды: при показе меню и при нажатии "Оплатить"']

  // ИСПРАВЛЕНО #9: Порядок сортировки блюд внутри категории
  sort_order      integer [default: 0, note: 'Порядок отображения внутри категории. Берётся из iiko при синхронизации. ORDER BY sort_order ASC']

  // Физические параметры
  weight          integer [note: 'Вес порции в граммах из iiko. NULL если не заполнено в iiko. Отображается в карточке: "250 г"']

  // ИСПРАВЛЕНО #13: КБЖУ из ТЗ блок 3
  calories        numeric(8,2) [note: 'ккал на 100г из iiko. NULL если не заполнено. Блок КБЖУ показываем только если calories IS NOT NULL']
  protein         numeric(8,2) [note: 'Белки на 100г из iiko. Обновляется при ежедневной полной синхронизации']
  fat             numeric(8,2) [note: 'Жиры на 100г из iiko']
  carbs           numeric(8,2) [note: 'Углеводы на 100г из iiko']

  deleted_at      timestamptz [note: 'Soft delete. При исчезновении из iiko: SET deleted_at = now(), is_available = false. Все SELECT: WHERE deleted_at IS NULL']

  indexes {
    // Супер-индекс: покрывает и поиск в категории, и статус стоп-листа, и порядок отображения
    (venue_id, category_id, is_available, sort_order) [name: 'idx_products_pwa_main']
    // Для синхронизации с iiko
    (venue_id, external_id) [unique, name: 'idx_products_external_unique']
  }
}

Table product_translations {
  product_id    uuid    [ref: > products.id, note: 'FK на товар']
  language_code varchar [ref: > languages.code, note: 'Код языка перевода']
  name          varchar [not null, note: 'Название товара. Используется в поиске. GIN-индекс на to_tsvector создаётся отдельной миграцией']
  description   text [note: 'Описание товара. NULL разрешён. Отображается в карточке товара при нажатии на него']

  indexes {
    (product_id, language_code) [pk]
    // ИСПРАВЛЕНО: Full-text search для поиска по названию (ТЗ блок 3)
    // Создаётся отдельно в миграции:
    // CREATE INDEX idx_product_fts ON product_translations USING GIN (to_tsvector('russian', name));
  }
}

Table product_modifiers {
  id          uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID модификатора']
  product_id  uuid    [not null, ref: > products.id, note: 'FK на товар. При загрузке карточки: SELECT * FROM product_modifiers WHERE product_id = ?']
  external_id uuid    [note: 'UUID модификатора в iiko. NULL если создан вручную. Матчинг при синхронизации аналогично products.external_id']
  extra_price numeric(12,2) [default: 0, note: '0 = бесплатный модификатор (напр. "Без сахара"). Прибавляется к цене. Снапшот хранится в order_item_modifiers']

  // ИСПРАВЛЕНО #19: Различаем обязательные и групповые модификаторы
  is_required  boolean [default: false, note: 'true = клиент не может добавить товар в корзину без выбора. Напр. группа "Размер": S/M/L обязательна']
  max_quantity integer [default: 1, note: '1 = чекбокс. >1 = можно выбрать несколько (напр. количество сиропов до 3). Валидировать на бэкенде']
  group_name   varchar [note: 'Название группы для UI: "Размер", "Добавки", "Основа". NULL = одиночный модификатор. Одинаковый group_name = одна визуальная группа']
}

Table modifier_translations {
  modifier_id   uuid    [ref: > product_modifiers.id, note: 'FK на модификатор']
  language_code varchar [ref: > languages.code, note: 'Код языка перевода']
  name          varchar [not null, note: 'Название добавки на данном языке. Напр. "Большой", "Large", "Без сахара"']

  indexes {
    (modifier_id, language_code) [pk]
  }
}

// ИСПРАВЛЕНО #12 (рекомендация): Upsell — "С этим также берут" (ТЗ блок 3)
Table product_recommendations {
  product_id         uuid    [ref: > products.id, note: 'Товар, на странице которого показывается блок. SELECT recommended_id WHERE product_id = ? ORDER BY sort_order']
  recommended_id     uuid    [ref: > products.id, note: 'Рекомендуемый товар. При запросе: проверять is_available = true — не показывать стоп-лист в upsell']
  sort_order         integer [default: 0, note: 'Порядок в блоке "С этим также берут". Обновляется drag&drop в админке. ORDER BY sort_order ASC']

  indexes {
    (product_id, recommended_id) [pk]
    (product_id, sort_order)     [name: 'idx_recs_sort']
  }

  Note: 'Настраивается вручную в админке. Блок "С этим также берут" на странице товара.'
}


// ─────────────────────────────────────────────────────────────────────
// БЛОК 5: ЗАКАЗЫ
// ─────────────────────────────────────────────────────────────────────

Table orders {
  id              uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID заказа. Показывается клиенту как номер: последние 4 символа UUID']
  venue_id        uuid    [not null, ref: > venues.id, note: 'FK на точку. Все заказы точки: WHERE venue_id = ? AND deleted_at IS NULL']
  user_id         uuid    [not null, ref: > users.id, note: 'FK на клиента. История клиента: WHERE user_id = ? ORDER BY created_at DESC LIMIT 20']

  status          order_status [default: 'NEW', note: 'Менять только через сервисные методы, не UPDATE напрямую — иначе пуши и логи не сработают']
  order_type      order_type   [default: 'takeaway', note: 'В MVP всегда takeaway. dine_in и delivery — для будущих версий']
  total_amount    numeric(12,2) [not null, note: 'Итог после скидки. Формула: SUM(unit_price*qty) + SUM(modifier.extra_price) - discount_amount. Фиксируется при создании']
  scheduled_at    timestamptz  [note: 'Предзаказ на время.NULL = "как можно скорее". Если задано: шаг 15 мин, только текущий или следующий день. Передаётся в iiko как время доставки']

  // ИСПРАВЛЕНО #14: Поля из ТЗ блок 4
  comment         text         [note: 'Комментарий клиента к заказу, "Без лука, пожалуйста". Передаётся в iiko в поле комментария. Показывается сотруднику в карточке заказа']
  cutlery_needed  boolean      [default: false, note: 'Чекбокс "Нужны приборы" в корзине. Передаётся в iiko как тег или комментарий']

  // ИСПРАВЛЕНО #20: Поля для логики отмены и возврата
  cancelled_reason text        [note: 'NULL если заказ не отменён. Заполняется автоматически ("Таймаут оплаты") или вручную сотрудником. Для аналитики потерь']
  refund_id        text        [note: 'ID возврата в платёжной системе. NULL пока возврата не было. При наличии: показываем статус "Возврат выполнен"']
  refunded_at      timestamptz [note: 'Момент выполнения возврата. NULL если возврата не было']

  // ИСПРАВЛЕНО #3: Прямая ссылка на платёж для ускорения JOIN
  payment_id      uuid         [ref: > payments.id, note: 'FK → payments.id (nullable, проставляется после создания платежа). NULL до инициации платежа. Порядок: создать заказ → создать payment → обновить orders.payment_id']

  // iiko-синхронизация
  pos_order_id    uuid         [note: 'UUID заказа от iiko. NULL до успешной отправки. Используется для отмены в iiko и запроса статуса готовности']
  pos_sync_status pos_sync_status [default: 'NOT_SENT', note: 'Воркер каждые 5 сек: WHERE pos_sync_status = NOT_SENT AND status = COOKING. Retry до 3 раз при ERROR']

  // ТЗ блок 9: поле orderSource для маркетинга
  order_source    varchar      [default: 'PWA', note: 'Передаётся в iiko как orderSource для маркетинговой аналитики. По ТЗ = название нашего приложения']

  // Промокод, применённый к заказу
  promocode_id    uuid         [ref: > promocodes.id, note: 'NULL если промокод не применялся. При применении: проверить валидность → INCREMENT used_count → сохранить ID']
  discount_amount numeric(12,2) [default: 0, note: 'Сумма скидки по промокоду. total_amount уже с учётом скидки. Поле хранится для прозрачности аналитики промокодов']

  created_at      timestamptz  [default: `now()`, note: 'Время создания заказа клиентом']
  deleted_at      timestamptz  [note: 'Soft delete. Заказы НИКОГДА не удалять физически — финансовая история. deleted_at только для технически ошибочных записей']

  indexes {
    // ИСПРАВЛЕНО #10: Главный составной индекс для фильтрации в админке
    (venue_id, status, created_at) [name: 'idx_orders_venue_status_date']
    (user_id, created_at)          [name: 'idx_orders_user_date', note: 'История заказов клиента']
    (payment_id)                   [name: 'idx_orders_payment']
    (pos_sync_status)              [name: 'idx_orders_sync', note: 'Для воркера повторных отправок в iiko']
  }
}

Table order_items {
  id          uuid    [primary key, default: `uuid_generate_v4()`]
  order_id    uuid    [not null, ref: > orders.id, note: 'FK на заказ. Все позиции заказа: WHERE order_id = ?']
  product_id  uuid    [not null, ref: > products.id, note: 'FK на товар. Использовать LEFT JOIN — товар может быть удалён из меню']
  quantity    integer [not null, note: 'Количество единиц. Минимум 1. Нельзя изменить после создания заказа']
  unit_price  numeric(12,2) [not null, note: 'Цена 1 единицы товара на момент оплаты (снапшот из products.price). Не зависит от будущих изменений цены']

  // Снапшот названия товара (на случай если товар удалят из меню)
  product_name_snapshot varchar [note: 'Название товара на языке клиента на момент заказа. ОБЯЗАТЕЛЬНО заполнять. Без снапшота история пустеет при удалении товара']

  indexes {
    (order_id) [name: 'idx_order_items_order']
  }
}

// ИСПРАВЛЕНО #6: Нормализованная таблица модификаторов заказа
// Вместо jsonb-снапшота — полноценные строки для аналитики
Table order_item_modifiers {
  id              uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID записи модификатора заказа']
  order_item_id   uuid    [not null, ref: > order_items.id, note: 'FK на позицию заказа. Все добавки позиции: WHERE order_item_id = ?']
  modifier_id     uuid    [ref: > product_modifiers.id, note: 'FK на модификатор. NULLABLE — модификатор может быть удалён из меню. Всегда LEFT JOIN, не INNER']
  name_snapshot   varchar [not null, note: 'Название добавки на момент заказа. ОБЯЗАТЕЛЬНО заполнять. Аналогично product_name_snapshot — защита от удаления']
  extra_price     numeric(12,2) [not null, default: 0, note: 'Цена добавки на момент заказа (снапшот). Итог позиции: (unit_price + SUM(extra_price)) * quantity']

  indexes {
    (order_item_id) [name: 'idx_oim_order_item']
  }

  Note: 'Нормализованная замена jsonb selected_modifiers. Позволяет строить аналитику по добавкам.'
}


// ─────────────────────────────────────────────────────────────────────
// БЛОК 6: ПЛАТЕЖИ
// ─────────────────────────────────────────────────────────────────────

Table payments {
  id          uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID платежа в нашей системе']
  order_id    uuid    [unique, not null, ref: - orders.id, note: '1:1 с заказом. UNIQUE гарантирует один платёж на заказ']
  external_id varchar [not null, note: 'ID транзакции в Тинькофф/ЮKassa. Используется для идентификации входящих webhook: WHERE external_id = ?']

  // ИСПРАВЛЕНО #5: Enum вместо varchar
  status      payment_status [default: 'PENDING', note: 'Обновлять ТОЛЬКО через webhook от банка или polling. Никогда не менять вручную без подтверждения банка']
  amount      numeric(12,2) [not null, note: 'Сумма платежа в рублях. Должна совпадать с orders.total_amount']

  // ИСПРАВЛЕНО #3: Обратная ссылка order.payment_id → payments.id проставляется после создания
  // Хранение деталей провайдера
  provider    varchar [note: 'Название провайдера: tinkoff, yookassa, sbp. Нужен при возврате — используем правильный API']
  paid_at     timestamptz [note: 'Момент фактического списания от банка (из данных webhook). NULL до SUCCESS. Не ставить своё время — только из данных банка']
  created_at  timestamptz [default: `now()`, note: 'Время создания платёжной записи (инициации платежа)']

  indexes {
    (order_id)    [name: 'idx_payments_order']
    external_id [unique, name: 'idx_payments_external_unique', note: 'Защита от повторной обработки одного и того же платежа от банка']
    (status)      [name: 'idx_payments_status']
  }
}


// ─────────────────────────────────────────────────────────────────────
// БЛОК 7: ЛОГИ СИНХРОНИЗАЦИИ С iiko
// ─────────────────────────────────────────────────────────────────────

// ИСПРАВЛЕНО #7: TTL = 14 дней (согласно ТЗ блок 9, было 30)
Table pos_sync_logs {
  id                uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID записи лога']
  venue_id          uuid    [not null, ref: > venues.id, note: 'FK на точку. Для дебага: WHERE venue_id = ? ORDER BY created_at DESC']
  order_id          uuid    [not null, ref: > orders.id, note: 'FK на заказ. Все попытки отправки заказа: WHERE order_id = ? ORDER BY created_at DESC']
  request_payload   jsonb   [note: 'Полное тело запроса отправленного в iiko. Логировать ДО отправки. НЕ включать api_key в payload']
  response_payload  jsonb   [note: 'Полный ответ iiko (JSON). NULL если нет ответа (таймаут). При ERROR статусе — главный источник для диагностики']
  http_status_code  integer [note: 'HTTP статус ответа iiko. 200 = успех → SENT. 500/503 = retry до 3 раз через 30 сек. 401 = невалидный api_key → оповестить владельца. NULL = таймаут']
  created_at        timestamptz [default: `now()`, note: 'Время попытки синхронизации. TTL 14 дней: pg_cron DELETE WHERE created_at < NOW() - INTERVAL 14 days']

  indexes {
    // ИСПРАВЛЕНО #18: Индексы для дебага
    (order_id, created_at)   [name: 'idx_pos_logs_order_date']
    (venue_id, created_at)   [name: 'idx_pos_logs_venue_date', note: 'Для фоновой очистки TTL 14 дней']
  }

  Note: 'TTL: DELETE WHERE created_at < NOW() - INTERVAL 14 days (pg_cron job).'
}


// ─────────────────────────────────────────────────────────────────────
// БЛОК 8: НАСТРОЙКИ iiko
// ─────────────────────────────────────────────────────────────────────

// ИСПРАВЛЕНО #1 и #25: Привязка к venue; api_key помечен как требующий шифрования
Table pos_settings {
  venue_id                uuid    [primary key, ref: - venues.id, note: '1:1 с точкой. Создаётся при подключении iiko в онбординге']
  api_key_encrypted       text    [not null, note: 'API-ключ iikoTransport. ЗАШИФРОВАТЬ через pgcrypto/Vault. В БД только blob. Ключ шифрования в env. Никогда не логировать']
  terminal_id             uuid    [not null, note: 'UUID терминала iiko, куда "падают" заказы. Если несколько терминалов — владелец выбирает нужный при настройке']
  external_org_id         uuid    [not null, note: 'UUID организации в облаке iiko (≠ наш organizations.id). Передаётся в каждый запрос к iikoCloud как organizationId']
  sync_interval_minutes   integer [default: 1,  note: 'Интервал полной синхронизации меню в минутах. Планировщик запускает каждые N минут. Полная (картинки) — раз в сутки']
  stoplist_poll_seconds   integer [default: 60, note: 'Интервал опроса стоп-листа в секундах. Отдельный частый воркер обновляет products.is_available всех товаров точки']

  // ИСПРАВЛЕНО: Timestamp последней синхронизации из ТЗ блок 9
  last_menu_sync_at       timestamptz [note: 'Timestamp последней успешной синхронизации меню. Обновлять при каждом успехе. Отображать в админке: "Обновлено 2 мин назад"']
  last_stoplist_sync_at   timestamptz [note: 'Timestamp последнего опроса стоп-листа. Если > 5 мин назад — алерт в Sentry. NULL = синхронизации ещё не было']

  updated_at  timestamptz [note: 'Дата последнего изменения настроек. Проставлять при каждом UPDATE']
}


// ─────────────────────────────────────────────────────────────────────
// БЛОК 9: ПРОМОКОДЫ
// ─────────────────────────────────────────────────────────────────────

// ИСПРАВЛЕНО #22: Новая таблица — промокоды из ТЗ блок 4
Table promocodes {
  id              uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID промокода']
  venue_id        uuid    [not null, ref: > venues.id, note: 'FK на точку. Промокод действует только в этой точке. UNIQUE(venue_id, code)']
  code            varchar [not null, note: 'Код промокода: WELCOME, SUMMER20. Нормализовывать: UPPER(TRIM(code)) при записи и при поиске']
  discount_type   discount_type [not null, note: 'PERCENT или FIXED. Определяет логику вычисления discount_amount в заказе']
  discount_value  numeric(12,2) [not null, note: 'PERCENT: 10 = 10%. FIXED: 150 = 150 руб. Для FIXED: discount = MIN(value, total) — защита от отрицательной суммы']
  min_order_amount numeric(12,2) [default: 0, note: 'Минимальная сумма заказа для применения промокода. 0 = без ограничений. Проверять до применения скидки']
  max_uses        integer  [note: 'Лимит использований. NULL = безлимит. При применении: SELECT FOR UPDATE → проверить → INCREMENT. Защита от race condition']
  used_count      integer  [default: 0, not null, note: 'Счётчик использований. INCREMENT в транзакции с SELECT FOR UPDATE. Никогда не уменьшать вручную']
  expires_at      timestamptz [note: 'Срок действия. NULL = бессрочно. Проверка: WHERE (expires_at IS NULL OR expires_at > NOW()) AND is_active = true']
  is_active       boolean  [default: true, note: 'false = промокод отключён вручную. При деактивации существующие заказы с промокодом не затрагиваются']
  created_at      timestamptz [default: `now()`, note: 'Дата создания промокода']

  indexes {
    (venue_id, code) [unique, name: 'idx_promocodes_venue_code']
    (venue_id, is_active, expires_at) [name: 'idx_promocodes_active']
  }
}


// ─────────────────────────────────────────────────────────────────────
// БЛОК 10: ОТЗЫВЫ И РЕЙТИНГИ
// ─────────────────────────────────────────────────────────────────────

// ИСПРАВЛЕНО #23: Таблица отзывов из ТЗ блок 6
Table order_reviews {
  id           uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID отзыва']
  order_id     uuid    [unique, not null, ref: - orders.id, note: '1:1 с заказом. UNIQUE гарантирует один отзыв на заказ. Предлагать оценку через пуш спустя 1 час']
  user_id      uuid    [not null, ref: > users.id, note: 'FK на клиента']
  venue_id     uuid    [not null, ref: > venues.id, note: 'FK на точку. Средний рейтинг: SELECT AVG(rating) WHERE venue_id = ?']
  rating       smallint [not null, note: 'Оценка 1-5. Валидировать: CHECK (rating BETWEEN 1 AND 5). 1 звезда = очень плохо, 5 = отлично']
  comment      text    [note: 'Текст отзыва. NULL если клиент поставил только оценку без текста. Показывать ТОЛЬКО в админке владельца — не публиковать в PWA']
  owner_reply  text    [note: 'Ответ владельца. NULL если не ответил. Заполняется в разделе "Отзывы" в админке. В MVP клиенту не доставляется']
  replied_at   timestamptz [note: 'Время ответа владельца. Проставлять: SET replied_at = now() вместе с owner_reply']
  created_at   timestamptz [default: `now()`, note: 'Время создания отзыва']

  indexes {
    (venue_id, created_at) [name: 'idx_reviews_venue_date']
    (user_id)              [name: 'idx_reviews_user']
  }

  Note: 'Внутренние отзывы — видны только владельцу в админке (ТЗ блок 6).'
}


// ─────────────────────────────────────────────────────────────────────
// БЛОК 11: PUSH-УВЕДОМЛЕНИЯ
// ─────────────────────────────────────────────────────────────────────

// ИСПРАВЛЕНО #24: Таблица браузерных push-подписок из ТЗ блок 6
Table push_subscriptions {
  id           uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID подписки']
  user_id      uuid    [not null, ref: > users.id, note: 'FK на пользователя']
  venue_id     uuid    [not null, ref: > venues.id, note: 'FK на точку. Маркетинговые пуши точки: WHERE venue_id = ? AND is_active = true']
  endpoint     text    [not null, note: 'URL push-сервиса браузера (Google FCM, Apple и др.). Получаем из ServiceWorker.pushManager.subscribe(). Уникален для устройства']
  p256dh       text    [not null, note: 'Публичный ключ шифрования устройства. Не секрет. Передаётся в web-push библиотеку для шифрования payload']
  auth         text    [not null, note: 'Auth-секрет устройства. Не секрет. Передаётся в web-push: sendNotification({endpoint, keys:{p256dh,auth}}, payload)']
  is_active    boolean [default: true, note: 'false = пользователь отозвал разрешение. Ставить false если браузер вернул 410 Gone при отправке пуша']
  created_at   timestamptz [default: `now()`, note: 'Время создания подписки']

  indexes {
    (user_id, venue_id) [name: 'idx_push_user_venue']
  }
}

// Telegram-уведомления для персонала (ТЗ блок 6)
Table telegram_bindings {
  id           uuid    [primary key, default: `uuid_generate_v4()`, note: 'UUID привязки']
  user_id      uuid    [not null, ref: > users.id, note: 'FK на пользователя. UNIQUE гарантирует один Telegram на аккаунт']
  telegram_chat_id bigint [not null, note: 'Числовой chat_id из Telegram Bot API (bigint — может быть > 2^31). Получаем из message.chat.id при команде /start']
  is_active    boolean [default: true, note: 'false = пользователь написал /stop боту. При отправке пушей: WHERE is_active = true']
  created_at   timestamptz [default: `now()`, note: 'Дата привязки аккаунта к боту']

  indexes {
    (user_id)          [unique, name: 'idx_tg_user']
    (telegram_chat_id) [name: 'idx_tg_chat']
  }

  Note: 'Привязка аккаунта к Telegram-боту для уведомлений персоналу и клиентам.'
}


