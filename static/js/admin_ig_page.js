function initFilterSearch() {
    var status = qs('#filter-status').value;
    optSearchParam = "&filterStatus="+status;
    initSearch();
}
