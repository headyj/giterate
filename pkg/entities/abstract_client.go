package entities

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/google/go-github/github"
	bitbucketv2 "github.com/ktrysmt/go-bitbucket"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

type RepositoryInfo struct {
	URL           string
	Path          string
	DefaultBranch string
}

type GitClient interface {
	GetGroupRepositories(path string, name string, recurse bool, cloneType string) []RepositoryInfo
	GetRepository(path string, name string, cloneType string) RepositoryInfo
}

type GitlabClient struct {
	Client *gitlab.Client
}

type BitbucketClient struct {
	Client *bitbucketv1.APIClient
}

type BitbucketClientV2 struct {
	Client *bitbucketv2.Client
}

type GithubClient struct {
	Client *github.Client
}

// Gitlab
func (gl GitlabClient) GetGroupRepositories(path string, name string, recurse bool, cloneType string) []RepositoryInfo {
	var repositoriesInfo []RepositoryInfo
	var finalPath string
	if path == "" {
		finalPath = name
	} else {
		finalPath = path + "/" + name
	}
	projects, _, err := gl.Client.Groups.ListGroupProjects(finalPath, &gitlab.ListGroupProjectsOptions{ListOptions: gitlab.ListOptions{PerPage: 10000}, IncludeSubgroups: gitlab.Bool(recurse)})
	if err != nil {
		log.Fatalf("Failed to get group %s: %v\n", finalPath, err)
	}
	for _, project := range projects {
		var repositoryLink string
		if cloneType == "ssh" {
			repositoryLink = project.SSHURLToRepo
		} else {
			repositoryLink = project.HTTPURLToRepo
		}
		repositoryInfo := RepositoryInfo{repositoryLink, project.PathWithNamespace, project.DefaultBranch}
		repositoriesInfo = append(repositoriesInfo, repositoryInfo)
	}
	return repositoriesInfo
}
func (gl GitlabClient) GetRepository(path string, name string, cloneType string) RepositoryInfo {
	project, _, err := gl.Client.Projects.GetProject(path+"/"+name, &gitlab.GetProjectOptions{})

	if err != nil {
		log.Fatalf("Failed to get repository %s: %v\n", path+"/"+name, err)
	}
	var repositoryLink string
	if cloneType == "ssh" {
		repositoryLink = project.SSHURLToRepo
	} else {
		repositoryLink = project.HTTPURLToRepo
	}
	return RepositoryInfo{repositoryLink, project.PathWithNamespace, project.DefaultBranch}
}
func GitlabAuth(service Service) GitClient {
	var git *gitlab.Client
	var err error
	if service.Username != "" && service.Password != "" {
		git, err = gitlab.NewBasicAuthClient(service.Username, service.Password, gitlab.WithBaseURL(service.BaseURL+service.APIURI))
	} else if service.APIToken != "" {
		git, err = gitlab.NewClient(service.APIToken, gitlab.WithBaseURL(service.BaseURL+service.APIURI))
	} else {
		git, err = gitlab.NewClient("", gitlab.WithBaseURL(service.BaseURL+service.APIURI))
	}
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}
	return GitlabClient{git}
}

