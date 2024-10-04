import datetime
import sys
sys.path.append("../Bots/")

# Import the module from the parent directory
import bot_utils as bu

def encrypt(name):
    ciph = ""
    for c in name:
        pass

    return name


def save_data(name, contents):
    with open(f"{encrypt(name)}.log", "a") as f:
        # header
        f.write(f"{datetime.datetime.today()}  ::_-> lines: {len(contents)}  \n")
        for index, line in enumerate(contents):
            f.write(f"{index + 1}: {line}\n")
        f.write("-------------\n")



def format_row_data_print(data):
    for index, entry in enumerate(data):
        print(f'[ENTRY]:-> {index + 1}. {entry}')
        
        
def format_lines_to_insertable_data(bid, phrases, author):
    return [(bid, phrase, author) for phrase in phrases]