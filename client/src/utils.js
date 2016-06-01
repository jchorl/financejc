export function toCurrency(num, currency) {
	return num.toLocaleString(window.navigator.language, { style: 'currency', currency: currency });
}

export function toDate(date) {
	let d = new Date(date);
	return d.toLocaleDateString();
}
