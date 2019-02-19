function apiGet(url, {onSuc, onErr, before, after}) {
    fetch(url, {
        credentials: 'same-origin'
    })
    .then(res => res.json())
    .then(json => {
        if (before) {
            before(json);
        }
        if (json.status == 'OK') {
            if (onSuc) {
                onSuc(json.data);
            }
        } else {
            if (onErr) {
                onErr(json.message);
            }
        }
        if (after) {
            after(json);
        }
    })
    .catch(err => console.error(err));
}

function apiPost(url, data, {onSuc, onErr, before, after}) {
    fetch(url, {
        method: 'POST',
        body: JSON.stringify(data),
        credentials: 'same-origin',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(res => res.json())
    .then(json => {
        if (before) {
            before(json);
        }
        if (json.status == 'OK') {
            if (onSuc) {
                onSuc(json.data);
            }
        } else {
            if (onErr) {
                onErr(json.message);
            }
        }
        if (after) {
            after(json);
        }
    })
    .catch(err => console.error(err));
}

qs = q => {
    if (q[0] == '.') {
        return document.getElementsByClassName(q.slice(1))[0];
    }
    return document.querySelector(q);
}