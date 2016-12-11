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

function pad(n) {
    return n<10 ? '0'+n : n;
}

export function queryByFieldAndVal(accountId, field, val) {
    return fetch(`/api/account/${accountId}/transactions/query?field=${field}&value=${val}`, {
        credentials: 'include'
    })
    .then(response => response.json());
}

export function toRFC3339(d) {
    return d.getUTCFullYear() + '-'
    + pad(d.getUTCMonth() + 1) + '-'
    + pad(d.getUTCDate());
}
