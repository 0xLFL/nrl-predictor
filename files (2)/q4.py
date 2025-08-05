#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import re
import psycopg2

def main():
    if len(sys.argv) != 2:
        print("Usage: ./q4.py <filter_expr>")
        sys.exit(1)

    filter_expr = sys.argv[1]

    conn = psycopg2.connect(dbname="mymyunsw")

    cur = conn.cursor()

    cur.close()
    conn.close()

if __name__ == "__main__":
    main()