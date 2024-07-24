import datetime


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
