# k8s-workload-balancer

实现一个集群资源负载均衡器，要求实现如下功能：

kubernetes 集群新增节点或节点宕机恢复后，运行中的 pod 不会自动迁移到新节点上，导致集群整体资源分配不均衡, 实现一个模块，定时将部分 pod 迁出, 交给 scheduler 重新调度，实现集群资源负载均衡。

## 整体工作流程

1. 初始化程序
2. 获取所有的 Pod 和 Node
3. 首先驱逐不符合节点亲和性的 Pod，无论此 Pod 是否存在一下的特殊情况
4. list 所有 Node，过滤掉 NotReady 的节点
5. 除去上述特殊情况不能驱逐的 Pod，按照 QoS 进行驱逐，Besteffort -> Burstable -> Guaranteed 的顺序进行驱逐，每次驱逐一半的 Pod

## 几个关键点：

1. 尽量 `Evict Pod`，不要使用 `Delete Pod`，因为 `Evict Pod` 可以保证 PDB。
2. 按照 QoS 进行驱逐，Besteffort -> Burstable -> Guaranteed 的顺序进行驱逐
3. 由于拿不到监控数据，很那判断节点负载到底是多少，因此只是驱逐一半 Pod
3. Node Selector 和 节点亲和性也需要考虑
4. 不要做 kubelet、kube-proxy 和 kube-scheduler 相关的事情
5. 发生过 OOM killer 的 Pod 要被调度到其他节点（未实现）
6. 设置连个阈值，如果节点负载达到这个阈值，就执行驱逐，如果没设置，按照默认方式进行即可（未实现）

## 几种无法驱逐的特殊情况：

1. DaemonSet 创建的 Pod 无法被驱逐
2. 开启 CPU manager 后，独占 CPU 的 Pod 不建议被驱逐（未实现）
3. 使用了 Local Storage 的 Pod 不建议被驱逐
4. 带存储的 statefulset 不建议被驱逐
6. Job 创建的、处于 Completed 状态 Pod 不建议被驱逐
7. Static Pod 不建议被驱逐
8. Critical Pod 不能被驱逐
9. Preemptor Pod 不能被驱逐
10. 还未调度的 Pod 不能被驱逐
