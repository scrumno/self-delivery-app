package payment_provider

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
