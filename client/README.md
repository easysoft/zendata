# ZD 客户端

## 开发

### 安装依赖

执行：

```
npm install
```

### 启动调试模式

启动默认的调试模式，此时会预先自动启动 UI 开发服务器和 ZD 开发服务器。

```
npm run start
```

参考：https://www.electronforge.io/cli#start

### 环境变量

在调试模式下可以通过制定环境变量来设置 ZD UI 服务访问地址和 ZD 可执行程序位置。

* `UI_SERVER_URL`：ZD UI 服务访问地址或静态资源文件目录，如果不指定会自动进入 `../ui/` 目录执行 `npm run serve` 获取开发服务器地址；
* `UI_SERVER_PORT`：ZD UI 服务端口，如果不指定则使用 `8000`；
* `SERVER_EXE_PATH`：ZD 服务可执行文件位置（相对于 `/client` 目录），如果不指定则会自动进入 `../cmd/server/` 目录下执行 `go run main.go -p 8085` 启动 ZD 服务器。
* `SERVER_CWD_PATH`：ZD 服务运行目录；
* `SKIP_SERVER`：跳过启动 ZD 服务，适用于 ZD 服务已经在外部启动的情况。

### 特殊调试模式

**模式一：使用外部 UI 服务**

```
SKIP_SERVER=1 UI_SERVER_URL=http://localhost:8000 npm run start
```

**模式二：使用本地 UI 静态文件目录**

```
UI_SERVER_URL=../ui/dist UI_SERVER_PORT=8000 npm run start
```

**模式三：自定义 ZD 服务执行文件路径**

```
SERVER_EXE_PATH=bin/darwin/zd npm run start
```

**模式四：跳过启动 ZD 服务，使用外部 ZD 服务**

```
SKIP_SERVER=1 npm run start
```

**模式五：综合使用外部 UI 服务和外部 ZD 服务**

```
UI_SERVER_URL=http://localhost:8000 SKIP_SERVER=1 npm run start
```

## 代码检查

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

打包之前确保如下目录有相应的资源文件：

* `ui/dist/`：包含 UI 服务相关的所有文件
* `client/bin/win32/`：包含适用于 Windows 的 ZD 程序；
* `client/bin/darwin/`：包含适用于 macOS 的 ZD 程序；
* `client/bin/linux/`：包含适用于 Linux 的 ZD 程序；

## 发布

```
npm run publish
```

参考：https://www.electronforge.io/cli#publish
