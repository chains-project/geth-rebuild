package buildconfig

type Spec interface {
	ToMap() map[string]string
	String() string
}
