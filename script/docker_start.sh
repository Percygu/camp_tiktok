name="tiktok_user_svr"
docker rm -f $name
docker build -t zhugeqing/tiktok_user_svr ../server/usersvr
docker run -d --name $name -v /etc/localtime:/etc/localtime -p 8000:8000 zhugeqing/tiktok_user_svr

name="tiktok_relation_svr"
docker rm -f $name
docker build -t zhugeqing/tiktok_relation_svr ../server/relationsvr
docker run -d --name $name -v /etc/localtime:/etc/localtime -p 8001:8001 zhugeqing/tiktok_relation_svr
docker system prune

name="tiktok_comment_svr"
docker rm -f $name
docker build -t zhugeqing/tiktok_comment_svr ../server/commentsvr
docker run -d --name $name -v /etc/localtime:/etc/localtime -p 8002:8002 zhugeqing/tiktok_comment_svr

name="tiktok_favorite_svr"
docker rm -f $name
docker build -t zhugeqing/tiktok_favorite_svr ../server/favoritesvr
docker run -d --name $name -v /etc/localtime:/etc/localtime -p 8003:8003 zhugeqing/tiktok_favorite_svr

name="tiktok_video_svr"
docker rm -f $name
docker build -t zhugeqing/tiktok_video_svr ../server/videosvr
docker run -d --name $name -v /etc/localtime:/etc/localtime -v tiktok-video:/root/tiktok/video -v tiktok-pic:/root/tiktok/pic -p 8004:8004 zhugeqing/tiktok_video_svr

name="tiktok_gateway_svr"
docker rm -f $name
docker build -t zhugeqing/gatewaysvr ../server/gatewaysvr
docker run -d --name $name -v /etc/localtime:/etc/localtime -v tiktok-video:/root/tiktok/video -p 8005:8005 zhugeqing/tiktok_gateway_svr
docker system prune
