package exceptions

type Exception interface {
	Error() string
	Code() int
}
