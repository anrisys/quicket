CREATE TABLE bookings (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    public_id CHAR(36) NOT NULL UNIQUE,

    event_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,

    seats BIGINT UNSIGNED NOT NULL,
    total_price DECIMAL(12,2) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'USD',

    status ENUM ('pending', 'success', 'failed', 'cancelled') 
        NOT NULL DEFAULT 'pending',

    payment_method ENUM (
        'credit_card',
        'bank_transfer',
        'e_wallet',
        'cash'
    ) NOT NULL,

    channel ENUM (
        'web',
        'mobile',
        'partner'
    ) NOT NULL,

    confirmed_at DATETIME(3) NULL,
    expired_at DATETIME(3) NOT NULL,

    metadata JSON NULL,
    /*
      examples:
      {
        "promo_code": "NEWUSER50",
        "partner_ref": "abc-123",
        "notes": "manual adjustment"
      }
    */

    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
        ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,

    FOREIGN KEY (event_id)
        REFERENCES events_snapshot (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    FOREIGN KEY (user_id)
        REFERENCES users_snapshot (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    INDEX idx_bookings_public_id (public_id),
    INDEX idx_bookings_event_id (event_id),
    INDEX idx_bookings_user_id (user_id),
    INDEX idx_bookings_status (status),
    INDEX idx_bookings_payment_method (payment_method),
    INDEX idx_bookings_channel (channel),
    INDEX idx_bookings_created_at (created_at),
    INDEX idx_bookings_confirmed_at (confirmed_at),
    INDEX idx_bookings_deleted_at (deleted_at)
) ENGINE = InnoDB;
