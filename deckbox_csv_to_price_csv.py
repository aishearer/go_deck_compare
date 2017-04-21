#!/usr/bin/python
import sys
import csv

file_path = sys.argv[1]

csv_file_in = open(file_path, 'rb')
reader = csv.reader(csv_file_in)
rows = []
for row in reader:
    if len(row) >= 2:
        if len(row[2]) > 8:
            row[2] = "0.01"
        rows.append([row[2].replace("$", ""), row[1]])

csv_file_in.close()

csv_file_out = open(file_path, 'wb')
writer = csv.writer(csv_file_out)
writer.writerows(rows)
csv_file_out.close()
