function appendCards(cards) {
    var cc = qs('#card-container');
    cards.forEach(card => {                
        var col = document.createElement('div');
        col.setAttribute('class', 'col-sm-4 col-md-3')
        col.innerHTML = `
        <div class="card">
            <img src="`+card.url+`" style="height:100%" onerror='if (this.src != "/static/default.png") this.src = "/static/default.png";'/>
            <div class="section dark center">
                <a href="https://www.instagram.com/`+card.ig_id+`" target="_blank">
                    @`+card.ig_id+`
                </a>
            </div>
        </div>`;
        cc.appendChild(col);
    });
}