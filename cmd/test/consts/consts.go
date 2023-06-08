package consts

import "bytes"

const (
	CommandTestFile  = "test/unittest/cmd.yaml"
	CommandTestFile2 = "test/unittest/cmd2.yaml"

	CommandTestFileProto    = "test/unittest/person.proto"
	CommandTestFileProtoOut = "test/unittest/out/GPBMetadata/Test/Unittest/Person.php"
)

var (
	Buf bytes.Buffer
)
