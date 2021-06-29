-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE blogss
(
    id SERIAL PRIMARY KEY,
    userid INT NOT NULL,
    title VARCHAR(100) NOT NULL,
    message VARCHAR(1000) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    category_id INT NOT NULL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE blogss;