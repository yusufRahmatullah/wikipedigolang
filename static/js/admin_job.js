var offset = 0;
var urlLocks = [];
var timeout = null;

window.onscroll = function(ev) {
    if ((window.innerHeight + window.pageYOffset) >= document.body.offsetHeight) {
        search();
    }
};

function initSearch() {
    offset = 0;
    qs('#job-table-body').innerHTML = '';
    urlLocks.length = 0;  // reset locks
    clearTimeout(timeout);
    timeout = setTimeout(search, 300);
}

function search() {
    var query = qs('#search').value;
    var sortBy = qs('#sort-by').value;
    var order = qs('#order-by').value;
    var filterStatus = qs('#filter-status').value;
    var url = '/api/jobs?offset='+offset+'&query='+query+'&sort='+sortBy+'&order='+order+'&filterStatus='+filterStatus;
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
            var data = d.jobs;
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

function reset() {
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

function errorMessage(msg) {
    var sn = qs('#error_message');
    sn.innerHTML = msg;
    sn.style.opacity = 1;
    sn.style.visibility = 'visible';
}

function hideErrorMessage() {
    var sn = qs('#error_message');
    sn.style.visibility = 'hidden';
    sn.style.opacity = 0;
    sn.innerHTML = '';
}

function loadAvaJobs() {
    var selectNode = qs("#jobName");
    fetch('/api/available_jobs', {
        method: 'GET',
        credentials: 'same-origin'
    }).then(resp => resp.json())
    .then(json => {
        var data = json.data;
        data.forEach(el => {
            var opt = document.createElement("option");
            opt.setAttribute("value", el.name);
            opt.innerHTML = el.name;
            selectNode.append(opt)
        });
        setParamsByJobs();
    })
}

function setParamsByJobs() {
    var jobName = qs('#jobName').value;
    switch(jobName) {
        case 'BanAccountJob':
        case 'MultiAccountJob':
        case 'SingleAccountJob':
        case 'SingleUpdaterJob':
        case 'UpdateIgMediaStatusJob':
            qs('#igId').setAttribute('class', '');
            qs('#igIdLabel').setAttribute('class', '');
            qs('#postId').setAttribute('class', 'hidden');
            qs('#postIdLabel').setAttribute('class', 'hidden');
            break;
        case 'AccountFromPostJob':
        case 'PostExtractionJob':
            qs('#igId').setAttribute('class', 'hidden');
            qs('#igIdLabel').setAttribute('class', 'hidden');
            qs('#postId').setAttribute('class', '');
            qs('#postIdLabel').setAttribute('class', '');
            break;
        case 'MediaFromPostJob':
            qs('#igId').setAttribute('class', '');
            qs('#igIdLabel').setAttribute('class', '');
            qs('#postId').setAttribute('class', '');
            qs('#postIdLabel').setAttribute('class', '');
            break;
        case 'UpdaterJob':
        default:
            qs('#igId').setAttribute('class', 'hidden');
            qs('#igIdLabel').setAttribute('class', 'hidden');
            qs('#postId').setAttribute('class', 'hidden');
            qs('#postIdLabel').setAttribute('class', 'hidden');
            break;
    }
}

function post() {
    var jobNameInput = qs("#jobName");
    var jobName = jobNameInput.value;
    var igIdInput = qs("#igId");
    var igId = '';
    if (igIdInput) {
        igId = igIdInput.value;
    }
    var postIdInput = qs('#postId');
    var postId = '';
    if (postIdInput) {
        postId = postIdInput.value;
    }
    var data = {
        name: jobName,
        params: {
            ig_id: igId,
            post_id: postId,
        }
    }
    apiPost('/api/jobs', data, {
        onSuc: d => {
            igIdInput.value = "";
        },
        onErr: e => {
            errorMessage(e);
        },
    });
}

function doAction(btn, action) {
    spinner(true);
    var id = btn.getAttribute('data-id');
    var data = {
        job_id: id,
        action: action,
    };
    apiPost('/api/jobs/action', data, {
        onSuc: d => {
            // remove row
            var td = qs('#action-'+id);
            var row = td.parentNode;
            row.parentNode.removeChild(row);
        },
        onErr: e => {
            errorMessage(e);
        },
        after: json => {
            spinner(false);
        },
    });
}

function generateButtonsText(status, id) {
    var reqCls = 'primary';
    var delCls = 'secondary';
    if (status == 'active') {
        reqCls = 'hidden';
    } else if (status == 'finished') {
        delCls = 'hidden';
    }
    return `
    <button class="`+reqCls+`" data-id="`+id+`" onclick="doAction(this, 'requeue');">Requeue</button>
    <button class="`+delCls+`" data-id="`+id+`" onclick="doAction(this, 'delete');">Delete</button>
    `;
}

function generateParamsText(params) {
    var retVal = '';
    if (params.ig_id) {
        retVal += 'IG ID: '+params.ig_id;
    }
    if (params.post_id) {
        retVal += '<br /> Post ID: '+params.post_id
    }
    return retVal;
}

function generateStatusText(card) {
    var retVal = card.status;
    if (card.status == 'postponed' && card.reason) {
        retVal += '<br/>Reason: '+card.reason;
    }
    return retVal
}

function appendCards(cards) {
    var tbl = qs('#job-table-body');
    cards.forEach(card => {
        var tr = document.createElement('tr');
        tr.innerHTML = `
        <tr id="job-row-`+card.id+`">
            <td data-label="Name">`+card.name+`</td>
            <td data-label="Params">`+generateParamsText(card.params)+`</td>
            <td data-label="Status" class="capitalize">`+generateStatusText(card)+`</td>
            <td data-label="Action" id="action-`+card.id+`">
                <div class="section button-group">
                    `+generateButtonsText(card.status, card.id)+`
                </div>
            </td>
        </tr>
        `;
        tbl.appendChild(tr);
    });
}