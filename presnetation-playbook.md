【start】

大家好，我是来自惠普公司的张周健，很高兴有这个机会，今天能在这里跟大家分享。





【weight split】
我们先来简单介绍一下关于这个demo的内部结构，在系统中有两个主要的service，一个是company service 一个是product service。 当一个请求调用get company 的时候，company svc 内部会去调用product service，获取的信息组合成json response 返回。当kubernetes 中存在有多个product serivce的pod， 或许有多个版本的时候，istio调用是使用round robin 轮询的方式来调用，流量会平均地分布在各个pod上。

那么现在我们有这样的需求场景，需要新上线一个金丝雀版本，比如我们这里的V2版本，要求一部分流量导入新的金丝雀版本，我们一般可以配置一些router规则来导流，那我们在istio应该怎样做呢，可以在product的virtual service yaml 文件中配置weight参数，来达到一样的效果。接下来我来具体演示一下。

# 拖动窗口到投影。三个窗口
为了演示方便，我们的系统的系统里面已经部署了两个版本，一个是v1就是我们的当前版本，另一个v2就是我们的金丝雀版本。在没有配置权重控制策略的访问情况，流量是平均分配的。
weight_split 
然后我们来应用具体权重分配的策略。
make weight_split/apply
再来看的话，就可以看策略已经生效，版本v1基本接受95%的流量，金丝雀版本基本上接受到5%的流量。
我们来看到具体的配置文件， 
cat virtual-service-weight-split.yaml
在route的destination下面有一个weight字段，所有weight的总和是100，按照这个权重配置来分配流量。
清除权重配置之后就会再次回到轮询访问方式，
make weight_split/clean

# 回到ppt， 转场到header split
以上呢就是通过配置权重来分配流量的方式，但是有时候我们需要更加精细化的控制，通过解析request的某些信息来控制流量的导向，将流量倒入到特定版本，

【header split】
这里我们通过解析请求的自定义 header 来实现，比如我们设定，从ios平台发出的请求，仅可以访问我们新上线的版本，这里是V2的product service，依旧可以通过istio virtualserivce 的配置来实现这个功能。前端的请求访问带上特定的header，自定义的x-request-platform，如果这个值是ios，就访问我们的v2版本。我们再来看一下具体的操作。

#拖动窗口到投影
在没有应用策略的情况下，带有ios header的请求是轮训访问product service的
header_split ios
make header_split/apply
在应用策略之后我们可以看到，所有带有ios header的请求全部都访问v2版本的product service，
 cat virtual-service-header-split.yaml
结合具体的配置文件来看，所有的ios header导入版本2，所有的android header导入版本1，而其他的，比如这个header 为空的情况，就会以轮询的方式访问。

# 回到ppt，承接security
通过以上这两个小demo，展示了 traffic splitting的功能，istio另外一个核心的功能就是security，istio security，在整个service mesh 里面可以在很多层面上去做，比如说我们可以在service之间使用mTLS，使服务之间的交互更安全，也可以通过pilot给proxy配置一些policy， 在proxy里做check，我们也可以给mixer配置一些policy，在mixer里面做check。在我们项目中，我们在proxy里面配置了一个JWT 认证policy，来验证jwt 是否是有效的，下面我们再通过一个小 demo来演示一下。

【policy jwt】
我们生成JWT的时候是用RS256签名的，那我们这里准了两个token，一个是正确的，另一个是用HS256签名的不正确的token，首先我们使用错误的token来访问api，在没有policy配置的情况下是可以访问的，
policy_jwt HS256
make policy_jwt/apply
这个policy的生效大约需要十几秒的时间，在这期间，我们先来看一下具体的配置
cat policy-jwt.yaml
大家看一下配置文件，有几个重要的地方需要注意，第一个这个issuer表示这个JWT的颁发者，audiences表示颁发给谁，
jwksUri是能够访问到我们用于验证签名的公钥，这里呢我们配置了一个公网地址，等待这个policy生效之后，一旦用无效的token来访问我们的api，就会返回401 unauthorized。
现在我们看到这个policy已经生效了，但是大家注意下，这里出现401之后还竟然出现了200，这个具体原因我们后面会做进一步解释。
那么我们再使用一个正确的token来访问api，是可以访问的
policy_jwt RS256
make policy_jwt/clean

