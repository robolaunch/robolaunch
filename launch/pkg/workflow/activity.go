package launchflow

import (
	"fmt"
	"strconv"

	"github.com/robolaunch/robolaunch/launch/pkg/account"
	"github.com/robolaunch/robolaunch/launch/pkg/helmops"
	"github.com/robolaunch/robolaunch/launch/pkg/kubeops"
)

func CreateUserSpace(l LaunchRequest) (string, error) {
	//TODO: Check namespace is available!
	err := kubeops.CheckNamespace(l.Namespace)
	if err != nil {
		// Create namespace
		_, err := kubeops.CreateNamespace(l.Namespace)
		if err != nil {
			return "", err
		}
		//Create role for user
		_, _, err = kubeops.CreateRole(l.Namespace)
		if err != nil {
			return "", err
		}

		//Create keycloak role
		_, err = account.CreateGroup(l.Namespace)
		if err != nil {
			return "", err
		}
		//Bind the role
		err = account.BindGroup(l.Username, l.Namespace)
		if err != nil {
			return "", err
		}
		return "Namespace created: " + l.Namespace, nil
	}
	return "Namespace avaliable", nil
	//TODO: Create Namespace & Role if not available
	//TODO: Create Group and bind them user

}

func CreateLaunch(l LaunchRequest) (LaunchStatus, error) {

	// Launch type not used right now!
	//TODO: Add following functions part as a workflow

	// Check namespace first

	udpPort, theiaPort, err := kubeops.CreateDeploymentService(l.Name, l.Namespace, l.IDToken)
	if err != nil {
		return LaunchStatus{}, err
	}
	return LaunchStatus{
		Name:           l.Name,
		Namespace:      l.Namespace,
		LaunchType:     "",
		WorkloadStatus: "RUNNING",
		NodeIp:         "", // TODO: Add Get Node IP ops function
		TheiaPort:      theiaPort,
		WebRpcPort:     udpPort,
	}, nil
}

func CreateLaunchHelm(l LaunchRequest) (LaunchStatus, error) {

	// Launch type not used right now!
	//TODO: Add following functions part as a workflow

	// Check namespace first
	httpPort, err := kubeops.GetUnallocatedPort(l.IDToken)
	if err != nil {
		return LaunchStatus{}, err
	}
	// webrtcPort, err := kubeops.GetUnallocatedPort(l.IDToken)
	// if err != nil {
	// 	return LaunchStatus{}, err
	// }

	resp, err := helmops.CreateRelease(l.IDToken, "default", l.Namespace, helmops.CreateReleaseBody{
		AppRepositoryResourceName:      "robot-helm-charts",
		AppRepositoryResourceNamespace: "default",
		ReleaseName:                    l.Name,
		ChartName:                      "jackal",
		Version:                        "0.1.0",
		// Values:                         "{\"launchName\": \"jackal-2\"\nhttpPort: " + string(httpPort) + "\n\"webrtcPort\": " + string(webrtcPort) + "}",

		Values: "launchName: jackal-2\nhttpPort: " + strconv.Itoa(int(httpPort)) + "\nwebrtcPort: " + strconv.Itoa(int(httpPort)),
	})
	if err != nil {
		return LaunchStatus{}, err
	}
	fmt.Println("launchName: jackal-2\nhttpPort: " + string(httpPort) + "\nwebrtcPort: " + string(httpPort))
	return LaunchStatus{
		Name:           resp.Data.Name,
		Namespace:      resp.Data.Namespace,
		LaunchType:     "",
		WorkloadStatus: "RUNNING",
		NodeIp:         "", // TODO: Add Get Node IP ops function
		TheiaPort:      0,
		WebRpcPort:     httpPort,
	}, nil
}

func DeleteLaunchHelm(l LaunchRequest) (LaunchStatus, error) {

	_, err := helmops.DeleteRelease(l.IDToken, "default", l.Namespace, l.Name)
	if err != nil {
		return LaunchStatus{}, err
	}
	return LaunchStatus{
		Name:           l.Name,
		Namespace:      l.Namespace,
		LaunchType:     "",
		WorkloadStatus: "DELETED",
		NodeIp:         "", // For a moment it would be static
		TheiaPort:      0,
		WebRpcPort:     0,
	}, nil

}

func DeleteLaunch(l LaunchRequest) (LaunchStatus, error) {

	err := kubeops.DeleteDeploymentService(l.Name, l.Namespace, l.IDToken)
	if err != nil {
		return LaunchStatus{}, err
	}
	return LaunchStatus{
		Name:           l.Name,
		Namespace:      l.Namespace,
		LaunchType:     "",
		WorkloadStatus: "DELETED",
		NodeIp:         "", // For a moment it would be static
		TheiaPort:      0,
		WebRpcPort:     0,
	}, nil

}

func ScaleOut(l LaunchRequest) (string, error) {
	err := kubeops.ScaleDeployment(l.Name, l.Namespace, 0, l.IDToken)
	if err != nil {
		return "", nil
	}
	return "STOPPED", nil
}

func ScaleDownHelm(l LaunchRequest) (string, error) {
	_, err := helmops.UpdateRelease(l.IDToken, "default", l.Namespace, l.Name, helmops.UpdateReleaseBody{
		AppRepositoryResourceName:      "robot-helm-charts",
		AppRepositoryResourceNamespace: "default",
		ReleaseName:                    l.Name,
		ChartName:                      "jackal",
		Version:                        "0.1.0",
		Values:                         "replicas: 0",
	})
	if err != nil {
		return "", err
	}
	return "STOPPED", nil

}

func ScaleUpHelm(l LaunchRequest) (string, error) {
	_, err := helmops.UpdateRelease(l.IDToken, "default", l.Namespace, l.Name, helmops.UpdateReleaseBody{
		AppRepositoryResourceName:      "robot-helm-charts",
		AppRepositoryResourceNamespace: "default",
		ReleaseName:                    l.Name,
		ChartName:                      "jackal",
		Version:                        "0.1.0",
		Values:                         "replicas: 1",
	})
	if err != nil {
		return "", err

	}
	return "RUNNING", nil
}

func ScaleUp(l LaunchRequest) (string, error) {
	err := kubeops.ScaleDeployment(l.Name, l.Namespace, 1, l.IDToken)
	if err != nil {
		return "", nil
	}
	return "RUNNING", nil
}
