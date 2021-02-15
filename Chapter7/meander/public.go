package meander

// Facade is an inteface containing Public method
type Facade interface {
	Public() interface{}
}

// Public is a function checking if o can be a Facade and if so returning
// their respective Public() result
func Public(o interface{}) interface{} {
	if p, ok := o.(Facade); ok {
		return p.Public()
	}
	return o
}
