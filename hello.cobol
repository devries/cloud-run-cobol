       IDENTIFICATION DIVISION.
       PROGRAM-ID. HELLOWRD.
       
       DATA DIVISION.
       WORKING-STORAGE SECTION.
       01 arg-value PIC X(25).

       PROCEDURE DIVISION.
       DISPLAY 1 UPON ARGUMENT-NUMBER.
       ACCEPT arg-value FROM ARGUMENT-VALUE.

       IF arg-value = SPACE THEN
          DISPLAY "HELLO WORLD"
       ELSE
          DISPLAY "HELLO " arg-value
       END-IF.
       STOP RUN.
