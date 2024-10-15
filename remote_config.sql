/*
 Navicat Premium Data Transfer

 Source Server         : remote_config
 Source Server Type    : SQLite
 Source Server Version : 3035005 (3.35.5)
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3035005 (3.35.5)
 File Encoding         : 65001

 Date: 15/10/2024 19:43:22
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for config
-- ----------------------------
DROP TABLE IF EXISTS "config";
CREATE TABLE "config"
(
    "id"      integer NOT NULL,
    "name"    TEXT    NOT NULL,
    "value"   TEXT    NOT NULL,
    "message" TEXT    NOT NULL,
    "secret"  TEXT    NOT NULL DEFAULT '',
    "enable"  integer NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

-- ----------------------------
-- Table structure for config_history
-- ----------------------------
DROP TABLE IF EXISTS "config_history";
CREATE TABLE "config_history"
(
    "id"          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "config_id"   INTEGER NOT NULL,
    "old_value"   TEXT    NOT NULL,
    "new_value"   TEXT    NOT NULL,
    "nickname"    TEXT    NOT NULL,
    "enable"      integer NOT NULL,
    "create_time" integer NOT NULL
);

-- ----------------------------
-- Table structure for sqlite_sequence
-- ----------------------------
DROP TABLE IF EXISTS "sqlite_sequence";
CREATE TABLE "sqlite_sequence"
(
    "name",
    "seq"
);

-- ----------------------------
-- Indexes structure for table config
-- ----------------------------
CREATE INDEX "key_name"
    ON "config" (
                 "name" ASC
        );

-- ----------------------------
-- Indexes structure for table config_history
-- ----------------------------
CREATE INDEX "key_config_id"
    ON "config_history" (
                         "config_id" ASC
        );

PRAGMA foreign_keys = true;
