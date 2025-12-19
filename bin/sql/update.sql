ALTER TABLE `sys_client`
    MODIFY COLUMN `id` varchar(36) NOT NULL COMMENT 'id';

ALTER TABLE `sys_config`
    MODIFY COLUMN `config_id` varchar(36) NOT NULL COMMENT '参数主键';

ALTER TABLE `sys_dept`
    MODIFY COLUMN `dept_id` varchar(36) NOT NULL COMMENT '部门id',
    MODIFY COLUMN `parent_id` varchar(36) DEFAULT '0' COMMENT '父部门id',
    MODIFY COLUMN `leader` varchar(36) DEFAULT NULL COMMENT '负责人';

ALTER TABLE `sys_dict_data`
    MODIFY COLUMN `dict_code` varchar(36) NOT NULL COMMENT '字典编码';

ALTER TABLE `sys_dict_type`
    MODIFY COLUMN `dict_id` varchar(36) NOT NULL COMMENT '字典主键';

ALTER TABLE `sys_logininfor`
    MODIFY COLUMN `info_id` varchar(36) NOT NULL COMMENT '访问ID';

ALTER TABLE `sys_menu`
    MODIFY COLUMN `menu_id` varchar(36) NOT NULL COMMENT '菜单ID',
    MODIFY COLUMN `parent_id` varchar(36) DEFAULT '0' COMMENT '父菜单ID';

ALTER TABLE `sys_notice`
    MODIFY COLUMN `notice_id` varchar(36) NOT NULL COMMENT '公告ID';

ALTER TABLE `sys_oper_log`
    MODIFY COLUMN `oper_id` varchar(36) NOT NULL COMMENT '日志主键';

ALTER TABLE `sys_oss`
    MODIFY COLUMN `oss_id` varchar(36) NOT NULL COMMENT '对象存储主键';

ALTER TABLE `sys_oss_config`
    MODIFY COLUMN `oss_config_id` varchar(36) NOT NULL COMMENT '主键';

ALTER TABLE `sys_post`
    MODIFY COLUMN `post_id` varchar(36) NOT NULL COMMENT '岗位ID',
    MODIFY COLUMN `dept_id` varchar(36) NOT NULL COMMENT '部门id';

ALTER TABLE `sys_role`
    MODIFY COLUMN `role_id` varchar(36) NOT NULL COMMENT '角色ID';

ALTER TABLE `sys_role_dept`
    MODIFY COLUMN `role_id` varchar(36) NOT NULL COMMENT '角色ID',
    MODIFY COLUMN `dept_id` varchar(36) NOT NULL COMMENT '部门ID';

ALTER TABLE `sys_role_menu`
    MODIFY COLUMN `role_id` varchar(36) NOT NULL COMMENT '角色ID',
    MODIFY COLUMN `menu_id` varchar(36) NOT NULL COMMENT '菜单ID';

ALTER TABLE `sys_social`
    MODIFY COLUMN `id` varchar(36) NOT NULL COMMENT '主键',
    MODIFY COLUMN `user_id` varchar(36) NOT NULL COMMENT '用户ID';

ALTER TABLE `sys_tenant`
    MODIFY COLUMN `package_id` varchar(36) COMMENT '套餐ID',
    MODIFY COLUMN `id` varchar(36) NOT NULL COMMENT 'id';

ALTER TABLE `sys_tenant_package`
    MODIFY COLUMN `package_id` varchar(36) NOT NULL COMMENT '租户套餐id';

ALTER TABLE `sys_user`
    MODIFY COLUMN `user_id` varchar(36) NOT NULL COMMENT '用户ID',
    MODIFY COLUMN `dept_id` varchar(36) DEFAULT NULL COMMENT '部门ID',
    MODIFY COLUMN `avatar` varchar(36) DEFAULT NULL COMMENT '头像地址';

ALTER TABLE `sys_user_post`
    MODIFY COLUMN `user_id` varchar(36) NOT NULL COMMENT '用户ID',
    MODIFY COLUMN `post_id` varchar(36) NOT NULL COMMENT '岗位ID';

ALTER TABLE `sys_user_role`
    MODIFY COLUMN `user_id` varchar(36) NOT NULL COMMENT '用户ID',
    MODIFY COLUMN `role_id` varchar(36) NOT NULL COMMENT '角色ID';
