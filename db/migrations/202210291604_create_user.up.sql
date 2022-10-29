CREATE TABLE user (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(191),
    last_name VARCHAR(191),
    email VARCHAR(191),
    crypted_password VARCHAR(255),
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    UNIQUE (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
