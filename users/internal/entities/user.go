package entities

type Address struct {
	Street  string
	City    string
	ZipCode string
}

type User struct {
	ID      string
	Name    string
	Email   string
	Address Address
}
