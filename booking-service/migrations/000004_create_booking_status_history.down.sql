CREATE TABLE booking_status_history (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    booking_id BIGINT UNSIGNED NOT NULL,

    from_status ENUM (
        'pending',
        'success',
        'failed',
        'cancelled'
    ) NULL,

    to_status ENUM (
        'pending',
        'success',
        'failed',
        'cancelled'
    ) NOT NULL,

    reason VARCHAR(255) NULL,
    changed_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    FOREIGN KEY (booking_id)
        REFERENCES bookings (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    INDEX idx_bsh_booking_id (booking_id),
    INDEX idx_bsh_to_status (to_status),
    INDEX idx_bsh_changed_at (changed_at)
) ENGINE = InnoDB;
