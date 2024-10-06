CREATE TABLE IF NOT EXISTS `order_receipt` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    order_id BIGINT(20) NOT NULL,
    order_address_id BIGINT(20) NOT NULL,
    total_price DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_quantity INT(11) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT(20),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT(20),
    deleted_at DATETIME DEFAULT NULL,
    deleted_by BIGINT(20),
    FOREIGN KEY (order_id) REFERENCES `order`(id),
    FOREIGN KEY (order_address_id) REFERENCES order_address(id)
);