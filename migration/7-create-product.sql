CREATE TABLE `product` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    store_id BIGINT(20) NOT NULL,
    warehouse_id BIGINT(20) NOT NULL,
    name VARCHAR(128) NOT NULL,
    price DECIMAL(12,2) DEFAULT 0 NOT NULL,
    is_actived BOOLEAN DEFAULT 0 NOT NULL,
    is_removed BOOLEAN DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT(20),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT(20),
    deleted_at DATETIME DEFAULT NULL,
    deleted_by BIGINT(20),
    FOREIGN KEY (store_id) REFERENCES store(id),
    FOREIGN KEY (warehouse_id) REFERENCES warehouse(id)
);