# Test

In this section we are going to test our ```worksync``` library by
creating a master/slave topology system. In the master pod we are going
to create 2 types of jobs (aka works). Moreover, in the slave pods we are
going to create a ```NewSyncer``` and do master jobs.

## master

In master, we are going to get a config file in order to create jobs.
The master configs are as below:

```yaml
database:
    dns: ""
jobs:
    - name: "remove orders"
      cron: "@every 5m"
      max_run: 1
    - name: "emails"
      cron: "* 1 * * *"
      max_run: 24
```

## slave

In slave pod, we are going to get jobs that are needed to be done
by that slave. The slaves configs are as below:

```yaml
database:
    dns: ""
job:
    name: "emails"
    cron: "@every 5m"
```

## docker images

In order to test ```worksync``` you can use the following images:

```shell
docker pull amirhossein21/workersync:v0.1-slave
docker pull amirhossein21/workersync:v0.1-master
```

### execute

First make sure to create a ```config.yaml``` file for container configs. Master and slave
have different config types. Make sure to create suitable config files for each based on
the templates in their directory.

```shell
docker run -d -v type=bind,source=$(pwd)/config.yaml,dest=/src/config.yaml amirhossein21/workersync:v0.1-master
```

## kubernetes

Use the manifests in [```k8s```](./k8s) directory in order to test ```worksync``` library on kubernetes
cluster.
