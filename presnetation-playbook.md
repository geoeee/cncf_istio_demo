【start】

大家好，我是来自惠普公司的张周健，很高兴有这个机会，今天能在这里跟大家分享。





【weight split】
我们先来简单介绍一下关于这个demo的内部结构，在系统中有两个主要的service，一个是company service 一个是product service。 当一个请求调用get company 的时候，company svc 内部会去调用product service，获取的信息组合成json response 返回。当kubernetes 中存在有多个product serivce的pod， 或许有多个版本的时候，istio调用是使用round robin 轮询的方式来调用，流量会平均地分布在各个pod上。

那么现在我们有这样的需求场景，需要新上线一个金丝雀版本，比如我们这里的V2版本，要求一部分流量导入新的金丝雀版本，传统的做法是在load balance上进行配置，使用istio的话，可在product的virtual service yaml 文件中配置weight参数，来达到一样的效果。接下来我来具体演示一下。

# 拖动窗口到投影。三个窗口
首先我们来看没有应用权重控制策略的访问情况，
weight_split 
执行调用之后，会在内部调用三百次product service，返回具体的版本和pods信息，
可以看到，访问product svc的流量是平均的。然后我们来应用具体的策略。
make weight_split/apply
再来看的话，就可以看策略已经生效，版本v1基本接受95%的流量，金丝雀版本基本上接受到5%的流量。
我们来看到具体的配置文件， 
cat virtual-service-weight-split.yaml
在route的destination下面有一个weight字段，所有weight的总和是100，按照这个权重配置来分配流量。
清除权重配置之后就会再次回到轮询访问方式，

# 回到ppt， 转场到header split
以上呢就是通过配置权重来分配流量的方式，但是有时候我们需要更加精细化的控制，通过解析request的某些信息来控制流量的导向，将流量倒入到特定版本，

【header split】
这里我们通过解析请求的特定header来实现，比如我们设定，从ios平台发出的请求，仅可以访问我们新上线的版本，这里是V2的product service，依旧可以通过istio virtualserivce 的配置来实现这个功能。前端的请求访问带上特定的header，自定义的x-request-platform，如果这个值是ios，就访问我们的v2版本。我们再来看一下具体的操作。

#拖动窗口到投影
在没有应用策略的情况下，带有ios header的请求是轮训访问product service的
header_split ios
make header_split/apply
在应用策略之后我们可以看到，所有带有ios header的请求全部都访问v2版本的product service，
 cat virtual-service-header-split.yaml
结合具体的配置文件来看，所有的ios header导入版本2，所有的android header导入版本1，而其他的，比如这个header 为空的情况，就会以轮询的方式访问。

# 回到ppt，承接security
通过以上这两个小demo，展示了 traffic splitting的功能，istio另外一个核心的功能就是security，istio security，在整个service mesh 里面可以在很多层面上去做，比如说我们可以在service之间建立安全的mTLS连接，把我们的传输数据进行加密，也可以在pilot里面下发一些policy到proxy， 在proxy这一层做check，我们也可以在mixer这一层配置policy，自mixer这一层做check，那我们在实际工作过程中，通过pilot配置了一个jwt的authentication的 policy，并下发到proxy，在实际收到请求的时候来验证jwt token是否是有效的，下面我们再通过一个小 demo来演示一下。

【policy jwt】
相信大家对于authorization的token验证已经很熟悉，我们准备了两个token，一个是以RS256的方式加密，另一个是HS256的方式加密，我们配置的加密方式是RS256，首先我们使用错误的token来访问我们的api，在没有policy配置的情况下是可以访问的，
policy_jwt HS256
make policy_jwt/apply
这个policy的生效大约需要十几秒的时间，我们结合文件来看一下具体的配置
cat policy-jwt.yaml
文件中配置的issuer表示这个JWT的颁发者，audiences表示颁发给谁，允许访问什么，而jwksUri，是我们配置的认证JWT签名的加密公钥的公网地址，等待这个policy生效之后，会将不符合policy的jwt token的请求当掉，返回401 unauthorized。
现在我们看到这个policy已经生效了，但是我们注意下这里有一段比较短暂的200 和401交替出现的情况，是因为proxy的work保证的是最终一致性，需要一个逐步更新的过程。
那么我们在用一个正确的token来访问api，是可以访问的
policy_jwt RS256
make policy_jwt/clean

