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
	NodePortList map[string]ServiceList
}

func GetNodePort(token string) (NodePortsList, error) {
	nodePortsList := NodePortsList{
		NodePortList: make(map[string]ServiceList),
	}
	// 通过配置文件创建Kubernetes客户端

	config, err := clientcmd.BuildConfigFromFlags("", token)
	if err != nil {
		//fmt.Println(err.Error())
		return nodePortsList, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		//fmt.Println(err.Error())
		return nodePortsList, err
	}

	//serviceList := make(map[string]ServiceList)
	//fmt.Println(ServiceList)
	// 列出 Namespace 下的所有 Service
	services, err := client.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		//fmt.Println(err.Error())
		return nodePortsList, err
	}
	for _, svc := range services.Items {
		if svc.Spec.Type != "NodePort" {
			continue //只处理NodePort类型的服务
		}
		// 确保每个命名空间对应的 ServiceList 已经初始化
		if _, ok := nodePortsList.NodePortList[svc.Namespace]; !ok {
			nodePortsList.NodePortList[svc.Namespace] = ServiceList{Namespace: svc.Namespace, Service: []ServicePorts{}}
		}

		// 创建当前服务的 ServicePorts 实例
		servicePorts := ServicePorts{
			ServiceName: svc.Name,
			ServicePort: []int32{},
		}

		// 收集 NodePort 并添加到 ServicePorts
		for _, port := range svc.Spec.Ports {
			if port.NodePort > 0 { // 确保只添加有效的 NodePort
				servicePorts.ServicePort = append(servicePorts.ServicePort, port.NodePort)
			}
		}

		// 获取命名空间对应的 ServiceList 的副本，更新后再放回映射
		serviceList := nodePortsList.NodePortList[svc.Namespace]
		serviceList.Service = append(serviceList.Service, servicePorts)
		nodePortsList.NodePortList[svc.Namespace] = serviceList
	}
	return nodePortsList, nil

}
