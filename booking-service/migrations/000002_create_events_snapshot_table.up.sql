CREATE TABLE events_snapshot (
    id BIGINT UNSIGNED PRIMARY KEY,
    public_id CHAR(36) NOT NULL UNIQUE,

    title VARCHAR(256) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description TEXT NULL,

    location_city VARCHAR(100) NOT NULL,
    location_country VARCHAR(100) NOT NULL,

    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,

    base_price DECIMAL(10,2) NOT NULL,
    max_seats BIGINT UNSIGNED NOT NULL,
    available_seats BIGINT UNSIGNED NOT NULL,

    attributes JSON NULL,
    /*
      examples:
      {
        "tags": ["music", "festival"],
        "is_featured": true,
        "age_limit": 18
      }
    */

    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL,
    deleted_at DATETIME(3) NULL,

    INDEX idx_events_public_id (public_id),
    INDEX idx_events_category (category),
    INDEX idx_events_city (location_city),
    INDEX idx_events_start_date (start_date),
    INDEX idx_events_deleted_at (deleted_at)
) ENGINE = InnoDB;
