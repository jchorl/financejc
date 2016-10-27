export function toCurrency(num, currency) {
	return num.toLocaleString(window.navigator.language, { style: 'currency', currency: currency });
}

export function toDate(date) {
	return date.toLocaleDateString(undefined, {timeZone: "UTC"});
}

export function toDecimal(whole, digitsAfterDecimal) {
  return whole / Math.pow(10, digitsAfterDecimal);
}

export function toWhole(decimal, digitsAfterDecimal) {
  return decimal * Math.pow(10, digitsAfterDecimal);
}
