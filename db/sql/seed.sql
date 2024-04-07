DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS tokens;

CREATE TABLE documents (
  id      BIGINT NOT NULL AUTO_INCREMENT,
  title   TEXT(255) NOT NULL,
  body    LONGTEXT NOT NULL,
  PRIMARY KEY (id),
  UNIQUE (title(255))
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci ENGINE = INNODB;

CREATE TABLE tokens (
  id         BIGINT NOT NULL AUTO_INCREMENT,
  token      TEXT(255) NOT NULL,
  docs_count INT NOT NULL,
  postings   LONGBLOB NOT NULL,
  PRIMARY KEY (id),
  UNIQUE (token(255))
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci ENGINE = INNODB;
