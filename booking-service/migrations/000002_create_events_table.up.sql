CREATE TABLE `events_snapshot` (
    `id`                BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `public_id`         CHAR(36) NOT NULL UNIQUE,
    `title`             VARCHAR(256) NOT NULL,
    `start_date`        DATETIME NOT NULL, 
    `end_date`          DATETIME NOT NULL, 
    `available_seats`   BIGINT UNSIGNED NOT NULL,
    `created_at`        DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at`        DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    `deleted_at`        DATETIME(3) NULL,
    INDEX `idx_events_deleted_at` (`deleted_at`),
    INDEX `idx_events_start_date` (`start_date`)
) ENGINE = InnoDB;