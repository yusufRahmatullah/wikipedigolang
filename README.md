# PG1-Go-Work

This is the code base of this [heroku app](http://pg1-go-work.herokuapp.com/igprofiles). PG1 means Playground-one because this application is a sandbox application to improve my programming skill. Go means this application is written using Golang. Work means this application contains worker app which runs in the background.

## Introduction
 
This application contains two processes, web and worker. The web is the sub-app for handling HTTP request from user meanwhile the worker is the sub-app for handling background process. This app is written in Golang using Gin framework and persists the data using MongoDB. This app is deployable on Heroku.

## Installation 

I develop this app using Golang version 1.11.2 for Linux/amd64. To run this app, we should install golang first. Golang dependencies are using dep. Install dep and run `dep ensure` to ensure the dependencies are ready. The vendor directory is not registered on .gitignore file. So, the dependencies automatically exists.

Because of this app is configured to be deployable in Heroku, don't forget to install heroku-cli too. Heroku is required to read the environment variables.

## Run in Local

First, ensure the MongoDB server has been run and the .env has been set. Copy `env.sample` and set the environment variables.

```bash
$ cp env.sample .env
```

Another dependecies is PhantomJS. Download PhantomJS executable because the worker uses PhantomJS to crawl the IG Profiles.

To run the app, just call the `./run_local.sh` script (ensure to enable execution on script run_local.sh).

## Deploy to Heroku

Read how to deploy an apps in Heroku. This app requires two buildpacks, `heroku/go` buildpack and PhantomJS buildpack. Add PhantomJS buildpack by run this command

```bash
$ heroku buildpacks:set  https://github.com/stomita/heroku-buildpack-phantomjs.git
```
