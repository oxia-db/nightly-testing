package operations

type Type string

const (
	TypeSequence = "sequence"
)

type Operation interface {
	Type() Type
}
