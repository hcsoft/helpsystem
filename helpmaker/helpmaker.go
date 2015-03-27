package  helpmaker

import (
	"database/sql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"fmt"
	"strings"
	dbutil "github.com/hcsoft/helpsystem/db"
	erutil "github.com/hcsoft/helpsystem/error"
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
	values := getPageData(db,r,params)
	r.HTML(200, "slide", values)
}

func getPageData(db *sql.DB , r render.Render, params martini.Params)  []map[string]interface{}{
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

			strs := strings.Split(url,",")

			urls := make ([]map[string]string,0)
			for _,v := range strs{

				ext :=strings.ToLower( v[len(v)-4:len(v)])
				fmt.Println(ext)
				item := make(map[string]string)
				item ["url"] = v;
				if ext ==".jpg" || ext == ".png" || ext ==".gif" {
					item["type"]="pic"
				}else{
					item["type"]="video"
				}
				urls = append(urls,item)
			}
			value["urls"] = urls
		}else{
			ext :=strings.ToLower( url[len(url)-4:len(url)])
			fmt.Println(ext)
			if ext ==".jpg" || ext == ".png" || ext ==".gif" {
				value["type"]="pic"
			}else{
				value["type"]="video"
			}
		}
	}
	return values;
}

func EditPages(db *sql.DB , r render.Render, params martini.Params){
	values := getPageData(db,r,params)
	r.HTML(200, "admin/cat/content", values)
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

