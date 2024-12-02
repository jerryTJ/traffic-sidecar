package tools

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	"github.com/jerryTJ/sidecar/internal/app"
	pb "github.com/jerryTJ/sidecar/web/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GetServerInfo(serverName string, deployVersion string, addr string, timeout time.Duration) app.ServerInfo {

	// 加载客户端证书和私钥
	clientCert, err := tls.LoadX509KeyPair("configs/ssl/client.crt", "configs/ssl/client.key")
	if err != nil {
		log.Printf("Failed to load client certificate: %v", err)
	}

	// 加载 CA 证书
	caCert, err := os.ReadFile("configs/ssl/ca.crt")
	if err != nil {
		log.Printf("Failed to read CA certificate: %v", err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		log.Printf("Failed to append CA certificate")
	}

	// 配置 TLS
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert}, // 客户端证书
		RootCAs:            certPool,                      // 服务端 CA
		InsecureSkipVerify: false,                         // 验证服务端证书
	}

	// 创建 TransportCredentials
	creds := credentials.NewTLS(tlsConfig)

	// 连接 gRPC 服务端
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// 创建服务客户端
	client := pb.NewServerInfoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	resp, err := client.GetColoringInfo(ctx, &pb.ServerRequest{Name: serverName, Version: deployVersion})
	if err != nil {
		log.Printf("Could not greet: %v", err)
	}

	return app.ServerInfo{Name: resp.Name, ChainID: resp.Chain, Color: resp.Color, Version: resp.Version}
}
