package models

// GenericResponse represents Common Response structure
type GenericResponse struct {
	Error   bool                   `json:"error"`
	Message string                 `json:"message"`
	Result  map[string]interface{} `json:"result"`
	Details string                 `json:"details"`
}

// Product model
type Product struct {
	ID          string `json:"id" bson:"id" mapstructure:"id"`
	Name        string `json:"name" bson:"name" mapstructure:"name"`
	Description string `json:"description" bson:"description" mapstructure:"description"`
	AvailQty    int    `json:"avail_qty" bson:"avail_qty" mapstructure:"avail_qty" `
	ReserveQty  int    `json:"reserve_qty" bson:"reserve_qty" mapstructure:"reserve_qty"`
}

// ProductModel model
type ProductModel struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

// CartModel model
type CartModel struct {
	ID       string         `json:"id"`
	Products []ProductModel `json:"products"`
}

// CheckoutModel model
type CheckoutModel struct {
	CartID string  `json:"cart_id"`
	Amount float64 `json:"amount"`
}

type OrderModel struct {
	OrderID      string `json:"order_id,omitempty"`
	OrderAddress string `json:"order_address,omitempty"`
}
