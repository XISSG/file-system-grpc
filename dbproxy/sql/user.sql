CREATE TABLE `tbl_user`
(
    `id`        BIGINT PRIMARY KEY ,
    `username`  VARCHAR(256) NOT NULL COMMENT '用户名',
    `password` varchar(256) NOT NULL COMMENT '密码',
    `create_at` datetime              DEFAULT NOW() COMMENT '创建日期',
    `update_at` datetime              DEFAULT NOW() on update CURRENT_TIMESTAMP() COMMENT '更新日期',
    `status`    int(11)      NOT NULL DEFAULT 0 COMMENT '状态（可用/禁用/已删除等状态）'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;