import sqlite3
import re
import sys


sys.path.append("../Bots/")

# Import the module from the parent directory
import bot_utils as bu

def is_valid_table_name(table_name):
    # Check if the table name contains only alphanumeric characters and underscores
    return re.match(r'^[A-Za-z0-9_]+$', table_name) is not None

# select all
def select_all_from_table(con, table_name):
    if not is_valid_table_name(table_name):
        raise ValueError(f"Invalid table name: {table_name}")

    cur = con.cursor()
    
    # Use a safe method to construct the query
    query = f'SELECT * FROM "{table_name}";'
    
    try:
        cur.execute(query)
        rows = cur.fetchall()
        bu.format_row_data_print(rows)
    except sqlite3.Error as e:
        print(f"An error occurred: {e}")

# select by ID
def getBotByID():
    pass
# select by BotName
def getBotByName():
    pass

# select trigger by ID
# select trigger by phrase
# select trigger by datetime

# select line by ID
# select line by phrase
# select line by datetime


def insert_trigger_words_in_table(db_name, table_name, data):
    if not is_valid_table_name(table_name):
        raise ValueError(f"Invalid table name: {table_name}")

    con = sqlite3.connect(db_name)

    cur = con.cursor()
    
    
    query = f'INSERT INTO {table_name} (bot_id, phrase, author) VALUES (?, ?, ?)'
    try:
        cur.executemany(query, data)
        res = cur.fetchone()
    except sqlite3.Error as e:
        print(f"An error occurred: {e}")
        con.close()
    
    con.commit()

    con.close()
    
    return res    

def insert_trigger_words_in_table(db_name, table_name, data):
    if not is_valid_table_name(table_name):
        raise ValueError(f"Invalid table name: {table_name}")

    con = sqlite3.connect(db_name)
    cur = con.cursor()
    
    
    query = f'INSERT INTO {table_name} (bot_id, phrase, author) VALUES (?, ?, ?)'
    try:
        cur.executemany(query, data)
        res = cur.fetchone()
    except sqlite3.Error as e:
        print(f"An error occurred: {e}")
        con.close()
    
    con.commit()
    
    con.close()

    return res  

# delete
# update



