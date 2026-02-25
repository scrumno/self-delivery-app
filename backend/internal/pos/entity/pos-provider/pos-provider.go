package pos_provider

// POSProvider — поддерживаемые кассовые / ресторанные системы.
type POSProvider string

const (
	POSProviderIiko    POSProvider = "iiko"
	POSProviderSyrve   POSProvider = "syrve" // бывший iiko cloud, отдельный API
	POSProviderRkeeper POSProvider = "rkeeper"
	POSProviderManual  POSProvider = "manual" // без интеграции, ручное управление меню
)
