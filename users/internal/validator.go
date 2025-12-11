package internal

type Validator interface {
	Validate(request any) error
}
