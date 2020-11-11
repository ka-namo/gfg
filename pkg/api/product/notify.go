package product

type StockChangedNotifier interface {
	StockChanged(sellerUUID, sellerReceiverID string, oldStock, newStock int, product string)
}
