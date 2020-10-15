# Giterate

 1. [Overview](#overview)
 1. [Installation](#installation)
 1. [Configuration](#configuration)
 1. [Usage](#usage)
 1. [Examples](#examples)
 1. [Roadmap](#roadmap)

## Overview <a name="overview"></a>
Giterate is a wrapper that will help you to easily deal with multiple git repositories from multiple sources.

With Giterate, you can:
- recursively clone repositories accross gitlab groups / bitbucket projects according to JSON or YAML configuration file
- execute almost all git commands on all (or filtered subset of) repositories
- execute custom git commands

## Installation <a name="installation"></a>
### Binaries
- Download the [latest release](https://github.com/headyj/giterate/releases) binary corresponding to your system
- <ins>For Linux</ins>: put the binary somewhere on your path (example: /usr/bin)
- <ins>For Windows</ins>: put the binary somewhere on your computer (example: C:/giterate) and [update the path variable](https://www.architectryan.com/2018/03/17/add-to-the-path-on-windows-10/)

## Configuration <a name="configuration"></a>
By default, giterate will use config.json or config.yaml file (in this order) in ~/.giterate folder
Each git provider will have the same available parameters:

- **BaseURL**: Base URL of the git provider
- **Name**: Name of the git provider (easier to filter with --provider _base URL or name_)
- **API**: Type of api (can be gitlab, bitbucketv1, bitbucketv2 or github)
- **ApiURI**: URI of the api
- **ApiToken**: API token (required with ssh authentication)
- **CloneType**: Type of authentication (can be ssh or http)
- **Username**: username for authentication
- **Password**: password for authentication
- **SSHPrivateKeyPath**: absolute path to private key (required with ssh authentication)
- **Destination**: default destination directory
- **CloneOptions**: global clone options array
    - **Key**: Name of the clone option
    - **Value**: Value of the clone option
- **Entities**: List of repositories and groups to clone
    - **Type**: Type of entity, can be repository or group
	- **Path**: path to the entity (e.g. group, project, etc)
	- **Name**: name to the entity (e.g. name of the group, name of the repository)
	- **Destination**: destination absolute path. If not given, will take the default destination + path of entity
    - **Recurse**: in case of group, will clone recursively
    - **CloneOptions**: clone options array
        - **Key**: Name of the clone option
        - **Value**: Value of the clone option

You can find an example of configuration file on this repository

## Usage <a name="usage"></a>
### Available commands
- [x] **clone**: clone repositories according to configuration file
    - if the repository already exists or is already clone, it will not be updated
    - parameters:
        - [ ] **--force**: will clean all an recreate from conf
        - [x] **-r, --repository _URL or path_**: target one or multiple repositories (chain multiple times)
        - [x] **-p, --provider _base URL or name_**: target one or multiple providers (chain multiple times)

- [x] **pull**: pull repositories on current branches according to configuration file
    - parameters
        - [x] **-r, --repository _URL or path_**: target one or multiple repositories (chain multiple times)
        - [x] **-p, --provider _base URL or name_**: target one or multiple providers (chain multiple times)

- [x] **fetch**: fetch repositories on current branches according to configuration file
    - parameters
        - [x] **-r, --repository _URL or path_**: target one or multiple repositories (chain multiple times)
        - [x] **-p, --provider _base URL or name_**: target one or multiple providers (chain multiple times)

- [x] **status**: check status of each git repositories according to configuration file
    - parameters
        - [x] -**f, --full**: show status of all repositories, even if there's no uncommited changes
        - [x] **-r, --repository _URL or path_**: target one or multiple repositories (chain multiple times)
        - [x] **-p, --provider _base URL or name_**: target one or multiple providers (chain multiple times)

- [x] **checkout**: checkout the configured/default branch on all repositories
    - parameters
        - [x] **--force**: reset uncommited changes
        - [x] **-r, --repository _URL or path_**: target one or multiple repositories (chain multiple times)
        - [x] **-p, --provider _base URL or name_**: target one or multiple providers (chain multiple times)

- [x] **commit**: check changes and ask for commit message in case of changes
    - if you don't provide any message, it will go to the next one without commiting
    - parameters
        - [ ] **-g, --global**: define a single message for all commits (you'll have to answer "yes" instead of providing a message)
        - [x] **-r, --repository _URL or path_**: target one or multiple repositories (chain multiple times)
        - [x] **-p, --provider _base URL or name_**: target one or multiple providers (chain multiple times)

- [x] **push**: push commited changes
    - if you don't provide any message, it will go to the next one without commiting
    - parameters
        - [ ] **--force**: push without asking
        - [x] **-r, --repository _URL or path_**: target one or multiple repositories (chain multiple times)
        - [x] **-p, --provider _base URL or name_**: target one or multiple providers (chain multiple times)

- [x] **providers**: list configured providers

- [x] **repositories**: list configured repositories
    - parameters
        - [x] **-p, --provider _base URL or name_**: target one or multiple providers (chain multiple times)

- [x] **exec**: execute a custom git command
    - parameters
        - [x] **-c, --command  _'command'_** : command to be executed
        - [x] **-r, --repository _URL or path_**: target one or multiple repositories (chain multiple times)
        - [x] **-p, --provider _base URL or name_**: target one or multiple providers (chain multiple times)

### Global parameters
- [x] **--config-file**: set json/yaml configuration file path
- [x] **--log-level**: set log level (info, warn, error, debug). default: info

### Available services
- Gitlab (https://www.github.com/xanzy/go-gitlab)
- Bitbucket V1 (https://www.github.com/gfleury/go-bitbucket-v1)
- Bitbucket V2 (https://www.github.com/ktrysmt/go-bitbucket)
- Github (https://www.github.com/google/go-github/github)

## Examples <a name="examples"></a>
Basic clone command
```bash
giterate clone
```

Use alternative config file
```bash
giterate clone --config-file ~/giterate-config.json
```

Pull from a specific provider name from the configuration file
```bash
giterate pull --provider gitlab-mycompany
```

Get the status of a specific provider URL from the configuration file
```bash
giterate status --provider https://gitlab.mycompany.com

Repository: https://gitlab.mycompany.com (/home/usr/giterate/mycompany)
Branch: master
Changes:
MM crypted.auto.tfvars
 M directory/file.json
 M directory2/other_file.php

Repository: https://gitlab.mycompany.com (/home/usr/giterate/mycompany)
Branch: develop
Changes:
 M file.txt
MM other_file.json
```

Commit
```bash
giterate commit

Repository: ssh://git@gitlab.*****.com:7999/path/repo.git (/home/usr/giterate/gitlab/repo)
Branch: production
Changes:
 M full/path/to/file.pp

Enter commit message (let empty to ignore): update puppet configuration for repo


Repository: https://git.****.com/path/other_repo.git (/home/usr/giterate/gitlab/other_repo)
Branch: master
Changes:
 M full/path/to/file.json

Enter commit message (let empty to ignore): update json configuration for other_repo
```

Execute a custom command on a subset of repositories
```bash
giterate exec -c 'reset --hard' -r /home/usr/giterate/bitbucketv1/repo1 -r /home/usr/giterate/bitbucketv1/repo2
```

## Roadmap <a name="roadmap"></a>
- implement tests
- implement parameters
