name="tiktok_user_svr"
docker build -t zhugeqing/tiktok_user_svr:latest ../server/usersvr
docker run -d --name $name -v /etc/localtime:/etc/localtime -p 8000:8000 zhugeqing/tiktok_user_svr:latest

name="tiktok_relation_svr"
docker build -t zhugeqing/tiktok_relation_svr:latest ../server/relationsvr
docker run -d --name $name -v /etc/localtime:/etc/localtime -p 8001:8001 zhugeqing/tiktok_relation_svr:latest
docker system prune

name="tiktok_comment_svr"
docker build -t zhugeqing/tiktok_comment_svr:latest ../server/commentsvr
docker run -d --name $name -v /etc/localtime:/etc/localtime -p 8002:8002 zhugeqing/tiktok_comment_svr:latest

name="tiktok_favorite_svr"
docker build -t zhugeqing/tiktok_favorite_svr:latest ../server/favoritesvr
docker run -d --name $name -v /etc/localtime:/etc/localtime -p 8003:8003 zhugeqing/tiktok_favorite_svr:latest

name="tiktok_video_svr"
docker build -t zhugeqing/tiktok_video_svr:latest ../server/videosvr
docker run -d --name $name -v /etc/localtime:/etc/localtime -v tiktok-video:/root/tiktok/video -v tiktok-pic:/root/tiktok/pic -p 8004:8004 zhugeqing/tiktok_video_svr:latest

name="tiktok_gateway_svr"
docker build -t zhugeqing/tiktok_gateway_svr:latest ../server/gatewaysvr
docker run -d --name $name -v /etc/localtime:/etc/localtime -v tiktok-video:/root/tiktok/video -p 8005:8005 zhugeqing/tiktok_gateway_svr:latest
docker system prune
