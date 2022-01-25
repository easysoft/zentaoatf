import moment from "moment";

export function momentTimeDef(tm) {
    return moment.unix(tm).format("YYYY-MM-DD HH:mm:ss")
}
export function percentDef(numb, total) {
    if (total == 0) return '0%'
    return Number(numb / total * 100).toFixed(2) + '%'
}

const osMap = {
    windows: 'Windows',
    linux: 'Linux',
    mac: 'Mac',
}
export function testEnvDef(code) {
    return osMap[code]
}
const testTypeMap = {
    func: 'Functional Testing',
    unit: 'Unit Testing',
    auto: 'Automated Testing',
}
export function testTypeDef(code) {
    return testTypeMap[code]
}
export function resultStatusDef(code) {
    if (code === true || code === 'pass')
    return '通过'

    return '失败'
}
