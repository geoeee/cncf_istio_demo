# CNCF-istio-demo

Script and config repo for istio demo in CNCF 2018 China

- [Prerequisites](#prerequisites)
- [Demo step](#demo-step)
    - [Weight Split](#weight-split)
    - [Header Split](#header-split)
    - [Policy JWT Token](#policy-jwt-token)
    - [Policy Whitelist](#policy-whitelist)
    - [Policy Authorization](#policy-authorization)
- [Repositories](#repositories)
- [Known Issues](#known-issues)

Maintainer: Joey Zhang[me@zhangzhoujian.com]

## Prerequisites

- Kubernetes environment

- istio release command: using version[1.0.2]

### install minkube kubernetes environment

NOTE: we use aliyun version minikube in ***China network***. kubectl :1.9+. Please use ***0.28.1*** version. Latest version doesnot support localkube. [Minikube-Aliyun-Kubernetes](https://yq.aliyun.com/articles/221687)

start command

`minikube start --bootstrapper=localkube --memory=8196 --logtostderr --v 0`

### install istio to Kubernetes

Go to install directory and run make coammnd. Waiting for pods of istio-system running. In istio-demo dir:

`
make istio/apply
`

check installed status. If all pods are running, will be ok.

`
make istio/check
`

## Demo step

install base insole demo app. check api get ok when deploy done. get json response.

`
cd $ISTIO_DEMO/install/insoledemo && make insole/apply
`

curl minikube ip and http port to check 

`
curl -X GET http://192.168.99.100:31380/company/api/v1/companies
`

### Weight Split

We config the weight split in product virtual service. Different version will use different weight. Can try curl 100 times and see product version split.

`
cd $ISTIO_DEMO/install/insoledemo && make weight_split/apply
`

Clean config after done.

`
cd $ISTIO_DEMO/install/insoledemo && make weight_split/clean
`

### Header Split

We config the header split in product virtual service. Different version will consume different header. Can try curl coammnd with header then see product version split. As config header android  will be route to v1, ios will be routed to v2, others will random.

```
- headers:
        x-request-platform:
          exact: android
    route:
    - destination:
        host: product
        subset: v1
  - match:
    - headers:
        x-request-platform:
          exact: ios
    route:
    - destination:
        host: product
        subset: v2
  - route:
    - destination:
        host: product
```

`
cd $ISTIO_DEMO/install/insoledemo && make header_split/apply
`

Clean config after done.

`
cd $ISTIO_DEMO/install/insoledemo && make header_split/clean
`

### Policy JWT Token

1. generate new key pairs, and host the public key in remote host.
2. config with policy 
3. generate new JWT token in jwt.io 


`
cd $ISTIO_DEMO/install/insoledemo && make policy_jwt/apply
`

Clean config after done.

`
cd $ISTIO_DEMO/install/insoledemo && make policy_jwt/clean
`

### Policy Whitelist

white list rule need about 5 minutes to take affect. Maybe need patient. Internal call of comany to product will be block by white list, beacause of white list in product is set to user.

`
cd $ISTIO_DEMO/install/insoledemo && make white_list/apply
`

Clean config after done.

`
cd $ISTIO_DEMO/install/insoledemo && make white_list/clean
`

### Policy Authorization

Open customized mixer adapter authz, will check token payload user_roles, api will be unauth except admin role. NEED restart policy pod to take affect.

`
cd $ISTIO_DEMO/install/insoledemo && make authz/apply
`

Clean config after done.

`
cd $ISTIO_DEMO/install/insoledemo && make authz/clean
`


## Repositories

test command in another repo [echo](https://github.com/zhangzhoujian/echo)

## Known Issues

1. white list take too long for take affect
2. authz need install first before istio start, or need restart policy pod to take affect
