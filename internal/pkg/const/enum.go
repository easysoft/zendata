package consts

type RunMode int

const (
	RunModeGen RunMode = iota
	RunModeServer
	RunModeServerRequest
	RunModeMockPreview
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
	Success         = ResponseCode{0, "Request Successfully"}
	CommErr         = ResponseCode{100, "Common Error"}
	ParamErr        = ResponseCode{200, "Parameter Error"}
	UnAuthorizedErr = ResponseCode{401, "UnAuthorized"}
)

type OpenApiDataType string

const (
	OpenApiDataTypeInteger  OpenApiDataType = "integer"
	OpenApiDataTypeLong     OpenApiDataType = "long"
	OpenApiDataTypeFloat    OpenApiDataType = "float"
	OpenApiDataTypeDouble   OpenApiDataType = "double"
	OpenApiDataTypeString   OpenApiDataType = "string"
	OpenApiDataTypeByte     OpenApiDataType = "byte"
	OpenApiDataTypeBinary   OpenApiDataType = "binary"
	OpenApiDataTypeBoolean  OpenApiDataType = "boolean"
	OpenApiDataTypeDate     OpenApiDataType = "date"
	OpenApiDataTypeDateTime OpenApiDataType = "dateTime"
	OpenApiDataTypePassword OpenApiDataType = "password"
)

type OpenApiDataFormat string

const (
	// integer
	OpenApiDataFormatInt32 OpenApiDataFormat = "int32"
	OpenApiDataFormatInt64 OpenApiDataFormat = "int64"

	// float
	OpenApiDataFormatFloat  OpenApiDataFormat = "float"
	OpenApiDataFormatDouble OpenApiDataFormat = "double"

	// date
	OpenApiDataFormatDate     OpenApiDataFormat = "date"
	OpenApiDataFormatDateTime OpenApiDataFormat = "date-time"

	OpenApiDataFormatByte     OpenApiDataFormat = "byte"
	OpenApiDataFormatBinary   OpenApiDataFormat = "binary"
	OpenApiDataFormatPassword OpenApiDataFormat = "password"
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

type ColumnType string

const (
	// number
	Bit       ColumnType = "bit"
	Tinyint   ColumnType = "tinyint"
	Smallint  ColumnType = "smallint"
	Mediumint ColumnType = "mediumint"
	Int       ColumnType = "int"
	Bigint    ColumnType = "bigint"
	Float     ColumnType = "float"
	Double    ColumnType = "double"

	// fixed-point
	Decimal ColumnType = "decimal"

	// character string
	Char       ColumnType = "char"
	Varchar    ColumnType = "varchar"
	Tinytext   ColumnType = "tinytext"
	Text       ColumnType = "text"
	Mediumtext ColumnType = "mediumtext"
	Longtext   ColumnType = "longtext"

	// binary data
	Tinyblob   ColumnType = "tinyblob"
	Blob       ColumnType = "blob"
	Mediumblob ColumnType = "mediumblob"
	Longblob   ColumnType = "longblob"
	Binary     ColumnType = "binary"
	Varbinary  ColumnType = "varbinary"

	// date and time type
	Date      ColumnType = "date"
	Time      ColumnType = "time"
	Year      ColumnType = "year"
	Datetime  ColumnType = "datetime"
	Timestamp ColumnType = "timestamp"

	// other type
	Enum               ColumnType = "enum"
	Set                ColumnType = "set"
	Geometry           ColumnType = "geometry"
	Point              ColumnType = "point"
	Linestring         ColumnType = "linestring"
	Polygon            ColumnType = "polygon"
	Multipoint         ColumnType = "multipoint"
	Multilinestring    ColumnType = "multilinestring"
	Multipolygon       ColumnType = "multipolygon"
	Geometrycollection ColumnType = "geometrycollection"
	Json               ColumnType = "json"
)

func (e ColumnType) String() string {
	return string(e)
}

type VarcharType string

const (
	Empty        VarcharType = ""
	Username     VarcharType = "username"
	Email        VarcharType = "email"
	Url          VarcharType = "url"
	Ip           VarcharType = "ip"
	Mac          VarcharType = "mac"
	CreditCard   VarcharType = "creditcard"
	IdCard       VarcharType = "idcard"
	MobileNumber VarcharType = "mobilenumber"
	TelNumber    VarcharType = "telnumber"
	Token        VarcharType = "token"
	Uuid         VarcharType = "uuid"
	JsonStr      VarcharType = "jsonstr"
	Md5          VarcharType = "md5"
	//UnixTime     VarcharType = "unixtime"
)

func (e VarcharType) String() string {
	return string(e)
}

type GeneratedBy string

const (
	ByRange GeneratedBy = "range"
	ByRefer GeneratedBy = "refer"
)

func (e GeneratedBy) String() string {
	return string(e)
}
