package gomuti

// DoFunc is a function signature that you can pass to Allow() in order
// to have fine-gained control over the behavior of your mocks. You are
// responsible for validating the number and type of parameters, and ensuring
// that the number of returned values is correct for the mocked method.
type DoFunc func(...interface{}) []interface{}