// Bitbucket V1
func (bb BitbucketClient) GetGroupRepositories(path string, name string, recurse bool, cloneType string) []RepositoryInfo {
	var repositoriesInfo []RepositoryInfo
	var finalPath string
	if path == "" {
		finalPath = name
	} else {
		finalPath = path + "/" + name
	}
	response, err := bb.Client.DefaultApi.GetRepositories(finalPath)
	if err != nil {
		log.Fatalf("Failed to get repositories from project %s: %v\n", path+"/"+name, err)
	}
	repositories, err := bitbucketv1.GetRepositoriesResponse(response)
	if err != nil {
		log.Fatalf("Cannot get repositories response: %s\n", err)
	}
	for _, repository := range repositories {
		var repositoryLink string
		for _, links := range repository.Links.Clone {
			if cloneType == "http" && links.Name == "http" {
				repositoryLink = links.Href
			}
			if cloneType == "ssh" && links.Name == "ssh" {
				repositoryLink = links.Href
			}
		}
		defaultBranch, err := bb.Client.DefaultApi.GetDefaultBranch(repository.Project.Key, repository.Slug)
		if err != nil {
			log.Fatalf("Cannot get repository default branch: %s\n", err)
		}
		var defaultBranchString = "master"
		for key, value := range defaultBranch.Values {
			if key == "displayId" {
				defaultBranchString = fmt.Sprintf("%v", value.(string))
			}
		}
		repositoryInfo := RepositoryInfo{repositoryLink, repository.Slug, defaultBranchString}
		repositoriesInfo = append(repositoriesInfo, repositoryInfo)
	}
	return repositoriesInfo
}
func (bb BitbucketClient) GetRepository(path string, name string, cloneType string) RepositoryInfo {
	response, err := bb.Client.DefaultApi.GetRepository(path, name)
	if err != nil {
		log.Fatalf("Failed to get repositories from project %s: %v\n", path+"/"+name, err)
	}
	repository, err := bitbucketv1.GetRepositoryResponse(response)
	if err != nil {
		log.Fatalf("Cannot get repository reponse: %s\n", err)
	}
	var repositoryLink string
	for _, links := range repository.Links.Clone {
		if cloneType == "http" && links.Name == "http" {
			repositoryLink = links.Href
		}
		if cloneType == "ssh" && links.Name == "ssh" {
			repositoryLink = links.Href
		}
	}
	defaultBranch, err := bb.Client.DefaultApi.GetDefaultBranch(repository.Project.Key, repository.Slug)
	if err != nil {
		log.Fatalf("Cannot get repository default branch: %s\n", err)
	}
	var defaultBranchString = "master"
	for key, value := range defaultBranch.Values {
		if key == "displayId" {
			defaultBranchString = fmt.Sprintf("%v", value.(string))
		}
	}
	return RepositoryInfo{repositoryLink, strings.ToLower(repository.Project.Key) + "/" + repository.Slug, defaultBranchString}
}
func BitbucketAuth(service Service) GitClient {
	var git *bitbucketv1.APIClient
	cfg := bitbucketv1.NewConfiguration(service.BaseURL + service.APIURI)
	ctx, _ := context.WithTimeout(context.Background(), 6000*time.Millisecond)
	if service.Username != "" && service.Password != "" {
		auth := bitbucketv1.BasicAuth{UserName: service.Username, Password: service.Password}
		ctx = context.WithValue(ctx, bitbucketv1.ContextBasicAuth, auth)
	} else if service.APIToken != "" {
		auth := bitbucketv1.APIKey{Key: service.APIToken, Prefix: ""}
		ctx = context.WithValue(ctx, bitbucketv1.ContextAPIKey, auth)
	}
	git = bitbucketv1.NewAPIClient(ctx, cfg)
	return BitbucketClient{git}
}

