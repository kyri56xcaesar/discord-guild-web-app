-- DROP TABLE bots;
CREATE TABLE IF NOT EXISTS bots (
    id integer primary key AUTOINCREMENT,
    bot_name char(20)
);

CREATE TABLE IF NOT EXISTS trigger_words (
    trigger_id integer primary key AUTOINCREMENT,
    bot_id integer ,
    phrase varchar(255),
    author char(20),
    data_datetime DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS typical_lines (
    line_id integer primary key AUTOINCREMENT,
    bot_id integer ,
    phrase varchar(255),
    author char(20),
    data_datetime DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sad_lines (
    line_id integer primary key AUTOINCREMENT,
    bot_id integer,
    phrase varchar(255),
    author char(20),
    data_datetime DATETIME DEFAULT CURRENT_TIMESTAMP
);