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
    // set card-img opacity
    var cardImg = qs('#card-img-'+nid);
    if (act == 'show') {
        cardImg.style.opacity = '1.0';
    } else {
        cardImg.style.opacity = '0.4';
    }
    // set card status
    if (act == 'show') {
        qs('#status-'+nid).innerHTML = 'shown';        
    } else {
        qs('#status-'+nid).innerHTML = 'hidden';
    }
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

function prepareModal(div) {
    var id = div.getAttribute('data-id');
    var status = qs('#status-'+id).innerHTML;
    var m = qs('#media');
    m.innerHTML = `
    <div id="btnGroup-`+id+`">
        `+generateButtonsText(status, id)+`
    </div>
    `;
}

function appendCards(cards) {
    var cc = qs('#card-container');
    cards.forEach(card => {
        var col = document.createElement('div');
        col.setAttribute('class', 'col-sm-4 col-md-3')
        var opacity = card.status == 'shown' ? '1.0' : '0.4';
        col.innerHTML = `
        <div class="hidden" id="status-`+card.id+`">`+card.status+`</div>
        <div class="card" data-id="`+card.id+`" id="card-`+card.id+`" onclick="prepareModal(this); showModal('`+card.url+`');" style="background: black;">
            <img id="card-img-`+card.id+`" src="`+ card.url + `" style="height:100%; opacity:`+opacity+`;" onerror='if (this.src != "/static/default.png") this.src = "/static/default.png";'/>
        </div>`;
        cc.appendChild(col);
    });
}