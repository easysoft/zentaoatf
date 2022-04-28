
export const ScriptFileNotExist = 'script_file_not_exist'

export const ScriptLanguages = [
    'Python',
    'PHP',
    'Ruby',
    'JavaScript',
    'Lua',
    'TCL',
    'Perl',
    'AutoIt'
]

export const TestTools = [
    'JUnit',
    'TestNG',
    'PHPUnit',
    'PyTest',
    'Jest',
    'CppUnit',
    'GTest',
    'QTest',

    'RobotFramework',
    'Cypress',
]
export const BuildTools = {
    junit: ['Maven'],
    testng: ['Maven'],
}

export const AutoTestTools = [
    'AutoIt',
    'Selenium',
    'Appium',
    'RobotFramework',
    'Cypress',
]

export const MonacoOptions = {
    colorDecorators: true,
    lineHeight: 24,
    tabSize: 2,
    autoIndent: true,
    formatOnPaste: true,
    formatOnType: true
}

export const disableStatusMap = [
    {
        label: 'all',
        value: '',
    },
    {
        label: 'enable',
        value: '1',
    },
    {
        label: 'disable',
        value: '0',
    },
]

export const ztfTestTypesDef = [{
        label: 'ZTF',
        value: 'ztf',
    },
]

export const unitTestTypesDef = [
    {
        label: 'JUnit',
        value: 'junit',
        cmd: 'ztf.exe junit -p product_id mvn clean package test',
    },
    {
        label: 'TestNG',
        value: 'testng',
        cmd: 'ztf.exe testng -p product_id mvn clean package test',
    },
    {
        label: 'PHPUnit',
        value: 'phpunit',
        cmd: 'ztf.exe phpunit -p product_id phpunit -c phpunit.xml --testdox-xml testresults.xml',
    },
    {
        label: 'PyTest',
        value: 'pytest',
        cmd: 'ztf.exe pytest -p product_id pytest --junitxml=testresult.xml',
    },
    {
        label: 'Jest',
        value: 'jest',
        cmd: 'ztf.exe jest -p product_id  ./node_modules/.bin/jest test',
    },
    {
        label: 'GTest',
        value: 'gtest',
        cmd: 'ztf  gtest -p product_id ./build/bin/unit_tests --gtest_output=xml:testresults.xml',
    },
    {
        label: 'QTest',
        value: 'qtest',
        cmd: 'ztf qtest -p product_id ./tests/QTestLibCocktailsTest -xunitxml -o testresults.xml',
    },
    {
        label: 'CppUnit',
        value: 'cppunit',
        cmd: ' ztf cppunit -p product_id ./MathTest',
    },
]

export const logLevelMap = {
    output: {code: 100},
    run: {code: 200},
    result: {code: 300},
    error: {code: 400},
}