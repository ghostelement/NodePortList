package getport

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	//v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type ServicePorts struct {
	ServiceName string
	ServicePort []int32
}
type ServiceList struct {
	Namespace string
	Service   []ServicePorts
}

type NodePortsList struct {
	NodePortList []ServiceList
}

func GetNodePort(token string) (NodePortsList, error) {
	NodePortsList := NodePortsList{}
	// 通过配置文件创建Kubernetes客户端

	config, err := clientcmd.BuildConfigFromFlags("", token)
	if err != nil {
		//fmt.Println(err.Error())
		return NodePortsList, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		//fmt.Println(err.Error())
		return NodePortsList, err
	}

	// 查询所有的 Namespace
	namespaces, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//fmt.Println(err.Error())
		return NodePortsList, err
	}

	for _, ns := range namespaces.Items {
		ServiceList := ServiceList{}
		ServiceList.Namespace = ns.Name
		// 列出 Namespace 下的所有 Service
		services, err := client.CoreV1().Services(string(ns.Name)).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			//fmt.Println(err.Error())
			return NodePortsList, err
		}
		for _, svc := range services.Items {
			ServicePorts := ServicePorts{}
			ServicePorts.ServiceName = svc.Name
			if svc.Spec.Type == "NodePort" {
				for _, port := range svc.Spec.Ports {
					ServicePorts.ServicePort = append(ServicePorts.ServicePort, port.NodePort)
				}
				ServiceList.Service = append(ServiceList.Service, ServicePorts)
			}
		}
		if ServiceList.Service != nil {
			NodePortsList.NodePortList = append(NodePortsList.NodePortList, ServiceList)
		}
	}

	return NodePortsList, nil

}
