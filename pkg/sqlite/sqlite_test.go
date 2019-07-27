package sqlite

import (
	"fmt"
	"testing"
)

func TestInsterData(t *testing.T) {
	type parseSqliteData struct {
		key      string
		nodeName string
	}
	dbCilent := InitKeyNodeTable()

	tests := []parseSqliteData{
		{"keya", "node1"},
		{"keyab", "node1"},
		{"keyac", "node2"},
		{"keyad", "node1"},
	}

	for _, lb := range tests {
		if _, err := dbCilent.KeyNodeInsert(lb.key, lb.nodeName, 0); err != nil {
			fmt.Println(err)
			t.Errorf("insert to db error: ")
		}
	}
	row,_ := KeyNodeCilent.KeyNodeSearch("jenkins", "node1")
	fmt.Println(row)
	//os.Remove("/Users/xiaoyangzhu/work/test/sqlite/test.db")
}
