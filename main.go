package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func getLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// 跳过不可用的接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 检查是否是IPv4地址并且是私有IP
			if ip != nil && ip.To4() != nil && isPrivateIP(ip) {
				return ip.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no private IP address found")
}

func isPrivateIP(ip net.IP) bool {
	privateBlocks := []*net.IPNet{
		// 10.0.0.0/8
		{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		// 172.16.0.0/12
		{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		// 192.168.0.0/16
		{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
	}

	for _, block := range privateBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

func main() {
	localIP, err := getLocalIP()
	if err != nil {
		fmt.Println("Error getting local IP address:", err)
		return
	}

	// 首页
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "." + r.URL.Path

		// 检查路径是否为文件夹
		info, err := os.Stat(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if info.IsDir() {
			// 如果是文件夹，则列出文件夹中的文件
			files, err := os.ReadDir(path)
			if err != nil {
				http.Error(w, "Unable to read directory", http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "<h1>Index of %s</h1><ul>", r.URL.Path)
			for _, file := range files {
				name := file.Name()
				fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>", filepath.Join(r.URL.Path, name), name)
			}
			fmt.Fprint(w, "</ul>")
		} else {
			// 如果是文件，则直接显示文件内容
			http.ServeFile(w, r, path)
		}
	})

	// about关于页面
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<html>
				<head>
					<title>About This Server</title>
				</head>
				<body>
					<h1>About This Server</h1>
					<p>This is a simple file server written in Go. It serves files from the current directory and allows users to browse and view text files via a web browser.</p>
					<h2>Features:</h2>
					<ul>
						<li>Lists files and directories in the current directory.</li>
						<li>Allows users to view the content of text files directly in the browser.</li>
						<li>Accessible via local network IP as well as localhost.</li>
					</ul>
					<h2>How to Access:</h2>
					<p>You can access the server using the following URLs:</p>
					<ul>
						<li><a href="http://127.0.0.1">http://127.0.0.1</a> (Localhost)</li>
						<li><a href="http://%s">http://%s</a> (LAN IP)</li>
					</ul>
					<p>Replace <code>LAN IP</code> with the actual IP address provided above.</p>
				</body>
			</html>
		`, localIP, localIP)
	})

	// 监听所有接口地址
	addr := "0.0.0.0:80"
	fmt.Printf("Server started at http://%s (local) and http://%s (LAN)\n", "127.0.0.1", localIP)
	http.ListenAndServe(addr, nil)
}
