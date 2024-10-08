CREATE TABLE IF NOT EXISTS `user_address` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id BIGINT(20) NOT NULL,
    full_address TEXT DEFAULT NULL,
    is_actived BOOLEAN DEFAULT 0 NOT NULL,
    is_removed BOOLEAN DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT(20),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT(20),
    deleted_at DATETIME DEFAULT NULL,
    deleted_by BIGINT(20),
    FOREIGN KEY (user_id) REFERENCES user(id)
);