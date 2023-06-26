#!/bin/bash

# 入口函数
function Init() {


        GetVar
        WriteToConf

        docker run -d --network host --restart unless-stopped --name $name -v /docker/redis/$name/conf:/etc/redis/  -v $name-data:/data zhugeqing/redis redis-server /etc/redis/redis.conf
        # 软链接
        ln -s /var/lib/docker/volumes/$name-data/_data /docker/redis/$name/data
}

# 获取变量
function GetVar() {
        echo "正在使用redLock节点脚本！"
        read -t 20 -p "请输入容器映射到主机的端口（容器内默认为6379，主机默认为6379）：" port 
        port=${port:-6379}
        echo "主机端口为${port}"

        read -t 20 -p "请输入容器运行时的容器名字（默认为redis）：" name
        name=${name:-redis}
        echo "容器名字为${name}"

        read -t 20 -p "请输入主节点密码（默认为123456）：" password
        password=${password:-123456}
        echo "主节点密码为${password}"

        # 创建需要映射的目录
        mkdir -p /docker/redis/$name/conf /docker/redis/$name/data
}


# 编写配置文件
function WriteToConf() {
cat > /docker/redis/$name/conf/redis.conf << EOF

# 容器内部端口
port $port
# daemonize yes 

logfile "redis-$port.log" 
dbfilename "dump-redis-$port.rdb"

# 访问redis-server密码
requirepass $password

# 主节点密码（也需要配置，不然当master无法成为新的master节点的slave节点）
masterauth $password

# 设置redis最大内存
maxmemory 500MB
# 设置redis内存淘汰策略
maxmemory-policy volatile-lru

# 开启AOF
appendonly yes
# 总是追加到AOF文件
appendfsync always
# 超过100MB就进行重写
auto-aof-rewrite-min-size 100mb
# 超过增长率就进行重写
auto-aof-rewrite-percentage 40
# 开启AOF-RDB 混合持久化
aof-use-rdb-preamble yes

# 配置慢查询日志（不超过100微秒就不会被记录）
slowlog-log-slower-than 100
# 最多记录100条，先进先出
slowlog-max-len 100
EOF
}

Init
