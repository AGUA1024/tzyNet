package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/crypto/ssh"
	"tzyNet/tCommon"
	"log"
	"os"
	"time"
)

// 正式服节点配置
var mpNodeTagToHost = map[string]string{}

// port
var nodePort = "8000"

// 镜像版本
<<<<<<< HEAD
var imageVersion = "v1.0.0.2"
var preVersion = "v1.0.0.1" //老版本名，用于删除老版本容器

// 主机身份
var hostUserName = "root"
var hostPassWord = "123456"

// etcd
var etcdHost = "127.0.0.1:2381"
=======
var imageVersion = "v1.0.0.5"
var preVersion = "v1.0.0.4" //老版本名，用于删除老版本容器

// 主机身份
var hostUserName = "root"
var hostPassWord = "5$vu7X9dCj&Zun9e"

// etcd
var etcdHost = "114.132.213.154:2381"
>>>>>>> origin/main

func main() {
	// 销毁旧版本
	destroyContainer(mpNodeTagToHost)
	// 部署新版本
	deployImage(mpNodeTagToHost)
}

func init() {
	mpEtcdRegPerList, ok := tCommon.GetYamlMapCfg("serverCfg", "server", "release").(map[string]any)
	if !ok {
		tCommon.Logger.SystemErrorLog("ERR_SERVER_RELEASE_CFG1")
	}

	for k, v := range mpEtcdRegPerList {
		var strV string
		strV, ok = v.(string)
		if !ok {
			tCommon.Logger.SystemErrorLog("ERR_SERVER_RELEASE_CFG2")
		}
		mpNodeTagToHost[k] = strV
	}
}

func deployImage(arrNode map[string]string) {
<<<<<<< HEAD
=======
	fmt.Println("docker run -d -p " + nodePort + ":80 -e  NODEMARK=" + "nodeTag" + " --name " + imageVersion + " agua1024/hdyx:" + imageVersion)
>>>>>>> origin/main
	for nodeTag, nodeHost := range arrNode {
		fmt.Printf("%s[%s]: %s\n", nodeTag, nodeHost, "Docker image updateing .....")

		// 配置SSH客户端
		config := &ssh.ClientConfig{
			User: hostUserName,
			Auth: []ssh.AuthMethod{
				ssh.Password(hostPassWord),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略主机密钥验证
		}

		// 建立SSH连接
		client, err := ssh.Dial("tcp", nodeHost+":22", config)
		if err != nil {
			log.Fatalf("Failed to dial: %s", err)
		}
		defer client.Close()

		// 创建SSH会话
		session, err := client.NewSession()
		if err != nil {
			log.Fatalf("Failed to create session: %s", err)
		}
		defer session.Close()

		// 设置会话输出
		session.Stdout = os.Stdout

		// 执行shell命令
		err = session.Run("docker run -d -p " + nodePort + ":80 -e  NODEMARK=" + nodeTag + " --name " + imageVersion + " agua1024/hdyx:" + imageVersion)
		if err != nil {
<<<<<<< HEAD
			log.Fatalf("Failed to run command: %s", err)
=======
			fmt.Println("Failed to run command: %s", err)
>>>>>>> origin/main
		}

		// 端口注册
		etcdCli, _ := clientv3.New(clientv3.Config{
			Endpoints:   []string{etcdHost},
			DialTimeout: 5 * time.Second,
		})

		// 注册服务
		key := "sevPort"
		_, err = etcdCli.Put(context.Background(), key, nodePort)
		if err != nil {
			log.Fatalf("Failed to put etcd kv: %s", err)
		}

		fmt.Println(nodeTag, "Docker update successfully!\n")
	}
}

func destroyContainer(arrNode map[string]string) {
	for nodeTag, nodeHost := range arrNode {
		fmt.Printf("%s[%s]: %s\n", nodeTag, nodeHost, "Docker destroyContainer .....")

		// 配置SSH客户端
		config := &ssh.ClientConfig{
			User: hostUserName,
			Auth: []ssh.AuthMethod{
				ssh.Password(hostPassWord),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略主机密钥验证
		}

		// 建立SSH连接
		client, err := ssh.Dial("tcp", nodeHost+":22", config)
		if err != nil {
			log.Fatalf("Failed to dial: %s", err)
		}
		defer client.Close()

		// 创建SSH会话
		session, err := client.NewSession()
		if err != nil {
			log.Fatalf("Failed to create session: %s", err)
		}
		defer session.Close()

		// 设置会话输出
		session.Stdout = os.Stdout

		// 执行shell命令
		err = session.Run("docker rm -f " + preVersion)
		if err != nil {
			log.Fatalf("Failed to run command: %s", err)
		}

		fmt.Println(nodeTag, "Docker destroyContainer successfully!\n")
	}
}
