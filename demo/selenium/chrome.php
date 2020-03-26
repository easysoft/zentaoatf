<?php
namespace Facebook\WebDriver;
use Facebook\WebDriver\Remote\RemoteWebDriver;
use Facebook\WebDriver\Remote\DesiredCapabilities;
use Facebook\WebDriver\Chrome\ChromeOptions;
include 'vendor/autoload.php';
if (isWindows()) {
    exec("CHCP 936");

	$command = 'start /B runtime\selenium\chrome80.exe >log.txt 2>&1';
	pclose(popen($command, 'w'));
}

$host = 'http://127.0.0.1:9515';

$options = new ChromeOptions();
$options->addArguments(['--no-sandbox']); // ['--headless', '--no-sandbox']
$desiredCapabilities = DesiredCapabilities::chrome();
$desiredCapabilities->setCapability(ChromeOptions::CAPABILITY, $options);

$driver = RemoteWebDriver::create($host, $desiredCapabilities);
$driver->get("https://www.baidu.com");
$html= $driver->getPageSource();
print_r("$html \n");

$keywordsInput = $driver->findElement(WebDriverBy::id('kw'));
$keywordsInput->clear();
$keywordsInput->sendKeys("禅道");

$submitButton = $driver->findElement(WebDriverBy::id('su'));
$submitButton->click();

$driver-> wait(10,500)-> until(WebDriverExpectedCondition::titleContains('禅道'));
$title = $driver->getTitle();
if (isWindows()) $title = iconv("UTF-8","GB2312", $title);
print(">> $title\n");

$driver->close();

function  isWindows() {
    return strtoupper(substr(PHP_OS,0,3))==='WIN';
}