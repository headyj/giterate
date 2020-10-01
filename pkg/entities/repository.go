package entities

import (
	"log"
)

type Repository struct {
	URL          string
	Destination  string
	CloneOptions []Option `json:"CloneOptions" yaml:"CloneOptions"`
}

func PopulateRepositories(services []Service) []Service {
	var repositories = []Repository{}
	for i, service := range services {
		repositories = nil
		var gitClient GitClient
		switch service.API {
		case "gitlab":
			gitClient = GitlabAuth(service)
		case "bitbucketv1":
			gitClient = BitbucketAuth(service)
		case "bitbucketv2":
			gitClient = BitbucketV2Auth(service)
		case "github":
			gitClient = GithubAuth(service)
		default:
			log.Fatalf("%s is not a valid API", service.API)
		}

		for _, jsonEntity := range service.Entities {
			switch jsonEntity.Type {
			case "group":
				repositoriesInfo := gitClient.GetGroupRepositories(jsonEntity.Path, jsonEntity.Name, jsonEntity.Recurse, service.CloneType)
				for _, repositoryInfo := range repositoriesInfo {
					repository := NewRepository(&service, repositoryInfo.URL, jsonEntity.Destination, repositoryInfo.Path, jsonEntity.CloneOptions)
					UpdateRepositories(&repositories, repository)
				}
			case "repository":
				repositoryInfo := gitClient.GetRepository(jsonEntity.Path, jsonEntity.Name, service.CloneType)
				repository := NewRepository(&service, repositoryInfo.URL, jsonEntity.Destination, repositoryInfo.Path, jsonEntity.CloneOptions)
				UpdateRepositories(&repositories, repository)
			default:
				log.Fatalf("%s is not a valid Type", jsonEntity.Type)
			}
		}
		(services)[i].Repositories = repositories
	}
	return services
}

func UpdateRepositories(repositories *[]Repository, newRepository Repository) {
	if len(*repositories) > 0 {
		_, exists := FindRepository(repositories, newRepository.URL)
		if !exists {
			*repositories = append(*repositories, newRepository)
		}
	} else {
		*repositories = append(*repositories, newRepository)
	}
}

func NewRepository(service *Service, URL string, destination string, path string, cloneOptions []Option) Repository {
	if destination == "" {
		destination = service.Destination
	}
	return Repository{URL, destination + "/" + path, cloneOptions}
}

func FindRepository(repositories *[]Repository, val string) (int, bool) {
	for i, repository := range *repositories {
		if repository.URL == val {
			return i, true
		}
	}
	return -1, false
}
