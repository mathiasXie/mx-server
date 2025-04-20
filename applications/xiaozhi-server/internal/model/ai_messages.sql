CREATE TABLE ai_messages (  
    id int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    user_id int NOT NULL COMMENT '此消息的用户id',
    device_id int NOT NULL COMMENT '此消息的设备id',
    messsage text NOT NULL COMMENT '消息内容',
    role varchar(100) NOT NULL COMMENT '消息由谁发出,USER,ASSISTANT',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) COMMENT '消息表';