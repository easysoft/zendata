package consts

type FieldType int

const (
	LIST FieldType = iota
	TIMESTAMP
	IP
	SESSION_ID
)

func (c FieldType) String() string {
	switch c {
	case LIST:
		return "list"
	case TIMESTAMP:
		return "timestamp"
	case IP:
		return "ip"
	case SESSION_ID:
		return "session"
	}

	return "n/a"
}

type RunMode int

const (
	RunModeGen RunMode = iota
	RunModeServer
	RunModeServerRequest
)

func (c RunMode) String() string {
	switch c {
	case RunModeGen:
		return "gen"
	case RunModeServer:
		return "server"
	}
	return "unknown"
}

type ResultStatus int

const (
	PASS ResultStatus = iota
	FAIL
)

func (c ResultStatus) String() string {
	switch c {
	case PASS:
		return "pass"
	case FAIL:
		return "fail"
	}

	return "UNKNOWN"
}

type ResponseCode struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

var (
	Success  = ResponseCode{0, "Request Successfully"}
	CommErr  = ResponseCode{100, "Common Error"}
	ParamErr = ResponseCode{200, "Parameter Error"}
)

type OpenApiDataType string

const (
	Integer  OpenApiDataType = "integer"
	Long     OpenApiDataType = "long"
	Float    OpenApiDataType = "float"
	Double   OpenApiDataType = "double"
	String   OpenApiDataType = "string"
	Byte     OpenApiDataType = "byte"
	Binary   OpenApiDataType = "binary"
	Boolean  OpenApiDataType = "boolean"
	Date     OpenApiDataType = "date"
	DateTime OpenApiDataType = "dateTime"
	Password OpenApiDataType = "password"
)

type OpenApiSchemaType string

const (
	SchemaTypeString  OpenApiSchemaType = "string" // this includes dates and files
	SchemaTypeNumber  OpenApiSchemaType = "number"
	SchemaTypeFloat   OpenApiSchemaType = "integer"
	SchemaTypeBoolean OpenApiSchemaType = "boolean"
	SchemaTypeArray   OpenApiSchemaType = "array"
	SchemaTypeObject  OpenApiSchemaType = "object"
)
