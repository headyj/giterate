package entities

import (
	"log"
)

type Repository struct {
	URL           string
	Destination   string
	DefaultBranch string
	CloneOptions  []Option `json:"CloneOptions" yaml:"CloneOptions"`
}

func PopulateRepositories(services []Service, arguments *Arguments) []Service {
	var repositories = []Repository{}
	for i, service := range services {
		_, exists := FindProviderArgs(&arguments.Providers, &service)
		if arguments.Providers == nil || exists {
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
				log.Fatalf("%s is not a valid API\n", service.API)
			}

			for _, jsonEntity := range service.Entities {
				switch jsonEntity.Type {
				case "group":
					repositoriesInfo := gitClient.GetGroupRepositories(jsonEntity.Path, jsonEntity.Name, jsonEntity.Recurse, service.CloneType)
					for _, repositoryInfo := range repositoriesInfo {

						repository := NewRepository(&service, repositoryInfo, jsonEntity)
						_, exists := FindRepositoryArgs(&arguments.Repositories, &repository)
						if arguments.Repositories == nil || exists {
							UpdateRepositories(&repositories, repository)
						}
					}
				case "repository":
					repositoryInfo := gitClient.GetRepository(jsonEntity.Path, jsonEntity.Name, service.CloneType)
					repository := NewRepository(&service, repositoryInfo, jsonEntity)
					_, exists := FindRepositoryArgs(&arguments.Repositories, &repository)
					if arguments.Repositories == nil || exists {
						UpdateRepositories(&repositories, repository)
					}
				default:
					log.Fatalf("%s is not a valid Type\n", jsonEntity.Type)
				}
			}
			(services)[i].Repositories = repositories
		}
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

func NewRepository(service *Service, repositoryInfo RepositoryInfo, jsonEntity Entity) Repository {
	if jsonEntity.Destination == "" {
		jsonEntity.Destination = service.Destination
	}
	return Repository{repositoryInfo.URL, jsonEntity.Destination + "/" + repositoryInfo.Path, repositoryInfo.DefaultBranch, jsonEntity.CloneOptions}
}

func FindRepository(repositories *[]Repository, val string) (int, bool) {
	for i, repository := range *repositories {
		if repository.URL == val {
			return i, true
		}
	}
	return -1, false
}

func FindProviderArgs(providers *[]string, service *Service) (int, bool) {
	for i, provider := range *providers {
		if service.Name != "" {
			if service.Name == provider {
				return i, true
			}
		} else {
			if service.BaseURL == provider {
				return i, true
			}
		}
	}
	return -1, false
}

func FindRepositoryArgs(repoArgs *[]string, repository *Repository) (int, bool) {
	for i, repoArg := range *repoArgs {
		if repository.Destination == repoArg {
			return i, true
		}
	}
	return -1, false
}
