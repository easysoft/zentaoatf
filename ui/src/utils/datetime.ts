import moment from "moment";

export function momentTimeDef(tm) {
    return moment.unix(tm).format("YYYY-MM-DD HH:mm:ss")
}
export function percentDef(numb, total) {
    if (total == 0) return '0%'
    return Number(numb / total * 100).toFixed(2) + '%'
}