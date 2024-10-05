CREATE DATABASE ecommerce;

USE ecommerce;

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

CREATE TABLE `auth_refresh_token` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id BIGINT(20) NOT NULL,
    token VARCHAR(60) NOT NULL,
    is_actived BOOLEAN DEFAULT 0 NOT NULL,
    expired_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES user(id)
);

CREATE TABLE `user_address` (
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

CREATE TABLE `store` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL,
    is_actived BOOLEAN DEFAULT 0 NOT NULL,
    is_removed BOOLEAN DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT(20),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT(20),
    deleted_at DATETIME DEFAULT NULL,
    deleted_by BIGINT(20)
);

CREATE TABLE `warehouse` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    store_id BIGINT(20) NOT NULL,
    name VARCHAR(128) NOT NULL,
    location VARCHAR(256),
    is_actived BOOLEAN DEFAULT 0 NOT NULL,
    is_removed BOOLEAN DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT(20),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT(20),
    deleted_at DATETIME DEFAULT NULL,
    deleted_by BIGINT(20),
    FOREIGN KEY (store_id) REFERENCES store(id)
);

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

CREATE TABLE `product_stock` (
    product_id BIGINT(20) PRIMARY KEY,
    stock_reserved INT(11) UNSIGNED,
    stock_on_hand INT(11) UNSIGNED,
    FOREIGN KEY (product_id) REFERENCES product(id)
);

CREATE TABLE `cart` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id BIGINT(20) NOT NULL,
    address_id BIGINT(20) NOT NULL,
    total_price DECIMAL(12,2),
    total_quantity INT(11),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT(20),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT(20),
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (address_id) REFERENCES user_address(id)
);

CREATE TABLE `cart_detail` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cart_id BIGINT(20) NOT NULL,
    product_id BIGINT(20) NOT NULL,
    quantity INT(11),
    price DECIMAL(12,2),
    total_price DECIMAL(12,2),
    FOREIGN KEY (cart_id) REFERENCES cart(id),
    FOREIGN KEY (product_id) REFERENCES product(id)
);

CREATE TABLE `order` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id BIGINT(20) NOT NULL,
    total_price DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_quantity INT(11) NOT NULL DEFAULT 0,
    order_code VARCHAR(20) NOT NULL,
    is_removed BOOLEAN DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT(20),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT(20),
    deleted_at DATETIME DEFAULT NULL,
    deleted_by BIGINT(20),
    FOREIGN KEY (user_id) REFERENCES user(id)
);

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

CREATE TABLE `order_receipt` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    order_id BIGINT(20) NOT NULL,
    address_id BIGINT(20) NOT NULL,
    total_price DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_quantity INT(11) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT(20),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT(20),
    deleted_at DATETIME DEFAULT NULL,
    deleted_by BIGINT(20),
    FOREIGN KEY (order_id) REFERENCES `order`(id),
    FOREIGN KEY (address_id) REFERENCES order_address(id)
);

CREATE TABLE `order_detail` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    order_id BIGINT(20) NOT NULL,
    order_receipt_id BIGINT(20) NOT NULL,
    product_id BIGINT(20) NOT NULL,
    product_price DECIMAL(12,2) NOT NULL DEFAULT 0,
    quantity INT(11) NOT NULL DEFAULT 0,
    total_price DECIMAL(12,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT(20),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT(20),
    deleted_at DATETIME DEFAULT NULL,
    deleted_by BIGINT(20),
    FOREIGN KEY (order_id) REFERENCES `order`(id),
    FOREIGN KEY (order_receipt_id) REFERENCES order_receipt(id),
    FOREIGN KEY (product_id) REFERENCES product(id)
);