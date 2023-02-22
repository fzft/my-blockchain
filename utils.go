package main

func filter(slice []*Transaction, f func(transaction *Transaction) bool) []*Transaction {
	var n []*Transaction
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func reducer(transactions []*Transaction) float64 {
	var sum float64
	for _, trans := range transactions {
		sum += trans.Amount
	}
	return sum
}