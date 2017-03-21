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
        ERRSUSP = -2,    // Must be a negative. Number of error...
        ERRSLEEP = -3,   // Must be a negative. Number of error...
        DEAD = 0,
        READY = 1,
        RUNNING = 2,
        BLOCKED = 3,
        SUSPEND = 4,
        SLEEP = 5,
        QUANTUM = 200    // ms.
};

typedef struct Thread {
        int id;
        ucontext_t uct;
        char *sp;
        int state;
        long ms;
        long sleep;
} Thread;

typedef struct Threads {
        Thread th[MAXTHREADS];
        int current_th;
} Threads;

/*
* Components of library.
*/
Threads g_ths;
int count_thread_id;
int threads_initialized;

static int
isrunning(Thread th)
{
        return th.state == RUNNING;
}

static int
issleep(Thread th)
{
        return th.state == SLEEP;
}

static long
getms(struct timeval t)
{
        return (t.tv_sec * 1000L) + (t.tv_usec / 1000L);
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
        return gettime() - th.ms < QUANTUM;
}

static void
print_state(Thread th)
{
        fprintf(stderr, "%s", "State : ");
        switch (th.state) {
        case DEAD :
                fprintf(stderr, "%s\n", "DEAD");
                break;
        case READY :
                fprintf(stderr, "%s\n", "READY");
                break;
        case RUNNING :
                fprintf(stderr, "%s\n", "RUNNING");
                break;
        case BLOCKED :
                fprintf(stderr, "%s\n", "BLOCKED");
                break;
        case SUSPEND :
                fprintf(stderr, "%s\n", "SUSPEND");
                break;
        case SLEEP :
                fprintf(stderr, "%s\n", "SLEEP");
                break;
        default :
                fprintf(stderr, "%s\n", "NOT STATE");
        }
}

static void
print_thread(Thread th)
{
        long ms , sleep = 0.0;

        if (isrunning(th)) {
                // Time of the thread running in cpu.
                ms = gettime() - th.ms;
        } else if (issleep(th)) {
                // Time of the thread sleeping...
                sleep = th.sleep - gettime();
        }
        fprintf(stderr, "Thread id : %2d | time : %4ld ms of CPU | time : %4ld sleep | sp : %9p | ", th.id, ms , sleep , th.sp);
        print_state(th);
}

static void
print_threads()
{
        fprintf(stderr, "%s\n", "");
        fprintf(stderr, "%s\n", "Prints Table of threads");
        fprintf(stderr, "%s\n", "=======================");
        fprintf(stderr, "Current thread id running... : %d\n", g_ths.th[g_ths.current_th].id);
        fprintf(stderr, "Position in array of current thread running : %d\n", g_ths.current_th);
        int i;
        for (i = 0; i < MAXTHREADS ; i++) {
                print_thread(g_ths.th[i]);
        }
        fprintf(stderr, "%s\n", "=======================");
        fprintf(stderr, "%s\n", "");
}

static int
is_thread_running_out_time()
{
        return isrunning(g_ths.th[g_ths.current_th]) && !istime(g_ths.th[g_ths.current_th]);
}

static int
f_stack(Thread th)
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

        for (i = 0 ; i < MAXTHREADS ; i++) {
                if (f_stack(g_ths.th[i])) {
                        free(g_ths.th[i].sp);
                        g_ths.th[i].sp = NULL;
                }
        }
}

static int
is_new_th(Thread th , int create)
{
        return th.state == DEAD && create;
}

static int
is_th_rdy(Thread th , int create)
{
        return th.state == READY && !create;
}

static int
istime_sleep (Thread th)
{
        return th.sleep - gettime() < 0;
}

static int
is_wake_up (Thread th , int create)
{
        return th.state == SLEEP && istime_sleep(th) && !create;
}


static int
isnext(Thread th , int create)
{
        return is_new_th(th , create) || is_th_rdy(th , create) || is_wake_up(th , create);
}

static int
th_alone (int create, int next_th)
{
        return !create && next_th < 0 && isrunning(g_ths.th[g_ths.current_th]);
}

void
debug_gen(int next_th , int old_th)
{
        fprintf(stderr, "next_th GENERAL : %d -> %s\n", next_th , "Scheduler");
        fprintf(stderr, "%s\n", "CURRENT THREAD");
        print_thread(g_ths.th[old_th]);
        fprintf(stderr, "%s\n", "");
        if (next_th >= 0) {
                fprintf(stderr, "%s\n", "NEXT THREAD");
                print_thread(g_ths.th[next_th]);
                fprintf(stderr, "%s\n", "");
        } else if (next_th == ERRSUSP) {
                fprintf(stderr, "%s\n", "THERE IS A THREAD SUSPEND");
        } else if (next_th == ERRSLEEP) {
                fprintf(stderr, "%s\n", "THERE IS A THREAD SLEEP");
        } else {
                fprintf(stderr, "%s\n", "NO THREADS READY AND NOT THREADS SUSPEND");
        }
        fprintf(stderr, "%s\n", "");
}

