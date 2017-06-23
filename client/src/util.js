export function toCurrencyString(value, currency, digitsAfterDecimal) {
    const formatter = new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency,
        minimumFractionDigits: digitsAfterDecimal
    });

    return formatter.format(value / Math.pow(10, digitsAfterDecimal));
}

export function toRFC3339(date) {
    return date.toISOString().substring(0, 10);
}

export function fromRFC3339(dateString) {
    return new Date(dateString);
}
