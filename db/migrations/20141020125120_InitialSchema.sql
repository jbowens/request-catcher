
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- TODO(jackson): Add headers, as blob, or in a separate table? idk.

CREATE TABLE `requests` (
  id              BIGINT NOT NULL AUTO_INCREMENT,
  host            VARCHAR(100) NOT NULL,
  when            TIMESTAMP NOT NULL,
  method          VARCHAR(20) NOT NULL,
  path            VARCHAR(255) NOT NULL,
  content_length  INT NOT NULL,
  remote_addr     VARCHAR(50) NOT NULL,
  body            TEXT NOT NULL,
  raw_request     TEXT NOT NULL
) ENGINE=InnoDB;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE `requests`;
