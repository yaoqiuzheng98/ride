#!/bin/bash
# 部署脚本：把二进制和 service 文件传到服务器并启动
# 部署前请先修改 deploy/ride-server.service 里的 DB_USER/DB_PASS

set -e

TARGET=root@38.207.185.207

echo ">>> 编译 Linux amd64 二进制..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ride-server main.go

echo ">>> 上传文件到 $TARGET:/tmp/ride-deploy/..."
ssh $TARGET "mkdir -p /tmp/ride-deploy"
scp ride-server $TARGET:/tmp/ride-deploy/
scp deploy/ride-server.service $TARGET:/tmp/ride-deploy/
scp 骑行日记v*.apk $TARGET:/tmp/ride-deploy/
scp -r web $TARGET:/tmp/ride-deploy/

echo ">>> 在服务器上安装..."
ssh $TARGET 'bash -s' <<'REMOTE'
set -e
# 先停服务，避免覆盖正在运行的二进制报 "Text file busy"
systemctl stop ride-server 2>/dev/null || true

# 安装二进制 + apk + web 资源
mkdir -p /opt/ride
cp /tmp/ride-deploy/ride-server /opt/ride/ride-server
cp /tmp/ride-deploy/骑行日记v*.apk /opt/ride/
cp -r /tmp/ride-deploy/web /opt/ride/web
chmod +x /opt/ride/ride-server
chown -R www-data:www-data /opt/ride 2>/dev/null || true

# 安装 systemd service
cp /tmp/ride-deploy/ride-server.service /etc/systemd/system/ride-server.service
systemctl daemon-reload
systemctl enable ride-server
systemctl start ride-server

sleep 2
systemctl status ride-server --no-pager -l | head -15
echo ">>> 部署完成。查看日志: journalctl -u ride-server -f"
REMOTE

BUILD=$(date +%s)
echo ""
echo ">>> 下载链接："
echo "https://hk-vmiss.dokodemo.top/ride/download?build=$BUILD"
