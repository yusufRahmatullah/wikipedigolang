function showModal(url) {
    qs('#media').style.background = 'url("'+url+'")';
    qs('#media').style.backgroundSize = 'contain';
    qs('#media').style.backgroundRepeat = 'no-repeat';
    qs('#media').style.backgroundPosition = 'center';
    qs('#mediaOverlay').style.display = 'block';
}

function closeModal() {
    qs('#media').style.background = '';
    qs('#mediaOverlay').style.display = 'none';
}