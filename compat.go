package sqltest

var (
	// DriversByProduct maps a Product to the Driver implementations that it is
	// compatible with.
	DriversByProduct = map[Product][]Driver{}

	// ProductsByDriver maps a Driver to the Product implementations that it is
	// compatible with.
	ProductsByDriver = map[Driver][]Product{}

	// CompatiblePairs contains all compatible driver/product pairs.
	CompatiblePairs []Pair
)

// Pair is a struct containing a driver and product that are compatible with
// each other.
type Pair struct {
	Driver  Driver
	Product Product
}

func init() {
	for _, d := range Drivers {
		for _, p := range Products {
			if p.IsCompatibleWith(d) {
				DriversByProduct[p] = append(DriversByProduct[p], d)
				ProductsByDriver[d] = append(ProductsByDriver[d], p)
				CompatiblePairs = append(CompatiblePairs, Pair{d, p})
			}
		}
	}
}
