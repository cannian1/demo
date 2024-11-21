CREATE TABLE book
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `title`        varchar(128) NOT NULL COMMENT '书籍名称',
    `author`       varchar(128) NOT NULL COMMENT '作者',
    `price`        int          NOT NULL DEFAULT '0' COMMENT '价格',
    `publish_date` datetime COMMENT '出版日期',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='书籍表';
