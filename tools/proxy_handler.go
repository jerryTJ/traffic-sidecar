package tools

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jerryTJ/sidecar/init/logger"
)

type ProxyHandler struct {
	TargetUrl string
	Addr      string
}

// ServeHTTP方法，绑定DefaultHandler
func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Info("target url " + ph.TargetUrl)
	// 解析目标URL
	targetURL, err := url.Parse(ph.TargetUrl)
	if err != nil {
		http.Error(w, "Invalid target URL:", http.StatusBadRequest)
		return
	}
	// 创建一个新的请求，复制原始请求的数据
	proxyReq, err := http.NewRequest(r.Method, targetURL.ResolveReference(r.URL).String(), r.Body)
	logger.Info("request url :" + targetURL.ResolveReference(r.URL).String())
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	// 查询路由配置信息
	logger.Info("source host :" + r.Host)
	serverName := strings.Split(r.Host, ":")[0]
	logger.Info("quest domain:" + serverName)
	// env from deployment
	deployVersion := os.Getenv("DEPLOY_VERDION")
	serverInfo := GetServerInfo(serverName, deployVersion, ph.Addr, 60)
	logger.Info("data domain:" + serverInfo.Domain)
	// 复制请求头
	proxyReq.Header = r.Header
	proxyReq.Header.Set("x-color", serverInfo.Color)
	proxyReq.Header.Set("x-chain", serverInfo.ChainID)
	proxyReq.Header.Set("x-version", serverInfo.Version)
	logger.Info("x-color:" + serverInfo.Color + ",x-chain" + serverInfo.ChainID + ",x-version:" + serverInfo.Version)

	// 发送代理请求
	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		http.Error(w, "Failed to reach target server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 将响应头写回给客户端
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 设置响应状态码
	w.WriteHeader(resp.StatusCode)
	logger.Info("end proxy status code" + resp.Status)
	// 将响应体写回给客户端
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println("Failed to copy response body:", err)
	}
}