//GitHub
func (gh GithubClient) GetGroupRepositories(path string, name string, recurse bool, cloneType string) []RepositoryInfo {
	var repositoriesInfo []RepositoryInfo
	opt := &github.RepositoryListByOrgOptions{Type: "public"}
	var finalPath string
	if path == "" {
		finalPath = name
	} else {
		finalPath = path + "/" + name
	}
	response, _, err := gh.Client.Repositories.ListByOrg(context.Background(), finalPath, opt)
	if err != nil {
		log.Fatalf("Cannot get repositories response: %s\n", err)
	}
	for _, repository := range response {
		var repositoryLink string
		if cloneType == "ssh" {
			repositoryLink = repository.GetSSHURL()
		} else {
			repositoryLink = repository.GetHTMLURL()
		}
		repositoryInfo := RepositoryInfo{repositoryLink, repository.GetFullName(), repository.GetDefaultBranch()}
		repositoriesInfo = append(repositoriesInfo, repositoryInfo)
	}
	return repositoriesInfo
}
func (gh GithubClient) GetRepository(path string, name string, cloneType string) RepositoryInfo {
	repository, _, err := gh.Client.Repositories.Get(context.Background(), path, name)
	if err != nil {
		log.Fatalf("Cannot get repository response: %s\n", err)
	}
	var repositoryLink string
	if cloneType == "ssh" {
		repositoryLink = repository.GetSSHURL()
	} else {
		repositoryLink = repository.GetHTMLURL()
	}
	return RepositoryInfo{repositoryLink, repository.GetFullName(), repository.GetDefaultBranch()}
}
func GithubAuth(service Service) GitClient {
	var git *github.Client
	if service.Username != "" && service.Password != "" {
		tp := github.BasicAuthTransport{
			Username: strings.TrimSpace(service.Username),
			Password: strings.TrimSpace(service.Password),
		}
		git = github.NewClient(tp.Client())
	} else if service.APIToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: service.APIToken},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		git = github.NewClient(tc)
	} else {
		git = github.NewClient(nil)
	}
	return GithubClient{git}
}

//Bitbucket V2
func (bb BitbucketClientV2) GetRepository(path string, name string, cloneType string) RepositoryInfo {
	project := bitbucketv2.Project{Name: path}
	repository := bitbucketv2.Repository{Slug: name, Project: project}
	repository2, err := repository.Get(&bitbucketv2.RepositoryOptions{})
	if err != nil {
		log.Fatalf("Cannot get repository response: %s\n", err)
	}
	var repositoryLink string
	for key, link := range repository2.Links {
		if key == "clone" {
			for _, v1 := range link.([]interface{}) {
				href := fmt.Sprintf("%v", v1.(map[string]interface{})["href"])
				name := fmt.Sprintf("%v", v1.(map[string]interface{})["name"])
				if cloneType == "http" && name == "https" {
					repositoryLink = href
				}
				if cloneType == "ssh" && name == "ssh" {
					repositoryLink = href
				}
			}
		}
	}
	return RepositoryInfo{repositoryLink, strings.ToLower(repository.Project.Key) + "/" + repository.Slug, "master"}
}
func (bb BitbucketClientV2) GetGroupRepositories(path string, name string, recurse bool, cloneType string) []RepositoryInfo {
	var repositoriesInfo []RepositoryInfo
	response, err := bb.Client.Repositories.ListForAccount(&bitbucketv2.RepositoriesOptions{Owner: name})
	if err != nil {
		log.Fatalf("Cannot get repositories response: %s\n", err)
	}
	for _, repository := range response.Items {
		var repositoryLink string
		for key, link := range repository.Links {
			if key == "clone" {
				for _, v1 := range link.([]interface{}) {
					href := fmt.Sprintf("%v", v1.(map[string]interface{})["href"])
					name := fmt.Sprintf("%v", v1.(map[string]interface{})["name"])
					if cloneType == "http" && name == "https" {
						repositoryLink = href
					}
					if cloneType == "ssh" && name == "ssh" {
						repositoryLink = href
					}
				}
			}
		}
		repositoryInfo := RepositoryInfo{repositoryLink, repository.Slug, "master"}
		repositoriesInfo = append(repositoriesInfo, repositoryInfo)
	}
	return repositoriesInfo
}
func BitbucketV2Auth(service Service) GitClient {
	var git *bitbucketv2.Client
	if service.Username != "" && service.Password != "" {
		git = bitbucketv2.NewBasicAuth(service.Username, service.Password)
	} else if service.APIToken != "" {
		log.Fatalf("Bitbucket V2 token authentication not implemented\n")
	} else {
		log.Fatalf("Bitbucket V2 without authentication not implemented\n")
	}

	git.SetApiBaseURL(service.BaseURL + service.APIURI)
	return BitbucketClientV2{git}
}
