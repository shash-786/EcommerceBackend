package database

import "errors"

var (
	ErrCantFindProduct        = errors.New("Can't Find The Product")
	ErrCantDecodeProduct      = errors.New("Can't Find the Product")
	ErrUserIdNotValid         = errors.New("User Invalid")
	ErrCantUpdateUser         = errors.New("Can't Add Product to Cart")
	ErrCantRemoveItemFromCart = errors.New("Can't Remove Item From Cart")
	ErrCantGetItem            = errors.New("Unable to get Item From Cart")
	ErrCantBuyCartItem        = errors.New("Cannot Update Purchase")
)

func AddProductToCart() error {
}

func RemoveProductFromCart() error {
}

func BuyItemFromCart() error {
}

func InstantBuy() error {
}
