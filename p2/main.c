#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <err.h>
#include <errno.h>
#include <string.h>

#include "threads.h"

enum {
        DEFSTACKSIZE = 32*1024,
        EXTRATH = 0,
        MAXTHREADS = 32,
};

int test = 1;

void
print_susp(int *list , int n)
{
        int i;

        fprintf(stdout, "Number of threads suspended : %d\n", n);
        for (i = 0; i < n; i++) {
                fprintf(stdout, "Tid Thread Suspended : %d\n", list[i]);
        }
}

void get_susp()
{
        int *suspends = NULL , nsusp = -1;

        if ((nsusp = suspendedthreads(&suspends)) > 0){
                print_susp(suspends , nsusp);
                free(suspends);
        } else {
                fprintf(stdout , "%s\n", "There is no threads suspended...");
        }
}


void
f1(void *arg)
{
        int i , id;
        char *argument;

        argument = (char*) arg;

        for (i = 1 ; i < 2 ; i++) {
                fprintf(stdout , "%s\n", "Hi f1!!");
                fprintf(stdout , "Argument is : %s\n", argument);
                if ((id = killthread(i)) < 0) {
                        fprintf(stderr , "Error killing thread id : %d\n", i);
                } else {
                        fprintf(stdout , "Thread kill id : %d\n", i);
                }
                sleepthread(1500);
                yieldthread();
                fprintf(stdout , "Current thread id : %d\n", curidthread());
        }
        fprintf(stdout, "%s\n", "### Exit ###");
        exitsthread();
        while (1){}; // Test exit...
}

void
f2(void *arg)
{
        int i , resume;
        char *argument;

        argument = (char*) arg;

        for (i = 0 ; i < 5 ; i++) {
                get_susp();
                fprintf(stdout , "%s\n", "Hi f2!!");
                fprintf(stdout , "Argument is : %s\n", argument);
                sleepthread(1500);
                if ((resume = resumethread(i) >= 0)) {
                        fprintf(stdout , "Thread resume id : %d\n", i);
                } else {
                        fprintf(stderr , "Thread no resume id : %d\n", i);
                }
                yieldthread();
                fprintf(stdout , "Current thread id : %d\n", curidthread());
        }
        get_susp();
        fprintf(stdout , "%s\n", "### Exit ###");
        exitsthread();
        while (1){}; // Test exit...
}


void
f3(void *arg)
{
        int i;
        char *argument;

        argument = (char*) arg;

        for (i = 0 ; i < 2 ; i++) {
                fprintf(stdout , "%s\n", "Hi f3!!");
                fprintf(stdout , "Argument is : %s\n", argument);
                sleepthread(1500);
                yieldthread();
                fprintf(stdout, "Current thread id : %d\n", curidthread());
        }
        suspendthread();
        fprintf(stdout , "%s\n", "Wake up! :)");
        sleepthread(1500);
        yieldthread();
        fprintf(stdout , "%s\n", "### Exit ###");
        exitsthread();
        while (1){}; // Test exit...
}

int
main(int argc, char const *argv[])
{
        initthreads();
        fprintf(stdout, "%s\n", "Hi main!!");
        createthread(f3 , "SUSPEND" , DEFSTACKSIZE);
        createthread(f3 , "SUSPEND" , DEFSTACKSIZE);
        sleepthread(10000);
        yieldthread();
        createthread(f2 , "RESUME THREADS" , DEFSTACKSIZE);
        createthread(f1 , "KILL THREADS" , DEFSTACKSIZE);
        get_susp();
        if (killthread(7) < 0) {
                fprintf(stderr, "%s\n", "No thread 7 could be killed.");
        }
        suspendthread();
        sleepthread(20000);
        sleepthread(1000);
        fprintf(stderr, "%s\n", "Wake up! :) main...");
        fprintf(stderr, "%s\n", "### Exit ### main...");
        get_susp();
        int i;
        for (i = 0 ; i < 10; i++){
                yieldthread();
                fprintf(stderr, "%s", ".");
                sleepthread(1000);
        }
        fprintf(stderr, "%s\n", "");
        killthread(0);
        exitsthread();
        while (1){}; // Test exit...
}
