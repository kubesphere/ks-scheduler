package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

const (
	DATA_PATH = "/home/data"
	CREATE_TABLE_SQL = `
		CREATE TABLE IF NOT EXISTS key_node 
		(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			podkey VARCHAR(64) NULL,
			nodename VARCHAR(64) NULL,
			count INT(10) NULL
		)
	;`
)

type KeyNodeTable struct {
	SQLiteDB *sql.DB `json:",omitempty"`
}

func InitKeyNodeTable() (KeyNodeTable) {
	db, err := sql.Open("sqlite3", os.Getenv(DATA_PATH))
	checkErr(err)
	CreateTable(db, CREATE_TABLE_SQL)
	KeyNodeCilent := KeyNodeTable{
		SQLiteDB: db,
	}

	return KeyNodeCilent
}

func CreateTable(db *sql.DB, sql string) {
	_, error := db.Exec(sql)
	checkErr(error)
}

func (kn *KeyNodeTable) KeyNodeInsert(key, nodeName string) (id int64) {
	stmt, err := kn.SQLiteDB.Prepare("INSERT INTO userinfo(podkey, nodename, count) values(?,?,?)")
	checkErr(err)

	res, err := stmt.Exec(key, nodeName, 0)
	checkErr(err)
	id, _ = res.LastInsertId()
	return id
}

func (kn *KeyNodeTable) KeyNodeSearch() {

	rows, err := kn.SQLiteDB.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var id int
		var podkey string
		var nodename string
		var count int
		err = rows.Scan(&id, &podkey, &nodename, &count)
		checkErr(err)
		fmt.Println(id)
		fmt.Println(podkey)
		fmt.Println(nodename)
		fmt.Println(count)
	}

}

//func (kn *KeyNodeTable) update() {
//
//}

//func sqlite() {
//	db, err := sql.Open("sqlite3", "/Users/xiaoyangzhu/work/test/sqlite/test.db")
//	checkErr(err)
//
//	sqlStmt := `
//	create table test (id integer not null primary key, name text);
//	delete from foo;
//	`
//
//	//insert
//	stmt, err := db.Prepare("INSERT INTO userinfo(username, department, created) values(?,?,?)")
//	checkErr(err)
//
//	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
//	checkErr(err)
//
//	id, err := res.LastInsertId()
//	checkErr(err)
//
//	fmt.Println(id)
//	//update
//	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
//	checkErr(err)
//
//	res, err = stmt.Exec("astaxieupdate", id)
//	checkErr(err)
//
//	affect, err := res.RowsAffected()
//	checkErr(err)
//
//	fmt.Println(affect)
//
//	//search
//	rows, err := db.Query("SELECT * FROM userinfo")
//	checkErr(err)
//
//	for rows.Next() {
//		var uid int
//		var username string
//		var department string
//		var created time.Time
//		err = rows.Scan(&uid, &username, &department, &created)
//		checkErr(err)
//		fmt.Println(uid)
//		fmt.Println(username)
//		fmt.Println(department)
//		fmt.Println(created)
//	}
//
//	//delete
//	stmt, err = db.Prepare("delete from userinfo where uid=?")
//	checkErr(err)
//
//	res, err = stmt.Exec(id)
//	checkErr(err)
//
//	affect, err = res.RowsAffected()
//	checkErr(err)
//
//	fmt.Println(affect)
//
//	db.Close()
//
//}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
