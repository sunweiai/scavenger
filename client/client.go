package client

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//func getConf() string {
//	kubeconfig := filepath.Join(
//		//os.Getenv("KUBECONFIG"),
//		homedir.HomeDir(), ".kube", "config",
//	)
//	return kubeconfig
//}

func RestClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
	}
	return clientSet
	// 使用当前上下文环境
	//kubeconfig := getConf()
	//config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//// 根据指定的 config 创建一个新的 clientSet
	//clientSet, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	//panic(err.Error())
	//	log.Fatal(err)
	//}
	//return clientSet
}
