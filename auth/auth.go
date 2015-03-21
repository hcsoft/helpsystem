package auth

import ("database/sql"
	"github.com/go-martini/martini"
	_ "github.com/mattn/go-adodb"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/render"
	"net/http"
	"fmt"
	dbutil "github.com/hcsoft/helpsystem/db"
	erutil "github.com/hcsoft/helpsystem/error"
)

func Login(session sessions.Session, db *sql.DB, r render.Render, req *http.Request) {
	userid := req.FormValue("userid")
	if userid == "" {
		r.HTML(200, "login", "请登录")
		return
	}
	fmt.Println(userid)
	password := req.FormValue("password")
	rows, err := db.Query("select * from auth_user where userid= ? ", userid)
	erutil.CheckErr(err)
	if rows.Next() {
		values := dbutil.GetOneResult(rows)
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

func Logout(session sessions.Session, r render.Render) {
	session.Delete("userid")
	r.HTML(200, "login", "登出成功")
}

func Auth(session sessions.Session, c martini.Context, r render.Render) {
	fmt.Println("auth..........")
	v := session.Get("userid")
	if v == nil {
		r.Redirect("/login")
	}else {
		c.Next();
	}
}
