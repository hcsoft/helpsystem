package main

import (
	"database/sql"
	"github.com/go-martini/martini"
	_ "github.com/mattn/go-adodb"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/render"
	"net/http"
	"fmt"
	"strings"
	"html/template"
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
	catid := req.FormValue("catid")

	ids := strings.Split(catid,",")
	cats := make(  map[string]map[string]interface{})
	for _,v := range ids{
		cat := getCats(v,db)
		for k, v := range cat {
			cats[k] = v
		}
	}
	ret["cats"] = cats

	r.HTML(200, "index-reveal", ret)
}
func pages(db *sql.DB , r render.Render, params martini.Params){
	id :=params["id"]
	rows, err := db.Query("select * from help_pages where catid= ?  order by idx",id)
	checkErr(err)
	values := getResultArray(rows)
	for _,value := range values{
		url :=value["url"].(string);
		fmt.Println(url);
		fmt.Println(strings.Index(url,","));

		if strings.Index(url,",") >0{
			value["isarray"] = true
			var urls  []interface{};

			strs := strings.Split(url,",")
			for _,str := range strs{
				fmt.Println(str);
				cols :=strings.Split(str,"|")
				fmt.Println(cols[0])
				fmt.Println(cols[1])
				colmap := make(map[string]interface{})
				if len(cols)>0 && cols[0] != "" {
					colmap["type"]=cols[0]
				}else{
					colmap["type"]="fragment"
				}
				if len(cols)>1 && cols[1] != "" {
					colmap["url"]=cols[1]
				}else{
					colmap["url"]=""
				}
				if len(cols)>2 && cols[2] != "" {
					colmap["css"]=template.CSS(cols[2])
				}else{
					colmap["css"]=template.CSS("width:100%;height:100%;")
				}
				if len(cols)>3  && cols[3] != ""  {
					colmap["animate"]=(cols[2])
				}else{
					colmap["animate"]=""
				}
				urls = append(urls,colmap)
			}
			value["urls"] = urls;
		}
	}
	r.HTML(200, "slide", values)
}

func getCats(catid string, db *sql.DB)map[string]map[string]interface{}{
	if catid == "0"{
		return getChildCats(catid,db)
	}else{
		cats := make(map[string]map[string]interface{})
		cats[catid] = make(map[string]interface{})
		rows, err := db.Query("select * from help_cat where id = ? order by ord ",catid)
		checkErr(err)
		values := getResultArray(rows)
		cats[catid]["data"] = values[0];
		cats[catid]["child"] = getChildCats(catid,db);
		return cats;
	}
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
