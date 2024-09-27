CREATE TABLE `devices` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '设备的Id',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '设备名',
  `server_addr` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '设备IP地址',
  `change_at` int NOT NULL COMMENT '新增或修改时间',
  `brand` int NOT NULL COMMENT '设备品牌',
  `port` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '设备登录端口号',
  `idc` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '设备隶属机房',
  `device_level` int NOT NULL COMMENT '设备层级',
  `status` tinyint NOT NULL COMMENT '设备状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;