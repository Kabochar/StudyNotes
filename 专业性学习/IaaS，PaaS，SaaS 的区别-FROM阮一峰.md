# IaaS，PaaS，SaaS 的区别

Iaas：提供底层基础资源。

PaaS：抽去硬件和系统 ，可以无缝拓展。关注业务逻辑。

SaaS：无需关心技术问题，拿来即用。



## 云服务三大类？

-   **IaaS**：基础设施服务，Infrastructure-as-a-service
-   **PaaS**：平台服务，Platform-as-a-service
-   **SaaS**：软件服务，Software-as-a-service

## 区别？差异？

案例模拟：披萨生意。

**（1）方案一：IaaS**

他人提供厨房、炉子、煤气，你使用这些基础设施，来烤你的披萨。

**（2）方案二：PaaS**

除了基础设施，他人还提供披萨饼皮。

你要做的就是设计披萨的味道（海鲜披萨或者鸡肉披萨），他人提供平台服务，让你把自己的设计实现。

**（3）方案三：SaaS**

他人直接做好了披萨，不用你的介入，到手的就是一个成品。你要做的就是把它卖出去，最多再包装一下，印上你自己的 Logo。

### 总结

![1555555389971](pics\1555555389971.png)

从左到右，自己承担的工作量（上图蓝色部分）越来越少，IaaS > PaaS > SaaS。

对应软件开发

![1555555462593](pics\1555555462593.png)

**SaaS 是软件的开发、管理、部署都交给第三方，不需要关心技术问题，可以拿来即用。**普通用户接触到的互联网服务，几乎都是 SaaS，下面是一些例子。

>   -   客户管理服务 Salesforce
>   -   团队协同服务 Google Apps
>   -   储存服务 Box
>   -   储存服务 Dropbox
>   -   社交服务 Facebook / Twitter / Instagram

**PaaS 提供软件部署平台（runtime），抽象掉了硬件和操作系统细节，可以无缝地扩展（scaling）。开发者只需要关注自己的业务逻辑，不需要关注底层。**下面这些都属于 PaaS。

>   -   Heroku
>   -   Google App Engine
>   -   OpenShift

**IaaS 是云服务的最底层，主要提供一些基础资源。**它与 PaaS 的区别是，用户需要自己控制底层，实现基础设施的使用逻辑。下面这些都属于 IaaS。

>   -   Amazon EC2
>   -   Digital Ocean
>   -   RackSpace Cloud

**参考链接**

-   [SaaS, PaaS and IaaS explained in one graphic](https://m.oursky.com/saas-paas-and-iaas-explained-in-one-graphic-d56c3e6f4606), by David Ng
-   [When to use SaaS, PaaS, and IaaS](https://www.computenext.com/blog/when-to-use-saas-paas-and-iaas/), by Eamonn Colman



## 备注

整理自：http://www.ruanyifeng.com/blog/2017/07/iaas-paas-saas.html

作者：[阮一峰](http://www.ruanyifeng.com/)

图片来源：http://www.ruanyifeng.com/blog/2017/07/iaas-paas-saas.html