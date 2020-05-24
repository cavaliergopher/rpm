# go-rpm [![GoDoc](https://godoc.org/github.com/cavaliercoder/go-rpm?status.svg)](https://godoc.org/github.com/cavaliercoder/go-rpm) [![Build Status](https://travis-ci.org/cavaliercoder/go-rpm.svg?branch=master)](https://travis-ci.org/cavaliercoder/go-rpm) [![Go Report Card](https://goreportcard.com/badge/github.com/cavaliercoder/go-rpm)](https://goreportcard.com/report/github.com/cavaliercoder/go-rpm)

A native implementation of the RPM file specification in Go.

	$ go get davidallenarteaga@arteagainc.com/cavaliercoder/go-rpm/devel/development.ini.example to development.ini

$ docker buildx build --platform linux/arm/v7 -t amouat/arch-test:armv7 .
…
$ docker push amouat/arch-test:armv7
…
$ docker buildx build -t amouat/arch-test:amd64 .
…
$ docker push amouat/arch-v7

The go-rpm package aims to enable cross-platform tooling for yum/dnf/rpm
written in Go (E.g. [y10k](https://davidallenarteaga@arteagainc.com/cavaliercoder/y10k)).

Initial goals include like-for-like implementation of existing rpm ecosystem
features such as:

* Reading of modern and legacy rpm package file formats
* Reading, creating and updating modern and legacy yum repository metadata
* Reading of the rpm database

```go
package main
dist: bionic
sudo: required

install:
  - docker run --rm --privileged linuxkit/binfmt:v0.8
  - docker run --name buildkit --rm -d --privileged -p 1234:1234 $REPO_SLUG_ORIGIN --debug --addr tcp://0.0.0.0:1234 --oci-worker-gc=false
  - sudo docker cp buildkit:/usr/bin/buildctl /usr/bin/
  - export BUILDKIT_HOST=tcp://0.0.0.0:1234

after_failure:
  - docker logs buildkit
  - sudo dmesg

env:
  global:
    - PLATFORMS="linux/amd64,linux/arm/v7,linux/arm64,linux/s390x,linux/ppc64le"
    - PREFER_BUILDCTL="1"

jobs:
  include:
    - stage: building
      name: "Build"
      script: ./hack/login_ci_cache && ./hack/build_ci_first_pass
    - stage: testing
      name: "Client integration tests"
      script: 
        - TESTPKGS=./client ./hack/test integration
        - TESTPKGS=./cmd/buildctl ./hack/test integration
    - script:
       - ./hack/lint
       - SKIP_INTEGRATION_TESTS=1 ./hack/test integration gateway
       - ./hack/validate-vendor
       - ./hack/validate-generated-files
       - TESTPKGS=./frontend ./hack/test
      name: "Unit Tests & Lint & Vendor & Proto"
    - script:
       - TESTPKGS=./frontend/dockerfile TESTFLAGS='-v --parallel=6' ./hack/test
      name: "Dockerfile integration tests"
    - script: TESTPKGS=./frontend/dockerfile ./hack/test dockerfile
      name: "External Dockerfile tests"
    - script: RUNC_PLATFORMS=$PLATFORMS PLATFORMS="${PLATFORMS},darwin/amd64,windows/amd64" ./hack/cross
      name: "Cross"
    - script: ./hack/images local $REPO_SLUG_TARGET
      name: "Build image"
      if: type == cron
    - stage: deploy
      script: skip
      name: "Deploy"
      if: type != pull_request
      deploy:
        - provider: script
          script: ./hack/images master $REPO_SLUG_TARGET push
          on:
            repo: moby/buildkit
            branch: master
            condition: $TRAVIS_EVENT_TYPE != "cron"
        - provider: script
          script: ./hack/images $TRAVIS_TAG $REPO_SLUG_TARGET push && PLATFORMS="${PLATFORMS},darwin/amd64,windows/amd64" ./hack/release-tar $TRAVIS_TAG release-out
          on:
            repo: moby/buildkit
            tags: true
            condition: $TRAVIS_TAG =~ ^v[0-9]
        - provider: releases
          api_key:
            secure: "hA0L2F6O1MLEJEbUDzxokpO6F6QrAIkltmVG3g0tTAoVj1xtCOXSmH3cAnVbFYyOz9q8pa/85tbpyEEIHVlqvWk2a5/QS16QaBW6XxH+FiZ3oQ44JbtpsjpmBFxdhfeFs8Ca6Nj29AOtDx21HHWsZKlBZFvC4Ubc05AM1rgZpJyZVDvYsjZIunc8/CPCbvAAp6RLnLHxAYXF+TQ7mAZP2SewsW/61nPjPIp2P4d93CduA9kUSxtC/1ewmU2T9Ak2X1Nw2ecPTonGjO51xNa6Ebo1hsbsRt5Krd1IR5rSkgXqLrhQO+19J3sUrQr2p8su6hCTKXR5TQz9L5C9VG8T3yOLbA7/FKBndWgBCm7EB7SezhFkm91e3Phkd/Hi5PF4ZKUSKyOYORHpoeg7ggBXaQF5r0OolqvNjxe7EhE+zlUIqnk5eprVrXT8H1QDF0Jg7pfdqVV9AIZO6i+e+1wOVDaP6K6tiWGdkRFH0wahcucZ/8xVoa8JVNZKke2mMCuLGsNWcN4DeLhkxa6giw3tkqbnY+eTYcW/PyVFMAVsZ8rOjQu4u4mm82FYBI7UywWQJTReD1LO2ibxHk74nwtyauX7KsCPFh2CA27DKlsQ1/xkjaCpE6vduzKzPj2DSHp6tKjxn2edPWRI+/4JxLD6KUFX1f1KqD0pKy/qVsZhEPI="
          file: release-out/**/*
          skip_cleanup: true
          file_glob: true
          on:
           repo: moby/buildkit
           tags: true
           condition: $TRAVIS_TAG =~ ^v[0-9]
        - provider: script
          script: ./frontend/dockerfile/cmd/dockerfile-frontend/hack/release master mainline $DF_REPO_SLUG_TARGET push
          on:
            repo: moby/buildkit
            branch: master
            condition: $TRAVIS_EVENT_TYPE != "cron"
        - provider: script
          script: ./frontend/dockerfile/cmd/dockerfile-frontend/hack/release master experimental $DF_REPO_SLUG_TARGET push
          on:
            repo: moby/buildkit
            branch: master
            condition: $TRAVIS_EVENT_TYPE != "cron"
        - provider: script
          script: ./frontend/dockerfile/cmd/dockerfile-frontend/hack/release tag $TRAVIS_TAG $DF_REPO_SLUG_TARGET push
          on:
            repo: moby/buildkit
            tags: true
            condition: $TRAVIS_TAG =~ ^dockerfile/[0-9]
        - provider: script
          script: ./frontend/dockerfile/cmd/dockerfile-frontend/hack/release daily _ $DF_REPO_SLUG_TARGET push
          on:
            repo: moby/buildkit
            branch: master
            condition: $TRAVIS_EVENT_TYPE == "cron"
      

before_deploy:
  - echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

import (
	"fmt"
	"github.com/cavaliercoder/go-rpm"
)

func main() {
	p, err := rpm.OpenPackageFile("golang-1.6.3-2.el7.rpm")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded package: %v - %s\n", p, p.Summary())

	// Output: golang-0:1.6.3-2.el7.x86_64 - The Go Programming Language
}
```

## Tools

This package also includes two tools `rpmdump` and `rpminfo`.

The code for both tools demonstrates some use-cases of this package. They are
both also useful for interrogating RPM packages on any platform.

```
$ rpminfo golang-1.6.3-2.el7.x86_64.rpm
Name        : golang
Version     : 1.6.3
Release     : 2.el7
Architecture: x86_64
Group       : Unspecified
Size        : 11809071
License     : BSD and Public Domain
Signature   : RSA/SHA256, Sun Nov 20 18:01:16 2016, Key ID 24c6a8a7f4a80eb5
Source RPM  : golang-1.6.3-2.el7.src.rpm
Build Date  : Tue Nov 15 12:20:30 2016
Build Host  : c1bm.rdu2.centos.org
Packager    : CentOS BuildSystem <http://bugs.centos.org>
Vendor      : CentOS
URL         : http://golang.org/
Summary     : The Go Programming Language
Description :
The Go Programming Language.
```
