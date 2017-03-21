#include <stdio.h>
#include <stdlib.h>
#include <ucontext.h>
#include <err.h>
#include <errno.h>
#include <string.h>
#include <sys/time.h>

#include "threads.h"

enum {
        MAXTHREADS = 32,
        DEAD = 0,
        READY = 1,
        RUNNING = 2,
        BLOCKED = 3,
        QUANTUM = 200 // ms.
};

typedef struct Thread {
        int id;
        ucontext_t uct;
        char *sp;
        int state;
        long ms;
} Thread;

typedef struct Threads {
        Thread th[MAXTHREADS];
        int n_th;
        int current_th;
} Threads;

/*
* Components of library.
*/
Threads g_threads;
int count_thread_id;
int threads_initialized;


static long
getms(struct timeval t)
{
        return (t.tv_sec * 1000) + (t.tv_usec / 1000);
}

static long
gettime()
{
        struct timezone zone;
        struct timeval now;

        if (gettimeofday(&now , &zone) < 0) {
                err (1 , "Error in function gettime : %s" , strerror(errno));
        }
        return getms(now);
}

static long
istime(Thread th)
{
        return  (gettime() - th.ms) <= QUANTUM;
}

static int
isrunning(Thread th)
{
        return th.state == RUNNING;
}

static int
is_thread_running_out_time()
{
        return isrunning(g_threads.th[g_threads.current_th]) &&
                !istime(g_threads.th[g_threads.current_th]);
}

static int
freeStack(Thread th)
{
        return th.state == DEAD && th.sp != NULL;
}

/*
* Test stack without free and
* free the stack's whose are not in use...
*/
static void
garbage_collector()
{
        int i;
        for (i = 0 ; i < MAXTHREADS ; i++){
                if (freeStack(g_threads.th[i])) {
                        free(g_threads.th[i].sp);
                        g_threads.th[i].sp = NULL;
                }
        }
}

/*
* getfollow_th get the follow thread with state that you define like argument, if is found
* the thread with that state the pos in thread array that have this state is returning
* if not threads have that state getfollow_th return -1 and not thread have that stated.
* inside the array...
*/
static int
getfollow_th(int state)
{
        enum {
                FOLLOW = 1,
        };
        int  i = g_threads.current_th + FOLLOW, next_th = -1, count = 0;

        while (next_th < 0 && count < MAXTHREADS) {

                if (i >= MAXTHREADS) {
                        // Circular iterative in threads array...
                        i = 0;
                }

                if (g_threads.th[i].state == state) {
                        /* Look for the follow thread with a state... */
                        next_th = i;
                }
                count++;
                i++;
        }
        return next_th;
}

static int
th_alone (int next_th)
{
        return next_th < 0 && isrunning(g_threads.th[g_threads.current_th]);
}

/*
* The scheduler change the context to the follow thread in wait ready to running.
* If you don't considerade the time to the follow thread scheduler change the context
* to the new thread althought the quantum was less than threshold..., scheduler always
* change the context to the follow thread ready to running without test nothing.
*/
static void
scheduler()
{
        int next_th = -1;
        int old_th = g_threads.current_th;

        garbage_collector();

        // Change context to follow thread ready...
        next_th = getfollow_th(READY);

        if (th_alone(next_th)) {
                // Only one thread is running anyone thread is ready...
                // we have to put next thread the same thread we were running.
                next_th = g_threads.current_th;
        }

        if (next_th < 0) {
                // There are not more thread to run... The library of
                // collaborative threads have finished succesfull.
                exit(0);
        }

        if (isrunning(g_threads.th[old_th])) {
                // The thread must be running to change the state to ready...
                g_threads.th[old_th].state = READY;
        }

        g_threads.th[next_th].state = RUNNING;
        g_threads.current_th = next_th;
        g_threads.th[next_th].ms = gettime(); // Update time...

        if (swapcontext(&g_threads.th[old_th].uct , &g_threads.th[next_th].uct) < 0) {
                err (1, "Error with swapcontext : %s\n" , strerror(errno));
        }
}

