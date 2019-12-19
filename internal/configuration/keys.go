package configuration

type Key string

const (
	KeyAddress Key = "address"
	KeyPort    Key = "port"
)

func (k Key) String() string {
	return string(k)
}
