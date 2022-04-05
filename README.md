# CI-KUBERNETES

Used to run in CI pipeline to perform actions on Kubernetes cluster

## Variables

`KUBE_CONFIG` - Base64 encoded kubernetes config file
`NAMESPACE` - Kubernetes namespace to run on
`ACTION` - Action to perform, `update`
`KIND` - Kind of object to update, `cronjob` `deployment`
`NAME` - Name of object to update
`IMAGE_NAME` - Name of image to use
`IMAGE_TAG` - Tag of image to use

Example:

```yaml
deploy-release:
  image: deyung/drone-kubernetes:latest
  stage: deploy
  only:
    - tags
  variables:
    NAMESPACE: development
    ACTION: update
    KIND: deployment
    NAME: $CI_PROJECT_NAME
    IMAGE_NAME: $CI_REGISTRY_IMAGE
    IMAGE_TAG: ${CI_COMMIT_TAG}
  script:
    - /app/ci-kubernetes
```
