<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/github_username/repo_name">
    <img src="logo.png" alt="Logo" width="125" height="125">
  </a>
<div align="center">
  <h3 align="center">KUclap Back-End</h3>

![Deploy-to-DO-production](https://github.com/KUclap/KUclap-review/workflows/Deploy-to-DO-production/badge.svg?branch=release&event=push)

</div>
</p>

<!-- ABOUT THE PROJECT -->

## Overview

This repository is KUclap back-end source code which is written in golang for implementing a web API and using it access to the database 🚀.

### Documentation

Postman Collection 📝
URL : https://www.getpostman.com/collections/79cb50bda1b010277ac9

### Built With 🔧

-   [Golang](https://golang.org/)
-   [mgo.v2](https://godoc.org/gopkg.in/mgo.v2)

## Getting Started

To get a local copy up and running follow these simple steps 🎉.

### Prerequisites

Install these prerequisites ✅ .

-   Go
-   Docker

### Installation

1. Clone the repo

```sh
git clone https://github.com/KUclap/KUclap-review.git
```

2. Install packages

```sh
go get ./...
```

<!-- USAGE EXAMPLES -->

## Development / Usage

Use `modd` for live reloading by follow this command 😎 .

```sh
make gomodd
```

## Deployment

### Staging

This command is for deploying to Heroku 🤒 (Stagging Environment).

```sh
make deploy-to-staging
```

### Pre-Production

Merge commits from master into pre-prod-release branch. The pipeline will deploy to Gandalf's server (DigitalOcean droplet) automatically 🤮.

```sh
git checkout pre-prod-release
git pull origin master
...
git push
```

### Production

Like Pre-Production 😬, Merge commits from master into release branch. The pipeline will deploy to Gandalf's server (DigitalOcean droplet) automatically 😳 .

```sh
git checkout release
git pull origin master
...
git push
```

## Note 🌶

-   For more details about commands, Please read `Makefile`.
-   `.github/workflows` is used for storing pipeline script for automated deployment.
-   You have to install `heroku` CLI for deploying image to staging.
-   Heroku only detects docker image (Dockerfile) which filename starting with 'D' capital letter.
