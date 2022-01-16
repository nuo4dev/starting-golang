package main

/**
接收客户端 request，并将 request 中带的 header 写入 response header
读取当前系统的环境变量中的 VERSION 配置，并写入 response header
Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
当访问 localhost/healthz 时，应返回 200
*/

import (
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/golang/glog"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server...")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
}

func ClientIP(r *http.Request) string {   
	xForwardedFor := r.Header.Get("X-Forwarded-For")   
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])   
	if ip != "" {
		return ip   
	}   
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))   
	if ip != "" {
		return ip   
	}   
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip   
	}   
	return ""
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		// io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
		w.Header().Add(k, strings.Join(v, ","))
	}
	os.Setenv("VERSION", "1.0.1")
	w.Header().Add("VERSION", os.Getenv("VERSION"))
	io.WriteString(w, "response completed\n")
	// get clientIP
	clientIp := ClientIP(r)
	log.Printf("Success! Response code: %d", 200)
	log.Printf("Success! clientip: %s", clientIp)
}