void
initthreads()
{
        if (threads_initialized) {
                fprintf(stderr , "%s\n", "Threads are just initialized...");
                return;
        }

        int i;
        for (i = 0; i < MAXTHREADS; i++) {
                g_threads.th[i].id = -1;
                g_threads.th[i].state = DEAD;
        }

        g_threads.th[g_threads.n_th].id = 0;
        // makecontext... Not necessary main is the flow of the thread 0...
        g_threads.th[g_threads.n_th].state = RUNNING;
        g_threads.th[g_threads.n_th].ms = gettime();
        g_threads.current_th = 0;       // Position of thread in data array of threads...
        g_threads.n_th++;

        threads_initialized = 1;        // library threads initializaded...
        count_thread_id = 0;            // Counter of numbers of tid assigned...

        if (getcontext(&(g_threads.th[g_threads.current_th].uct)) < 0) {
                err (1 , "Error getting context: %s" , strerror(errno));
        }
}

static void
print_state(Thread th)
{
        printf("%s", "State : ");
        switch (th.state) {
        case DEAD :
                printf("%s\n", "DEAD");
                break;
        case READY :
                printf("%s\n", "READY");
                break;
        case RUNNING :
                printf("%s\n", "RUNNING");
                break;
        case BLOCKED :
                printf("%s\n", "BLOCKED");
                break;
        default :
                printf("%s\n", "NOT STATE");
        }
}


static void
print_thread(Thread th)
{
        long ms = 0.0;

        if (isrunning(th)) {
                // Time of the thread running in cpu
                // when we call to print_thread.
                ms = gettime() - th.ms;
        }
        printf("Thread id : %2d | time : %3ld ms of CPU | sp : %9p | ", th.id, ms , th.sp);
        print_state(th);
}

static void
print_threads()
{
        printf("%s\n", "");
        printf("%s\n", "Prints Table of threads");
        printf("%s\n", "=======================");
        printf("Current thread id running... : %d\n", g_threads.th[g_threads.current_th].id);
        printf("Position in array of current thread running : %d\n", g_threads.current_th);
        printf("Number of threads : %d\n", g_threads.n_th);
        int i;
        for (i = 0; i < MAXTHREADS ; i++) {
                print_thread(g_threads.th[i]);
        }
        printf("%s\n", "=======================");
        printf("%s\n", "");
}

void
print_attr_library()
{
        printf("Counter thread id last number id was assigend : %d\n", count_thread_id);
        printf("Library threads initialized : %d\n", threads_initialized);
        print_threads();
}

int
createthread(void (*mainf)(void*), void *arg, int stacksize)
{
        char *stack_f = NULL;
        int next_th = -1;

        next_th = getfollow_th(DEAD);
        if (next_th < 0) {
                fprintf(stderr, "%s\n", "We can't create a new thread anyone thread is dead...");
                return -1;
        }

        stack_f = malloc(stacksize);
        if (stack_f == NULL) {
                fprintf(stderr, "%s%s\n", "Error allocating memory to thread stack" , strerror(errno));
                return -1;
        }

        // Set attributes of the new thread created...
        g_threads.th[next_th].id = ++count_thread_id;
        g_threads.th[next_th].state = READY;
        g_threads.th[next_th].sp = stack_f;
        g_threads.th[next_th].ms = gettime();

        // Create context to thread.
        if (getcontext(&g_threads.th[next_th].uct) < 0) {
                fprintf(stderr, "Error creating context to the thread : %s\n", strerror(errno));
                return -1;
        }
        g_threads.th[next_th].uct.uc_stack.ss_sp = stack_f;
        g_threads.th[next_th].uct.uc_stack.ss_size = stacksize;
        g_threads.th[next_th].uct.uc_link = NULL;
        makecontext(&g_threads.th[next_th].uct, (void *) mainf, 1 , arg , NULL);

        // New thread was created we return its id.
        g_threads.n_th++;
        return g_threads.th[next_th].id;
}

void
exitsthread()
{
        g_threads.th[g_threads.current_th].state = DEAD;
        g_threads.th[g_threads.current_th].id = -1;
        // The thread was deleted we pass the context
        // to follow thread that have to run.
        g_threads.n_th--;
        scheduler();
}

void
yieldthread()
{
        if (is_thread_running_out_time()) {
                scheduler();
        }
}

int
curidthread()
{
        return g_threads.th[g_threads.current_th].id;
}
