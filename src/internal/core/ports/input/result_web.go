package driverport

type ResultAPI interface {
	GetResult() error
	GetResultHistory() error
}
