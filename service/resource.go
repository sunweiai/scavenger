package service

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

type KindInfo struct {
	sourceName string
	kind       string
	nameSpace  string
}

// 获取pod的相关信息，并通过pod信息查询到部署此pod的资源类型
func (kindInfo *KindInfo) GetPodType(client *kubernetes.Clientset, nameSpace, podName string) *KindInfo {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Runtime err caught: %v", r)
		}
	}()
	pod, err := client.CoreV1().Pods(nameSpace).Get(context.TODO(), podName, metav1.GetOptions{})
	fmt.Printf("namespace:%s,podname:%s\n", nameSpace, podName)
	if err != nil {
		panic("Get pod error!")
	}
	kindInfo.kind = pod.OwnerReferences[0].Kind
	kindInfo.sourceName = pod.OwnerReferences[0].Name
	kindInfo.nameSpace = pod.Namespace
	//fmt.Printf("类型：%s,名称：%s\n", kindInfo.kind, kindInfo.sourceName)
	return kindInfo
}

// 如果是replicaSet类型，则继续通过replicaSet查询到deployment的信息
func getRepType(clienSet *kubernetes.Clientset, sourceName, namespace string) KindInfo {
	var kindInfo KindInfo
	api := clienSet.AppsV1()
	replicaSet, err := api.ReplicaSets(namespace).Get(context.TODO(), sourceName, metav1.GetOptions{})
	if err != nil {
		log.Fatal("Get pod error!")
	}
	kindInfo.kind = replicaSet.OwnerReferences[0].Kind
	kindInfo.sourceName = replicaSet.OwnerReferences[0].Name
	kindInfo.nameSpace = replicaSet.Namespace
	//fmt.Printf("类型:%s,名称:%,命名空间:%s\n", kindInfo.kind, kindInfo.sourceName, kindInfo.nameSpace)
	return kindInfo
}

// 根据资源类型进行删除
func DeleteSource(clienSet *kubernetes.Clientset, sourceName, kind, namespace string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Runtime err caught: %v", r)
		}
	}()
	switch kind {
	case "ReplicaSet":
		kindInfo := getRepType(clienSet, sourceName, namespace)
		ifDeploy := GetDeployment(clienSet, kindInfo.sourceName, kindInfo.nameSpace)
		// 判断资源类型是否存在，避免重复删除
		if ifDeploy {
			DeleteDeployment(clienSet, kindInfo.sourceName, kindInfo.nameSpace)
		}

	case "DaemonSet":
		ifDaemon := GetDaemonset(clienSet, sourceName, namespace)
		if ifDaemon {
			DeleteDaemonSet(clienSet, sourceName, namespace)
		}
	}
}

// 删除deployment资源
func DeleteDeployment(clientSet *kubernetes.Clientset, deploymentName, namespace string) {
	// 通过appsv1去访问核心api资源,并获取deployment列表
	api := clientSet.AppsV1()

	fmt.Printf("删除成功: deployment: %s,namespace: %s \n", deploymentName, namespace)
	fmt.Printf("++++++++\n")

	err := api.Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		panic("delete deployment failed")

	}
}

// 删除DaemonSet资源
func DeleteDaemonSet(clientSet *kubernetes.Clientset, deamonSetName, namespace string) {
	api := clientSet.AppsV1()
	fmt.Printf("删除成功: deamonSet name: %s,namespace: %s\n", deamonSetName, namespace)
	fmt.Printf("++++++++\n")
	err := api.DaemonSets(namespace).Delete(context.TODO(), deamonSetName, metav1.DeleteOptions{})
	if err != nil {
		panic("delete daemonset failed")
	}
}

func GetDeployment(clientSet *kubernetes.Clientset, deploymentName, namespace string) bool {
	api := clientSet.AppsV1()
	_, err := api.Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return false
	}
	return true
}

func GetDaemonset(clientSet *kubernetes.Clientset, daemonsetName, namespace string) bool {
	api := clientSet.AppsV1()
	_, err := api.DaemonSets(namespace).Get(context.TODO(), daemonsetName, metav1.GetOptions{})
	if err != nil {
		return false
	}
	return true
}
