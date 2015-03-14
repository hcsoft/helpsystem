package admin
import (
	"database/sql"
	"github.com/go-martini/martini"
	_ "github.com/mattn/go-adodb"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/render"
	erutil "helpsystem/error"
	"helpsystem/helpmaker"
)

func Router( router martini.Router) {
	router.Get("", func(r render.Render, session sessions.Session) {
			r.HTML(200, "admin/index", session.Get("username"))
		})
	router.Get("/index", func(r render.Render, session sessions.Session) {
			r.HTML(200, "admin/index", session.Get("username"))
		})
	router.Get("/EditPages/:id",helpmaker.EditPages)
	router.Get("/helpmanager", func(r render.Render, session sessions.Session, db *sql.DB) {
			r.HTML(200, "admin/helpmanager", helpmaker.GetCats("0", db))
		})
	router.Get("/helpcatsave/:id/:parentid/:ord/:name", func(r render.Render, params martini.Params, db *sql.DB) string {
			id := params["id"]
			parentid := params["parentid"]
			name := params["name"]
			ord := params["ord"]
			rows , err := db.Query("select * from help_cat where id= ? ", id)
			erutil.CheckErr(err)
			if rows.Next(){
				_ , err := db.Exec("update help_cat set parentid=? ,name=? ,ord=? where id= ? ",parentid,name , ord, id)
				erutil.CheckErr(err)
				return "保存成功";
			}else{
				_ , err := db.Exec("insert into help_cat (id,name,parentid,ord)values(?,?,?,?) ",id,name, parentid,ord)
				erutil.CheckErr(err)
				return "保存成功"
			}
			return "保存失败"
		})
	router.Get("/helpcatdel/:id", func(r render.Render, params martini.Params, db *sql.DB) string {
			id := params["id"]
			_ , err := db.Exec("delete help_cat  where id = ? ",id)
			erutil.CheckErr(err)
			return "删除成功"
		})
}

