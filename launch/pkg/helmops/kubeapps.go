package helmops

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

var KubeappsHost = os.Getenv("KUBEAPPS_SERVER_IP")

type AppRepoResponse struct {
	AppRepo AppRepoItems `json:"appRepository"`
	Code    int          `json:"code,omitempty"`
	Message string       `json:"message,omitempty"`
}

type RefreshAppRepoResponse struct {
	AppRepo AppRepository `json:"appRepository"`
	Code    int           `json:"code"`
	Message string        `json:"message"`
}

type AppRepoItems struct {
	Items []AppRepository `json:"items"`
}

type AppRepository struct {
	Metadata AppRepositoryMetadata `json:"metadata"`
	Spec     AppRepositorySpec     `json:"spec"`
}

type AppRepositorySpec struct {
	URL string `json:"url"`
}

type AppRepositoryMetadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type DeleteAppRepositoryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CreateReleaseBody struct {
	AppRepositoryResourceName      string `json:"appRepositoryResourceName"`
	AppRepositoryResourceNamespace string `json:"appRepositoryResourceNamespace"`
	ChartName                      string `json:"chartName"`
	ReleaseName                    string `json:"releaseName"`
	Version                        string `json:"version"`
	Values                         string `json:"values"`
}

type CreateReleaseResponse struct {
	Data    ReleaseInfo `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
}

type ReleaseInfo struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Info      ReleaseInfoDetails `json:"info"`
	Version   int                `json:"version"`
}

type ReleaseInfoDetails struct {
	FirstDeployed string `json:"first_deployed"`
	LastDeployed  string `json:"last_deployed"`
	Description   string `json:"description"`
	Status        string `json:"status"`
}

type RegisterAppRepositoryBody struct {
	AppRepository RegisterAppRepositoryBodyDetails `json:"appRepository"`
}

type RegisterAppRepositoryBodyDetails struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	RepoURL     string `json:"repoURL"`
}

type RegisterAppRepositoryResponse struct {
	AppRepository AppRepository `json:"appRepository"`
	Code          int           `json:"code"`
	Message       string        `json:"message"`
}

type UpdateReleaseBody struct {
	AppRepositoryResourceName      string `json:"appRepositoryResourceName"`
	AppRepositoryResourceNamespace string `json:"appRepositoryResourceNamespace"`
	ChartName                      string `json:"chartName"`
	ReleaseName                    string `json:"releaseName"`
	Version                        string `json:"version"`
	Values                         string `json:"values"`
}

type UpdateReleaseResponse struct {
	Data    ReleaseInfo `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
}

/*
 * App Repository Information
 */

func GetAppRepository(token string, cluster string, namespace string, name string) (AppRepository, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", KubeappsHost+"/api/v1/clusters/"+cluster+"/namespaces/"+namespace+"/apprepositories", nil)
	if err != nil {
		return AppRepository{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return AppRepository{}, err
	}

	var appRepositoriesResponse AppRepoResponse
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AppRepository{}, err
	}

	err = json.Unmarshal(body, &appRepositoriesResponse)
	if err != nil {
		return AppRepository{}, err
	}

	for _, repo := range appRepositoriesResponse.AppRepo.Items {
		if repo.Metadata.Name == name {
			return repo, nil
		}
	}

	if appRepositoriesResponse.Code != 0 {
		return AppRepository{}, errors.New(appRepositoriesResponse.Message)
	}
	return AppRepository{}, errors.New("app repository not found on this namespace, please register the app repository first")

}

/*
 * App Repository will be added/registered to X namespace with this function.
 */

