#-*- coding: utf-8
#!/usr/bin/python

'''

title=use ztf to run selenium test
cid=1
pid=1

1. check webpage title >> 禅道_百度搜索

'''

from selenium import webdriver
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
import sys,io,platform

## for Chinese display
if(platform.system()=='Windows'):
   import sys,io
   sys.stdout = io.TextIOWrapper(sys.stdout.buffer,encoding='utf8')

option = webdriver.ChromeOptions()
# option.add_argument('--headless')
driver = webdriver.Chrome(options=option)
driver.get('https://www.baidu.com/')

keywordsInput = driver.find_element_by_id("kw")
keywordsInput.clear()
keywordsInput.send_keys("禅道")
submitButton = driver.find_element_by_id("su")
submitButton.click()

element = WebDriverWait(driver, 5, 0.5).until(EC.title_contains("禅道"),message='超时啦!')

print(driver.title)

driver.quit()