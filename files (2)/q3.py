#!/usr/bin/env python3
import sys
import re
import psycopg2

def main():
    if len(sys.argv) != 2:
        print("Usage: ./q3.py <zID>")
        sys.exit(1)
    zid = sys.argv[1]
    if zid[0] == 'z':
        zid = zid[1:8]
    digits = re.compile("^\d{7}$")
    if not digits.match(zid):
        print("Invalid zID")
        exit(1)
    conn = psycopg2.connect(dbname="mymyunsw")
    cur = conn.cursor()

    cur.close()
    conn.close()

if __name__ == "__main__":
    main()