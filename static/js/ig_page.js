var offset = 0;
var urlLocks = [];
var cs = document.currentScript;

window.onscroll = function(ev) {
    if ((window.innerHeight + window.pageYOffset) >= document.body.offsetHeight) {
        search();
    }
};

function initSearch() {
    offset = 0;
    qs('#card-container').innerHTML = '';
    urlLocks.length = 0;  // reset locks
    search();
}

function getScriptAttr(attr) {
    return cs.getAttribute(attr) || "";
}

function search() {
    var query = qs('#search').value;
    var sortBy = qs('#sort-by').value;
    var order = qs('#order-by').value;
    var srchUrl = getScriptAttr('search');
    var pageCode = getScriptAttr('page');
    var url = '/api/'+pageCode+srchUrl+'?offset='+offset+'&query='+query+'&sort='+sortBy+'&order='+order;
    if (urlLocks.includes(url)) {
        return
    }
    urlLocks.push(url);
    spinner(true);
    apiGet(url, {
        before: j => {
            var currentQuery = qs('#search').value;
            if (currentQuery !== j.data.query) {
                return
            }
            spinner(false);
        },
        onSuc: d => {
            var page = getScriptAttr('page');
            var data;
            if (page == 'igmedias') {
                data = d.medias;
            } else if (page == 'igprofiles') {
                data = d.profiles;
            }
            if (offset == 0 && data == null) {
                notFound(true);
            } else {
                notFound(false);
                appendCards(data);
                offset += data.length;
            }
            urlLocks.length = 0;
        },
        onErr : e => {
            notFound(true);            
        },
    });
}

function initCount() {
    var page = getScriptAttr('page');
    var countUrl = getScriptAttr('count');
    var subtitle = getScriptAttr('subtitle');
    apiGet('/api/'+page+countUrl, {
        onSuc: d => {
            qs('#subtitle').innerHTML = 'IGO '+subtitle+' <sub>('+d+' '+subtitle.toLowerCase()+')</sub>';
        }
    });
}

function reset() {
    initCount();
    offset = 0;
    qs('#search').value = '';
    initSearch();
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
