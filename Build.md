# 构建指南

推荐在`Ubuntu 18.04.6`下构建，其他版本未测试

## 安装依赖

### 安装node

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
nvm install v14.21.3 --default
npm install -g yarn
```

### 安装Go

```bash
wget https://golang.google.cn/dl/go1.21.5.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH="$HOME/go/bin:/usr/local/go/bin:$PATH"
export GOPROXY='https://goproxy.cn,direct'
```

### 安装系统依赖

```bash
apt install -y zip unzip sqlite3 wine64
``` 

### 安装rsrc

```bash
wget https://github.com/akavel/rsrc/releases/download/v0.10.2/rsrc_linux_amd64
chmod +x rsrc_linux_amd64
mv rsrc_linux_amd64 /usr/bin/rsrc
```

## 构建

```bash
make default
```
