import moment from "moment";
import {AutoTestTools, TestTools, BuildTools} from "@/utils/const";

function addItems(item, list, map) {
    const lowerCase = item.toLowerCase()
    list.push(lowerCase)
    map[lowerCase] = item
}

export function getUnitTestFrameworks(): any {
    const list = new Array<string>()
    const map = {}
    TestTools.forEach((item) => {
        addItems(item, list, map)
    })

    return {list: list, map: map}
}
export function getUnitTestTools(): any {
    const data = {}
    const map = {}

    Object.keys(BuildTools).forEach((key) => {
        if (! (key in data)) data[key] = []

        BuildTools[key].forEach(item => {
            const lowerCase = item.toLowerCase()
            data[key].push(lowerCase)
            map[lowerCase] = item
        })
    })

    return {data: data, map: map}
}

export function getAutoTestTools(): any {
    const list = new Array<string>()
    const map = {}
    AutoTestTools.forEach((item) => {
        addItems(item, list, map)
    })

    return {list: list, map: map}
}

const execByMap = {
    case: 'by_case',
    module: 'by_module',
    suite: 'by_suite',
    task: 'by_task',
}
export const testToolMap = {
    junit: 'JUnit',
    testng: 'TestNG',
    phpunit: 'PHPUnit',
    pytest: 'PyTest',
    jest: 'Jest',
    cppunit: 'CppUnit',
    gtest: 'GTest',
    qtest: 'QTest',
    goTest: 'GoTest',
    allure: 'Allure',

    robotframework: 'RobotFramework',
    cypress: 'Cypress',

    playwright: 'Playwright',
    puppeteer: 'Puppeteer',

    autoit: 'AutoIt',
    selenium: 'Selenium',
    appium: 'Appium',
}
export function execByDef(record) {
    if (record.execBy) return execByMap[record.execBy] ? execByMap[record.execBy]: '';
    else return testToolMap[record.testTool] ? testToolMap[record.testTool]: '';
}
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
    mac: 'mac',
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
    if (code === true || code === 'pass') {
        return 'pass'
    } else {
        return 'fail'
    }
}

export function expectDesc(str) {
    return str === '' ? 'pass' : str
}
export function actualDesc(str) {
    return str === 'N/A' ? '' : str
}
