#如何在任意目录中执行ztf命令？
## 方法一，将ztf加入系统环境变量中：
### Windows系统
1. 在命令行输入sysdm.cpl，打开系统属性窗口；
2. 依次点击"高级"标签、"环境变量"按钮，打开环境变量编辑窗口；
3. 在上部"用户变量"列表中，点击"编辑"按钮修改Path变量；若无Path变量，则点击"新建"按钮；
4. 填入或追加ztf.exe文件所在的目录绝对路径，Win7中为单行编辑模式，路径间用英文分号隔开；
5. 重新打开命令行窗口，使设置生效。

### Linux/Mac系统
1. 编辑用户目录下的.bash_profile文件；
2. 在文件末尾添加export PATH=$PATH:<ztf文件所在目录绝对路径>；
3. 执行source ~/.bash_profile使设置生效。

## 方法二，拷贝ztf到PATH变量所在的目录：
### Windows系统
手动拷贝至ztf.exe命令至C:\Windows，或
以管理员身份打开cmd行，执行copy /Y ztf.exe C:\Windows

### Linux/Mac系统
在命令行执行cp -f ztf /usr/local/bin/