function appendCards(cards) {
    var cc = qs('#card-container');
    cards.forEach(card => {
        var col = document.createElement('div');
        col.setAttribute('class', 'col-sm-4 col-md-3')
        col.innerHTML = `
        <div class="card" onclick="showModal('`+card.url+`');">
            <img src="`+ card.url + `" style="height:100%" onerror='if (this.src != "/static/default.png") this.src = "/static/default.png";'/>
        </div>`;
        cc.appendChild(col);
    });
}