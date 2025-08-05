#!/usr/bin/env python3
import sys
import re
import psycopg2
from helpers import get_student, get_program, get_stream

def main():
    argc = len(sys.argv)
    if argc < 2:
      print(f"Usage: {sys.argv[0]} zID [Program Stream]")
      exit(1)
    zid = sys.argv[1]
    if zid[0] == 'z':
        zid = zid[1:8]
    digits = re.compile("^\d{7}$")
    if not digits.match(zid):
        print("Invalid zID")
        exit(1)

    prog_code = None
    strm_code = None

    if argc >= 3:
        prog_code = sys.argv[2]
    if argc >= 4:
        strm_code = sys.argv[3]

    conn = psycopg2.connect("dbname=mymyunsw")
    cur = conn.cursor()
    
    stu_info = get_student(conn,zid)
    if not stu_info:
        print(f"Invalid student id {zid}")
        exit(1)
    #print(stu_info) # debug

    if prog_code:
        prog_info = get_program(conn,prog_code)
        if not prog_info:
            print(f"Invalid program code {prog_code}")
            exit(1)
            #print(prog_info)  #debug

    if strm_code:
        strm_info = get_stream(conn,strm_code)
        if not strm_info:
            print(f"Invalid stream code {strm_code}")
            exit(1)
    # your code goes here
    

    cur.close()
    conn.close()
if __name__ == "__main__":
    main()




