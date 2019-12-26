# OpenOLT Scale Tester
This is used for scale testing of OpenOLT Agent + BAL

# Design
The proposed design is here https://docs.google.com/document/d/1hQxbA8jvG1BEHeeLkM5L3sYrvgJUX7z2Tk7j1sPdd1s

# How to build
```shell
make build
```

# How to run
Make sure openolt-agent and dev_mgmt_daemon are running. Then run the below command from `openolt-scale-tester` folder.

```shell
DOCKER_HOST_IP=<your-host-ip-here> docker-compose -f compose/openolt-scale-tester.yml up -d
```
