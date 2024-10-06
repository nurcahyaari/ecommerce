CREATE TABLE `order_address` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    order_id BIGINT(20) NOT NULL,
    user_id BIGINT(20) NOT NULL,
    full_address TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT(20),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT(20),
    deleted_at DATETIME DEFAULT NULL,
    deleted_by BIGINT(20),
    FOREIGN KEY (user_id) REFERENCES `user`(id),
    FOREIGN KEY (order_id) REFERENCES `order`(id)
);