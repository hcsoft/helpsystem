$(".helpmanager li.dd-item > div").click(function(){
    $(this).parent().toggleClass("expend");
});

function add(flag,obj){
    var item = $(obj).parent().parent().parent();
    var subitem = $(obj).parent().parent();
    var text = $("#cattemplate").html();
    text = text.replace("<<id>>",guid());
    var ord = item.data("ord");
    if(!$.isNumeric(ord)){
        ord = 0 ;
    }
    if(flag==1){
        text = text.replace("<<parentid>>",item.data("parentid"));
        text = text.replace("<<ord>>",""+(ord-1));
        item.before(text)
    }else if(flag==2){
        text =  text.replace("<<parentid>>",item.data("parentid"));
        text = text.replace("<<ord>>",""+(ord+1));
        item.after(text)
    }else if(flag==3){
        text = text.replace("<<parentid>>",item.data("id"));
        item.addClass("expend");
        if(subitem.next("ol.dd-list").length==0){
            subitem.after('<ol class="dd-list"></ol>');
            ord = 0 ;
        }else{
            ord =  subitem.next("ol.dd-list").find(">li:last-child").data("ord");
            if(!$.isNumeric(ord)){
                ord = 0 ;
            }
            ord = ord+1;
        }
        text = text.replace("<<ord>>",""+(ord));
        subitem.next("ol.dd-list").append(text)
    }
    event.stopPropagation();
}

function addfirst(obj){
    var olparent = $(obj).parent().parent().next().children();
    var ord = 0;
    if(olparent.find(">ol.dd-list").length == 0){
        olparent.append('<ol class="dd-list"></ol>');
    }else{
        ord =  olparent.find(">ol.dd-list").find(">li:last-child").data("ord");
        if(!$.isNumeric(ord)){
            ord = 0 ;
        }
        ord = ord+1;
    }
    var text = $("#cattemplate").html();
    text = text.replace("<<id>>",guid());
    text = text.replace("<<parentid>>","0");
    text = text.replace("<<ord>>",""+ord);
    olparent.find(">ol.dd-list").append(text)
    event.stopPropagation();
}

function namechanged(obj){
    $(obj).parent().parent().parent().addClass("changed");
}

function save(obj){
    var item = $(obj).parent().parent();
    saveCat(item,function(){
        item.removeClass("changed");
    });
    event.stopPropagation();
}
function saveCat(item,func){
    var parent = item.parent("ol").parent("li.dd-item");
    if (parent.length>0){
        saveCat(parent,function(){
            parent.removeClass("changed");
            dosavecat(item,func);
        })
    }else{
        dosavecat(item,func);
    }

}
function dosavecat(item,func){
    var name = item.find(">div.dd2-content>.itemname>input").val();
    $.ajax({
        url: "/admin/helpcatsave/"+item.data("id")+"/"+item.data("parentid")+"/"+item.data("ord")+"/"+name,
        contentType:"html/text"
    }).done(function(data) {
        console.log(data)
        if(func)
            func()
    });
}

function del(obj){
    console.log("1111111111111")
    BootstrapDialog.confirm({
        title: '删除确认',
        message: '是否确定删除目录! 将会连同子目录一起删除!',
        btnOKLabel: '删除',
        btnCancelLabel: '取消',
        callback:function(result){
            if(result) {
                var item = $(obj).parent().parent().parent();
                delCat(item,function(){
                    item.remove();
                });
            }
        }
    });
    event.stopPropagation();
}
function delCat(item,func){
    var childs = item.find(">ol>li");
    var childol = item.find(">ol");
    if (childs.length>0){
        for(var idx = 0 ; idx < childs.length;idx++){
            delCat($(childs[idx]),function(){
                $(childs[idx]).remove();
                dodelcat(item,func);
                if(item.find(">ol>li").length==0){
                    childol.remove();
                }
            })
        }
    }else{
        dodelcat(item,func);
    }

}
function dodelcat(item,func){
    $.ajax({
        url: "/admin/helpcatdel/"+item.data("id"),
        contentType:"html/text"
    }).done(function(data) {
        console.log(data)
        if(func)
            func()
    });
}

function showContent(obj){
    var item = $(obj).parent().parent().parent();
    $.ajax({
        url: "/admin/EditPages/"+item.data("id"),
        contentType:"html/text"
    }).done(function(data) {
        console.log(data)
    });
    event.stopPropagation();
}