DROP TABLE IF EXISTS `t_rule`;
CREATE TABLE `t_rule`
(
    `id`            INT(10) unsigned NOT NULL auto_increment comment,
    `create_time`   DATETIME NOT NULL default current_timestamp comment,
    `update_time`   DATETIME NOT NULL default current_timestamp ON update current_timestamp comment,
    `rule_id`       VARCHAR(64) NOT NULL comment 'ID of the rule',
    `cidr`          VARCHAR(255) NOT NULL comment 'Allow traffic inbound CIDR',
    `destination`   VARCHAR(255) NOT NULL comment 'DNAT destination',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`uuid`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  COLLATE = utf8_general_ci