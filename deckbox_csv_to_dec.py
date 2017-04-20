import sys
import csv

file_path = sys.argv[1]

csv_file_in = open(file_path, 'rb')
reader = csv.reader(csv_file_in)
rows = []
for row in reader:
    if len(row) >= 2:
        rows.append([row[0], row[2]])

csv_file_in.close()

rows = rows[1:]

with open(file_path + '.dec', 'wb') as output_file:
	for row in rows:
		output_file.write(row[0] + ' ' + row[1] + '\n')

