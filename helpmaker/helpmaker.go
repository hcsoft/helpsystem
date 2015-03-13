package  helpmaker

import (
	"database/sql"
	"github.com/go-martini/martini"
	_ "github.com/mattn/go-adodb"
	"github.com/martini-contrib/render"
	"fmt"
	"strings"
	"html/template"
	dbutil "helpsystem/db"
	erutil "helpsystem/error"
)

func Cats( db *sql.DB , r render.Render, params martini.Params) {
	ret := make(map[string]interface{})
	catid := params["catid"]
	ret["cats"] = GetCats(catid,db)
	r.HTML(200, "index-reveal", ret)
}

func GetCats(  catid string,db *sql.DB)   map[string]map[string]interface{} {
	if catid==""{
		catid = "0"
	}
	ids := strings.Split(catid,",")
	cats := make(  map[string]map[string]interface{})
	for _,v := range ids{
		cat := GetCat(v,db)
		for k, v := range cat {
			cats[k] = v
		}
	}
	return cats
}

func Pages(db *sql.DB , r render.Render, params martini.Params){
	id :=params["id"]
	rows, err := db.Query("select * from help_pages where catid= ?  order by idx",id)
	erutil.CheckErr(err)
	values := dbutil.GetResultArray(rows)
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

func GetCat(catid string, db *sql.DB)map[string]map[string]interface{}{
	if catid == "0"{
		return GetChildCats(catid,db)
	}else{
		cats := make(map[string]map[string]interface{})
		cats[catid] = make(map[string]interface{})
		rows, err := db.Query("select * from help_cat where id = ? order by ord ",catid)
		erutil.CheckErr(err)
		values := dbutil.GetResultArray(rows)
		cats[catid]["data"] = values[0];
		cats[catid]["child"] = GetChildCats(catid,db);
		return cats;
	}
}

func GetChildCats(catid interface{}, db *sql.DB) map[string]map[string]interface{}{
	rows, err := db.Query("select * from help_cat where parentid = ? order by ord ",catid)
	erutil.CheckErr(err)
	values := dbutil.GetResultArray(rows)
	cats := make(map[string]map[string]interface{})
	for _, v := range values {
		id := v["id"].(string)
		cats[id] = make(map[string]interface{})
		cats[id]["data"] = v;
		cats[id]["child"] = GetChildCats(id,db);
	}
	return cats;
}

