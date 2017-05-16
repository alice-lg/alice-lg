

/*
 * Fetch query params from location
 */
export function queryParams() {
    if (!window && !window.location && !window.location.search) {
        return {}
    }
    let search = window.location.search.slice(1); // omit ?
    let tokens = search.split("&");
    let params = {};
    for (let t of tokens) {
        let kv = t.split("=", 2)
        params[kv[0]] = kv[1];
    }
    return params;
}


