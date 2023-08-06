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
      max_run: 0
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

