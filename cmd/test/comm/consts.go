package comm

const (
	TruncateTable = `truncate table %s;`

	CreateTableTempl = `CREATE TABLE IF NOT EXISTS %s (
		id bigint auto_increment,
		content varchar(1000) not null unique,
		primary key(id)
	) engine=innodb default charset=utf8 auto_increment=1;`

	InsertTemplate = "INSERT INTO %s (content) VALUES %s;"
)
