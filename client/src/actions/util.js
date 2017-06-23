export function handleErrors(resp) {
    if (resp.ok) {
        return resp;
    }

    throw new Error(resp.statusText);
}
