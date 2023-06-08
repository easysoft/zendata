package consts

import "bytes"

const (
	CommandTestFile  = "test/unittest/cmd.yaml"
	CommandTestFile2 = "test/unittest/cmd2.yaml"

	CommandTestFileTables    = "test/unittest/tables.sql"
	CommandTestFileTablesOut = "test/unittest/out/biz_task.yaml"

	CommandTestFileArticle       = "test/unittest/article.txt"
	CommandTestFileArticleOut    = "test/unittest/out/article.yaml"
	CommandTestFileArticleConfig = "test/unittest/article.yaml"

	CommandTestFileProto    = "test/unittest/person.proto"
	CommandTestFileProtoOut = "test/unittest/out/GPBMetadata/Test/Unittest/Person.php"
)

var (
	Buf bytes.Buffer
)
