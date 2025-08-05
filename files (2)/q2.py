#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import psycopg2

def main():
    if len(sys.argv) != 2:
        print("Usage: ./q2.py <SubjectCode>")
        sys.exit(1)
    subject_code = sys.argv[1]

    conn = psycopg2.connect(dbname="mymyunsw")
    cur = conn.cursor()

    cur.close()
    conn.close()

if __name__ == "__main__":
    main()