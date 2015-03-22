package  db

import(
	"database/sql"
	cn "golang.org/x/text/encoding/simplifiedchinese"
//	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"bytes"
	"fmt"
)


/*获得数据库的map类型的array*/
func GetResultArray(rows *sql.Rows) []map[string]interface{} {
	cols, _ := rows.Columns()
	count := len(cols)
	var ret []map[string]interface{};
	for rows.Next() {
		row := make(map[string]interface{})
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)
		for i, _ := range cols {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		for i, s := range cols {
			var v interface{}

			val := values[i]

			b, ok := val.([]byte)
			if (ok) {
				data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(b), cn.GBK.NewDecoder()))
				v= string(data)
			} else {
				v = val
			}
			row[s] = v
		}
		ret = append(ret, row);
	}
	return ret;
}

/*获得数据库的map类型单一结果*/
func GetOneResult(rows *sql.Rows) map[string]interface{} {
	cols, _ := rows.Columns()
	count := len(cols)
	fmt.Println("aaaaaaaaaaa")
	row := make(map[string]interface{})
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for i, _ := range cols {
		valuePtrs[i] = &values[i]
	}
	rows.Scan(valuePtrs...)
	fmt.Println("aaaaaaaaaaa")

	for i, s := range cols {
		var v interface{}

		val := values[i]

		b, ok := val.([]byte)
		if (ok) {
			data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(b), cn.GBK.NewDecoder()))
			v= string(data)
		} else {
			v = values[i]
		}
		fmt.Println("aaaaaaaaaaa")
		row[s] = v
	}
	return row;
}
