package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Definir o nome do namespace que você deseja verificar
	namespace := "userdept"

	// Obter a lista de recursos disponíveis no namespace
	resourcesBytes, err := exec.Command("kubectl", "api-resources", "--namespaced=true", "--verbs=list", "-o", "name").Output()
	if err != nil {
		fmt.Println("Erro ao obter lista de recursos:", err)
		return
	}
	resources := strings.TrimSuffix(string(resourcesBytes), "\n")
	resourceList := strings.Split(resources, "\n")

	// Definir lista de recursos que queremos verificar
	targetResources := []string{"deployments.apps", "services", "secrets", "ingress", "configmaps"}

	path := "kubernetes"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			fmt.Println("Erro ao criar a pasta:", err)
			return
		}
		fmt.Println("Pasta criada com sucesso!")
	} else {
		fmt.Println("A pasta", path, "existe")
	}

	// Verificar cada recurso usando kubectl e salvá-lo em um arquivo YAML separado, se for um recurso de destino
	for _, resource := range resourceList {
		resourceName := strings.TrimPrefix(resource, "name/")
		if contains(targetResources, resourceName) {
			fmt.Printf("Verificando recurso %s ...\n", resourceName)

			cmd := exec.Command("kubectl", "get", resource, "-n", namespace, "-o", "yaml")
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Erro ao verificar recurso %s: %s\n", resourceName, err)
				continue
			}

			yaml := stdout.Bytes()
			err = ioutil.WriteFile(fmt.Sprintf("%s.yaml", "kubernetes/"+resourceName), yaml, os.ModePerm)
			if err != nil {
				fmt.Printf("Erro ao salvar recurso %s em arquivo YAML: %s\n", resourceName, err)
				continue
			}
		}
	}
}

func contains(list []string, item string) bool {
	for _, listItem := range list {
		if item == listItem {
			return true
		}
	}
	return false
}