void
debug_alone(int next_th , int old_th)
{
        fprintf(stderr, "next_th ALONE : %d -> %s\n", next_th , "Scheduler.");
        print_thread(g_ths.th[old_th]);
        print_thread(g_ths.th[next_th]);
        fprintf(stderr, "%s\n", "");
}

void
debug_exit(int next_th , int old_th)
{
        fprintf(stderr, "next_th EXIT : %d -> %s\n", next_th , "Scheduler.");
        print_thread(g_ths.th[old_th]);
        fprintf(stderr, "%s\n", "");
}

static void
errsusp()
{
        fprintf(stderr, "%s", "There is no more threads to run ");
        fprintf(stderr, "%s", "we have threads SUSPEND ");
        fprintf(stderr, "%s\n", "error while running thread library.");
}

/*
 * checkexit will exit if next is anyone to get change context by the mechanism.
 * is a special case.
 * All threads suspended or dead... ERRSUSP exit...
 * Anyone thread ready or sleep to get run... normal exit...
 * Thread alone hi has to follow run...
 */
void
exitsch(int create , int next_th)
{
        if (!create && next_th == -1) {
                // There are not more thread to run... The library of
                // collaborative threads have finished succesfull.
                garbage_collector();
                exit(0);
        } else if (!create && next_th == ERRSUSP) {
                // There are not more thrads to run... and there are
                // threads suspended exit with Error status by to
                // suspened threads.
                errsusp();
                garbage_collector();
                exit(1);
        }
}

/*
* scheluder, test the state of follow thread...
* if follow thread is called for createthread scheluder pass a thread DEAD
* or -1 if he can't get the follow thread dead because anyone thread is dead.
*
* if the thread is called by changecontext(), you can some possibilities...
*
* if the follow thread is READY, scheluder return this thread.
* if the follow thread is SUSPEND, appoint this thread like ERRSUSP and follow
* iterate by the array, if I found a follows threads READY, I return this threads, if I
* found a thread SLEEP the next_th is appoint like ERRSLEEP.
*
* if the thread is SLEEP, I appoint the thread like sleep, if I found a thread READY
* return the thread ready else I follow looking for a thread ready or finalizing sleep.
*/
static int
scheluder(int create)
{
        enum {
                FOLLOW = 1,
        };
        int  pos = g_ths.current_th + FOLLOW, next_th = -1, count = 0;

        while (next_th < 0 && count < MAXTHREADS) {
                if (pos >= MAXTHREADS) {
                        // Circular iterative in threads array...
                        pos = 0;
                }
                if (isnext(g_ths.th[pos] , create)) {
                        // Gotten the follow thread to running.
                        next_th = pos;
                } else if (next_th == ERRSLEEP || g_ths.th[pos].state == SLEEP) {
                        // There aren't anyone thread ready, and are threads sleep.
                        // We must test until threads sleeps wake up!...
                        next_th = ERRSLEEP;
                        count = 0;
                } else if (next_th != ERRSLEEP && g_ths.th[pos].state == SUSPEND) {
                        // At least a thread is suspend... and no threads are sleeping or running by now...
                        next_th = ERRSUSP;
                }
                count++;
                pos++;
        }
        if (th_alone(create , next_th)) {
                // This thread is the Only one thread is running anyone thread is ready...
                // We have to put next thread the same thread we were running.
                next_th = g_ths.current_th;
        }
        exitsch(create , next_th);
        return next_th;
}

/*
 * The changecontext() change the context to the follow thread in wait ready to running.
 * If you don't considerade the time to the follow thread changecontext() change the context
 * to the new thread althought the quantum was less than threshold..., changecontext() always
 * change the context to the follow thread ready to running without test nothing.
 */
