/*
    Cinnamon is a UCI chess engine
    Copyright (C) Giuseppe Cannella

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

/// stacktrace.h (c) 2008, Timo Bingmann from http://idlebox.net/
/// published under the WTFPL v2.0

#ifndef _STACKTRACE_H_
#define _STACKTRACE_H_
#if defined(_WIN32) || !defined(DEBUG_MODE)

static inline void print_stacktrace() { }

#else

#include <stdio.h>
#include <stdlib.h>
#include <execinfo.h>
#include <cxxabi.h>


/** Print a demangled stack backtrace of the caller function to FILE* out. */
static inline void print_stacktrace(FILE *out = stdout, unsigned int max_frames = 512) {
    fprintf(out, "stack trace:\n");
    // storage array for stack trace address data
    void **addrlist = (void **) malloc(max_frames + 1);
    // retrieve current stack addresses
    int addrlen = backtrace(addrlist, max_frames + 1);
    if (addrlen == 0) {
        fprintf(out, "  <empty, possibly corrupt>\n");
        free(addrlist);
        return;
    }
    // resolve addresses into strings containing "filename(function+address)",
    // this array must be free()-ed
    char **symbollist = backtrace_symbols(addrlist, addrlen);
    // allocate string which will be filled with the demangled function name
    size_t funcnamesize = 256;
    char *funcname = (char *) malloc(funcnamesize);
    // iterate over the returned symbol lines. skip the first, it is the
    // address of this function.
    for (int i = 1; i < addrlen; i++) {
        char *begin_name = 0, *begin_offset = 0, *end_offset = 0;
        // find parentheses and +address offset surrounding the mangled name:
        // ./module(function+0x15c) [0x8048a6d]
        for (char *p = symbollist[i]; *p; ++p) {
            if (*p == '(') {
                begin_name = p;
            } else if (*p == '+') {
                begin_offset = p;
            } else if (*p == ')' && begin_offset) {
                end_offset = p;
                break;
            }
        }
        if (begin_name && begin_offset && end_offset && begin_name < begin_offset) {
            *begin_name++ = '\0';
            *begin_offset++ = '\0';
            *end_offset = '\0';
            // mangled name is now in [begin_name, begin_offset) and caller
            // offset in [begin_offset, end_offset). now apply
            // __cxa_demangle():
            int status;
            char *ret = abi::__cxa_demangle(begin_name, funcname, &funcnamesize, &status);
            if (status == 0) {
                funcname = ret;        // use possibly realloc()-ed string
                fprintf(out, "  %s+%s\n", funcname, begin_offset);
            } else {
                // demangling failed. Output function name as a C function with
                // no arguments.
                fprintf(out, "  %s : %s()+%s\n", symbollist[i], begin_name, begin_offset);
            }
        } else {
            // couldn't parse the line? print the whole line.
            fprintf(out, "  %s\n", symbollist[i]);
        }
    }
    free(funcname);
    free(symbollist);
    free(addrlist);
}

#endif
#endif // _STACKTRACE_H_