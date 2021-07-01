# semver

> semantic versions in pipelines


`semver` was created to serve as reliable pipeline helper. Born out of
necessity and distrust towards regex.


## Installation

```
go get github.com/krysopath/semver@v0.1.2
```

## Usage

Parse tags from git:
```
$ git tag | ./semver -sort | jq .[-1]
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

Parse many versions and sort them:
```
$ echo 'v3 v1.1.1-pre0+999 v3.1.1-dest0 v2 v3.1.1' | semver -sort | jq
[
  {
    "canonical": "v1.1.1-pre0",
    "major": "v1",
    "majorminor": "v1.1",
    "prerelease": "-pre0",
    "build": "+999",
    "source": "v1.1.1-pre0+999"
  },
  {
    "canonical": "v2.0.0",
    "major": "v2",
    "majorminor": "v2.0",
    "prerelease": "",
    "build": "",
    "source": "v2"
  },
  {
    "canonical": "v3.0.0",
    "major": "v3",
    "majorminor": "v3.0",
    "prerelease": "",
    "build": "",
    "source": "v3"
  },
  {
    "canonical": "v3.1.1-dest0",
    "major": "v3",
    "majorminor": "v3.1",
    "prerelease": "-dest0",
    "build": "",
    "source": "v3.1.1-dest0"
  },
  {
    "canonical": "v3.1.1",
    "major": "v3",
    "majorminor": "v3.1",
    "prerelease": "",
    "build": "",
    "source": "v3.1.1"
  }
]
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
