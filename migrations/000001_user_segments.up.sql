CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE segments
(
    id SERIAL PRIMARY KEY,
    slug VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE user_segments
(
    user_id INTEGER REFERENCES users(id),
    segment_slug VARCHAR(100) REFERENCES segments(slug),
    expiration_date TIMESTAMP,
    PRIMARY KEY (user_id, segment_slug)
);

CREATE TABLE user_segment_history (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    segment_slug TEXT NOT NULL,
    operation TEXT NOT NULL,
    operation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);