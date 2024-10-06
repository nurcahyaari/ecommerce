CREATE TABLE IF NOT EXISTS `product_stock` (
    product_id BIGINT(20) PRIMARY KEY,
    stock_reserved INT(11) UNSIGNED,
    stock_on_hand INT(11) UNSIGNED,
    FOREIGN KEY (product_id) REFERENCES product(id)
);