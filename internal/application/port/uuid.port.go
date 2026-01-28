package port

type UUIDGenerator interface {
	Generate() string
}