static void
changecontext()
{
        int next_th = -1 , create = 1;
        int old_th = g_ths.current_th;

        next_th = scheluder(!create); // Politic Round-Robin to next_th...
        if (g_ths.th[old_th].state != DEAD) {
                // Garbage collector when changecontext from a thread only
                // with not exitsthreads or killthread call.
                garbage_collector();
        }
        if (isrunning(g_ths.th[old_th])) {
                // The thread must be running to change the state to ready...
                g_ths.th[old_th].state = READY;
        }
        g_ths.th[next_th].state = RUNNING;
        g_ths.current_th = next_th;
        g_ths.th[next_th].ms = gettime(); // Update time...
        if (swapcontext(&g_ths.th[old_th].uct , &g_ths.th[next_th].uct) < 0) {
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
                g_ths.th[i].id = -1;
                g_ths.th[i].state = DEAD;
        }
        g_ths.current_th = 0;           // Position of thread in data array of threads...
                                        // First thead.
        g_ths.th[g_ths.current_th].id = 0;
        // makecontext... Not necessary main is the flow of the thread 0...
        g_ths.th[g_ths.current_th].state = RUNNING;
        g_ths.th[g_ths.current_th].ms = gettime();
        threads_initialized = 1;    // library threads initializaded...
        count_thread_id = 0;        // Counter of numbers of tid assigned...
        if (getcontext(&(g_ths.th[g_ths.current_th].uct)) < 0) {
                err (1 , "Error getting context: %s" , strerror(errno));
        }
}

void
print_attr_library()
{
        fprintf(stderr, "Counter thread id last number id was assigend : %d\n", count_thread_id);
        fprintf(stderr, "Library threads initialized : %d\n", threads_initialized);
        print_threads();
}

void
debug_create(int next_th)
{
        fprintf(stderr, "%s%d\n", "Next thread to create -> CREATE THREAD : " , next_th);
        if (next_th >= 0) {
                print_thread(g_ths.th[next_th]);
        }
        fprintf(stderr , "%s\n", "");
}

int
createthread(void (*mainf)(void*), void *arg, int stacksize)
{
        char *stack_f = NULL;
        int next_th = -1, create = 1;

        next_th = scheluder(create);
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
        count_thread_id++;
        g_ths.th[next_th].id = count_thread_id;
        g_ths.th[next_th].state = READY;
        g_ths.th[next_th].sp = stack_f;
        g_ths.th[next_th].ms = gettime();
        // Create context to thread.
        if (getcontext(&g_ths.th[next_th].uct) < 0) {
                fprintf(stderr, "Error creating context to the thread : %s\n", strerror(errno));
                return -1;
        }
        g_ths.th[next_th].uct.uc_stack.ss_sp = stack_f;
        g_ths.th[next_th].uct.uc_stack.ss_size = stacksize;
        g_ths.th[next_th].uct.uc_link = NULL;
        makecontext(&g_ths.th[next_th].uct, (void *) mainf, 1 , arg , NULL);
        // New thread was created we return its id.
        return g_ths.th[next_th].id;
}

void
exitsthread()
{
        g_ths.th[g_ths.current_th].id = -1;
        g_ths.th[g_ths.current_th].state = DEAD;
        // The thread was deleted we pass the context
        // to follow thread that have to run.
        changecontext();
}

void
yieldthread()
{
        if (is_thread_running_out_time()) {
                changecontext();
        }
}

int
curidthread()
{
        return g_ths.th[g_ths.current_th].id;
}

void
suspendthread()
{
        g_ths.th[g_ths.current_th].state = SUSPEND;
        changecontext();
}

int
isresume(Thread th , int id)
{
        return th.id == id && th.state == SUSPEND;
}

void
resume(Thread *th)
{
        th->state = READY;
}

int
resumethread(int id)
{
        int found = -1 , pos = 0;

        while (found < 0 && pos < MAXTHREADS) {
                if (isresume(g_ths.th[pos] , id)) {
                        resume(&g_ths.th[pos]);
                        found = 0;
                } else {
                        pos++;
                }
        }
        return found;
}

int
suspendedthreads(int **list)
{
        int i = -1 , nsusp = -1;

        (*list) = malloc(sizeof(int) * MAXTHREADS);
        if (*list == NULL){
                err (1, "Not malloc to vector of suspended threads... : %s" , strerror(errno));
        }
        nsusp = 0;
        for (i = 0; i < MAXTHREADS; i++) {
                if (g_ths.th[i].state == SUSPEND) {
                        (*list)[nsusp] = g_ths.th[i].id;
                        nsusp++;
                 }
        }
        return nsusp;
}

int
iskill(Thread th , int id)
{
        return th.id == id && th.state != DEAD;
}

void
_kill (Thread *th)
{
        th->state = DEAD;
        th->id = -1;
}

int
killthread(int id)
{
        int tid = -1 , pos = 0;

        while (tid < 0 && pos < MAXTHREADS) {
                if(iskill(g_ths.th[pos] , id)){
                        tid = g_ths.th[pos].id;
                        _kill(&g_ths.th[pos]);
                } else {
                        pos++;
                }
        }
        if (id == tid) {
                // I kill me..
                changecontext();
        }
        return tid;
}

void
sleepthread(int msec)
{
        g_ths.th[g_ths.current_th].sleep = gettime() + (long) msec;
        g_ths.th[g_ths.current_th].state = SLEEP;
        changecontext();
}
