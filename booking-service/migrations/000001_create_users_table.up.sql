CREATE TABLE `users_snapshot` (
    `id`          BIGINT UNSIGNED PRIMARY KEY,
    `public_id`   CHAR(36) NOT NULL UNIQUE,
    `created_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at`  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    `deleted_at`  DATETIME(3) NULL,
    INDEX `idx_users_deleted_at` (`deleted_at`)
) ENGINE = InnoDB;