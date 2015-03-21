package  db

import(
	"database/sql"
	_ "github.com/mattn/go-adodb"
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
			row[s] = values[i]
		}
		ret = append(ret, row);
	}
	return ret;
}

/*获得数据库的map类型单一结果*/
func GetOneResult(rows *sql.Rows) map[string]interface{} {
	cols, _ := rows.Columns()
	count := len(cols)
	row := make(map[string]interface{})
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for i, _ := range cols {
		valuePtrs[i] = &values[i]
	}
	rows.Scan(valuePtrs...)
	for i, s := range cols {
		row[s] = values[i]
	}
	return row;
}
