# 创建需要用到的目录
mkdir -p /root/tiktok/video
mkdir -p /root/tiktok/pic

# 安装minio
docker run -p 9000:9000 -p 9090:9090 \
     --net=host \
     --name minio \
     -d --restart=always \
     -e "MINIO_ACCESS_KEY=admin123" \
     -e "MINIO_SECRET_KEY=admin123" \
     minio/minio server \
     /data --console-address ":9090" -address ":9000"


# 安装mysql
docker run --name camps_mysql -e MYSQL_ROOT_PASSWORD=123456 -d -e MYSQL_DATABASE=camps_tiktok -p 8086:3306 mysql:8.0

# 安装consul
docker run --name consul -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp consul consul agent -dev -client=0.0.0.0

# 安装redis
. redis.sh