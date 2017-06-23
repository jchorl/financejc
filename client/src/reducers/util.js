export function dateStringToDate(dateString) {
    let date = new Date(dateString);
    return new Date(date.getTime() + date.getTimezoneOffset() * 60000);
}
