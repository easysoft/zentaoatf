import moment, {utc} from "moment";

export function momentUtcDef(tm) {
    return moment.parseZone(tm).format("YYYY-MM-DD HH:mm:ss")
}
export function momentUnixDef(tm) {
    return moment.unix(tm).format("YYYY-MM-DD HH:mm:ss")
}
export function momentTime(tm) {
    return moment.parseZone(tm).format("HH:mm:ss")
}

export function momentUnixDefFormat(tm, format) {
    console.log(111, tm, format)
    return moment.unix(tm).format(format)
}

export function percentDef(numb, total) {
    if (total == 0) return '0%'
    return Number(numb / total * 100).toFixed(2) + '%'
}