<h1 align="center" style="margin: 30px 0 30px; font-weight: bold; font-size: 30px">Ovra-Zero</h1>
<h4 align="center">基于Go-Zero实现的若依服务端脚手架（支持多租户）</h4>



## 平台简介
+ 提供了完整的权限系统、多租户支持、RBAC 权限控制、菜单管理等功能，适合快速搭建企业级后台管理系统。
+ 前端项目：基于 [Ruoyi-Plus-Vben5](https://gitee.com/dapppp/ruoyi-plus-vben5.git)
+ 后端项目：基于 [RuoYi-Vue-Plus](https://gitee.com/dromara/RuoYi-Vue-Plus.git) 的功能模型，使用 Go-Zero 重写。

## 在线体验
> 账号密码：admin/admin123
+ 演示地址：
    + https://vben5.ovra.dev/ (vben5版本)
+ 文档地址：
    + https://ovra.dev/

## 内置功能
1.  用户管理：用户是系统操作者，该功能主要完成系统用户配置。
2.  部门管理：配置系统组织机构（公司、部门、小组），树结构展现支持数据权限。
3.  岗位管理：配置系统用户所属担任职务。
4.  菜单管理：配置系统菜单，操作权限，按钮权限标识等。
5.  角色管理：角色菜单权限分配、设置角色按机构进行数据范围权限划分。
6.  字典管理：对系统中经常使用的一些较为固定的数据进行维护。
7.  参数管理：对系统动态配置常用参数。
8.  操作日志：系统正常操作日志记录和查询；系统异常信息日志记录和查询。
9.  登录日志：系统登录日志记录查询包含登录异常。
10. 文件管理：文件上传下载、文件配置管理。
11. 租户管理：多租户支持，租户维度数据隔离、租户启停管理。
12. 租户套餐管理：支持租户套餐管理、模块功能控制、租户容量限制。

## 待完成功能
+ [X] 代码生成，需结合代码生成器（采用GoLand插件形式，仅实现后端生成功能）
+ [ ] 报表大屏可视化
+ [ ] OAuth2.0

## 快速启动（开发环境）
+ 后端需安装 Go 1.24+，数据库为 MySQL（推荐 8.0+）
+ clone下来代码需要对应修改etc下的配置文件中数据库与redis地址
+ 服务注册以来etcd(etcd安装教程查看 https://ovra.dev/docs/tutorial-install/etcd-install)

### 后端运行
```shell
# 克隆后端代码
git clone https://github.com/cls-cloud/atlas-zero.git
make init && make build-all && make back-all
```
### 网关运行
> 网关使用traefik，下载地址：https://github.com/traefik/traefik/releases
```shell
# 下载系统对应版本的traefik，解压后放在bin/traefik目录下
./bin/traefik/traefik --configfile=./bin/traefik/traefik.yaml
```

## 联系方式 / 技术交流

- **Telegram**：[@ovra12](https://t.me/ovra12)
- **QQ**：2579260178（备注：`ovra-zero`）
- **邮箱**：ut1221@icloud.com（标题：`[ovra-zero] : 简要说明问题`）

