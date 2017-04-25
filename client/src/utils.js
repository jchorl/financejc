export function toCurrency(num, currency) {
    return num.toLocaleString(window.navigator.language, { style: 'currency', currency: currency });
}

export function toDate(date) {
    return date.toLocaleDateString(undefined, { timeZone: 'UTC' });
}

export function toDecimal(whole, digitsAfterDecimal) {
    return whole / Math.pow(10, digitsAfterDecimal);
}

export function toWhole(decimal, digitsAfterDecimal) {
    return Math.round(decimal * Math.pow(10, digitsAfterDecimal));
}

export function queryByFieldAndVal(accountId, field, val) {
    return fetch(`/api/account/${accountId}/transactions/query?field=${field}&value=${val}`, {
        credentials: 'include'
    })
    .then(response => response.json());
}

export function toRFC3339(d) {
    return d.toJSON().slice(0, 10);
}

// currentDateAsUtc returns the current local date as a UTC date
// e.g. if it is 1/1/2000 locally and 2/1/2000 UTC, this returns
// a date with UTC timezone with date 1/1/2000.
export function currentDateAsUtc() {
    let local = new Date();
    local.setMinutes(local.getMinutes() - local.getTimezoneOffset());
    return local;
}

export function formDateValueToDate(val) {
    return new Date(val);
}
