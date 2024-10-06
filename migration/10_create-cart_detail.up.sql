CREATE TABLE IF NOT EXISTS `cart_detail` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    cart_id BIGINT(20) NOT NULL,
    product_id BIGINT(20) NOT NULL,
    quantity INT(11),
    price DECIMAL(12,2),
    total_price DECIMAL(12,2),
    FOREIGN KEY (cart_id) REFERENCES cart(id),
    FOREIGN KEY (product_id) REFERENCES product(id)
);