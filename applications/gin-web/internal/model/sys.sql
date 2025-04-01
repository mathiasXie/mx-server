
CREATE TABLE `sys_users` (
                             `id` int NOT NULL AUTO_INCREMENT,
                             `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户名',
                             `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '密码',
                             `role_id` int DEFAULT NULL COMMENT '所属权限组ID',
                             `realname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '真实姓名',
                             `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮箱',
                             `mobile` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '手机号码',
                             `status` tinyint(1) DEFAULT NULL COMMENT '状态：1：正常; 0：禁用',
                             `is_super` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否为超级管理员 1:是 0:否',
                             `last_time` timestamp NULL DEFAULT NULL COMMENT '最后一次操作时间',
                             `last_ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '最后一次操作IP',
                             `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                             `updated_at` timestamp NULL DEFAULT NULL,
                             `deleted_at` timestamp NULL DEFAULT NULL,
                             PRIMARY KEY (`id`),
                             UNIQUE KEY `unique_user` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='后台用户表';

CREATE TABLE `sys_roles` (
                             `id` bigint NOT NULL AUTO_INCREMENT,
                             `parent_id` bigint DEFAULT NULL COMMENT '上级角色ID',
                             `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
                             `sequence` int DEFAULT NULL COMMENT '排序',
                             `status` tinyint DEFAULT NULL COMMENT '状态 1=启用;0=禁用',
                             `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '角色描述',
                             `deleted_at` timestamp NULL DEFAULT NULL,
                             `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                             `updated_at` timestamp NULL DEFAULT NULL,
                             PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='权限角色';

CREATE TABLE `sys_role_routes`  (
                                    `id` bigint NOT NULL AUTO_INCREMENT,
                                    `role_id` bigint NOT NULL COMMENT '角色id',
                                    `route_id` bigint NOT NULL COMMENT '路由功能表',
                                    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 66 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '角色权限表' ROW_FORMAT = Dynamic;

CREATE TABLE `sys_routes`  (
                               `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
                               `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '路由标识',
                               `label` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '名称',
                               `permission` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '权限标识,用于前端按钮级别的鉴权',
                               `parent_id` bigint NULL DEFAULT 0 COMMENT '父级菜单ID',
                               `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '路径;对应前端vue路由',
                               `api_path` varchar(250) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT 'api地址,不包含management,用于中间件鉴权,如果此功能使用到多个api用逗号隔开',
                               `component` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '组件',
                               `sequence` int NULL DEFAULT 1 COMMENT '排序',
                               `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '菜单图标',
                               `style` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '样式',
                               `type` tinyint NULL DEFAULT 1 COMMENT '类型（1=菜单;2=按钮）',
                               `status` tinyint NULL DEFAULT 1 COMMENT '1=启用;0=禁用',
                               `global` tinyint NULL DEFAULT 2 COMMENT '公共资源 1是,2否 无需分配所有人就可以访问的',
                               `display` tinyint NULL DEFAULT 1 COMMENT '0=隐藏;1=显示',
                               `created_at` timestamp NULL DEFAULT (now()) COMMENT '创建时间',
                               `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
                               `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                               PRIMARY KEY (`id`) USING BTREE,
                               UNIQUE INDEX `sys_routes_name_uindex`(`name` ASC, `deleted_at` ASC) USING BTREE,
                               INDEX `INX_STATUS`(`global` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 56 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '菜单' ROW_FORMAT = DYNAMIC;
