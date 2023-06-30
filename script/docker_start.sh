name="usersvr"
docker rm -f $name
docker build -t zhugeqing/go-geek:v2 .
docker run -d --name $name -v /etc/localtime:/etc/localtime -v /geek:/geek -p 10086:10086 zhugeqing/go-geek:v2
docker system prune

name="commentsvr"
docker rm -f $name
docker build -t zhugeqing/go-geek:v2 .
docker run -d --name $name -v /etc/localtime:/etc/localtime -v /geek:/geek -p 10086:10086 zhugeqing/go-geek:v2
docker system prune


name=""
docker rm -f $name
docker build -t zhugeqing/go-geek:v2 .
docker run -d --name $name -v /etc/localtime:/etc/localtime -v /geek:/geek -p 10086:10086 zhugeqing/go-geek:v2
docker system prune