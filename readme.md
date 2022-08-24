# semver

> semantic versions in pipelines


`semver` was created to serve as reliable pipeline helper. Born out of
necessity and distrust towards regex.


## Installation

```
go install github.com/krysopath/semver/cmd/semver@v1
```

or

```
GO111MODULE=on go get github.com/krysopath/semver/cmd/semver@v1
```

## Usage

Parse tags from git:
```
$ git tag | head -n1 | semver | jq
{
  "canonical": "v0.2.0",
  "major": "v0",
  "majorminor": "v0.2",
  "prerelease": "",
  "build": "",
  "source": "v0.2.0"
}
```

Parse semantic version from stdin:
```
$ echo "v0.1.23-alpha2+9999" | semver | jq
{
  "build": "+9999",
  "canonical": "v0.1.23-alpha2",
  "major": "v0",
  "majorminor": "v0.1",
  "prerelease": "-alpha2"
}
```


When run inside Gitlab Runners, can execute without stdin:
```
$ export CI_COMMIT_TAG=v1.22.4-some+123
$ semver | jq
{
  "build": "+123",
  "canonical": "v1.22.4-some",
  "major": "v1",
  "majorminor": "v1.22",
  "prerelease": "-some"
}
```

When you dont like json, but shell:
```
$ eval $(echo v3.2.1-yolo22 | ./semver -format=eval)
$ echo $MAJORMINOR
v3.2
```


## In pipelines

```
# lets use the great jq to create nobrain parseable 
export SEMVER="$(echo $CI_COMMIT_TAG | semver)"

# lets record the -prerelease+build as SUFFIX
export SUFFIX="$(echo $SEMVER | jq -r .prerelease)$(echo $SEMVER | jq -r .build)"

# lets create the semvers, appending SUFFIX (for edgecase, you know)
export MAJORMINOR="$(echo $SEMVER | jq -r .majorminor)$SUFFIX"
export MAJOR="$(echo $SEMVER | jq -r .major)$SUFFIX"
export CANONICAL="$(echo $SEMVER | jq -r .canonical)"

# tag one distinct artifact with these three versions
docker tag $CI_IMAGE_TAG "$CI_REGISTRY_IMAGE:$MAJORMINOR"
docker tag $CI_IMAGE_TAG "$CI_REGISTRY_IMAGE:$MAJOR"
docker tag $CI_IMAGE_TAG "$CI_REGISTRY_IMAGE:$CANONICAL"
```

> docker e.g. as standin for any artifact


> Above snippet can be used as is in gitlab