func RegisterAppRepository(token string, cluster string, namespace string, appRepository RegisterAppRepositoryBody) (AppRepository, error) {

	client := &http.Client{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(appRepository)
	if err != nil {
		return AppRepository{}, err
	}

	req, err := http.NewRequest("POST", KubeappsHost+"/api/v1/clusters/"+cluster+"/namespaces/"+namespace+"/apprepositories", &buf)
	if err != nil {
		return AppRepository{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return AppRepository{}, err
	}

	var registerAppRepositoryResponse RegisterAppRepositoryResponse
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AppRepository{}, err
	}

	err = json.Unmarshal(body, &registerAppRepositoryResponse)
	if err != nil {
		return AppRepository{}, err
	}

	if registerAppRepositoryResponse.Code != 0 {
		return AppRepository{}, errors.New(registerAppRepositoryResponse.Message)
	}

	return registerAppRepositoryResponse.AppRepository, nil

}

/*
 * When a launch is updated (on ChartMuseum), it'll be updated on Kubeapps with this function
 */

func RefreshAppRepository(token string, cluster string, namespace string, name string) (AppRepository, error) {

	client := &http.Client{}

	req, err := http.NewRequest("POST", KubeappsHost+"/api/v1/clusters/"+cluster+"/namespaces/"+namespace+"/apprepositories/"+name+"/refresh", nil)
	if err != nil {
		return AppRepository{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return AppRepository{}, err
	}

	var appRepositoryResp RefreshAppRepoResponse
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AppRepository{}, err
	}

	err = json.Unmarshal(body, &appRepositoryResp)
	if err != nil {
		return AppRepository{}, err
	}

	if appRepositoryResp.Code != 0 {
		return AppRepository{}, errors.New(appRepositoryResp.Message)
	}

	return appRepositoryResp.AppRepo, nil

}

/*
 * App Repository deletion. Needed for test cleanups. Returns a boolean.
 */

func DeleteAppRepository(token string, cluster string, namespace string, name string) (bool, error) {

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", KubeappsHost+"/api/v1/clusters/"+cluster+"/namespaces/"+namespace+"/apprepositories/"+name, nil)
	if err != nil {
		return false, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	var deleteAppRepositoryResponse DeleteAppRepositoryResponse
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == 200 {
		return true, nil
	} else {
		err = json.Unmarshal(body, &deleteAppRepositoryResponse)
		if err != nil {
			return false, err
		}

		return false, errors.New(deleteAppRepositoryResponse.Message)
	}

}

/*
 * Release (Launch Instance) creation. Send the values as an empty string to use default configuration.
 */

func CreateRelease(token string, cluster string, namespace string, release CreateReleaseBody) (CreateReleaseResponse, error) {

	client := &http.Client{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(release)
	if err != nil {
		return CreateReleaseResponse{}, err
	}

	req, err := http.NewRequest("POST", KubeappsHost+"/api/kubeops/v1/clusters/"+cluster+"/namespaces/"+namespace+"/releases", &buf)
	if err != nil {
		return CreateReleaseResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return CreateReleaseResponse{}, err
	}

	var createReleaseResp CreateReleaseResponse
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CreateReleaseResponse{}, err
	}

	err = json.Unmarshal(body, &createReleaseResp)
	if err != nil {
		return CreateReleaseResponse{}, err
	}

	if createReleaseResp.Code != 0 {
		// returns 401 if cannot be created
		return CreateReleaseResponse{}, errors.New(createReleaseResp.Message)
	}

	return createReleaseResp, nil
}

/*
 * Release (Launch Instance) update.
 */

func UpdateRelease(token string, cluster string, namespace string, name string, release UpdateReleaseBody) (UpdateReleaseResponse, error) {

	client := &http.Client{}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(release)
	if err != nil {
		return UpdateReleaseResponse{}, err
	}

	req, err := http.NewRequest("PUT", KubeappsHost+"/api/kubeops/v1/clusters/"+cluster+"/namespaces/"+namespace+"/releases/"+name, &buf)
	if err != nil {
		return UpdateReleaseResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return UpdateReleaseResponse{}, err
	}

	var updateReleaseResp UpdateReleaseResponse
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UpdateReleaseResponse{}, err
	}

	err = json.Unmarshal(body, &updateReleaseResp)
	if err != nil {
		return UpdateReleaseResponse{}, err
	}

	if updateReleaseResp.Code != 0 {
		return UpdateReleaseResponse{}, errors.New(updateReleaseResp.Message)
	}

	return updateReleaseResp, nil

}

/*
 * Release (Launch Instance) deletion.
 */

func DeleteRelease(token string, cluster string, namespace string, name string) (bool, error) {

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", KubeappsHost+"/api/kubeops/v1/clusters/"+cluster+"/namespaces/"+namespace+"/releases/"+name, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if string(body) != "OK" {
		return false, errors.New(string(body))
	}

	return true, nil
}
