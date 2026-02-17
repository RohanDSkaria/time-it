CREATE TABLE current_entry (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    task TEXT NOT NULL,
    start_time INTEGER NOT NULL
);
