CREATE TABLE `ai_roles` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
  `user_id` int NOT NULL COMMENT '用户id',
  `role_name` varchar(100) NOT NULL COMMENT '角色名称',
  `role_desc` text NOT NULL COMMENT '角色描述',
  `llm` varchar(100) NOT NULL COMMENT '角色使用的大模型提供商',
  `llm_model_id` varchar(100) NOT NULL COMMENT '角色使用的大模型id',
  `tts` varchar(100) NOT NULL COMMENT '角色使用的语音合成提供商',
  `tts_voice_id` varchar(100) NOT NULL COMMENT '角色使用的语音合成声音id',
  `language` varchar(40) DEFAULT NULL COMMENT '角色使用的语言',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色表'