package buildspec

type Spec interface {
	ToMap() map[string]string
	String() string
}