# 转场 pilot相关
这个demo展示的是如何通过简单配置来应用jwt token的验证，这个是跟pilot相关的功能，我们接下来继续了解下pilot组件的细节 

# 喜哥mixer 讲的差不多
adapter是什么呢，如何处理这些资源，如何传入adapter，根据mixer的check结果来决定请求在进入业务pod之前的访问权限，
通过两个demo来演示一下mixer policy的功能。

我们这里有两个demo，第一个demo演示mixer内置的listchecker，和listentry来演示白名单功能，第二个demo演示我们通过自定义authz adapter来演示给予RBAC的访问控制功能
【white list】
我们的应用场景就是想限制对prod service的内部访问，只有user service 可以访问prod。
# 切换到 命令行
white_list
在没有配置被名单之前，prod service是可以被company serivce 调用的，我们可以看到company service可以从product service获取数据，接下来我们先应用白名单机制的来阻止company 访问product， 
 make white_list/apply
mixer policy的生效时间大约需要三十秒，生效之后，company从product获取的数据会变成null，我们利用等待policy生效的时间来看下具体的配置。
cat white-list.yaml
listchecker是一个适配器，我们在这里配置了一个静态的白名单列表，第二listcheker是一个模板，这里面配置了listchecker要处理的数据，第三个rule是配置了这个policy在什么情况下被执行。最终达到的目的是，
如果匹配到了目标service是product，就应用list check来检查发起访问的service是否在白名单中。

# 看到白名单生效了，
我们可以看到白名单机制生效之后，company service不在白名单里，所以就无法获取product的数据了。

# 转场 authz
我们这个demo演示的是通过mixer内置的adapter，使用白名单来实现访问控制，接下来我们演示通过，自定义的adapter，authz，可以根据JWT中的role进行RBAC的检查，来达到访问控制的目的

【authz】

我们实现了一个了一个authz adapter作为插件加载在mixer当中，这个插件调用我们开发的authz service来进行RBAC检测，达到访问控制的目的。


每一个访问company的request，都会首先在mixer当中的authz adapter进行check，而具体的check逻辑是在内部的authz service中，这里我们准备了两个role，一个admin，一个user，这两个role配置在JWT的payload中，在authzzervice中的校验逻辑也很简单，admin通过，user则401。OK，我们具体来看一下。

在没有配置AUTHZ 的rule之前，request并不会收到mixer的authz adapter的check，所以user role的token也可以获取数据，
authz_role user
我们来应用这个rule，所有的请求都会经过authz的check。
 make authz/apply
在等待生效的时间内，我们来看下具体的配置文件，
cat authz-config.yaml

这里的authz是我们自定义的adapter，因为所有的功能都在authzservice当中实现，所以这里我们只需要配置authzservice的访问地址，authorization是一个模板，定义了authz adapter需要处理的数据，我们这里最需要就是从authorization header里面拿到JWT传给authz service进行处理，至于rule和之前类似，是设置了policy被执行的条件。

# 生效了authz，
我们看到，使用user role已经无法访问服务，那么我们使用adminrole再来尝试一下，
authz_role admin
adminrole经过adapter，而后经过authzserivce 的验证，是可以继续访问服务的。

那么以上两个demo简单的展示了，istio是怎样使用mixer来做访问控制的

#喜哥切入开始metrics
# opentracing

<!-- 
【istio-proxy admin api】

Istio-proxy本身提供了一系列的管理api来提供运行时的动态配置，这里我们带大家简略看一下如何使用istio-proxy内部的api来修改log的level，可以帮助我们在开发中更好的定位问题等等。
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
curl -X POST localhost:15000/logging?filter=info -->








