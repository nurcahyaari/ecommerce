
CREATE TABLE `user` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(128) NOT NULL,
    email VARCHAR(128),
    phone VARCHAR(15),
    `password` VARCHAR(128),
    is_actived BOOLEAN DEFAULT 0 NOT NULL,
    is_removed BOOLEAN DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT(20),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT(20),
    deleted_at DATETIME DEFAULT NULL,
    deleted_by BIGINT(20),
    UNIQUE (email),
    UNIQUE (phone)
);