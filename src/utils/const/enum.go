package constant

type FieldType int
const (
	LIST FieldType = iota
	TIMESTAMP
	IP
	SESSION
)
func (c FieldType) String() string {
	switch c {
	case LIST:
		return "list"
	case TIMESTAMP:
		return "timestamp"
	case IP:
		return "ip"
	case SESSION:
		return "session"
	}

	return "n/a"
}
