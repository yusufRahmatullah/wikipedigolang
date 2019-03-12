function doAction(btn, action) {
    var igid = btn.getAttribute('data-igId');
    var id = btn.getAttribute('data-id');
    var data = {
        ig_id: igid,
        action: action,
    };
    apiPost('/api/igprofiles/action', data, {
        onSuc: d => {
            toggleAction(igid, id, action);
        },
        onErr: e => {
            alert(e);
        },
    });
}

function toggleAction(igid, id, act) {
    // set buttons
    var bg = qs('#btnGroup-'+id);
    var status = 'active';
    if (act == 'ban') {
        status = 'banned';
    } else if (act == 'asMulti') {
        status = 'multi';
    } else if (act == 'update') {
        window.location.reload(true);
        return
    }
    bg.innerHTML = generateButtonsText(status, igid, id);
    // set status section
    var ss = qs('#status-'+id);
    ss.innerHTML = status;
}

function generateButtonsText(status, igid, id) {
    var actDis = '';
    var actCls = 'primary';
    var banDis = '';
    var banCls = 'secondary';
    var asMultiDis = '';
    var asMultiCls = 'tertiary';
    var updDis = '';
    var updCls = 'primary';
    if (status == 'active') {
        actCls = 'hidden';
        actDis = 'disabled';
    } else if (status == 'banned') {
        banCls = 'hidden';
        banDis = 'disabled';
        updCls = 'hidden';
        updDis = 'disabled';
    } else if (status == 'multi') {
        asMultiCls = 'hidden';
        asMultiDis = 'disabled';
        updCls = 'tertiary';
    } else if (status == 'banned_multi') {
        banCls = 'hidden';
        banDis = 'disabled';
    }
    return `
        <button class="`+updCls+` data-id="`+id+`" data-igId="`+igid+`" onclick="doAction(this, 'update');" `+updDis+`>Update</button>
        <button class="`+actCls+`" data-id="`+id+`" data-igId="`+igid+`" onclick="doAction(this, 'activate');" `+actDis+`>Activate</button>
        <button class="`+banCls+`" data-id="`+id+`" data-igId="`+igid+`" onclick="doAction(this, 'ban');" `+banDis+`>Ban</button>
        <button class="`+asMultiCls+`" data-id="`+id+`" data-igId="`+igid+`" onclick="doAction(this, 'asMulti');" `+asMultiDis+`>As Multi</button>
    `;
}

function appendCards(cards) {
    var cc = qs('#card-container');
    cards.forEach(card => {                
        var col = document.createElement('div');
        col.setAttribute('class', 'col-sm-4 col-md-3')
        col.innerHTML = `
        <div class="card">
            <img src="`+card.pp_url+`" style="height:100%" onerror='if (this.src != "/static/default.png") this.src = "/static/default.png";'/>
            <div class="section double-padded center">`+card.name+`</div>
            <div class="section dark center">
                <a href="/admin/igprofile/`+card.ig_id+`" target="_blank">
                    @`+card.ig_id+`
                </a>
            </div>
            <div class="section center">
                Status: <span id="status-`+card.id+`">`+card.status+`</span>
            </div>
            <div class="section button-group" id="btnGroup-`+card.id+`">
                `+generateButtonsText(card.status, card.ig_id, card.id)+`
            </div>
        </div>`;
        cc.appendChild(col);
    });
}