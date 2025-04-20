CREATE TABLE `ai_devices` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
  `device_id` varchar(100) NOT NULL COMMENT '设备id',
  `role_id` int DEFAULT NULL COMMENT '绑定的角色ID',
  `device_mac` varchar(17) NOT NULL COMMENT '设备网络mac地址',
  `device_name` varchar(100) NOT NULL COMMENT '设备名称',
  `token` varchar(100) NOT NULL COMMENT '设备token',
  `user_id` int NOT NULL COMMENT '用户id',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `bind_code` int DEFAULT NULL COMMENT '绑定码',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_mac` (`device_mac`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='设备表'