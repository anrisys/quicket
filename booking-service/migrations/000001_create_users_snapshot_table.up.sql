CREATE TABLE users_snapshot (
    id BIGINT UNSIGNED PRIMARY KEY,
    public_id CHAR(36) NOT NULL UNIQUE,

    email VARCHAR(255) NULL,
    full_name VARCHAR(255) NULL,

    registered_at DATETIME(3) NOT NULL,

    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
        ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,

    INDEX idx_users_public_id (public_id),
    INDEX idx_users_registered_at (registered_at),
    INDEX idx_users_deleted_at (deleted_at)
) ENGINE = InnoDB;