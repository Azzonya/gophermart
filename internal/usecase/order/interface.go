package order

type OrderServiceI interface {
	IsLuhnValid(orderNumber string) bool
}
