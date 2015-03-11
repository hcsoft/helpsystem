package main

import (
	"database/sql"
	"github.com/go-martini/martini"
	_ "github.com/mattn/go-adodb"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/render"
	"net/http"
	"fmt"
)

func main() {
	m := martini.Classic()
	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("helpsystem_session", store))
	m.Use(render.Renderer())
	//数据库
	db, err := sql.Open("adodb", "Provider=SQLNCLI11;DataTypeCompatibility=80;Server=127.0.0.1;UID=sa;PWD=11111111;Database=helpsystem;")
	checkErr(err)
	db.SetMaxOpenConns(100)
	m.Map(db)
	m.Post("/login", login)
	m.Post("/logout", logout)
	//静态内容
	m.Use(martini.Static("static"))
	//需要权限的内容
	m.Get("/admin", func(r render.Render) {
			r.Redirect("/login.html")
		})
	m.Group("/admin", func(r martini.Router) {
			r.Get("/:id", test)
		}, auth)
	m.Run()
}

func test(db *sql.DB) string {
	rows, err := db.Query("select * from auth_user ")
	checkErr(err)
	fmt.Printf("abc")
	values := getResultArray(rows)
	txt := values[0]["username"].(string)
	fmt.Printf(txt);
	return txt
}

/*获得数据库的map类型的array*/
func getResultArray(rows *sql.Rows) []map[string]interface{} {
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

/*获得数据库的map类型的array*/
func getOneResult(rows *sql.Rows) map[string]interface{} {
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


func login(session sessions.Session, db *sql.DB, r render.Render,req *http.Request) {
	userid := req.FormValue("userid")
	fmt.Println(userid)
	password := req.FormValue("password")
	rows, err := db.Query("select * from auth_user where userid= ? ", userid)
	checkErr(err)
	if rows.Next() {
		values := getOneResult(rows)
		if values["password"] == password {
			session.Set("userid", values["userid"])
			session.Set("username", values["username"])
			r.Redirect("/admin/index.html")
		}else {
			r.Redirect("/login.html&msg=密码错误")
		}
	}else {
		r.Redirect("/login.html&msg=用户名错误")
	}
}

func logout(session sessions.Session) string {
	session.Delete("userid")
	return "删除成功"
}

func auth(session sessions.Session, c martini.Context, r render.Render) {
	v := session.Get("userid")
	if v == nil {
		r.Redirect("/login.html")
	}else {
		c.Next();
	}
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
