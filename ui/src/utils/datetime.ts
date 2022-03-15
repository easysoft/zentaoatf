import moment, {utc} from "moment";

export function momentUtcDef(tm) {
    return moment.utc(tm).format("YYYY-MM-DD HH:mm:ss")
}
export function momentUnixDef(tm) {
    return moment.unix(tm).format("YYYY-MM-DD HH:mm:ss")
}

export function percentDef(numb, total) {
    if (total == 0) return '0%'
    return Number(numb / total * 100).toFixed(2) + '%'
}