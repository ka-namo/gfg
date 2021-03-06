package product

// StockChangedNotifier is a notifier when product stock is changed.
type StockChangedNotifier interface {
	// StockChanged notifies or gives warning through different media to the Seller,
	// when Product stock is changed.
	StockChanged(sellerUUID, sellerReceiverID string, oldStock, newStock int, product string)
}
