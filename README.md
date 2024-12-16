[![status-badge](https://dev-ci.shijicloud.com/woodpecker/api/badges/3/status.svg)](https://dev-ci.shijicloud.com/woodpecker/repos/3)

# scavenger
专注于k8s或者k3s集群的pod资源CPU、内存和句柄的占用监控，设置最高阈值，其中任意一种资源占用超过阈值则删除对应的deployment资源；

### 项目提交规范

工具: **Sourcetree**

#### 开发流程简介

1. 使用Sourcetree拉取项目到本地.
2. 利用Sourcetree的GitFlow功能进行GitFlow的初始化.
3. 基于`develop`分支建立新的`feature/新功能`分支.
4. 在`feature/新功能`分支进行开发.
5. 开发完后, 将修改提交到本地`feature/新功能`分支.
6. 将远端的`develop`分支pull到`feature/新功能`分支(如有冲突, 请解决).
7. 将本地`feature/新功能`分支推送到远端.
8. 去远端https://dev-scm.shijicloud.com/Kunlun/k-octopus提PR, 并写明修改内容.
9. 通知负责人审核并合并PR.
10. 删除本地`feature/新功能`分支. 顺便将远端的`develop`分支pull到本地`develop`分支, 以保持代码是最新的.

具体内容请参考: https://kunlun.shijicloud.com/docs/public/kunlun-kie/KunlunGitSpecification/index.html#/GitToolSourcetree

#### 分支命名
- 新功能分支  
  `feature/云效号_功能描述`  
  例, `feature/KLOT-1232_clue_list`
- bug修复分支  
  `bugfix/KLOT-2333_bug描述`  
  例, `bugfix/KLOT-2333_user_not_found`