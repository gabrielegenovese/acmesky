import random
from datetime import datetime
import time

    
def str_time_prop(start, end, time_format, prop):
    """Get a time at a proportion of a range of two formatted times.

    start and end should be strings specifying times formatted in the
    given format (strftime-style), giving an interval [start, end].
    prop specifies how a proportion of the interval to be taken after
    start.  The returned time will be in the specified format.
    """

    stime = time.mktime(time.strptime(start, time_format))
    etime = time.mktime(time.strptime(end, time_format))

    ptime = stime + prop * (etime - stime)

    return time.strftime(time_format, time.localtime(ptime))


def random_date(start, end, prop):
    return str_time_prop(start, end, '%m/%d/%Y %I:%M %p', prop)

    
random.randint(3, 9)

f = open("flights.txt", "w")
for i in range(35, 1000):
    first = random.randint(1, 20)
    second = random.randint(1, 20)
    while(first == second):
        second = random.randint(1, 20)
    depart = random_date("9/9/2024 1:30 PM", "1/1/2025 4:50 AM", random.random())
    arrival = datetime.strptime(random_date(depart, "1/1/2026 4:50 AM", random.random()), '%m/%d/%Y %I:%M %p')
    f.write(f"({i}, {first}, {second}, {datetime.strptime(depart, '%m/%d/%Y %I:%M %p')}, {arrival}, {random.randint(50,100)}, {round(random.uniform(40,500), 2)}),\n")

f.write(f"(1000, {first}, {second}, {datetime.strptime(depart, '%m/%d/%Y %I:%M %p')}, {arrival}, {random.randint(50,100)}, {round(random.uniform(40,500), 2)});\n")
f.close()
