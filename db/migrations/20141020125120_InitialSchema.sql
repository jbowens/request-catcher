
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE `requests` (
  `id`              BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `host`            VARCHAR(100) NOT NULL,
  `when`            TIMESTAMP NOT NULL,
  `method`          VARCHAR(20) NOT NULL,
  `path`            VARCHAR(255) NOT NULL,
  `content_length`  INT NOT NULL,
  `remote_addr`     VARCHAR(50) NOT NULL,
  `body`            TEXT NOT NULL,
  `raw_request`     TEXT NOT NULL,
  `cleared`         TINYINT(1) NOT NULL DEFAULT 0
) ENGINE=InnoDB;

CREATE INDEX `host_when` ON `requests` (`host`, `when`);

CREATE TABLE `request_headers` (
  `id`              BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `request_id`      BIGINT NOT NULL,
  `key`             VARCHAR(100) NOT NULL,
  `value`           TEXT NOT NULL
) ENGINE=InnoDB;

CREATE INDEX `header_request` ON `request_headers` (`request_id`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE `requests`;
