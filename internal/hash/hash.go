package hash

type HashStrategy interface {
	Encode(str string) (string, error)
	Compare(enconded, str string) bool
}
