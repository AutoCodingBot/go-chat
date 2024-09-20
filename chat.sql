/*
 Navicat Premium Data Transfer

 Source Server         : MySQL_from_Docker
 Source Server Type    : MySQL
 Source Server Version : 80037
 Source Host           : localhost:3306
 Source Schema         : chat

 Target Server Type    : MySQL
 Target Server Version : 80037
 File Encoding         : 65001

 Date: 20/09/2024 17:53:14
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for group_members
-- ----------------------------
DROP TABLE IF EXISTS `group_members`;
CREATE TABLE `group_members`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` bigint UNSIGNED NULL DEFAULT NULL,
  `user_id` int NULL DEFAULT NULL COMMENT '\'用户ID\'',
  `group_id` int NULL DEFAULT NULL COMMENT '\'群组ID\'',
  `nickname` varchar(350) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'昵称',
  `mute` smallint NULL DEFAULT NULL COMMENT '\'是否禁言\'',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_group_members_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_group_members_group_id`(`group_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '群组成员表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of group_members
-- ----------------------------
INSERT INTO `group_members` VALUES (1, '2024-07-12 18:17:38.770', '2024-07-12 18:17:38.770', 0, 5, 1, 'eric', 0);
INSERT INTO `group_members` VALUES (2, '2024-07-12 18:17:56.512', '2024-07-12 18:17:56.512', 0, 6, 1, 'samual', 0);

-- ----------------------------
-- Table structure for group_messages
-- ----------------------------
DROP TABLE IF EXISTS `group_messages`;
CREATE TABLE `group_messages`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
  `group_id` int NULL DEFAULT NULL COMMENT '群id,自增,非uuid',
  `from_user_id` int NULL DEFAULT NULL COMMENT '发送人ID',
  `content` varchar(250) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '消息内容',
  `url` varchar(350) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'文件或者图片地址\'',
  `pic` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '缩略图',
  `content_type` smallint NULL DEFAULT NULL COMMENT '\'消息内容类型：1文字，2语音，3图片\'',
  `created_at` datetime(2) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(2) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` bigint UNSIGNED NULL DEFAULT NULL COMMENT '删除时间戳',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_messages_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_messages_from_user_id`(`from_user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 46 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '群消息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of group_messages
-- ----------------------------
INSERT INTO `group_messages` VALUES (38, 1, 5, '滴滴', '', '', 1, '2024-07-30 08:50:26.04', '2024-07-30 08:50:26.04', 0);
INSERT INTO `group_messages` VALUES (39, 1, 5, 'we are united together', '', '', 1, '2024-08-22 11:12:30.89', '2024-08-22 11:12:30.89', 0);
INSERT INTO `group_messages` VALUES (40, 1, 6, 'That\'s for sure', '', '', 1, '2024-08-22 11:12:42.44', '2024-08-22 11:12:42.44', 0);
INSERT INTO `group_messages` VALUES (41, 1, 5, 'eee', '', '', 1, '2024-08-22 13:45:42.29', '2024-08-22 13:45:42.29', 0);
INSERT INTO `group_messages` VALUES (42, 1, 5, '1', '', '', 1, '2024-08-22 13:58:04.11', '2024-08-22 13:58:04.11', 0);
INSERT INTO `group_messages` VALUES (43, 1, 5, 'ttt', '', '', 1, '2024-08-22 13:59:29.32', '2024-08-22 13:59:29.32', 0);
INSERT INTO `group_messages` VALUES (44, 1, 5, '3', '', '', 1, '2024-08-22 14:00:17.90', '2024-08-22 14:00:17.90', 0);
INSERT INTO `group_messages` VALUES (45, 1, 5, 'd', '', '', 1, '2024-08-22 14:10:22.02', '2024-08-22 14:10:22.02', 0);
INSERT INTO `group_messages` VALUES (46, 1, 5, 'working?', '', '', 1, '2024-08-23 18:24:19.87', '2024-08-23 18:24:19.87', 0);

-- ----------------------------
-- Table structure for groups
-- ----------------------------
DROP TABLE IF EXISTS `groups`;
CREATE TABLE `groups`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NULL DEFAULT NULL COMMENT '\'群主ID\'',
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'群名称',
  `avatar` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `notice` varchar(350) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'群公告',
  `uuid` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'uuid\'',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` bigint UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_groups_user_id`(`user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '群组表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of groups
-- ----------------------------
INSERT INTO `groups` VALUES (1, 5, '木兰辞', '', '', 'affec7bb-8565-4cad-a05f-63a3334c53a2', '2024-07-12 18:17:38.752', '2024-07-12 18:17:38.752', 0);

-- ----------------------------
-- Table structure for user_friends
-- ----------------------------
DROP TABLE IF EXISTS `user_friends`;
CREATE TABLE `user_friends`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` int NULL DEFAULT NULL COMMENT '用户ID',
  `friend_id` int NULL DEFAULT NULL COMMENT '好友ID',
  `relation_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '好友关系ID;比如:tom_sam;tom指的是userName',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` bigint UNSIGNED NULL DEFAULT NULL COMMENT '删除时间戳',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_friends_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_user_friends_friend_id`(`friend_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '好友信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_friends
-- ----------------------------
INSERT INTO `user_friends` VALUES (7, 5, 6, 'eric_sam', '2024-07-29 12:50:47.034', '2024-07-29 12:50:47.034', 0);
INSERT INTO `user_friends` VALUES (16, 5, 14, 'eric_john', '2024-07-30 08:23:38.957', '2024-07-30 08:23:38.957', 0);
INSERT INTO `user_friends` VALUES (17, 5, 5, '', '2024-08-22 14:28:48.276', '2024-08-22 14:28:48.276', 0);

-- ----------------------------
-- Table structure for user_messages
-- ----------------------------
DROP TABLE IF EXISTS `user_messages`;
CREATE TABLE `user_messages`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
  `conversation_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '对话标识id,比如ericsam,concat两个用户昵称(降序)',
  `from_user_id` int NOT NULL COMMENT '发送人ID',
  `to_user_id` int NOT NULL COMMENT '发送对象ID',
  `content` varchar(2500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '消息内容',
  `url` varchar(350) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'文件或者图片地址\'',
  `pic` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '缩略图',
  `content_type` smallint NULL DEFAULT NULL COMMENT '\'消息内容类型：1文字，2语音，3视频\'',
  `created_at` datetime(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` bigint UNSIGNED NULL DEFAULT NULL COMMENT '删除时间戳',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_messages_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_messages_from_user_id`(`from_user_id` ASC) USING BTREE,
  INDEX `idx_messages_to_user_id`(`to_user_id` ASC) USING BTREE,
  INDEX `idx_user_messages_from_user_id`(`from_user_id` ASC) USING BTREE,
  INDEX `idx_user_messages_to_user_id`(`to_user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 152 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '消息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_messages
-- ----------------------------
INSERT INTO `user_messages` VALUES (106, 'eric_sam', 5, 6, 'working all right?', '', '', 1, '2024-07-20 19:17:16.846', '2024-07-20 19:17:16.846', 0);
INSERT INTO `user_messages` VALUES (107, 'eric_sam', 5, 6, 'oooo', '', '', 1, '2024-07-20 20:37:13.133', '2024-07-20 20:37:13.133', 0);
INSERT INTO `user_messages` VALUES (108, 'eric_sam', 5, 6, 'huoji?', '', '', 1, '2024-07-20 20:57:33.474', '2024-07-20 20:57:33.474', 0);
INSERT INTO `user_messages` VALUES (109, 'eric_sam', 6, 5, 'knock knock', '', '', 1, '2024-07-22 07:35:09.138', '2024-07-22 07:35:09.138', 0);
INSERT INTO `user_messages` VALUES (110, 'eric_sam', 6, 5, 'p', '', '', 1, '2024-07-22 07:36:07.122', '2024-07-22 07:36:07.122', 0);
INSERT INTO `user_messages` VALUES (111, 'eric_sam', 6, 5, 'e', '', '', 1, '2024-07-22 07:52:25.808', '2024-07-22 07:52:25.808', 0);
INSERT INTO `user_messages` VALUES (112, 'eric_sam', 6, 5, 'knock!', '', '', 1, '2024-07-22 07:59:17.435', '2024-07-22 07:59:17.435', 0);
INSERT INTO `user_messages` VALUES (113, 'eric_sam', 6, 5, 'lll', '', '', 1, '2024-07-22 08:02:22.117', '2024-07-22 08:02:22.117', 0);
INSERT INTO `user_messages` VALUES (114, 'eric_sam', 6, 5, '[[[', '', '', 1, '2024-07-22 08:02:37.187', '2024-07-22 08:02:37.187', 0);
INSERT INTO `user_messages` VALUES (115, 'eric_sam', 6, 5, '\\\\', '', '', 1, '2024-07-22 08:02:54.790', '2024-07-22 08:02:54.790', 0);
INSERT INTO `user_messages` VALUES (116, 'eric_sam', 5, 6, '', 'e8720bc7-8f6d-4172-9f95-23d0ef55154e.jpg', '', 3, '2024-07-22 09:34:40.022', '2024-07-22 09:34:40.022', 0);
INSERT INTO `user_messages` VALUES (117, 'eric_sam', 6, 5, 'hello there', '', '', 1, '2024-07-22 12:22:24.478', '2024-07-22 12:22:24.478', 0);
INSERT INTO `user_messages` VALUES (119, 'eric_sam', 5, 6, 'take a shot\n', '', '', 1, '2024-07-29 11:38:35.541', '2024-07-29 11:38:35.541', 0);
INSERT INTO `user_messages` VALUES (120, 'eric_sam', 5, 6, 'again', '', '', 1, '2024-07-29 11:38:45.454', '2024-07-29 11:38:45.454', 0);
INSERT INTO `user_messages` VALUES (121, 'sam_sam', 6, 6, 'didi', '', '', 1, '2024-07-29 15:17:15.480', '2024-07-29 15:17:15.480', 0);
INSERT INTO `user_messages` VALUES (122, 'eric_sam', 5, 6, 'hi', '', '', 1, '2024-07-29 15:58:50.452', '2024-07-29 15:58:50.452', 0);
INSERT INTO `user_messages` VALUES (123, 'eric_sam', 6, 5, '鸟人', '', '', 1, '2024-07-29 16:16:08.048', '2024-07-29 16:16:08.048', 0);
INSERT INTO `user_messages` VALUES (124, 'eric_sam', 5, 6, 'zai de ', '', '', 1, '2024-07-29 16:16:20.974', '2024-07-29 16:16:20.974', 0);
INSERT INTO `user_messages` VALUES (125, 'eric_john', 5, 14, 'Hi,john', '', '', 1, '2024-07-29 17:20:02.291', '2024-07-29 17:20:02.291', 0);
INSERT INTO `user_messages` VALUES (126, 'eric_john', 5, 14, 'd', '', '', 1, '2024-07-30 08:20:04.863', '2024-07-30 08:20:04.863', 0);
INSERT INTO `user_messages` VALUES (127, 'eric_john', 5, 14, 'ooo', '', '', 1, '2024-07-30 08:23:46.871', '2024-07-30 08:23:46.871', 0);
INSERT INTO `user_messages` VALUES (128, 'eric_john', 5, 14, 'it worled?', '', '', 1, '2024-07-30 08:49:27.107', '2024-07-30 08:49:27.107', 0);
INSERT INTO `user_messages` VALUES (129, 'eric_sam', 5, 6, 'p', '', '', 1, '2024-07-30 17:51:12.686', '2024-07-30 17:51:12.686', 0);
INSERT INTO `user_messages` VALUES (130, 'eric_john', 5, 14, 'dfd?', '', '', 1, '2024-08-11 12:10:51.808', '2024-08-11 12:10:51.808', 0);
INSERT INTO `user_messages` VALUES (131, 'eric_sam', 6, 5, 'ei?', '', '', 1, '2024-08-22 11:12:02.871', '2024-08-22 11:12:02.871', 0);
INSERT INTO `user_messages` VALUES (132, 'eric_sam', 5, 6, '', '549c0e14-5114-4fbb-9df2-0af06f35d694.png', '', 3, '2024-08-22 11:13:40.084', '2024-08-22 11:13:40.084', 0);
INSERT INTO `user_messages` VALUES (133, 'eric_sam', 5, 6, '', 'e41bfb8d-c1c7-452c-b775-16e60618be8a.mp4', '', 5, '2024-08-22 12:49:20.270', '2024-08-22 12:49:20.270', 0);
INSERT INTO `user_messages` VALUES (134, 'eric_sam', 5, 6, '', '2b2b5f02-f2d0-4882-acd0-bef988925bc5.mp4', '', 5, '2024-08-22 12:50:45.829', '2024-08-22 12:50:45.829', 0);
INSERT INTO `user_messages` VALUES (135, 'eric_john', 5, 14, '3333', '', '', 1, '2024-08-22 13:56:48.704', '2024-08-22 13:56:48.704', 0);
INSERT INTO `user_messages` VALUES (136, 'eric_john', 5, 14, '3332', '', '', 1, '2024-08-22 13:57:14.007', '2024-08-22 13:57:14.007', 0);
INSERT INTO `user_messages` VALUES (137, 'eric_john', 5, 14, '1024\n', '', '', 1, '2024-08-22 14:40:24.847', '2024-08-22 14:40:24.847', 0);
INSERT INTO `user_messages` VALUES (138, 'eric_sam', 5, 6, 'test!!', '', '', 1, '2024-08-23 14:32:33.374', '2024-08-23 14:32:33.374', 0);
INSERT INTO `user_messages` VALUES (139, 'eric_sam', 5, 6, 'Awake ?', '', '', 1, '2024-08-24 13:10:29.788', '2024-08-24 13:10:29.788', 0);
INSERT INTO `user_messages` VALUES (140, 'eric_sam', 5, 6, '', '332ed682-5fe4-4d6a-9d9b-1124bc3a18cc.mp4', '', 5, '2024-08-24 13:11:09.480', '2024-08-24 13:11:09.480', 0);
INSERT INTO `user_messages` VALUES (141, 'eric_sam', 5, 6, 'test\n', '', '', 1, '2024-09-13 09:05:29.539', '2024-09-13 09:05:29.539', 0);
INSERT INTO `user_messages` VALUES (142, 'eric_sam', 5, 6, '111', '', '', 1, '2024-09-13 09:05:44.714', '2024-09-13 09:05:44.714', 0);
INSERT INTO `user_messages` VALUES (143, 'eric_sam', 5, 6, '22', '', '', 1, '2024-09-13 09:05:46.559', '2024-09-13 09:05:46.559', 0);
INSERT INTO `user_messages` VALUES (144, 'eric_sam', 5, 6, '33', '', '', 1, '2024-09-13 09:05:49.042', '2024-09-13 09:05:49.042', 0);
INSERT INTO `user_messages` VALUES (145, 'eric_sam', 5, 6, '100', '', '', 1, '2024-09-13 09:13:43.901', '2024-09-13 09:13:43.901', 0);
INSERT INTO `user_messages` VALUES (146, 'eric_sam', 5, 6, '收到', '', '', 1, '2024-09-13 09:23:35.576', '2024-09-13 09:23:35.576', 0);
INSERT INTO `user_messages` VALUES (147, 'eric_sam', 5, 6, 'hit ?', '', '', 1, '2024-09-13 09:28:06.574', '2024-09-13 09:28:06.574', 0);
INSERT INTO `user_messages` VALUES (148, 'eric_eric', 5, 5, 'hello', '', '', 1, '2024-09-13 09:33:24.445', '2024-09-13 09:33:24.445', 0);
INSERT INTO `user_messages` VALUES (149, 'eric_eric', 5, 5, '', '13d76b07-dee2-42ac-b768-2ae2934de9f5.jpg', '', 3, '2024-09-13 09:35:15.166', '2024-09-13 09:35:15.166', 0);
INSERT INTO `user_messages` VALUES (150, 'eric_sam', 5, 6, 'dddd', '', '', 1, '2024-09-13 09:43:45.076', '2024-09-13 09:43:45.076', 0);
INSERT INTO `user_messages` VALUES (151, 'eric_sam', 5, 6, 'u know,for search', '', '', 1, '2024-09-13 09:45:13.824', '2024-09-13 09:45:13.824', 0);
INSERT INTO `user_messages` VALUES (152, 'eric_sam', 5, 6, 'd', '', '', 1, '2024-09-13 09:47:29.302', '2024-09-13 09:47:29.302', 0);

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'id',
  `uuid` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'uuid',
  `username` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '\'用户名\'',
  `nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '昵称',
  `email` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '邮箱',
  `password` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `avatar` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '\'头像\'',
  `create_at` datetime(3) NULL DEFAULT NULL,
  `update_at` datetime(3) NULL DEFAULT NULL,
  `delete_at` bigint NULL DEFAULT NULL,
  `jwt` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `username`(`username` ASC) USING BTREE,
  UNIQUE INDEX `idx_uuid`(`uuid` ASC) USING BTREE,
  UNIQUE INDEX `Emali`(`email` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 14 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (5, '28353ed6-5966-4804-9c52-9b00abd4401e', 'eric', 'newMe', '163@163.com', '$2a$10$cJc3LdVegpT4sJ.QvjJJu.gDOVyC53y8K98jRQldQOW6.dK4.E3gC', '6b476aca-3c11-4bdb-b9cb-5b78fec2eddd.jpg', '2024-07-12 10:09:14.109', '2024-07-30 15:48:45.909', 0, NULL);
INSERT INTO `users` VALUES (6, 'c0fd020f-fa4d-4379-8f8b-f84a07fadd6a', 'sam', 'samual', '88@163.com', '$2a$10$XJos2XseT9hVMI8/u2mjpeNX2bRFf/81aJLahPryYmCxtgqe6m1jW', '', '2024-07-12 11:15:18.005', '2024-07-12 18:06:19.555', 0, NULL);
INSERT INTO `users` VALUES (13, 'a1c6d308-46cc-46a7-bb44-94dadbcf6b89', 'alpha', 'alpha beta', '2222@222.com', '$2a$10$gb/Kqz3yAfTKAKC4Encp5eAsNnnHzsulIKfLP5K1YqdrSwC4YxQXm', '', '2024-07-19 13:42:12.686', NULL, 0, NULL);
INSERT INTO `users` VALUES (14, 'b99b7135-e25c-4f23-b6a6-7d17335ec9f6', 'john', 'john', '', '$2a$10$2T5voPj6n9suYi6HHPxrN.BssNAJcYCdCXnfkejRLOeMwf4aDKzba', '', '2024-07-22 08:04:32.463', NULL, 0, '');

SET FOREIGN_KEY_CHECKS = 1;
