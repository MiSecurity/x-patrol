Core is a lightweight wrapper of sql.DB.

[![CircleCI](https://circleci.com/gh/go-xorm/core/tree/master.svg?style=svg)](https://circleci.com/gh/go-xorm/core/tree/master)

# Open
```Go
db, _ := core.Open(db, connstr)
```

# SetMapper
```Go
db.SetMapper(SameMapper())
```

## Scan usage

### Scan
```Go
rows, _ := db.Query()
for rows.Next() {
    rows.Scan()
}
```

### ScanMap
```Go
rows, _ := db.Query()
for rows.Next() {
    rows.ScanMap()
```

### ScanSlice

You can use `[]string`, `[][]byte`, `[]interface{}`, `[]*string`, `[]sql.NullString` to ScanSclice. Notice, slice's length should be equal or less than select columns.

```Go
rows, _ := db.Query()
cols, _ := rows.Columns()
for rows.Next() {
    var s = make([]string, len(cols))
    rows.ScanSlice(&s)
}
```

```Go
rows, _ := db.Query()
cols, _ := rows.Columns()
for rows.Next() {
    var s = make([]*string, len(cols))
    rows.ScanSlice(&s)
}
```

### ScanStruct
```Go
rows, _ := db.Query()
for rows.Next() {
    rows.ScanStructByName()
    rows.ScanStructByIndex()
}
```

## Query usage
```Go
rows, err := db.Query("select * from table where name = ?", name)

user = User{
    Name:"lunny",
}
rows, err := db.QueryStruct("select * from table where name = ?Name",
            &user)

var user = map[string]interface{}{
    "name": "lunny",
}
rows, err = db.QueryMap("select * from table where name = ?name",
            &user)
```

## QueryRow usage
```Go
row := db.QueryRow("select * from table where name = ?", name)

user = User{
    Name:"lunny",
}
row := db.QueryRowStruct("select * from table where name = ?Name",
            &user)

var user = map[string]interface{}{
    "name": "lunny",
}
row = db.QueryRowMap("select * from table where name = ?name",
            &user)
```

## Exec usage
```Go
db.Exec("insert into user (`name`, title, age, alias, nick_name,created) values (?,?,?,?,?,?)", name, title, age, alias...)

user = User{
    Name:"lunny",
    Title:"test",
    Age: 18,
}
result, err = db.ExecStruct("insert into user (`name`, title, age, alias, nick_name,created) values (?Name,?Title,?Age,?Alias,?NickName,?Created)",
            &user)

var user = map[string]interface{}{
    "Name": "lunny",
    "Title": "test",
    "Age": 18,
}
result, err = db.ExecMap("insert into user (`name`, title, age, alias, nick_name,created) values (?Name,?Title,?Age,?Alias,?NickName,?Created)",
            &user)
```