package main
 
import (
	"context"
	"fmt"
	pb "grpcdemo/protos"
	"net"
 
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)
 
func main() {
 
	GPRCServer()
 
	// http
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
 
	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})
	app.Run(iris.Addr(":8080"))
}
 
func GPRCServer() {
	// 监听本地端口
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		return
	}
	s := grpc.NewServer()                  // 创建GRPC
	pb.RegisterGreeterServer(s, &server{}) // 在GRPC服务端注册服务
 
	reflection.Register(s)
	fmt.Println("grpc serve 9090")
	err = s.Serve(listener)
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to serve: %v", err))
	}
 
}
 
type server struct{}
 
func NewServer() *server {
	return &server{}
}
 
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	msg := "Resuest By:" + in.Name + " Response By :" + LocalIp()
	fmt.Println("GRPC Send: ", msg)
	return &pb.HelloReply{Message: msg}, nil
}
 
func LocalIp() string {
	addrs, _ := net.InterfaceAddrs()
	var ip string = "localhost"
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return ip
}