# 转场 pilot相关
这个demo展示的是如何通过简单配置来应用jwt token的验证，这个是跟pilot相关的功能，我们接下来继续了解下pilot组件的细节 

# 喜哥mixer 讲的差不多
adapter是什么呢，如何处理这些资源，如何传入adapter，根据mixer的check结果来决定请求在进入业务pod之前的访问权限，
通过两个demo来演示一下mixer check的功能。

【white list】
# 切换到 命令行
在没有配置被名单之前，系统内部的各个service之间是可以任意调用的，我们可以看到company service可以从product service获取数据，接下来我们先应用白名单机制的来阻止company 访问product， 
white_list
 make white_list/apply
mixer check的生效时间大约需要一分钟，products的数据会表示成null，我们来看下配置。
cat white-list.yaml
在这一个文件内，共配置三个yaml实体，第一个是listchecker，配置了白名单的静态数据，还有一个balcklist的flag来表示使用白名单机制还是黑名单机制，第二个是listentry，进来的request的来源的字段属性，第三个是具体的rule，如果匹配到了目标service是product，就应用list check来检查是否在被名单中。

# 看到白名单生效了，
我们可以看到白名单机制生效之后，company service不在白名单里，所以就无法获取product的数据了。

# 转场 authz
mixer的list adapter呢是官方内置实现的adapter，接下来我们修改istio代码之后，自定义的扩展adapter，authz，配合authz-service的实现，可以根据JWT token中的role进行RBAC的检查，来达到角色权限控制的目的

【authz】

每一个进入company的request，都会进入mixer当中的authz adapter进行check，而具体的check逻辑是现在内部的authz service中，这里我们准备了两个role，一个admin，一个user，这两个role配置在JWT token的payload中，在authzzervice中的校验逻辑也很简单，admin通过，user则401。OK，我们具体来看一下。

在没有配置AUTHZ 的rule之前，request并不会收到mixer的authz adapter的check，所以user role的token也可以获取数据，
authz_role user
我们来应用这个rule，进来的流量加上authz的check。
 make authz/apply
在等待生效的时间内，我们来看下具体的配置文件，
cat authz-config.yaml

首先我们来看我们要配置的rule， match的部分就是进入company service 的所有request，其中一个要配置的是handler，并且传入对应的instance，authz的handler就是第二个yaml实体，这里我们除了声明这个handler之外，还配置了所需的authzservice的具体变量，demo namespace下的authzservice，之后在handler当中我们会给authzservice post一个认证请求来验证具体逻辑，最上面一个也就是instance的定义，由istio来实例化这个对象，填充对应的attributes传给handler。

# 生效了authz，
我们看到，使用user role已经无法访问服务，那么我们使用adminrole再来尝试一下，
authz_role admin
adminrole经过adapter，而后经过authzserivce 的验证，是可以继续访问服务的。

那么以上两个demo简单的展示了，istio所提供mixer adapter 关于check的一些特性，接下来我们来看下istio mixer的另一些方面。

#喜哥切入开始metrics
# opentracing


【istio-proxy admin api】

Istio-proxy 或者说 就是envoy本身提供了一系列的管理api来提供运行时的动态配置，这里我们带大家简略看一下如何使用istio-proxy内部的api来修改log的level，可以帮助我们在开发中更好的定位问题等等。
首先我们看到的是在保持不断被访问状态下的authzservice的istio-proxy的container log，他的默认log level是info，通过调用本地15000端口的api，直接修改log level，就可以看到debug level的log，也可以通过同样的方式来关闭。这是istio或者说envoy提供的方便开发的功能。

# 打开authz_role admin 一直调用
authz_role admin
# 打开 authz service 的istio-proxy container的log
 kubectl logs -n demo -f -c istio-proxy $(kubectl get pods -n demo | awk ‘/authz/ {print $1}’)
# 打开 authz service 的istio-proxy container的交互式命令行
 kubectl exec -it -n demo -c istio-proxy $(kubectl get pods -n demo | awk ‘/authz/ {print $1}’) bash
# 查看 管理api
 curl -X POST localhost:15000/admin
# 产看logging的组件
curl -X POST localhost:15000/logging
# 修改logging 的组件log level
curl -X POST localhost:15000/logging?filter=debug
curl -X POST localhost:15000/logging?filter=info








