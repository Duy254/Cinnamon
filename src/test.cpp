/*
Copyright (C) 2008-2010
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

#include "stdafx.h"
#ifdef TEST_MODE
#include <time.h>
#include "maindefine.h"
#include "gen.h"
#include "butterfly.h"
#include "utility.h"
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include "utility.h"
#include <sys/stat.h>
#include "extern.h"
#include "eval.h"
#include "zobrist.h"
#ifdef _MSC_VER
#include <windows.h>
#endif
#define LUNG_FILE 300

typedef struct {
  char fen[LUNG_FILE][200];
  int side[LUNG_FILE];
} Mnode;
Mnode M;

typedef struct teststructt {
  char val[200];
  double eval;
} teststruct;

void
test_epd ( char *testfile ) {
  printf ( "\n ****** START TEST %s max time %d sec *******", testfile, MAX_TIME_MILLSEC / 1000 );
  struct timeb start, end;
  ftime ( &start );
  int count = 0;
  int foundit = 0;
  num_tot_moves = 0;
  FILE *stream;
  char line[2001];
  int side;
  strcpy ( line, testfile );
  stream = fopen ( line, "r" );
  if ( !stream ) {
    memset ( line, 0, sizeof ( line ) );
    strcpy ( line, "../" );
    strcat ( line, testfile );
    stream = fopen ( line, "r" );
  }
  myassert ( stream, "test error file not found" );
  while ( fgets ( line, sizeof ( line ), stream ) != NULL ) {
    //strcpy (line,"r3kr2/pppb1p2/2n3p1/3Bp2p/4P2N/2P5/PP3PPP/2KR3R b q - bm O-O-O; ");
    printf ( "\n%s", line );
    side = loadfen ( line );
    count++;
    print (  );
    do_move ( side );
    if ( strstr ( test_ris, test_found ) ) {
      foundit++;
      printf ( "\nOK" );
    }
    else {
      printf ( "\nKO" );
    }
    if ( mate ) {
      printf ( " MATE" );
    }
    printf ( " RESULT: (%s %s) %s \nfound %d/%d", test_found, test_ris, line, foundit, count );
  }
  fclose ( stream );
  ftime ( &end );
  printf ( "\ntime:  %d total nodes per whole test %I64u", diff_time ( end, start ), num_tot_moves );
  printf ( "\n ****** END TEST  %s found %d/%d *******", testfile, foundit, count );
}


void
test (  ) {
  //MAX_TIME_MILLSEC = 12000;
  //MAX_DEPTH_TO_SEARCH = 6;
  test_epd ( "wac.epd" );
  test_epd ( "sbd.epd" );
  test_epd ( "kaufman.epd" );
  test_epd ( "zugzwang.epd" );
  test_epd ( "bk.epd" );
  test_epd ( "mate.epd" );	//generated by http://www.frayn.net/beowulf/matetest.zip
  test_epd ( "arasan12.epd" );
}

#endif
