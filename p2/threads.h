void initthreads(void);
int createthread(void (*mainf)(void*), void *arg, int stacksize);
void exitsthread(void);
void yieldthread(void);
int curidthread(void);
void print_attr_library(void);

void suspendthread(void);
int resumethread(int id);
int suspendedthreads(int **list);
int killthread(int id);
void sleepthread(int msec);
