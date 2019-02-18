function doAction(btn, action) {
    var igid = btn.getAttribute('data-igId');
    var data = {
        ig_id: igid,
        action: action,
    };
    apiPost('/api/igprofiles/action', data, {
        onSuc: d => {
            toggleAction(igid, action);
        },
        onErr: e => {
            alert(e);
        },
    });
}

function toggleAction(nid, act) {
    // set buttons
    var bg = qs('#btnGroup-'+nid);
    var status = 'active';
    if (act == 'ban') {
        status = 'banned';
    } else if (act == 'asMulti') {
        status = 'multi';
    }
    bg.innerHTML = generateButtonsText(status, nid);
    // set status section
    var ss = qs('#status-'+nid);
    ss.innerHTML = status;
}

function generateButtonsText(status, igid) {
    var actDis = '';
    var actCls = 'primary';
    var banDis = '';
    var banCls = 'secondary';
    var asMultiDis = '';
    var asMultiCls = 'tertiary';
    if (status == 'active') {
        actCls = 'hidden';
        actDis = 'disabled';
    } else if (status == 'banned') {
        banCls = 'hidden';
        banDis = 'disabled';
    } else if (status == 'multi') {
        asMultiCls = 'hidden';
        asMultiDis = 'disabled';
    } else if (status == 'banned_multi') {
        banCls = 'hidden';
        banDis = 'disabled';
    }
    return `
        <button class="`+actCls+`" data-igId="`+igid+`" onclick="doAction(this, 'activate');" `+actDis+`>Activate</button>
        <button class="`+banCls+`" data-igId="`+igid+`" onclick="doAction(this, 'ban');" `+banDis+`>Ban</button>
        <button class="`+asMultiCls+`" data-igId="`+igid+`" onclick="doAction(this, 'asMulti'); `+asMultiDis+`">As Multi</button>`
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
                <a href="`+'https://www.instagram.com/'+card.ig_id+`" target="_blank">
                    @`+card.ig_id+`
                </a>
            </div>
            <div class="section center">
                Status: <span id="status-`+card.ig_id+`">`+card.Status+`</span>
            </div>
            <div class="section button-group" id="btnGroup-`+card.ig_id+`">
                `+generateButtonsText(card.Status, card.ig_id)+`
            </div>
        </div>`;
        cc.appendChild(col);
    });
}