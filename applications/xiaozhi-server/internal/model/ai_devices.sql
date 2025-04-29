CREATE TABLE `ai_devices` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
  `device_mac` varchar(17) NOT NULL COMMENT '设备网络mac地址',
  `version` varchar(100) NOT NULL COMMENT '设备版本号',
  `board_type` varchar(100) NOT NULL COMMENT '设备主板',
  `board_ssid` varchar(100) NOT NULL COMMENT '设备连接的ssid',
  `board_ip` varchar(100) NOT NULL COMMENT '设备内网ip',
  `language` varchar(100) NOT NULL COMMENT '语言',
  `chip_model_name` varchar(100) NOT NULL COMMENT '设备芯片',
  `ip` varchar(100) NOT NULL COMMENT '外网ip',
  `user_id` int NOT NULL COMMENT '用户id',
  `role_id` int DEFAULT NULL COMMENT '绑定的角色ID',
  `bind_code` int DEFAULT NULL COMMENT '绑定码',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_mac` (`device_mac`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='设备表'