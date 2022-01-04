# ZTF 客户端

## 开发

### 安装依赖

执行：

```
npm install
```

### 启动调试模式

```
npm run start
```

参考：https://www.electronforge.io/cli#start

### 环境变量

在调试模式下可以通过制定环境变量来设置 ZTF UI 服务访问地址和 ZTF 可执行程序位置。

* `UI_SERVER_PATH`: ZTF UI 服务访问地址，如果不指定会自动进入 `../ui/` 目录执行 `npm run serve` 获取开发服务器地址；
* `SERVER_EXE_PATH`：ZTF 服务可执行文件位置，如果不指定则会自动进入 `../cmd/server/` 目录下执行 `go run main.go -p 8085` 启动 ZTF 服务器。

### 代码检查

```
npm run lint
```

## 构建

```
npm run make
```

参考：https://www.electronforge.io/cli#make

## 打包

```
npm run package
```

参考：https://www.electronforge.io/cli#package

## 发布

```
npm run publish
```

参考：https://www.electronforge.io/cli#publish
