package sqltest

// Pair is a struct containing a driver and product that are compatible with
// each other.
type Pair struct {
	Driver  Driver
	Product Product
}

// CompatiblePairs returns all compatible driver/product pairs for
// the given set of products.
//
// If no products are provided compatible pairs are given for all products.
func CompatiblePairs(products ...Product) []Pair {
	if len(products) == 0 {
		products = Products
	}

	var pairs []Pair

	for _, d := range Drivers {
		if d.IsAvailable() {
			for _, p := range products {
				if p.IsCompatibleWith(d) {
					pairs = append(pairs, Pair{d, p})
				}
			}
		}
	}

	return pairs
}
