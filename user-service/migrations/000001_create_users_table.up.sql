CREATE TABLE users (
    `id`          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `public_id`   CHAR(36) NOT NULL UNIQUE,
    `email`       VARCHAR(255) NOT NULL UNIQUE,
    `password`    VARCHAR(255) NOT NULL,
    `role`        ENUM('user', 'organizer', 'admin') NOT NULL DEFAULT 'user',
    `created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    `deleted_at`  DATETIME(3) NULL,
    INDEX `idx_users_deleted_at` (`deleted_at`)
) ENGINE = INNODB;