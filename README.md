
### Some Docker Basics

#### get container PID from host point of view
``` bash
docker inspect -f '{{.State.Pid}}' b7b9bbce4a8e
> 7053

```

#### get PID inside of container namespace

``` bash
more /proc/7053/status | grep "NSpid" # 7053 is the container PID from host point of view
> NSpid:	7053	1                 # inside of container namepsaces. 7053 mapped to 1
```


### Namespaces

#### UTS Namespace

UTS Namespace is used to isolate node name and domain name
