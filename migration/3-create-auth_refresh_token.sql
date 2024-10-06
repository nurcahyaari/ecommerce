CREATE TABLE `auth_refresh_token` (
    id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    user_id BIGINT(20) NOT NULL,
    token VARCHAR(60) NOT NULL,
    is_actived BOOLEAN DEFAULT 0 NOT NULL,
    expired_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES user(id)
);