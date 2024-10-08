CREATE TABLE IF NOT EXISTS groups (
    id         serial         PRIMARY KEY,
    name       VARCHAR(100)   UNIQUE NOT NULL,
    members    VARCHAR(100)[],
    created_at TIMESTAMP             DEFAULT now(),
    updated_at TIMESTAMP             DEFAULT now()
);

CREATE TABLE IF NOT EXISTS songs (
    id          serial       PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    text        TEXT         NOT NULL,
    released_at VARCHAR(100) NOT NULL,
    link        VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP    DEFAULT now(),
    updated_at  TIMESTAMP    DEFAULT now()
);

CREATE TABLE IF NOT EXISTS group_song (
    id         serial       PRIMARY KEY,
    group_id   INTEGER      NOT NULL,
    song_id    INTEGER      NOT NULL,

    FOREIGN KEY (group_id) REFERENCES groups (id),
    FOREIGN KEY (song_id)  REFERENCES songs  (id) ON DELETE CASCADE,
    UNIQUE (group_id, song_id)
);
