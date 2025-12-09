package types

type BoolType string

const (
	BoolTypeTrue  BoolType = "true"
	BoolTypeFalse BoolType = "false"
	BoolTypeAll   BoolType = "all"
)

func (b BoolType) Bool() bool {
	switch b {
	case BoolTypeTrue:
		return true
	case BoolTypeFalse:
		return false
	default:
		return false
	}

}
