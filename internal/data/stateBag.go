package data

type StateBag struct {
	ProjectId string
	Region    string
}

func NewStateBag() StateBag {
	return StateBag{}
}
