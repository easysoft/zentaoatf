#!/usr/bin/env php
<?php
/**

title=use ztf to run selenium test
cid=0
pid=0

1. check webpage title >> 禅道_百度搜索

*/

namespace Facebook\WebDriver;
use Facebook\WebDriver\Remote\RemoteWebDriver;
use Facebook\WebDriver\Remote\DesiredCapabilities;
use Facebook\WebDriver\Chrome\ChromeOptions;
include 'vendor/autoload.php';

/* launch build-in selenium driver to test */
if (isWindows())
{
	$command = 'start ' . dirname(__FILE__, 3) . '\runtime\selenium\chrome80.exe >log.txt 2>&1';
    //	exec("CHCP 936");
} else // for no-windows system, pls download chrome driver from https://chromedriver.storage.googleapis.com/index.html
{
    $command = 'nohup ' . dirname(__FILE__, 3) . '/runtime/selenium/chrome80 >log.txt 2>&1 &';
}
pclose(popen($command, 'r'));
sleep(1);

$host = 'http://127.0.0.1:9515';

$options = new ChromeOptions();
$options->addArguments(['--no-sandbox']); // ['--headless', '--no-sandbox']
$desiredCapabilities = DesiredCapabilities::chrome();
$desiredCapabilities->setCapability(ChromeOptions::CAPABILITY, $options);

$driver = RemoteWebDriver::create($host, $desiredCapabilities);
$driver->get("https://www.baidu.com");
$html= $driver->getPageSource();
// print_r("$html \n");

$keywordsInput = $driver->findElement(WebDriverBy::id('kw'));
$keywordsInput->clear();
$keywordsInput->sendKeys("禅道");

$submitButton = $driver->findElement(WebDriverBy::id('su'));
$submitButton->click();

$driver-> wait(10, 500)-> until(WebDriverExpectedCondition::titleContains('禅道'));

$title = $driver->getTitle();
//if (isWindows()) $title = iconv("UTF-8","GB2312", $title);
print("$title\n");

$driver->close();

if (isWindows())
{
    exec('taskkill /F /im chrome80.exe');
} else
{
    exec('ps -ef | grep chrome80 | grep -v grep | xargs kill -9 2>/dev/null');
}

function  isWindows()
{
    return strtoupper(substr(PHP_OS, 0, 3)) === 'WIN';
}
