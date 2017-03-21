#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <err.h>
#include <errno.h>
#include <string.h>

#include "threads.h"

enum {
        DEFSTACKSIZE = 4*1024,
        MAXTHREADS = 32,
        EXTRATH = 10,
};

int test = 1;

void
f1(void *arg)
{
        int i;
        char *argument;

        argument = (char*) arg;
        for (i = 0 ; i < 100000 ; i++) {
                printf("%s\n", "Hi f1!!");
                printf("Argument is : %s\n", argument);
//                sleep(1);
                yieldthread();
                printf("********* Current thread id : %d **********\n", curidthread());
//                print_attr_library();
        }
        //        print_attr_library();
        exitsthread();
}

void
f2(void *arg)
{
        int i;
        char *argument;

        argument = (char*) arg;
        for (i = 0 ; i < 100000 ; i++) {
                printf("%s\n", "Hi f2!!");
                printf("Argument is : %s\n", argument);
//                sleep(1);
                yieldthread();
                printf("********* Current thread id : %d **********\n", curidthread());
//                print_attr_library();
        }
//        print_attr_library();
        exitsthread();
}


void
f3(void *arg)
{
        int i;
        char *argument;

        argument = (char*) arg;
        for (i = 0 ; i < 100000 ; i++) {
                printf("%s\n", "Hi f3!!");
                printf("Argument is : %s\n", argument);
//                sleep(1);
                yieldthread();
                printf("********* Current thread id : %d **********\n", curidthread());
//                print_attr_library();
        }
//        print_attr_library();
        exitsthread();
}

void
test0()
{
        initthreads(); // Already initialized the threads...
        // Test yield con el thread 0...
        yieldthread();
//        print_attr_library();
        sleep(1);
        yieldthread();
//        print_attr_library();
        printf("%s\n", "Test yield with thread 0 ok!!");
}

// Test max Threads...
void
test_max_th()
{
        int id_th_created = -1;
        int i;

        for (i = 0; i < MAXTHREADS + EXTRATH ; i++) {
                if (i%2 == 0) {
                        if ((id_th_created = createthread(f1 , "My firsts functions..." , DEFSTACKSIZE)) < 0) {
                                printf("%s\n", "Error creaing thread.");
                        }
                } else {
                        if ((id_th_created = createthread(f2 , "My seconds functions..." , DEFSTACKSIZE)) < 0) {
                                printf("%s\n", "Error creaing thread.");
                        }
                }
                if (id_th_created >= 0)
                        fprintf(stdout, "Thread with id %d created.\n", id_th_created);
        }
}

// Test more threads...
void
test_more_th()
{
        int id_th_created = -1;
        int i;

        for (i = 0 ; i < MAXTHREADS + EXTRATH ; i++) {
                if ((id_th_created = createthread(f3 , "My thirds functions... :)" , DEFSTACKSIZE)) < 0) {
                        printf("%s\n", "Error creaing thread.");
                } else {
                        printf("Thread with id : %d created.\n", id_th_created);
                }
        }
}

int
main(int argc, char const *argv[])
{
        int i;

        initthreads();
        if (test) {
                test0();
                test_max_th();
        }

        // Wait in main to follow test...
        for (i = 0 ; i < 12 ; i++) {
                fprintf(stderr, "In main : %d \n" , i);
                sleep(1);
                yieldthread();
        }
        if (test){
                test_more_th();
        }
        for (i = 0 ; i < 2 ; i++) {
                if (test) {
                        printf("%s\n", "Hi main... eyyy!!! we are been executing other things!!!!");
                }
                sleep(1);
                yieldthread();
        }
//        print_attr_library();
        printf("%s\n" , "Bye main!");
        exitsthread();
        while (1){}; // Test exit...
}
