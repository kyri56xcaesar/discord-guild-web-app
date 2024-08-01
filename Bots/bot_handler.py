import sqlite3
import logging

import database_handler as dh


logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

create_tables_sql = "init_create_tables.sql"
init_populate_sql = "populate_current_tables.sql"


def main():

    con = sqlite3.connect("myBots_databases.db")


    init_database_system(con)

    try:
        dh.select_all_from_table(con, "bots")
    except ValueError as ve:
        logger.info(f'Invalid table name: {ve}')
    finally:
        con.close()
    

    return 0



def init_database_system(con):
    
    cur = con.cursor()


    with open(create_tables_sql, "r", encoding='utf-8') as f:
        sql = f.read()
    with open(init_populate_sql, "r", encoding='utf-8') as f2:
        sql2 = f2.read()
        
    try:
        cur.executescript(sql)
        cur.executescript(sql2)
        
    
    except sqlite3.OperationalError as e:
        logger.error(f'SQL error {e}')
    
            
    con.commit()
    


if __name__ == "__main__":
    logger.info(main())