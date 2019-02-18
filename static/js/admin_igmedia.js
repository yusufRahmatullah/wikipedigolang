function doAction(btn, action) {
    var id = btn.getAttribute('data-id');
    var data = {
        id: id,
        action: action,
    };
    apiPost('/api/igmedias/action', data, {
        onSuc: d => {
            toggleAction(id, action);
        },
        onErr: e => {
            alert(e);
        },
    });
}

function toggleAction(nid, act) {
    // set buttons
    var bg = qs('#btnGroup-'+nid);
    var status = 'shown';
    if (act == 'hide') {
        status = 'hidden';
    }
    bg.innerHTML = generateButtonsText(status, nid);
    // set status section
    var ss = qs('#status-'+nid);
    ss.innerHTML = status;
}

function generateButtonsText(status, id) {
    var actDis = '';
    var actCls = 'primary';
    var banDis = '';
    var banCls = 'secondary';
    if (status == 'shown') {
        actCls = 'hidden';
        actDis = 'disabled';
    } else if (status == 'hidden') {
        banCls = 'hidden';
        banDis = 'disabled';
    }
    return `
        <button class="`+actCls+`" data-id="`+id+`" onclick="doAction(this, 'show');" `+actDis+`>Show</button>
        <button class="`+banCls+`" data-id="`+id+`" onclick="doAction(this, 'hide');" `+banDis+`>Hide</button>`
}

function appendCards(cards) {
    var cc = qs('#card-container');
    cards.forEach(card => {                
        var col = document.createElement('div');
        col.setAttribute('class', 'col-sm-4 col-md-3')
        col.innerHTML = `
        <div class="card">
            <img src="`+card.url+`" style="height:100%" onerror='if (this.src != "/static/default.png") this.src = "/static/default.png";'/>
            <div class="section dark center">
                <a href="`+'https://www.instagram.com/'+card.ig_id+`" target="_blank">@`+card.ig_id+`</a>
            </div>
            <div class="section center">
                Status: <span id="status-`+card.id+`">`+card.status+`</span>
            </div>
            <div class="section button-group" id="btnGroup-`+card.id+`">
                `+generateButtonsText(card.status, card.id)+`
            </div>
        </div>`;
        cc.appendChild(col);
    });
}