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
	m.Any("/login", login)
	m.Get("/logout", logout)
	m.Get("/", index)
	m.Get("/pages/:id", pages)

	//静态内容
	m.Use(martini.Static("static"))
	//需要权限的内容
	m.Get("/admin", auth, func(r render.Render, session sessions.Session) {
			r.HTML(200, "admin/index", session.Get("username"))
		})

	m.Group("/admin", func(r martini.Router) {
			m.Get("/index", func(r render.Render, session sessions.Session) {
					r.HTML(200, "admin/index", session.Get("username"))
				})
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

func index( db *sql.DB , r render.Render, req *http.Request) {
	ret := make(map[string]interface{})
	rows, err := db.Query("select * from help_pages ")
	checkErr(err)
	values := getResultArray(rows)
	ret["test"] = values
	catid := req.FormValue("catid")
	if catid == "" {
		catid = "0"
	}
	ret["cats"] = getChildCats(catid,db)

	r.HTML(200, "index-reveal", ret)
}
func pages(db *sql.DB , r render.Render, params martini.Params){
	id :=params["id"]
	rows, err := db.Query("select * from help_pages where catid= ? ",id)
	checkErr(err)
	values := getResultArray(rows)
	r.HTML(200, "slide", values)
}

func getChildCats(catid interface{}, db *sql.DB) map[string]map[string]interface{}{
	rows, err := db.Query("select * from help_cat where parentid = ? order by ord ",catid)
	checkErr(err)
	values := getResultArray(rows)
	cats := make(map[string]map[string]interface{})
	for _, v := range values {
		id := v["id"].(string)
		cats[id] = make(map[string]interface{})
		cats[id]["data"] = v;
		cats[id]["child"] = getChildCats(id,db);
	}
	return cats;
}


func login(session sessions.Session, db *sql.DB, r render.Render, req *http.Request) {
	userid := req.FormValue("userid")
	if userid == "" {
		r.HTML(200, "login", "请登录")
		return
	}
	fmt.Println(userid)
	password := req.FormValue("password")
	rows, err := db.Query("select * from auth_user where userid= ? ", userid)
	checkErr(err)
	if rows.Next() {
		values := getOneResult(rows)
		if values["password"] == password {
			session.Set("userid", values["userid"])
			session.Set("username", values["username"])
			r.Redirect("/admin")
			//			r.HTML(200, "admin/index", nil)
		}else {
			r.HTML(200, "login", "密码错误")
		}
	}else {
		r.HTML(200, "login", "用户名错误")
	}
}

func logout(session sessions.Session, r render.Render) {
	session.Delete("userid")
	r.HTML(200, "login", "登出成功")
}

func auth(session sessions.Session, c martini.Context, r render.Render) {
	fmt.Println("auth..........")
	v := session.Get("userid")
	if v == nil {
		r.Redirect("/login")
	}else {
		c.Next();
	}
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
