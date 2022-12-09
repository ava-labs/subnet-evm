package limitorders

type LimitOrder struct {
	id                int64
	PositionType      string
	UserAddress       string
	BaseAssetQuantity int
	Price             float64
	Status            string
	Salt              string
	Signature         []byte
	RawOrder          interface{}
	RawSignature      interface{}
}

type InMemoryDatabase interface {
	GetAllOrders() []*LimitOrder
	Add(order LimitOrder)
	Delete(signature []byte)
}

type inMemoryDatabase struct {
	orderMap map[string]*LimitOrder
}

func NewInMemoryDatabase() *inMemoryDatabase {
	orderMap := map[string]*LimitOrder{}
	return &inMemoryDatabase{orderMap}
}

func (db *inMemoryDatabase) GetAllOrders() []*LimitOrder {
	allOrders := []*LimitOrder{}
	for _, order := range db.orderMap {
		allOrders = append(allOrders, order)
	}
	return allOrders
}

func (db *inMemoryDatabase) Add(order LimitOrder) {
	db.orderMap[string(order.Signature)] = &order
}

// Deletes silently
func (db *inMemoryDatabase) Delete(signature []byte) {
	delete(db.orderMap, string(signature))
}
