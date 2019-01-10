# api-routerd
![N|Solid](https://ibin.co/4R6Hzr2H7l4A.png)

A RestAPI MicroService GateWay for Linux

[![Build Status CircleCI](https://circleci.com/gh/RestGW/api-routerd.svg?style=svg)](https://circleci.com/gh/RestGW/api-routerd)
[![Build Status](https://travis-ci.org/RestGW/api-routerd.svg?branch=master)](https://travis-ci.org/RestGW/api-routerd)
[![HitCount](http://hits.dwyl.io/ssahani/RestGW/api-routerd.svg)](http://hits.dwyl.io/ssahani/RestGW/api-routerd)
[![CodeFactor](https://www.codefactor.io/repository/github/restgw/api-routerd/badge)](https://www.codefactor.io/repository/github/restgw/api-routerd)
[![codebeat badge](https://codebeat.co/badges/1bdd48c6-4cc1-4255-a11b-9807473e9c3d)](https://codebeat.co/projects/github-com-restgw-api-routerd-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/RestGW/api-routerd)](https://goreportcard.com/report/github.com/RestGW/api-routerd)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/5094d66b86f44522b3de7381a0bba5a1)](https://app.codacy.com/app/ssahani/api-routerd?utm_source=github.com&utm_medium=referral&utm_content=RestGW/api-routerd&utm_campaign=Badge_Grade_Dashboard)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Coverage Status](https://coveralls.io/repos/github/RestGW/api-routerd/badge.svg?branch=master)](https://coveralls.io/github/RestGW/api-routerd?branch=master)


api-routerd is a cloud-enabled, mobile-ready, a super light weight remote management tool which uses REST API for real time configuration and performance as well as health monitoring for systems (containers) and applications. It provides fast API based monitoring without affecting the system it's running on.

- Platform independent REST APIs can be accessed via any application (curl, chrome, PostMan ...) from any OS (Linux, IOS, Android, Windows ...)
- An [Iphone App Demo](https://www.linkedin.com/feed/update/urn:li:activity:6486243669560127488) using REST APIS
- Minimal data transfer using JSON.
- Plugin based Architechture. See how to write plugin section for more information.

# Features!

- systemd socket activation support
- systemd
  - systemd informations
  - services (start, stop, restart, status)
  - service properties for example CPUShares
  - See service logs.
- networkd config
  - .network
  - .netdev
  - .link
- configure hostnamed
- configure users using login (list-sessions, list-users and terminate-user etc)
- configure timdate
- configure nameserver ```/etc/resolv.conf```
- configure timesynd
- configure journald.conf
- configure system.conf
- configure coredump.conf
- configure systemd-resolved.conf
- configure kernel modules (modprobe, lsmod, rmmod)
- configure network (netlink)
  - Link: mtu, up, down
  - Create bridge and enslave links
  - Adddress: Set, Get, Delete
  - Gateway: Default Gateway Add and Delete
- configure group add/delete/modify
- configure users add/delete/modify (requires newuser)
- configure sysctl add/delete/modify and apply
- see information from /proc such as netstat, netdev, memory and much more
- configure ```/proc/sys/net``` (core/ipv4/ipv6), VM
- See ethtool information and configure offload features
- See sudoers and sshd conf

### api-routerd JSON APIs

 Refer spreadsheet [APIs](https://docs.google.com/spreadsheets/d/e/2PACX-1vTl2Vmp-BdTE5Vgi_PiW-qKPJnbLxdSso9kT2GAkAxCu_iWrw3_PZLlEuyXz0lbFgd7DoofXlmmb3dP/pubhtml
)

### Tech

api-routerd uses a number of open source projects to work properly:

* [logrus](github.com/sirupsen/logrus)
* [gorilla mux](github.com/gorilla/mux)
* [netlink](github.com/vishvananda/netlink)
* [gopsutil](github.com/shirou/gopsutil)
* [coreos go-systemd](github.com/coreos/go-systemd)
* [dbus](github.com/godbus/dbus)
* [ethtool](github.com/safchain/ethtool)
* [BurntSushi toml](github.com/BurntSushi/toml)
* [go-ini](github.com/go-ini/ini)


And of course api-routerd itself is open source with a [public repository][git-repo-url]
 on GitHub.

### Development

Want to contribute? Great!

### Installation

First configure your ```$GOPATH```. If you have already done this skip this step.

```sh
# keep in ~/.bashrc
```

```sh
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
export OS_OUTPUT_GOPATH=1
```

Clone inside src dir of ```$GOPATH```. In my case

```sh
$ pwd
/home/sus/go/src
```

### Installation

```sh
$ go get github.com/RestGW/api-routerd
$ git clone https://github.com/RestGW/api-routerd
$ cd api-routerd
$ pwd
/home/sus/go/src/api-routerd
$ go build -v

$ sudo ./api-routerd
INFO[0000] api-routerd: v0.1 (built go1.11.4)
INFO[0000] Start Server at 0.0.0.0:8080
INFO[0000] Starting api-routerd in plain text mode

```

### How to configure IP and Port ?

Conf dir: ```/etc/api-routerd/```
Conf File: ```api-routerd.toml```

```sh
$ cat /etc/api-routerd/api-routerd.toml
[Network]
IPAddress="0.0.0.0"
Port="8080"
```

### How to configure users ?

Add user name and authentication string in space separated lines

```sh
# cat /etc/api-routerd/api-routerd-auth.conf
Susant secret
Max bbbb
Joy ccccc
```

### How to configure TLS ?

Generate private key (.key)

```sh
# Key considerations for algorithm "RSA" ≥ 2048-bit
$ openssl genrsa -out server.key 2048
Generating RSA private key, 2048 bit long modulus (2 primes)
.......................+++++
.+++++
e is 65537 (0x010001)

openssl genrsa -out server.key 2048
```

Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)

```sh
$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:

```

Place ```server.crt``` and ```server.key``` in the dir ```/etc/api-routerd/tls```

```sh
[root@Zeus tls]# ls
server.crt  server.key
[root@Zeus tls]# pwd
/etc/api-routerd/tls

```

Use case: https

```sh
$ curl --header "X-Session-Token: secret" --request GET https://localhost:8080/api/network/ethtool/vmnet8/get-link-features -k --tlsv1.2

```
## Use cases

Refer usecase document [use cases](https://github.com/RestGW/api-routerd/blob/master/examples.md)

## How to write your own plugin ?

api-routerd is designed with plugin based architecture in mind and can act as a thin client. You can always add and remove modules to it with minimal effort

* Choose namespace under cmd directory ( network, proc, system etc) whare you want to put your module.
* Write sub router see for example ```api-routerd/cmd/system/login```
* Write your module ```module.go``` and  ```module_router.go```
* Write ```RegisterRouterModule```
* Register ```RegisterRouterModule``` with parent router for example for ```login``` registered with
  ```RegisterRouterSystem``` under ```system``` namespace as ```login.RegisterRouterLogin(n)```

### Todos

 - Write Tests
 - Networkd
 - iptables

License
----

Apache 2.0


**Free Software, Hell Yeah!**

[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen. Thanks SO - http://stackoverflow.com/questions/4823468/store-comments-in-markdown-syntax)

   [git-repo-url]: <git@github.com:RestGW/api-routerd.git>
