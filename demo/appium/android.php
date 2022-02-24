#!/usr/bin/env php
<?php
/**

title=use ztf to run appium test
cid=0
pid=0

1. check image element displayed attribute is >> true

*/

require_once('vendor/autoload.php');

use Facebook\WebDriver\Remote\DesiredCapabilities;
use Facebook\WebDriver\Remote\RemoteWebDriver;
use Facebook\WebDriver\WebDriverBy;

class AndroidTest
{
    protected $webDriver;

    public function exec()
    {
        /* you need to start appium service on below host and port */
        $url = "http://172.16.13.1:4723/wd/hub";

        $capabilities = new DesiredCapabilities();
        $capabilities->setCapability("deviceName", "redmi");
        $capabilities->setCapability("platformName", "Android");

        // use android apk remote url
        $capabilities->setCapability("app", "https://applitools.bintray.com/Examples/eyes-android-hello-world.apk");

        // or use local apk path on host that appium run on
		// $capabilities->setCapability("app", '/Users/aaron/testing/res/eyes-android-hello-world.apk');

        $driver = RemoteWebDriver::create($url, $capabilities);

        $driver->findElement(WebDriverBy::id("random_number_check_box"))->click();
        $driver->findElement(WebDriverBy::id("click_me_btn"))->click();

        $image = $driver->findElement(WebDriverBy::id("image"));
        print($image->getAttribute('displayed') . "\n");

        $driver->quit();
    }
}

$test = new AndroidTest();
$test->exec();
