package kubeops

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	deploy "k8s.io/api/apps/v1"
	v1ns "k8s.io/api/core/v1"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//Internal functions
//Same cluster configuration setted

var APIServer = "https://kubernetes.default.svc.cluster.local:443"

//TODO: Set Configurable API Server adress with config-map or env-variable default one should be "https://kubernetes.default.svc.cluster.local:443"

func GetAdminKubeClient() (*kubernetes.Clientset, error) {
	// TODO: Set Configurable CA file Default one should be service account path!
	tlsClientConfig := rest.TLSClientConfig{}
	optionalCA := os.Getenv("CA_PATH")
	optionalToken := os.Getenv("TOKEN_PATH")
	optionalAPI := os.Getenv("KUBE_HOST")
	fmt.Println(optionalToken)
	tlsClientConfig.CAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	var config *rest.Config
	var err error
	if optionalCA != "" && optionalToken != "" && optionalAPI != "" {
		token, err := os.ReadFile(optionalToken)
		if err != nil {
			return nil, err
		}
		tlsClientConfig.CAFile = optionalCA
		config = &rest.Config{
			Host:            optionalAPI,
			BearerToken:     string(token),
			TLSClientConfig: tlsClientConfig,
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, err
}

func GetKubeClient(token string) (*kubernetes.Clientset, error) {
	// TODO: Set Configurable CA file Default one should be service account path!
	tlsClientConfig := rest.TLSClientConfig{}
	tlsClientConfig.CAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	optionalCA := os.Getenv("CA_PATH")
	optionalAPI := os.Getenv("KUBE_HOST")
	if optionalCA != "" {
		tlsClientConfig.CAFile = optionalCA
	}
	if optionalAPI != "" {
		APIServer = optionalAPI
	}
	config := &rest.Config{
		Host:            APIServer,
		BearerToken:     token,
		TLSClientConfig: tlsClientConfig,
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, err
}

func CreateDeploymentService(name string, namespace string, token string) (int32, int32, error) {
	client, err := GetKubeClient(token)
	if err != nil {
		return 0, 0, err
	}
	replicas := int32(1)
	dp := client.AppsV1().Deployments(namespace)

	//Create Service first after assign nodeport paramater to neko env var
	//Service Definition Template
	svc := client.CoreV1().Services(namespace)
	service := v1ns.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1ns.ServiceSpec{
			Selector: map[string]string{
				"robot": name + "robolaunch",
			},
			Type: v1ns.ServiceTypeNodePort,
			Ports: []v1ns.ServicePort{
				{
					Protocol:   v1ns.ProtocolTCP,
					Port:       8080,
					Name:       "http",
					TargetPort: intstr.FromInt(8080),
				},
				{
					Protocol:   v1ns.ProtocolUDP,
					Port:       31555,
					Name:       "neko-webrtc",
					TargetPort: intstr.FromInt(31555),
				},
				{
					Protocol:   v1ns.ProtocolTCP,
					Port:       3000,
					Name:       "theia",
					TargetPort: intstr.FromInt(3000),
				},
			},
		},
	}
	createdSvc, err := svc.Create(context.TODO(), &service, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Service did not created: %v", err)
		return 0, 0, err
	}

	var udpPort int32
	var theiaPort int32

	for _, port := range createdSvc.Spec.Ports {
		if port.Name == "neko-webrtc" {
			udpPort = port.NodePort
		}
		if port.Name == "theia" {
			theiaPort = port.NodePort
		}
	}

	// update service with node port details! only for demo
	//fetch service again!

	createdSvc.Spec.Ports[1].Port = udpPort
	createdSvc.Spec.Ports[1].TargetPort = intstr.FromInt(int(udpPort))

	_, err = svc.Update(context.TODO(), createdSvc, metav1.UpdateOptions{})
	if err != nil {
		fmt.Printf("Service didn't  updated!:%v\n", err)
		return 0, 0, err
	}
	//Deployment definition template
	deployment := deploy.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: deploy.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"robot": name + "robolaunch",
				},
			},
			Replicas: &replicas,
			Template: v1ns.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"robot": name + "robolaunch",
					},
				},
				Spec: v1ns.PodSpec{
					Containers: []v1ns.Container{
						{
							Name:  "neko",
							Image: "m1k1o/neko:firefox",
							Stdin: true,
							TTY:   true,
							Ports: []v1ns.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8080,
									Protocol:      v1ns.ProtocolTCP,
								},
								{
									Name:          "neko-webrtc",
									ContainerPort: udpPort,
									Protocol:      v1ns.ProtocolUDP,
								},
							},
							Env: []v1ns.EnvVar{
								{Name: "NEKO_BIND", Value: "0.0.0.0:8080"},
								{Name: "NEKO_UDP_PORT", Value: strconv.Itoa(int(udpPort)) + "-" + strconv.Itoa(int(udpPort))},
								{Name: "NEKO_EPR", Value: strconv.Itoa(int(udpPort)) + "-" + strconv.Itoa(int(udpPort))},
								{Name: "NEKO_ICELITE", Value: "1"},
								{Name: "NEKO_SCREEN", Value: "1920x1080@30"},
							},
						},
						{

							Name:  "theia",
							Image: "theiaide/sadl",
							Stdin: true,
							TTY:   true,
							Ports: []v1ns.ContainerPort{
								{
									Name:          "theia",
									ContainerPort: 3000,
									Protocol:      v1ns.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
	}

	//Create deployment!
	_, err = dp.Create(context.TODO(), &deployment, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Deployment did not created: %v", err)
		return 0, 0, err
	}

	return udpPort, theiaPort, nil

}

func DeleteDeploymentService(name string, namespace string, token string) error {
	fmt.Printf("Name: %v\nNamespace: %v\n", name, namespace)
	client, err := GetKubeClient(token)
	if err != nil {
		return err
	}
	deploy := client.AppsV1().Deployments(namespace)

	svc := client.CoreV1().Services(namespace)

	err = deploy.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	err = svc.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func ScaleDeployment(name string, namespace string, replicas int32, token string) error {
	client, err := GetKubeClient(token)
	if err != nil {
		return err
	}
	deployment, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	deployment.Spec.Replicas = &replicas
	_, err = client.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func GetUnallocatedPort() (int32, error) {
	client, err := GetAdminKubeClient()
	if err != nil {
		return 0, err
	}

	services, err := client.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0, err
	}

	var allocatedPorts []int32

	for _, service := range services.Items {
		for _, port := range service.Spec.Ports {
			if port.NodePort != 0 {
				allocatedPorts = append(allocatedPorts, port.NodePort)
			}
		}
	}

	min := 30000
	max := 32000

	for i := 0; i < 100; i++ {
		source := rand.NewSource(time.Now().UnixMilli())
		randGen := rand.New(source)
		randomPort := randGen.Intn(max-min) + min
		if !contains(allocatedPorts, randomPort) {
			return int32(randomPort), nil
		}
	}

	return 0, errors.New("ports in cluster are allocated, try again")
}

func contains(s []int32, num int) bool {
	for _, v := range s {
		if int(v) == num {
			return true
		}
	}
	return false
}

//TODO: Create edit deployment method to scale up & scale down operations.

//TODO: Implement namespace check function
func CheckNamespace(namespace string) error {
	client, err := GetAdminKubeClient()
	if err != nil {
		return err
	}
	res, err := client.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if res.Name != namespace {
		return errors.New("namespace not found")
	}
	return nil
}

//TODO: Implement namespace create function
func CreateNamespace(namespace string, args ...metav1.CreateOptions) (*v1ns.Namespace, error) {
	//namespace template
	ns := v1ns.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	// Default create option
	option := metav1.CreateOptions{}
	if len(args) == 1 {
		option = args[0]
	}
	client, err := GetAdminKubeClient()
	if err != nil {
		return nil, err
	}

	nsObj, err := client.CoreV1().Namespaces().Create(context.Background(), &ns, option)
	if err != nil {
		return nil, err
	}
	return nsObj, nil
}

// The following function can create role & role binding
func CreateRole(namespace string, args ...metav1.CreateOptions) (*v1.Role, *v1.RoleBinding, error) {
	client, err := GetAdminKubeClient()
	if err != nil {
		fmt.Printf("Client connection failed: %v", err)
		return nil, nil, err
	}
	option := metav1.CreateOptions{}
	if len(args) == 1 {
		option = args[0]
	}
	role := &v1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespace + "_role",
			Namespace: namespace,
		},
		Rules: []v1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"pods", "services", "deployments"},
				Verbs:     []string{"get", "list", "create", "update"},
			},
		},
	}

	rbind := &v1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespace + "_role",
			Namespace: namespace,
		},
		Subjects: []v1.Subject{
			{
				Kind:     "Group",
				Name:     namespace + "_role",
				APIGroup: "rbac.authorization.k8s.io",
			},
		},
		RoleRef: v1.RoleRef{
			Kind:     "Role",
			Name:     namespace + "_role",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	// Create user role
	clusterRole, err := client.RbacV1().Roles(namespace).Create(context.TODO(), role, option)
	if err != nil {
		return nil, nil, err
	}

	//Create user role binding
	clusterRoleBind, err := client.RbacV1().RoleBindings(namespace).Create(context.TODO(), rbind, option)
	if err != nil {
		return nil, nil, err
	}
	return clusterRole, clusterRoleBind, nil
}
