package models

type Order struct {
	Model
	TransactionId   string      `json:"transaction_id" gorm:"null"`
	UserId          uint        `json:"user_id"`
	Code            string      `json:"code"`
	AmbassadorEmail string      `json:"ambassador_email"`
	FirstName       string      `json:"-"`
	LastName        string      `json:"-"`
	Name            string      `json:"name"`
	Email           string      `json:"email"`
	Address         string      `json:"address" gorm:"null"`
	City            string      `json:"city" gorm:"null"`
	Country         string      `json:"country" gorm:"null"`
	Zip             string      `json:"zip" gorm:"null"`
	Complete        bool        `json:"complete" gorm:"default:false"`
	Total           float64     `json:"total"`
	OrderItems      []OrderItem `json:"order_items" gorm:"foreignKey:OrderId"`
}

func (order *Order) GetFullname() string {
	return order.FirstName + " " + order.LastName
}

func (order *Order) GetTotal() float64 {
	var total float64
	for _, item := range order.OrderItems {
		total += item.Price * float64(item.Quantity)
	}
	return total
}
