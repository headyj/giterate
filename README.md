# Project
Giterate is a wrapper that allows you to clone and pull multiple repositories from multiple providers with a single command

## Configuration
By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate folder
Each git provider will have the same available parameters:

- BaseURL: Base URL of the git provider
- API: Type of api (can be gitlab, bitbucketv1, bitbucketv2 or github)
- ApiURI: URI of the api
- ApiToken: API token (required with ssh authentication)
- CloneType: Type of authentication (can be ssh or http)
- Username: username for authentication
- Password: password for authentication
- SSHPrivateKeyPath: absolute path to private key (required with ssh authentication)
- Destination: default destination directory
- Entities: List of repositories and groups to clone
    - Type: Type of entity, can be repository or group
	- Path: path to the entity (e.g. group, project, etc)
	- Name: name to the entity (e.g. name of the group, name of the repository)
	- Destination: destination absolute path. If not given, will take the default destination + path of entity
    - Recurse: in case of group, will clone recursively
    - CloneOptions: clone options array
        - Key: Name of the clone option
        - Value: Value of the clone option

You can find an example of configuration file on this repository

## Available services
- Gitlab (https://www.github.com/xanzy/go-gitlab)
- Bitbucket V1 (https://www.github.com/gfleury/go-bitbucket-v1)
- Bitbucket V2 (https://www.github.com/ktrysmt/go-bitbucket)
- Github (https://www.github.com/google/go-github/github)

## Available commands
- [x] clone: clone repositories according to configuration file
    - if the repository already exists or is already clone, it will not be updated
    - parameters:
        - [ ] --force: will clean all an recreate from conf

- [x] pull: pull repositories on current branches according to configuration file
    - if a new repository has been added to the configuration/to the git provider, it will not be cloned

- [x] status: check status of each git repositories according to configuration file
    - parameters
        - [x] --full: show status of all repositories, even if there's no uncommited changes

- [x] checkout: checkout the configured/default branch on all repositories
    - parameters
        - [x] --force: reset uncommited changes

- [x] commit: check changes and ask for commit message in case of changes
    - if you don't provide any message, it will go to the next one without commiting
    - parameters
        - [ ] --target: target one or multiple repositories
        - [ ] -g, --global: define a single message for all commits (you'll have to answer "yes" instead of providing a message)

- [x] push: push commited changes
    - if you don't provide any message, it will go to the next one without commiting
    - parameters
        - [ ] --force: push without asking
        - [ ] --target: target one or multiple repositories

## Global parameters
- [x] --config-file: set json/yaml configuration file path
- [x] --log-level: set log level ("info", "warn", "error", "debug"). default: "info"

## Roadmap
- implement tests
- implement parameters
