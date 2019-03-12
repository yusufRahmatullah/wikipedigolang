var offset = 0;
var urlLocks = [];
var cs = document.currentScript;

window.onscroll = function (ev) {
    if ((window.innerHeight + window.pageYOffset) >= document.body.offsetHeight) {
        search();
    }
};

function getScriptAttr(attr) {
    return cs.getAttribute(attr) || "";
}

function search() {
    var srchUrl = getScriptAttr('search');
    var url = '/api/igmedias/' + srchUrl + '&offset=' + offset;
    if (urlLocks.includes(url)) {
        return
    }
    urlLocks.push(url);
    spinner(true);
    apiGet(url, {
        before: j => {
            spinner(false);
        },
        onSuc: d => {
            var data = d.medias;
            if (offset == 0 && data == null) {
                notFound(true);
            } else {
                notFound(false);
                appendCards(data);
                offset += data.length;
            }
            urlLocks.length = 0;
        },
        onErr: e => {
            notFound(true);
        },
    });
}

function reset() {
    offset = 0;
    search();
}

function notFound(v) {
    var n = qs('#not-found');
    if (v) {
        n.style.visibility = 'visible';
    } else {
        n.style.visibility = 'hidden';
    }
}

function spinner(b) {
    var n = qs('.spinner');
    n.style.visibility = b ? 'visible' : 'hidden';
}