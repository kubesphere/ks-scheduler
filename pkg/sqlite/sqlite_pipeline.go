package sqlite

import (
	"database/sql"
	log "github.com/golang/glog"
	_ "github.com/mattn/go-sqlite3"
)

const (
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

var KeyNodeCilent = InitKeyNodeTable()

func InitKeyNodeTable() (KeyNodeTable) {
	db, err := sql.Open("sqlite3", "/Users/xiaoyangzhu/work/test/sqlite/test.db")
	err = db.Ping()
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

func (kn *KeyNodeTable) KeyNodeInsert(key, nodeName string, count int) (id int64, err error) {
	tx, err := kn.SQLiteDB.Begin()
	stmt, err := tx.Prepare("INSERT INTO key_node(podkey, nodename, count) VALUES (?,?,?)")
	if err != nil {
		return 1, err
	}
	defer stmt.Close()

	stmt.Exec(key, nodeName, count)
	checkErr(err)
	tx.Commit()
	return id, nil
}

func (kn *KeyNodeTable) KeyNodeSearch(podKey string, nodeName string) ([]KeyNodeTable, error) {
	rows, err := kn.SQLiteDB.Query("SELECT *  FROM key_node WHERE podkey = ? AND nodename = ?", podKey,nodeName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []KeyNodeTable{}

	for rows.Next() {
		var id, count int
		var podkey, nodename string

		err = rows.Scan(&id, &podkey, &nodename, &count)
		checkErr(err)

		res = append(res, KeyNodeTable{
			Id:     id,
			PodKey: podkey,
			Count: count,
		})
	}
	return res, nil

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
		log.Fatal(err)
		panic(err)
	}
}
