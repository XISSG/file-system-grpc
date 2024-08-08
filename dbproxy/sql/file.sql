CREATE TABLE `tbl_file`
(
    `id`            bigint        NOT NULL,
    `file_checksum` varchar(256)  NOT NULL DEFAULT '' COMMENT '文件的校验和',
    `file_name`     varchar(256)  NOT NULL DEFAULT '' COMMENT '文件名',
    `file_size`     bigint(20)             DEFAULT 0 COMMENT '文件大小',
    `file_addr`     varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
    `user_id`       bigint        NOT NULL COMMENT '文件所有者id',
    `create_at`     datetime               DEFAULT NOW() COMMENT '创建日期',
    `update_at`     datetime               DEFAULT NOW() on update CURRENT_TIMESTAMP() COMMENT '更新日期',
    `status`        bigint       NOT NULL DEFAULT 0 COMMENT '状态（可用/禁用/已删除等状态）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_file_hash` (`file_checksum`),
    KEY `idx_status` (`status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;