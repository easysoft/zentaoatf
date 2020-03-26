<?php
namespace Facebook\WebDriver;
exec("CHCP 936");
use Facebook\WebDriver\Remote\RemoteWebDriver;
use Facebook\WebDriver\Remote\DesiredCapabilities;
use Facebook\WebDriver\Chrome\ChromeOptions;
include 'vendor\autoload.php';

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
$title = iconv("UTF-8","GB2312",$driver->getTitle()); 
print(">> $title\n");

$driver->close();
