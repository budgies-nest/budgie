package helpers

import (
	"bufio"
	"os"
	"runtime"
	"strings"

	"github.com/budgies-nest/budgie/enums/base"
	"github.com/budgies-nest/budgie/enums/environments"
)

func DetectContainerEnvironment() string {
	switch runtime.GOOS {
	case "linux":
		return detectLinuxContainer()
	case "windows":
		return detectWindowsContainer()
	case "darwin":
		return detectMacOSContainer()
	default:
		return "Unknown OS"
	}
}

func detectLinuxContainer() string {
	// Vérifier Docker
	if _, err := os.Stat("/.dockerenv"); !os.IsNotExist(err) {
		return "Docker"
	}

	// Vérifier Kubernetes
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return "Kubernetes"
	}

	// Vérifier les cgroups
	if file, err := os.Open("/proc/1/cgroup"); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "docker") {
				return "Docker"
			}
			if strings.Contains(line, "containerd") {
				return "Containerd"
			}
		}
	}

	return "Local"
}

// To be tested on Windows
func detectWindowsContainer() string {
	if os.Getenv("CONTAINER_SANDBOX_MOUNT_POINT") != "" {
		return "Docker"
	}

	hostname, err := os.Hostname()
	if err == nil && len(hostname) == 12 {
		return "Docker"
	}

	return "Local"
}

// To be tested on macOS
func detectMacOSContainer() string {
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return "Kubernetes"
	}

	if os.Getenv("DOCKER_DESKTOP") != "" {
		return "Docker"
	}

	return "Local"
}

func GetModelRunnerBaseUrl() string {
	// Detect if running in a container or locally
	if DetectContainerEnvironment() == environments.Local {
		return base.DockerModelRunnerLocalURL
	}
	return base.DockerModelRunnerContainerURL
}
