function appendCards(cards) {
    var cc = qs('#card-container');
    cards.forEach(card => {
        var col = document.createElement('div');
        col.setAttribute('class', 'col-sm-4 col-md-3');
        col.innerHTML = `
        <div class="card">
            <img src="`+card.pp_url+`" style="height:100%" onerror='if (this.src != "/static/default.png") this.src = "/static/default.png";'/>
            <div class="section double-padded center">
                `+card.name+`
            </div>
            <a href="/igprofile/`+card.ig_id+`" target="_blank">
                <div class="section dark center">
                    @`+card.ig_id+`
                </div>
            </a>
        </div>`;
        cc.appendChild(col);
    });
}