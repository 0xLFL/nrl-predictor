#!/usr/bin/env python3
import sys
import psycopg2

def main():
    conn = psycopg2.connect(dbname="mymyunsw")
    cur = conn.cursor()

    cur.close()
    conn.close()

if __name__ == "__main__":
    main()