function loadHtml(url,selector){
    $.ajax({
        url: url,
        contentType:"html/text"
    }).done(function(data) {
        $(selector).html(data);
    });
}

function guid() {
    function s4() {
        return Math.floor((1 + Math.random()) * 0x10000)
            .toString(16)
            .substring(1);
    }
    return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
        s4() + '-' + s4() + s4() + s4();
}